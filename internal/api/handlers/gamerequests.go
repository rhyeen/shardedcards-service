package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// CreateGameRequest parameters from the CreateGame call
type CreateGameRequest struct {
	Type string `json:"type"`
}

// NewCreateGameRequest extracts the CreateGameRequest
func NewCreateGameRequest(r *http.Request, p httprouter.Params) (CreateGameRequest, error) {
	var request CreateGameRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return request, errors.New("invalid request")
	}
	return request.validate()
}

func (request CreateGameRequest) validate() (CreateGameRequest, error) {
	if request.Type == "" {
		return request, errors.New("must provide type")
	}
	return request, nil
}
