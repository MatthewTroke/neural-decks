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
  id: string;
  type: CardType;
  card_value: string;
}

interface Collection {
  cards: Card[];
}

interface Player {
  score: number;
  role: PlayerRole;
  user_id: string;
  name: string;
  image?: string; // Optional image field
  deck: Card[]; // Array of cards in the player's deck
  is_card_czar: boolean;
  was_card_czar: boolean;
  placed_card?: Card | null; // Nullable placed card
  is_round_winner: boolean;
  is_game_winner: boolean;
}

interface Game {
  id: string;
  name: string;
  collection: Collection;
  winner_count: number;
  max_player_count: number;
  status: GameStatus;
  players: Player[];
  white_cards: Card[];
  black_card: Card | null;
  round_status: RoundStatus;
  round_winner: Player;
  current_game_round: number;
  last_vacated_at: Date | null;
  vacated: boolean;
  created_at: Date;
  updated_at: Date;
}
