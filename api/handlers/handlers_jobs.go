// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/swinslow/peridot-db/pkg/datastore"
)

// ========== HANDLER for /repopulls/{id}/jobs

func (env *Env) jobsSubHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// we only take GET or POST requests
	switch r.Method {
	case "GET":
		env.jobsSubGetHelper(w, r)
	case "POST":
		env.jobsSubPostHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) jobsSubGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least viewer
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access; get repopull id from vars
	repopullID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid repopull ID"}`)
		return
	}

	// get jobs from database
	jobs, err := env.db.GetAllJobsForRepoPull(repopullID)
	if err != nil {
		fmt.Fprintf(w, `{"error": "Database retrieval error"}`)
		return
	}

	// create map so we return a JSON object
	jobsMap := map[string][]*datastore.Job{}
	jobsMap["jobs"] = jobs
	js, err := json.Marshal(jobsMap)
	if err != nil {
		fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) jobsSubPostHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	// must be at least operator
	user := extractUser(w, r, datastore.AccessOperator)
	if user == nil {
		return
	}

	// sufficient access; get repopull id from vars
	repopullID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid repopull ID"}`)
		return
	}

	// parse JSON request
	js := map[string]interface{}{}
	err = json.NewDecoder(r.Body).Decode(&js)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid JSON request"}`)
		return
	}

	// and extract data
	agentID, ok := js["agent_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'agent_id'"}`)
		return
	}
	priorJobIDs := []uint32{}
	pjids, ok := js["priorjob_ids"]
	if ok {
		for _, pjid := range pjids.([]interface{}) {
			pjf, ok2 := pjid.(float64)
			if !ok2 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error": "Invalid non-integer value for 'priorjob_ids'"}`)
				return
			}

			priorJobIDs = append(priorJobIDs, uint32(pjf))
		}
	} else {
		priorJobIDs = []uint32{}
	}

	// not looking for startedAt, finishedAt, status, health or output,
	// because the API user cannot set those

	// we WILL look for is_ready, but it will be a second call to the database
	var isReady bool
	isReadyStr, ok := js["is_ready"]
	if ok {
		isReady = isReadyStr.(bool)
	}

	config, ok := js["config"]
	// expect to be present, at least as an empty object
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing required value for 'config'"}`)
		return
	}
	jcfg, err := env.helperParseJobConfig(config)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Error parsing value for 'config': %v"}`, err)
		return
	}

	// finally, add the new job
	newID, err := env.db.AddJobWithConfigs(repopullID, uint32(agentID.(float64)), priorJobIDs, jcfg.KV, jcfg.CodeReader, jcfg.SpdxReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Unable to create job: %v"}`, err)
		return
	}

	// and if is_ready was set, do the update call too
	if isReady {
		err := env.db.UpdateJobIsReady(newID, true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "Created job with ID %d but unable to set job as ready"}`, newID)
			return
		}
	}

	// success!
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": %d}`, newID)
}

// ========== HANDLER for /jobs/{id}

func (env *Env) jobsOneHandler(w http.ResponseWriter, r *http.Request) {
	// responses will be JSON format
	w.Header().Set("Content-Type", "application/json")

	// check valid request types
	switch r.Method {
	case "GET":
		env.jobsOneGetHelper(w, r)
	case "PUT":
		env.jobsOnePutHelper(w, r)
	case "DELETE":
		env.jobsOneDeleteHelper(w, r)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (env *Env) jobsOneGetHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessViewer)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	jobID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid ID"}`)
		return
	}

	// get job from database
	job, err := env.db.GetJobByID(jobID)
	if err != nil {
		fmt.Fprintf(w, `{"error": "Database retrieval error"}`)
		return
	}

	// create map so we return a JSON object
	jsData := struct {
		Job *datastore.Job `json:"job"`
	}{Job: job}
	js, err := json.Marshal(jsData)
	if err != nil {
		fmt.Fprintf(w, `{"error": "JSON marshalling error"}`)
		return
	}
	w.Write(js)
}

func (env *Env) jobsOnePutHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessOperator)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	jobID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid ID"}`)
		return
	}

	// check job exists in database
	_, err = env.db.GetJobByID(jobID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "Unknown job ID"}`)
		return
	}

	// parse JSON request
	js := map[string]interface{}{}
	err = json.NewDecoder(r.Body).Decode(&js)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid JSON request"}`)
		return
	}

	// and extract data; currently, can only update is_ready
	newIsReadyStr, ok := js["is_ready"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "No new value specified for is_ready"}`)
		return
	}
	newIsReady := newIsReadyStr.(bool)

	// modify the job data
	err = env.db.UpdateJobIsReady(jobID, newIsReady)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Unable to update job"}`)
		return
	}

	// success!
	w.WriteHeader(http.StatusNoContent)
}

func (env *Env) jobsOneDeleteHelper(w http.ResponseWriter, r *http.Request) {
	// get user and check access level
	user := extractUser(w, r, datastore.AccessAdmin)
	if user == nil {
		return
	}

	// sufficient access
	// extract ID for request
	jobID, err := extractIDasU32(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Missing or invalid ID"}`)
		return
	}

	// delete the job
	err = env.db.DeleteJob(jobID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Unable to delete job"}`)
		return
	}

	// success!
	w.WriteHeader(http.StatusNoContent)
}

// ========== HELPERS for parsing job details

func (env *Env) helperParseJobConfig(config interface{}) (*datastore.JobConfig, error) {
	jcfg := &datastore.JobConfig{KV: map[string]string{}}

	// and extract config sub-values
	configCast := config.(map[string]interface{})
	configKV, ok := configCast["kv"]
	if ok {
		for k, v := range configKV.(map[string]interface{}) {
			jcfg.KV[k] = v.(string)
		}
	}

	var err error
	jcfg.CodeReader, err = env.helperParseJobPathConfigs(configCast, "codereader")
	if err != nil {
		return nil, err
	}

	jcfg.SpdxReader, err = env.helperParseJobPathConfigs(configCast, "spdxreader")
	if err != nil {
		return nil, err
	}

	return jcfg, nil
}

func (env *Env) helperParseJobPathConfigs(configCast map[string]interface{}, which string) (map[string]datastore.JobPathConfig, error) {
	configReader, ok := configCast[which]
	if !ok {
		// that's fine
		return nil, nil
	}

	cr := map[string]datastore.JobPathConfig{}
	for k, v := range configReader.(map[string]interface{}) {
		// k is unique key for this codereader config record
		// v is JSON object with one key, which should be either "path" or "priorjob_id"
		recordCast, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Error parsing value for '%s' with key %s", which, k)
		}
		// we'll do a range but there should only be one key
		seenFirst := false
		var path string
		var priorJobID uint32
		for subk, subv := range recordCast {
			if seenFirst == true {
				return nil, fmt.Errorf("More than one sub-key present for '%s' with key %s", which, k)
			}
			switch subk {
			case "path":
				path, ok = subv.(string)
				if !ok {
					return nil, fmt.Errorf("Invalid non-string value for '%s' with key 'path'", which)
				}
			case "priorjob_id":
				pjidCast, ok := subv.(float64)
				if !ok {
					return nil, fmt.Errorf("Invalid non-integer value for '%s' with key 'priorjob_id'", which)
				}
				priorJobID = uint32(pjidCast)
			default:
				return nil, fmt.Errorf("Invalid sub-key '%s' for '%s' with key %s", subk, which, k)
			}
			seenFirst = true
		}

		// now actually create and add the jobPathConfig record
		cr[k] = datastore.JobPathConfig{Value: path, PriorJobID: priorJobID}
	}

	return cr, nil
}
