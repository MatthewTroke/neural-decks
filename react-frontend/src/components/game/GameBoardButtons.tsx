import { Button } from "@/components/ui/button";
import { useAuth } from "@/context/AuthContext";

interface ContinueRoundButtonProps {
  game: Game;
  handleContinueRound: () => void;
}

function ContinueRoundButton(props: ContinueRoundButtonProps) {
  const isGameInProgress = props.game.status === "InProgress";
  const isGameRoundOver =
    props.game.round_status === "JudgeChoseWinningCard";

  const shouldRenderContinueRoundButton = isGameInProgress && isGameRoundOver;

  if (!shouldRenderContinueRoundButton) {
    return null;
  }

  return (
    <Button onClick={props.handleContinueRound} variant="secondary" size="sm">
      Continue round
    </Button>
  );
}

interface JoinGameButtonProps {
  game: Game;
  handleJoinGame: () => void;
}

function JoinGameButton(props: JoinGameButtonProps) {
  const { user } = useAuth();

  const isUserInGame = props.game.players.some(
    (player: Player) => player.user_id === user?.user_id
  );
  const shouldRenderJoinGameButton = !isUserInGame;

  if (!shouldRenderJoinGameButton) {
    return null;
  }

  return (
    <Button
      onClick={() => props.handleJoinGame()}
      variant="secondary"
      size="sm"
    >
      Join game
    </Button>
  );
}

interface BeginGameButtonProps {
  game: Game;
  handleBeginGame: () => void;
}

function BeginGameButton(props: BeginGameButtonProps) {
  const isGameReadyToBegin = props.game.players.length > 1;
  const isGameInSetupState = props.game.status === "Setup";

  const shouldRenderBeginGameButton = isGameReadyToBegin && isGameInSetupState;

  if (!shouldRenderBeginGameButton) {
    return null;
  }

  return (
    <Button onClick={props.handleBeginGame} variant="secondary" size="sm">
      Begin game
    </Button>
  );
}

interface GameButtonsProps {
  game: Game;
  handleBeginGame: () => void;
  handleJoinGame: () => void;
  handleContinueRound: () => void;
}

export function GameButtons(props: GameButtonsProps) {
  const { user } = useAuth();
  const isGameReadyToBegin = props.game.players.length > 1;
  const isGameInSetupState = props.game.status === "Setup";
  const isUserInGame = props.game.players.some(
    (player: Player) => player.user_id === user?.user_id
  );

  const isGameInProgress = props.game.status === "InProgress";
  const isGameRoundOver =
    props.game.round_status === "JudgeChoseWinningCard";

  const shouldRenderContinueRoundButton = isGameInProgress && isGameRoundOver;
  const shouldRenderBeginGameButton = isGameReadyToBegin && isGameInSetupState;
  const shouldRenderJoinGameButton = !isUserInGame;

  if (
    !shouldRenderBeginGameButton &&
    !shouldRenderContinueRoundButton &&
    !shouldRenderJoinGameButton
  ) {
    return null;
  }

  return (
    <div className="flex justify-between items-start mb-4">
      <JoinGameButton game={props.game} handleJoinGame={props.handleJoinGame} />
      <BeginGameButton
        game={props.game}
        handleBeginGame={props.handleBeginGame}
      />
      <ContinueRoundButton
        game={props.game}
        handleContinueRound={props.handleContinueRound}
      />
    </div>
  );
}
