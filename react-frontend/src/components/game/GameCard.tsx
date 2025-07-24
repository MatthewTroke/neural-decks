import { SelectableCard } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { Crown } from "lucide-react";

interface CardProps {
  cardId: string;
  value: string;
  onCardClick: (cardId: string) => void;
  isDisabled: boolean;
  isBlack?: boolean;
  isWinningCard?: boolean;

  /**
   * Is the card active/selected, changes the background color to primary
   */
  selected?: boolean;
}

export default function GameCard(props: CardProps) {
  const {
    isDisabled,
    isBlack,
    value,
    cardId,
    isWinningCard,
    onCardClick,
    selected,
  } = props;

  return (
    <SelectableCard
      onClick={() => (isDisabled ? () => {} : onCardClick(cardId))}
      key={cardId}
      className={cn(
        "aspect-[3/4] w-36 h-48 flex items-start p-3 transition-all relative",
        !isDisabled && "cursor-pointer hover:ring-2 hover:ring-primary",
        isDisabled && "cursor-not-allowed opacity-70",
        isBlack
          ? "bg-black text-white"
          : "bg-white text-black border-2 border-black",
        isWinningCard ? "bg-primary text-white" : "",
        { "bg-primary": selected }
      )}
      aria-disabled={isDisabled ? true : false}
      aria-selected={selected}
    >
      {isWinningCard && (
        <span className="absolute -right-3 -top-4 z-10 rotate-30 text-yellow-400 animate-bounce">
          <Crown />
        </span>
      )}
      <div className="w-full">
        <p
          className={cn(
            "text-sm font-bold leading-tight tracking-tight",
            isBlack ? "text-white" : "text-black"
          )}
        >
          {value}
        </p>
      </div>

      {/* Small logo at bottom right */}
      <div
        className={cn(
          "absolute bottom-2 right-2 text-[10px] font-semibold",
          isBlack ? "text-white" : "text-black"
        )}
      >
        Neural Decks
      </div>
    </SelectableCard>
  );
}
