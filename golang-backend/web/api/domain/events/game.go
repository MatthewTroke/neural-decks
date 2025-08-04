package events

type GameEventType string

const (
	EventGameBegins               GameEventType = "GameBegins"
	EventJoinedGame               GameEventType = "JoinedGame"
	EventCardPlayed               GameEventType = "CardPlayed"
	EventRoundContinued           GameEventType = "RoundContinued"
	EventCardCzarChoseWinningCard GameEventType = "CardCzarChoseWinningCard"
	EventShuffle                  GameEventType = "Shuffle"
	EventDealCards                GameEventType = "DealCards"
	EventDrawBlackCard            GameEventType = "DrawBlackCard"
	EventSetCardCzar              GameEventType = "SetCardCzar"
	EventTimerUpdate              GameEventType = "TimerUpdate"
	EventGameWinner               GameEventType = "GameWinner"
	EventClockUpdate              GameEventType = "ClockUpdate"
	EventEmojiClicked             GameEventType = "EmojiClicked"
)
