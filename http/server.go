package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"unicorn"
)

// OrderIDHeader Is the name of the HTTP Header which contains the request id.
// Exported so that it can be changed by developers
var OrderIDHeader = "X-Unicorn-Order-Id"

var (
	ErrNoAmount        = errors.New("no unicorn amount provided for the order")
	ErrInvalidAmount   = errors.New("invalid amount of unicorns")
	ErrOrderIDNotFound = errors.New("coult not find your order ID")
)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type UnicornsResponse struct {
	Pending  int               `json:"pending"`
	OrderID  string            `json:"orderId,omitempty"`
	Unicorns []unicorn.Unicorn `json:"unicorns,omitempty"`
}

func HandleGetUnicorns(svc unicorn.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getOrderID(r)

		if id == "" {
			HandleNewOrder(svc)(w, r)
			return
		}

		if !svc.Validate(id) {
			raise(w, ErrOrderIDNotFound, http.StatusNotFound)
			return
		}

		pending, unicorns := svc.Pool(id)

		response := UnicornsResponse{
			Pending:  pending,
			OrderID:  string(id),
			Unicorns: unicorns,
		}

		reply(w, http.StatusOK, &response)
	}
}

func HandleNewOrder(svc unicorn.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		amount, err := getAmount(r)
		if err != nil {
			raise(w, err, http.StatusBadRequest)
			return
		}

		id, err := svc.OrderUnicorns(amount)
		if err != nil {
			raise(w, fmt.Errorf("could not order unicorns: %w", err), http.StatusServiceUnavailable)
			return
		}

		setOrderID(w, id)

		pending, unicorns := svc.Pool(id)

		response := UnicornsResponse{
			Pending:  pending,
			OrderID:  string(id),
			Unicorns: unicorns,
		}

		reply(w, http.StatusOK, &response)
	}
}

func getAmount(r *http.Request) (int, error) {
	s := r.URL.Query().Get("amount")
	if s == "" {
		return 0, ErrNoAmount
	}

	amount, err := strconv.Atoi(s)
	if err != nil {
		return 0, ErrInvalidAmount
	}

	return amount, nil
}

// setOrderID writes the order ID to the responde headers.
func setOrderID(w http.ResponseWriter, id unicorn.OrderID) {
	if id == "" {
		return
	}

	w.Header().Add(OrderIDHeader, string(id))
}

// getOrderID retrives an order ID from the headers.
func getOrderID(r *http.Request) unicorn.OrderID {
	id := r.Header.Get(OrderIDHeader)
	if id == "" {
		return ""
	}

	return unicorn.OrderID(id)
}

// reply replies to the resquest with the body and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
func reply(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(body)
}

// raise replies to the request with the service error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
func raise(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(&ErrorResponse{
		Error: err.Error(),
	})
}
