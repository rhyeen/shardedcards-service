package handlers

import (
	"context"
	"net/http"

	"github.com/rhyeen/shardedcards-service/internal/api"
	"github.com/rhyeen/shardedcards-service/internal/services/gameservice"
	"github.com/rhyeen/shardedcards-service/internal/stores/storeerror"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/rhyeen/shardedcards-service/internal/models/game"
	"github.com/rhyeen/shardedcards-service/internal/models/gameuser"
)

// GameService see Service for more details
type GameService interface {
	CreateGame(ctx context.Context, params gameservice.CreateGameParams) (game.Game, error)
}

// GameHandler is the handler for game API
type GameHandler struct {
	GameService GameService
}

// CreateGame see Service for more details
func (h GameHandler) CreateGame(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := NewCreateGameRequest(r, p)
	if err != nil {
		api.RespondWith(r, w, http.StatusBadRequest, err, err)
		return
	}
	ctx := r.Context()
	authData, err := api.GetDataFromContext(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, errors.Wrap(err, "failed to get auth data"))
		return
	}
	game, err := h.GameService.CreateGame(ctx, gameservice.CreateGameParams{
		User: gameuser.User{
			ID: authData.UserID,
		},
		Type: request.Type,
	})
	if castErr, ok := err.(*storeerror.DupEntry); ok {
		api.RespondWith(r, w, http.StatusBadRequest, castErr, err)
		return
	}
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, err)
		return
	}
	api.RespondWith(r, w, http.StatusOK, map[string]string{"id": game.ID}, nil)
}
