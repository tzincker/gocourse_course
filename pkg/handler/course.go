package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tzincker/go_lib_response/response"
	"github.com/tzincker/gocourse_course/internal/course"
)

func NewCourseHTTPServer(ctx context.Context, endpoints course.Endpoints) http.Handler {

	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateCourse,
		encodeResponse,
		opts...,
	)).Methods("POST")

	router.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllCourses,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	router.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	return router
}

func decodeCreateCourse(_ context.Context, r *http.Request) (any, error) {
	var req course.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetAllCourses(_ context.Context, r *http.Request) (any, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := course.GetAllReq{
		Name:  v.Get("name"),
		Limit: limit,
		Page:  page,
	}

	return req, nil
}

func decodeGetCourse(_ context.Context, r *http.Request) (any, error) {
	p := mux.Vars(r)
	req := course.GetReq{
		ID: p["id"],
	}

	return req, nil
}

func decodeUpdateCourse(_ context.Context, r *http.Request) (any, error) {
	p := mux.Vars(r)
	id := p["id"]

	var req course.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.ID = id
	return req, nil
}

func decodeDeleteCourse(_ context.Context, r *http.Request) (any, error) {
	p := mux.Vars(r)
	req := course.DeleteReq{
		ID: p["id"],
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp any) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utd-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utd-8")

	resp, ok := err.(response.Response)

	if !ok {
		newResponse := response.BadRequest("error parsing body")
		w.WriteHeader(newResponse.StatusCode())
		_ = json.NewEncoder(w).Encode(newResponse)
		return
	}

	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)

}
