import { useAuth } from "@/context/AuthContext";
import GameCard from "./GameCard";
import { Card } from "../ui/card";

export default function Board(props: {
  game: Game;
  handlePickWinningCard: (cardId: string) => void;
}) {
  let { user } = useAuth();

  let player = props.game.Players.find((p: any) => p.UserID === user?.user_id);
  // let isCardCzar = player?.IsCardCzar;
  // let isWhiteCardDisabled =
  //   !isCardCzar ||
  //   !(
  //     props.game.Status === "InProgress" &&
  //     props.game.RoundStatus === "CardCzarPickingWinningCard"
  //   );

  let hasRoundWinner =
    props.game.RoundWinner &&
    props.game.Status === "InProgress" &&
    props.game.RoundStatus === "CardCzarChoseWinningCard";

  let winningCard = null;

  if (hasRoundWinner) {
    winningCard = props.game.RoundWinner.PlacedCard;
  }

  const onWhiteCardClick = (card: Card) => {
    if (!player) {
      return;
    }

    if (props.game.Status === "Setup") {
      return;
    }

    if (props.game.RoundStatus !== "CardCzarPickingWinningCard") {
      return;
    }

    if (!player.IsCardCzar) {
      return;
    }

    props.handlePickWinningCard(card.ID);
  };

  return (
    <div>


      <Card key={props.game.ID} className="p-6">
        <div className="flex justify-between items-start">
          <div className="space-y-1">
            <div className="flex items-start">
              <div className="flex flex-wrap gap-4 justify-center sm:justify-start">
                <RenderBlackCard card={props.game.BlackCard} />
                {props.game.WhiteCards.map((card) => (
                  <GameCard
                    key={card.ID}
                    cardId={card.ID}
                    value={card.CardValue}
                    onCardClick={() => onWhiteCardClick(card)}
                    isDisabled={false}
                    isWinningCard={winningCard?.ID === card.ID ? true : false}
                  />
                ))}
              </div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
}

interface RenderBlackCardProps {
  card: Card | null;
}

function RenderBlackCard(props: RenderBlackCardProps) {
  const { card } = props;

  if (!card) {
    return null;
  }

  return (
    <GameCard
      key={card.ID}
      cardId={card.ID}
      value={card.CardValue}
      onCardClick={() => {}}
      isDisabled
      isBlack
    />
  );
}
