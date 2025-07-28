import GameCard from "@/components/game/GameCard";
import { useAuth } from "@/context/AuthContext";
import { cn } from "@/lib/utils";
import { Check, X } from "lucide-react";
import { useState } from "react";
import { Button } from "../ui/button";

export default function Hand(props: {
  game: Game;
  handlePlayCard: (card: Card) => void;
}) {
  const [selectedCardId, setSelectedCardId] = useState<string | undefined>(
    undefined
  );
  const { user } = useAuth();

  const player = props.game.players.find((p) => p.user_id === user?.user_id);
  const playerIsCardCzar =
    props.game.players.find((p) => p.is_card_czar)?.user_id === player?.user_id;

  if (!player) {
    return <div>no cards!</div>;
  }

  if (props.game.status === "Setup") {
    return <div>Waiting for game to start...</div>;
  }

  const disabled = Boolean(player.placed_card) || playerIsCardCzar;

  const onSelectCard = (cardId: string) => {
    const card = player.deck.find((card: Card) => card.id === cardId);

    setSelectedCardId((id) => (id === card?.id ? undefined : card?.id));

    // if (card) {
    //   props.handlePlayCard(card);
    // }
  };

  const onChooseChard = () => {
    const card = player.deck.find((card: Card) => card.id === selectedCardId);

    console.log({ card });

    if (card) {
      props.handlePlayCard(card);
    }
  };

  const rotations = [
    "rotate-295",
    "rotate-315",
    "rotate-335",
    "-rotate-5",
    "rotate-15",
    "rotate-35",
    "rotate-55",
  ];

  return (
    <div className="absolute bottom-0 flex flex-col justify-center items-center">
      <h3 className="text-center sm:text-left text-lg font-semibold mb-3">
        Your Hand
      </h3>
      <div className="relative lg:hidden flex flex-wrap gap-4 justify-center sm:justify-start">
        {player.deck?.map((card: Card) => (
          <>
            <GameCard
              key={card.id}
              onCardClick={onSelectCard}
              cardId={card.id}
              value={card.card_value}
              isDisabled={disabled}
              selected={selectedCardId === card.id}
            />
          </>
        ))}
      </div>

      <div className="hidden lg:block">
        {player.deck?.map((card: Card, index: number) => (
          <div
            className={cn(
              "absolute",
              "bottom-[16rem]",
              "-ml-8",
              "origin-[40px_500px]",
              "hover:z-10",
              rotations[index]
            )}
          >
            <GameCard
              key={card.id}
              onCardClick={onSelectCard}
              cardId={card.id}
              value={card.card_value}
              isDisabled={disabled}
              selected={selectedCardId === card.id}
            />
          </div>
        ))}
      </div>

      <div className="flex gap-2 mb-8">
        {playerIsCardCzar ? (
          <span className="text-center">
            You are the Card Czar! Choose a winning card.
          </span>
        ) : (
          <>
            <Button disabled={!selectedCardId} onClick={onChooseChard}>
              <Check />
              Choose
            </Button>
            <Button
              variant="outline"
              disabled={!selectedCardId}
              onClick={() => setSelectedCardId(undefined)}
            >
              <X />
              Cancel
            </Button>
          </>
        )}
      </div>
    </div>
  );
}
