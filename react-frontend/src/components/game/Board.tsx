import { useAuth } from "@/context/AuthContext";
import GameCard from "./GameCard";

export default function Board(props: {
  game: Game;
  handlePickWinningCard: (cardId: string) => void;
}) {
  const { user } = useAuth();

  let player = props.game.players.find((p: any) => p.user_id === user?.user_id);
  // let isCardCzar = player?.is_card_czar;
  // let isWhiteCardDisabled =
  //   !isCardCzar ||
  //   !(
  //     props.game.status === "InProgress" &&
  //     props.game.round_status === "CardCzarPickingWinningCard"
  //   );

  let hasRoundWinner =
    props.game.round_winner &&
    props.game.status === "InProgress" &&
    props.game.round_status === "CardCzarChoseWinningCard";

  let winningCard = null;

  if (hasRoundWinner) {
    winningCard = props.game.round_winner.placed_card;
  }

  const onWhiteCardClick = (card: Card) => {
    if (!player) {
      return;
    }

    if (props.game.status === "Setup") {
      return;
    }

    if (props.game.round_status !== "CardCzarPickingWinningCard") {
      return;
    }

    if (!player.is_card_czar) {
      return;
    }

    props.handlePickWinningCard(card.id);
  };

  return (
    <div>
      {/* <Card key={props.game.ID} className="p-6"> */}
      <div className="flex justify-between items-start">
        <div className="space-y-1">
          <div className="flex items-start">
            <div className="flex flex-wrap gap-4 justify-center sm:justify-start">
              <RenderBlackCard card={props.game.black_card} />
              {props.game.white_cards.map((card) => (
                <GameCard
                  key={card.id}
                  cardId={card.id}
                  value={card.card_value}
                  onCardClick={() => onWhiteCardClick(card)}
                  isDisabled={false}
                  isWinningCard={winningCard?.id === card.id ? true : false}
                />
              ))}
            </div>
          </div>
        </div>
      </div>
      {/* </Card> */}
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
      key={card.id}
      cardId={card.id}
      value={card.card_value}
      onCardClick={() => {}}
      isDisabled
      isBlack
    />
  );
}
