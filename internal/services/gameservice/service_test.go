package gameservice

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/rhyeen/shardedcards-service/internal/models/game"

	"github.com/stretchr/testify/require"
)

var gameService GameService
var ctx context.Context

func TestMain(m *testing.M) {
	gameService = GameService{}
	ctx = context.Background()
	result := m.Run()
	os.Exit(result)
}
func TestCreateGame(t *testing.T) {
	cases := []struct {
		name         string
		params       CreateGameParams
		expectedErr  error
		expectedGame game.Game
	}{}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			game, err := gameService.CreateGame(ctx, tc.params)
			errExpected := testErrorAgainstCase(t, err, tc.expectedErr)
			if errExpected {
				return
			}
			require.Equal(t, game, tc.expectedGame)
		})
	}
}

// returns true if tcErr was expected
func testErrorAgainstCase(t *testing.T, err error, tcErr error) bool {
	if tcErr != nil {
		require.EqualError(t, err, tcErr.Error())
		return true
	}
	require.NoError(t, err)
	return false
}
