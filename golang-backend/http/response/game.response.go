package response

import "cardgame/domain/aggregates"

type GetGamesResponse []*aggregates.Game
type GetGameResponse *aggregates.Game
