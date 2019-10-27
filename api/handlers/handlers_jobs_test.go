package handlers

import (
	"net/http"
	"testing"
	"time"

	hu "github.com/swinslow/peridot-api/test/handlerutils"
	"github.com/swinslow/peridot-db/pkg/datastore"
)

// ===== GET /repopulls/2/jobs =====

func TestCanGetJobsSubHandlerAsViewer(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repopulls/2/jobs", "", "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.jobsSubHandler), "/repopulls/{id}/jobs")
	hu.ConfirmOKResponse(t, rec)

	wanted := `{"jobs": [
		{"id":2, "repopull_id":2, "agent_id":4, "started_at":"2019-05-02T14:07:00Z", "finished_at":"2019-05-02T14:07:30Z", "status":"stopped", "health":"ok", "output":"successfully retrieved repo", "is_ready":true, "config":{}},
		{"id":5, "repopull_id":2, "agent_id":1, "started_at":"2019-05-02T14:07:00Z", "finished_at":"2019-05-02T14:08:00Z", "status":"stopped", "health":"ok", "output":"found 57 files with short-form license IDs in 182 files", "is_ready":true, "config":{}},
		{"id":6, "repopull_id":2, "agent_id":1, "priorjob_ids":[5], "started_at":"2019-05-02T14:09:00Z", "finished_at":"2019-05-02T14:09:10Z", "status":"stopped", "health":"ok", "output":"wrote attributions", "is_ready":true, "config":{}},
		{"id":7, "repopull_id":2, "agent_id":5, "started_at":"2019-05-02T14:09:30Z", "finished_at":"0001-01-01T00:00:00Z", "status":"running", "health":"degraded", "output":"unable to retrieve some dependencies", "is_ready":true, "config":{}},
		{"id":8, "repopull_id":2, "agent_id":6, "priorjob_ids":[5, 7], "started_at":"0001-01-01T00:00:00Z", "finished_at":"0001-01-01T00:00:00Z", "status":"startup", "health":"ok", "is_ready":true, "config":{"kv": {"prefer": "primary"}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 7}}}}
	]}`
	hu.CheckResponse(t, rec, wanted)
}

func TestCannotGetJobsSubHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "GET", "/repopulls/2/jobs", ``, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.jobsSubHandler), "/repopulls/{id}/jobs")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "GET", "/repopulls/2/jobs", ``, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.jobsSubHandler), "/repopulls/{id}/jobs")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

// ===== POST /repopulls/2/jobs =====

func TestCanPostJobsSubHandlerAsOperator(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/repopulls/3/jobs", `{"agent_id": 5, "priorjob_ids": [3], "is_ready":false, "config": {"kv": {"hello":"world"}, "codereader": {"godeps": {"priorjob_id": 7}}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 7}}}}`, "operator")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.jobsSubHandler), "/repopulls/{id}/jobs")
	//log.Printf("rec: %#v\n", rec)
	hu.ConfirmCreatedResponse(t, rec)

	wanted := `{"id": 9}`
	hu.CheckResponse(t, rec, wanted)

	// and verify state of database now
	jobs, err := env.db.GetAllJobsForRepoPull(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(jobs) != 2 {
		t.Errorf("expected %d, got %d", 2, len(jobs))
	}
	newJob, err := env.db.GetJobByID(9)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	wantedJob := &datastore.Job{
		ID:          9,
		RepoPullID:  3,
		AgentID:     5,
		PriorJobIDs: []uint32{3},
		StartedAt:   time.Time{},
		FinishedAt:  time.Time{},
		Status:      datastore.StatusStartup,
		Health:      datastore.HealthOK,
		Output:      "",
		IsReady:     false,
		Config: datastore.JobConfig{
			KV: map[string]string{"hello": "world"},
			CodeReader: map[string]datastore.JobPathConfig{
				"godeps": datastore.JobPathConfig{PriorJobID: 7},
			},
			SpdxReader: map[string]datastore.JobPathConfig{
				"primary": datastore.JobPathConfig{Value: "/path/wherever"},
				"godeps":  datastore.JobPathConfig{PriorJobID: 7},
			},
		},
	}
	helperCompareJobs(t, wantedJob, newJob)
}

func TestCannotPostJobsSubHandlerAsOtherUser(t *testing.T) {
	// as commenter
	rec, req, env := setupTestEnv(t, "POST", "/repopulls/3/jobs", `{"agent_id": 5, "priorjob_ids": [3], "is_ready":false, "config": {"kv": {"hello":"world"}, "codereader": {"godeps": {"priorjob_id": 7}}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 7}}`, "commenter")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.jobsSubHandler), "/repopulls/{id}/jobs")
	hu.ConfirmAccessDenied(t, rec)

	// as viewer
	rec, req, env = setupTestEnv(t, "POST", "/repopulls/3/jobs", `{"agent_id": 5, "priorjob_ids": [3], "is_ready":false, "config": {"kv": {"hello":"world"}, "codereader": {"godeps": {"priorjob_id": 7}}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 7}}`, "viewer")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.jobsSubHandler), "/repopulls/{id}/jobs")
	hu.ConfirmAccessDenied(t, rec)
}

func TestCannotPostJobsSubHandlerAsBadUser(t *testing.T) {
	rec, req, env := setupTestEnv(t, "POST", "/repopulls/3/jobs", `{"agent_id": 5, "priorjob_ids": [3], "is_ready":false, "config": {"kv": {"hello":"world"}, "codereader": {"godeps": {"priorjob_id": 7}}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 7}}`, "disabled")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.jobsSubHandler), "/repopulls/{id}/jobs")
	hu.ConfirmAccessDenied(t, rec)

	rec, req, env = setupTestEnv(t, "POST", "/repopulls/3/jobs", `{"agent_id": 5, "priorjob_ids": [3], "is_ready":false, "config": {"kv": {"hello":"world"}, "codereader": {"godeps": {"priorjob_id": 7}}, "spdxreader": {"primary": {"path": "/path/wherever"}, "godeps": {"priorjob_id": 7}}`, "invalid")
	hu.ServeHandler(rec, req, http.HandlerFunc(env.jobsSubHandler), "/repopulls/{id}/jobs")
	hu.ConfirmInvalidAuth(t, rec, ErrAuthGithub)
}

func helperCompareJobs(t *testing.T, expected *datastore.Job, got *datastore.Job) {
	if expected.ID != got.ID {
		t.Errorf("expected %#v, got %#v", expected.ID, got.ID)
	}

	if expected.RepoPullID != got.RepoPullID {
		t.Errorf("expected %#v, got %#v", expected.RepoPullID, got.RepoPullID)
	}

	if expected.AgentID != got.AgentID {
		t.Errorf("expected %#v, got %#v", expected.AgentID, got.AgentID)
	}

	if len(expected.PriorJobIDs) != len(got.PriorJobIDs) {
		t.Errorf("expected %#v, got %#v", len(expected.PriorJobIDs), len(got.PriorJobIDs))
	} else {
		for i := range expected.PriorJobIDs {
			if expected.PriorJobIDs[i] != got.PriorJobIDs[i] {
				t.Errorf("for index %d, expected %#v, got %#v", i, expected.PriorJobIDs[i], got.PriorJobIDs[i])
			}
		}
	}

	if expected.StartedAt != got.StartedAt {
		t.Errorf("expected %#v, got %#v", expected.StartedAt, got.StartedAt)
	}

	if expected.FinishedAt != got.FinishedAt {
		t.Errorf("expected %#v, got %#v", expected.FinishedAt, got.FinishedAt)
	}

	if expected.Status != got.Status {
		t.Errorf("expected %#v, got %#v", expected.Status, got.Status)
	}

	if expected.Health != got.Health {
		t.Errorf("expected %#v, got %#v", expected.Health, got.Health)
	}

	if expected.Output != got.Output {
		t.Errorf("expected %#v, got %#v", expected.Output, got.Output)
	}

	if expected.IsReady != got.IsReady {
		t.Errorf("expected %#v, got %#v", expected.IsReady, got.IsReady)
	}

	if len(expected.Config.KV) != len(got.Config.KV) {
		t.Errorf("expected %#v, got %#v", len(expected.Config.KV), len(got.Config.KV))
	} else {
		for kExp, vExp := range expected.Config.KV {
			vGot, ok := got.Config.KV[kExp]
			if !ok {
				t.Errorf("key %v in expected, not in got", kExp)
			} else {
				if vExp != vGot {
					t.Errorf("expected %#v, got %#v", vExp, vGot)
				}
			}
		}
		for kGot := range got.Config.KV {
			_, ok := expected.Config.KV[kGot]
			if !ok {
				t.Errorf("key %v in got, not in expected", kGot)
			}
		}
	}

	if len(expected.Config.CodeReader) != len(got.Config.CodeReader) {
		t.Errorf("expected %#v, got %#v", len(expected.Config.CodeReader), len(got.Config.CodeReader))
	} else {
		for kExp, vExp := range expected.Config.CodeReader {
			vGot, ok := got.Config.CodeReader[kExp]
			if !ok {
				t.Errorf("key %v in expected, not in got", kExp)
			} else {
				if vExp.Value != vGot.Value {
					t.Errorf("expected %#v, got %#v", vExp.Value, vGot.Value)
				}
				if vExp.PriorJobID != vGot.PriorJobID {
					t.Errorf("expected %#v, got %#v", vExp.PriorJobID, vGot.PriorJobID)
				}
			}
		}
		for kGot := range got.Config.CodeReader {
			_, ok := expected.Config.CodeReader[kGot]
			if !ok {
				t.Errorf("key %v in got, not in expected", kGot)
			}
		}
	}

	if len(expected.Config.SpdxReader) != len(got.Config.SpdxReader) {
		t.Errorf("expected %#v, got %#v", len(expected.Config.SpdxReader), len(got.Config.SpdxReader))
	} else {
		for kExp, vExp := range expected.Config.SpdxReader {
			vGot, ok := got.Config.SpdxReader[kExp]
			if !ok {
				t.Errorf("key %v in expected, not in got", kExp)
			} else {
				if vExp.Value != vGot.Value {
					t.Errorf("expected %#v, got %#v", vExp.Value, vGot.Value)
				}
				if vExp.PriorJobID != vGot.PriorJobID {
					t.Errorf("expected %#v, got %#v", vExp.PriorJobID, vGot.PriorJobID)
				}
			}
		}
		for kGot := range got.Config.SpdxReader {
			_, ok := expected.Config.SpdxReader[kGot]
			if !ok {
				t.Errorf("key %v in got, not in expected", kGot)
			}
		}
	}
}
