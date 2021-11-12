package handlers

import (
	"Task/constants"
	repo "Task/repository"
	"encoding/json"
	"net/http"

	"Task/helpers"

	"github.com/go-chassis/openlog"
	"github.com/gorilla/mux"
)

//
type TaskHandler struct {
	Repo *repo.Repository
}

//
type Response struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

var payloadErrorResp = Response{Status: 500, Msg: constants.UnableToReadPayload}

//
func Recover(rw *http.ResponseWriter) {
	if r := recover(); r != nil {
		writeErrorResp(rw, constants.InternalServerError)
	}
}

//
func writeResp(rw *http.ResponseWriter, resp *Response) {
	(*rw).Header().Set("Content-Type", "application/json")
	(*rw).WriteHeader(resp.Status)
	json.NewEncoder(*rw).Encode(*resp)
	return
}

//
func writeErrorResp(rw *http.ResponseWriter, err string) {
	(*rw).Header().Set("Content-Type", "application/json")
	status := helpers.ErrorStatusMap[err]
	resp := Response{Status: status, Msg: err}
	writeResp(rw, &resp)
	return
}

//
func (t *TaskHandler) Create(rw http.ResponseWriter, req *http.Request) {
	defer Recover(&rw)
	openlog.Info("got request to create record")
	record := make(map[string]interface{})
	err := json.NewDecoder(req.Body).Decode(&record)
	if err != nil {
		openlog.Error(err.Error())
		writeResp(&rw, &payloadErrorResp)
		return
	}
	res, err := t.Repo.CreateRecord(record)
	if err != nil {
		writeErrorResp(&rw, err.Error())
		return
	}
	response := Response{Status: 201, Data: res, Msg: constants.CreationSuccess}
	writeResp(&rw, &response)
	return
}

//
func (t *TaskHandler) FetchByID(rw http.ResponseWriter, req *http.Request) {
	defer Recover(&rw)
	params := mux.Vars(req)
	id := params["id"]
	openlog.Info("got request to fetch record " + id)
	res, err := t.Repo.FetchRecordByID(id)
	if err != nil {
		writeErrorResp(&rw, err.Error())
		return
	}
	response := Response{Status: 200, Data: res, Msg: constants.FetchSuccess}
	writeResp(&rw, &response)
	return
}

//
func (t *TaskHandler) DeleteByID(rw http.ResponseWriter, req *http.Request) {
	defer Recover(&rw)
	params := mux.Vars(req)
	id := params["id"]
	openlog.Info("got request to delete record " + id)
	res, err := t.Repo.DeleteRecordByID(id)
	if err != nil {
		writeErrorResp(&rw, err.Error())
		return
	}
	response := Response{Status: 200, Data: res, Msg: constants.DeletionSuccess}
	writeResp(&rw, &response)
	return
}

//
func (t *TaskHandler) UpdateByID(rw http.ResponseWriter, req *http.Request) {
	defer Recover(&rw)
	params := mux.Vars(req)
	id := params["id"]
	openlog.Info("got request to update record " + id)
	updatePayload := make(map[string]interface{})
	err := json.NewDecoder(req.Body).Decode(&updatePayload)
	if err != nil {
		openlog.Error(err.Error())
		writeResp(&rw, &payloadErrorResp)
		return
	}

	res, err := t.Repo.UpdateRecordByID(id, updatePayload)
	if err != nil {
		writeErrorResp(&rw, err.Error())
		return
	}
	response := Response{Status: 200, Data: res, Msg: constants.UpdateSuccess}
	writeResp(&rw, &response)
	return
}

//
func (t *TaskHandler) FetchAll(rw http.ResponseWriter, req *http.Request) {
	defer Recover(&rw)
	openlog.Info("got request to fetch records")
	page := req.URL.Query().Get("page")
	size := req.URL.Query().Get("size")
	filter := req.URL.Query().Get("filter")
	sort := req.URL.Query().Get("sort")
	res, err := t.Repo.FetchAllRecords(filter, sort, page, size)
	if err != nil {
		writeErrorResp(&rw, err.Error())
		return
	}
	response := Response{Status: 200, Data: res, Msg: constants.FetchAllSuccess}
	writeResp(&rw, &response)
	return
}
