import GameCard from "@/components/game/GameCard";
import { useAuth } from "@/context/AuthContext";
import { cn } from "@/lib/utils";
import { useState } from "react";

export default function Hand(props: {
  game: Game;
  handlePlayCard: (card: Card) => void;
}) {
  const [selectedCardId, setSelectedCardId] = useState<string | undefined>(
    undefined
  );
  const { user } = useAuth();

  const player = props.game.Players.find((p) => p.UserID === user?.user_id);

  if (!player) {
    return <div>no cards!</div>;
  }

  if (props.game.Status === "Setup") {
    return <div>Waiting for game to start...</div>;
  }

  const disabled = Boolean(player.PlacedCard);

  const onCardClick = (cardId: string) => {
    const card = player.Deck.find((card: Card) => card.ID === cardId);

    setSelectedCardId(card?.ID);

    if (card) {
      props.handlePlayCard(card);
    }
  };

  const rotations = [
    "rotate-300",
    "rotate-320",
    "rotate-340",
    "rotate-0",
    "rotate-20",
    "rotate-40",
    "rotate-60",
  ];

  return (
    <div className="relative">
      <h3 className="block md:hidden text-lg font-semibold mb-3">Your Hand</h3>
      <div className="relative md:hidden flex flex-wrap gap-4 justify-center sm:justify-start">
        {player.Deck?.map((card) => (
          <>
            <GameCard
              key={card.ID}
              onCardClick={onCardClick}
              cardId={card.ID}
              value={card.CardValue}
              isDisabled={disabled}
              selected={selectedCardId === card.ID}
            />
          </>
        ))}
      </div>

      <div className="hidden md:block fixed bottom-[500px] left-1/2">
        <div className="rotate-[-5deg] -ml-20">
          {player.Deck?.map((card, index) => (
            <div
              className={cn(
                "absolute",
                "w-[144px]",
                "h-[192px]",
                "origin-[40px_500px]",
                "hover:z-10",
                rotations[index]
              )}
            >
              <GameCard
                key={card.ID}
                onCardClick={onCardClick}
                cardId={card.ID}
                value={card.CardValue}
                isDisabled={disabled}
                selected={selectedCardId === card.ID}
              />
            </div>
          ))}
        </div>
      </div>

      {disabled && (
        <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
          <span className="text-white text-xl font-bold">
            Card Czar Picking
          </span>
        </div>
      )}
    </div>
  );
}
