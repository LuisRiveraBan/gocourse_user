package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LuisRiveraBan/go_lib_response/response"
	user "github.com/LuisRiveraBan/gocourse_user/internal"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func NewUserHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {
	r := mux.NewRouter()

	otps := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	// Cambia NewServer por el método ServeHTTP
	r.Methods("POST").Path("/users").Handler(httptransport.NewServer(
		endpoint.Endpoint(endpoints.CreateUser),
		decodeCreateUser,
		encodeResponse,
		otps...,
	))

	r.Methods("GET").Path("/users").Handler(httptransport.NewServer(
		endpoint.Endpoint(endpoints.ListUsers),
		decodeGetAllUser,
		encodeResponse,
		otps...,
	))

	r.Methods("GET").Path("/users/{id}").Handler(httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetUserByID),
		decodeGetUser,
		encodeResponse,
		otps...,
	))

	r.Methods("DELETE").Path("/users/{id}").Handler(httptransport.NewServer(
		endpoint.Endpoint(endpoints.DeleteUser),
		decodeDeleteUser,
		encodeResponse,
		otps...,
	))

	r.Methods("PATCH").Path("/users/{id}").Handler(httptransport.NewServer(
		endpoint.Endpoint(endpoints.UpdateUser),
		decodeUpdateUser,
		encodeResponse,
		otps...,
	))

	/*r.HandleFunc("/users", enpoints.ListUsers).Methods("GET")
	r.HandleFunc("/users/{id}", enpoints.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id}", enpoints.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}", enpoints.UpdateUser).Methods("PUT")*/

	return r
}

func decodeCreateUser(_ context.Context, r *http.Request) (interface{}, error) {
	// Implement the logic to decode the incoming JSON request into a Create request struct
	var req user.Create
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}
	return req, nil
}

func decodeGetUser(_ context.Context, r *http.Request) (interface{}, error) {
	// Implement the logic to decode the incoming request parameters into a GetUser request struct
	p := mux.Vars(r)
	req := user.GetReq{
		ID: p["id"],
	}
	return req, nil
}

func decodeGetAllUser(_ context.Context, r *http.Request) (interface{}, error) {
	// Implement the logic to decode the incoming request parameters into a GetAllUsers request struct
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := user.GetAllReq{
		FirstName: v.Get("first_name"),
		LastName:  v.Get("last_name"),
		Limit:     limit,
		Page:      page,
	}
	return req, nil
}

func decodeUpdateUser(_ context.Context, r *http.Request) (interface{}, error) {
	// Implement the logic to decode the incoming request parameters into a UpdateUser request struct
	var req user.Update

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}
	p := mux.Vars(r)
	req.ID = p["id"]

	return req, nil

}

func decodeDeleteUser(_ context.Context, r *http.Request) (interface{}, error) {
	// Implement the logic to decode the incoming request parameters into a DeleteUser request struct
	p := mux.Vars(r)
	req := user.DeleteReq{
		ID: p["id"],
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, res interface{}) error {
	r := res.(response.Response)
	// Implement the logic to encode the response into JSON format
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	// Implement the logic to encode errors into JSON format
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
