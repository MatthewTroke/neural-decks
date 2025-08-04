package response

import "cardgame/internal/domain/aggregates"

type GetGamesResponse []*aggregates.Game
type GetGameResponse *aggregates.Game
