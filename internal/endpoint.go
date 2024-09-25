package user

import (
	"context"
	"errors"
	"github.com/LuisRiveraBan/go_lib_response/response"
	"github.com/LuisRiveraBan/gocourse_meta/meta"
	/*"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"*/)

type (
	// Controller is a function that handles HTTP requests and writes responses
	Controller func(ctx context.Context, request interface{}) (interface{}, error)
	// Service is the interface that defines the methods required by the service
	// Endpoints represent a collection of functions that implement the service's endpoints'
	Endpoints struct {
		GetUserByID Controller
		CreateUser  Controller
		UpdateUser  Controller
		DeleteUser  Controller
		ListUsers   Controller
		//... more endpoints...
		// ... add more methods to this interface...
	}

	// User is a simple struct representing a user
	Create struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	GetReq struct {
		ID string
	}

	DeleteReq struct {
		ID string
	}

	GetAllReq struct {
		Page      int
		Limit     int
		FirstName string
		LastName  string
	}

	Update struct {
		ID        string
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}
	// ErrResponse is a simple struct to represent an error response
	ErrResponse struct {
		Message string `json:"message"`
	}

	Response struct {
		Status int `json:"status"`
		//omitempty con eso le decimos que si viene vacio lo omita
		Data interface{} `json:"data,omitempty"`
		Err  string      `json:"error,omitempty"`
		Meta *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {

	// Implement the Controller functions for each endpoint here
	// For example:
	// endpoints.GetUserByID = func(w http.ResponseWriter, r *http.Request) {
	//     // Implement logic to get user by ID
	//     //...
	// }

	//... more endpoints...
	//... add more Controller functions...
	return Endpoints{
		//... more endpoints...
		ListUsers:   makeListEndpoint(s, config),
		GetUserByID: makeToFindEndpoint(s),
		CreateUser:  makeCreateEndpoint(s),
		UpdateUser:  makeUpdateEndpoint(s),
		DeleteUser:  makeDeleteEndpoint(s),
		//... more endpoints...
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(Create)

		if req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email, req.Phone)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", user, nil), nil
	}
}

func makeListEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Extract query parameters from the request URL
		req := request.(GetAllReq)

		filters := Filters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		count, err := s.Count(ctx, filters)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		meta, err := meta.NewMeta(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		users, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", users, meta), nil
	}
}
func makeToFindEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetReq)

		user, err := s.Get(ctx, req.ID)

		if err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.NotFound(err.Error())
		}
		return response.OK("sucess", user, nil), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Update)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}
		err := s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email, req.Phone)

		if err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteReq)

		err := s.Delete(ctx, req.ID)

		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", nil, nil), nil

	}
}
