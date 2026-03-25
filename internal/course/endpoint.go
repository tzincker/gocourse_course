package course

import (
	"context"
	"errors"

	"github.com/tzincker/go_lib_response/response"
	"github.com/tzincker/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request any) (any, error)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	GetReq struct {
		ID string
	}

	GetAllReq struct {
		Name  string
		Limit int
		Page  int
	}

	UpdateReq struct {
		ID        string
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	DeleteReq struct {
		ID string
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {

		req := request.(CreateReq)

		if req.Name == "" {
			return nil, response.BadRequest("name is required")
		}

		if req.StartDate == "" {
			return nil, response.BadRequest("start_date is required")
		}

		if req.EndDate == "" {
			return nil, response.BadRequest("end_date is required")
		}

		course, err := s.Create(ctx, req.Name, req.StartDate, req.EndDate)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", course, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetReq)

		course, err := s.Get(ctx, req.ID)
		if err != nil {
			if _, ok := errors.AsType[*ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", course, nil), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetAllReq)
		filters := Filters{
			Name: req.Name,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		courses, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", courses, meta), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {

		req := request.(UpdateReq)

		err := s.Update(ctx, req.ID, &req.Name, &req.StartDate, &req.EndDate)

		if err != nil {
			if _, ok := errors.AsType[*ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(DeleteReq)
		err := s.Delete(ctx, req.ID)

		if err != nil {
			if _, ok := errors.AsType[*ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}
