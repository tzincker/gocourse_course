package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/tzincker/go_lib_response/response"
	"github.com/tzincker/gocourse_course/internal/course"
)

func NewCourseHTTPServer(ctx context.Context, endpoints course.Endpoints) http.Handler {

	router := gin.Default()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.POST("/courses", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateCourse,
		encodeResponse,
		opts...,
	)))

	router.GET("/courses", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllCourses,
		encodeResponse,
		opts...,
	)))

	router.GET("/courses/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse,
		encodeResponse,
		opts...,
	)))

	router.PATCH("/courses/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse,
		encodeResponse,
		opts...,
	)))

	router.DELETE("/courses/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse,
		encodeResponse,
		opts...,
	)))

	return router
}

func ginDecode(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "params", c.Params)
	c.Request = c.Request.WithContext(ctx)
}

func decodeCreateCourse(_ context.Context, r *http.Request) (any, error) {
	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	var req course.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetAllCourses(_ context.Context, r *http.Request) (any, error) {
	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

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

func decodeGetCourse(ctx context.Context, r *http.Request) (any, error) {
	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	params := ctx.Value("params").(gin.Params)
	req := course.GetReq{
		ID: params.ByName("id"),
	}

	return req, nil
}

func decodeUpdateCourse(ctx context.Context, r *http.Request) (any, error) {
	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	params := ctx.Value("params").(gin.Params)
	id := params.ByName("id")

	var req course.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.ID = id
	return req, nil
}

func decodeDeleteCourse(ctx context.Context, r *http.Request) (any, error) {
	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	params := ctx.Value("params").(gin.Params)
	req := course.DeleteReq{
		ID: params.ByName("id"),
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

func authorization(token string) error {
	if token != os.Getenv("TOKEN") {
		return errors.New("invalid token")
	}
	return nil
}
