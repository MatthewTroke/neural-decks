type CardType = "Black" | "White";

type RoundStatus =
  | "Waiting"
  | "PlayersPickingCard"
  | "CardCzarPickingWinningCard"
  | "CardCzarChoseWinningCard";

type GameStatus = "Setup" | "InProgress" | "Finished";

interface User {
  name: string;
  email: string;
  user_id: string;
  iss: string;
  sub: string;
  aud: string[];
  exp: number;
  iat: number;
}

interface Card {
  ID: string;
  Type: CardType;
  CardValue: string;
}

interface Collection {
  Cards: Card[];
}

interface Player {
  Score: number;
  Role: PlayerRole;
  UserID: string;
  Name: string;
  Image?: string; // Optional image field
  Deck: Card[]; // Array of cards in the player's deck
  IsCardCzar: boolean;
  WasCardCzar: boolean;
  PlacedCard?: Card | null; // Nullable placed card
  IsRoundWinner: boolean;
  IsGameWinner: boolean;
}

interface Game {
  ID: string;
  Name: string;
  Collection: Collection;
  WinnerCount: number;
  MaxPlayerCount: number;
  Status: GameStatus;
  Players: Player[];
  WhiteCards: Card[];
  BlackCard: Card | null;
  RoundStatus: RoundStatus;
  RoundWinner: Player;
  CurrentGameRound: number;
  LastVacatedAt: Date | null;
  Vacated: boolean;
  CreatedAt: Date;
  UpdatedAt: Date;
}
