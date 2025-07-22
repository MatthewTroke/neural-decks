import { useAuth } from "@/context/AuthContext";
import GameCard from "@/components/game/GameCard";

export default function Hand(props: {
  game: Game;
  handlePlayCard: (card: Card) => void;
}) {
  let { user } = useAuth();

  let player = props.game.Players.find((p: any) => p.UserID === user?.user_id);

  if (!player) {
    return <div>no cards!</div>;
  }

  if (props.game.Status === "Setup") {
    return <div>Waiting for game to start...</div>;
  }

  let disabled = Boolean(player.PlacedCard);

  let onCardClick = (cardId: string) => {
    let card = player.Deck.find((card: Card) => card.ID === cardId);

    if (card) {
      props.handlePlayCard(card);
    }
  };

  return (
    <div className="relative">
      <h3 className="text-lg font-semibold mb-3">Your Hand</h3>
      <div className="flex flex-wrap gap-4 justify-center sm:justify-start">

        {player.Deck.map((card: any) => (
          <GameCard key={card.ID} onCardClick={onCardClick} cardId={card.ID} value={card.CardValue} isDisabled={disabled} />
        ))}
      </div>
      {disabled && (
        <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
          <span className="text-white text-xl font-bold">Card Czar Picking</span>
        </div>
      )}
    </div>
  );
}
