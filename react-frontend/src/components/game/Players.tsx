import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";
import { useAuth } from "@/context/AuthContext";
import { Card } from "../ui/card";
import { Timer, Trophy, Users } from "lucide-react";
import { useEffect, useState } from "react";
import PlayerBadge from "@/components/game/players/Badge";
import { Badge } from "@/components/ui/badge";

export default function Players(props: {
  game: Game;
  handleJoinGame: () => void;
  handleBeginGame: () => void;
  handleContinueRound: () => void;
}) {
  const [showAnimation, setShowAnimation] = useState(false);

  // Reset animation when a new winner is chosen
  useEffect(() => {
    let hasRoundWinner = props.game.round_winner;

    if (!hasRoundWinner) {
      return;
    }

    const winner = props.game.players.find(
      (p) => p.user_id === props.game.round_winner.user_id
    );

    if (winner) {
      setShowAnimation(true);
      const timer = setTimeout(() => setShowAnimation(false), 2000);
      return () => clearTimeout(timer);
    }
  }, [props.game.players]);

  return (
    <div>
      <div className="flex justify-between items-start mb-4">
        <div className="space-y-1">
          <h3 className="font-semibold">Game #{props.game.id}</h3>
          <Badge
          // variant={
          //   game.status === "In Progress" ? "default" : "secondary"
          // }
          >
            {props.game.status}
          </Badge>
          <Badge
          // variant={
          //   game.status === "In Progress" ? "default" : "secondary"
          // }
          >
            {props.game.round_status}
          </Badge>
        </div>

        <JoinGameButton
          game={props.game}
          handleJoinGame={props.handleJoinGame}
        />
        <BeginGameButton
          game={props.game}
          handleBeginGame={props.handleBeginGame}
        />
        <ContinueRoundButton
          game={props.game}
          handleContinueRound={props.handleContinueRound}
        />
      </div>

      <div className="bg-[#121a2f] rounded-xl p-4 border border-gray-800">
        <h2 className="text-lg font-medium mb-3 flex items-center gap-2">
          <Users className="h-5 w-5" /> Players
        </h2>
        <div className="space-y-2">
          {props.game.players.map((player: Player) => (
            <div
              key={player.user_id}
              className="flex justify-between items-center"
            >
              <div className="flex items-center gap-2">
                <div className="w-6 h-6 rounded-full bg-gray-700 flex items-center justify-center text-xs">
                  {player.name.at(0)}
                </div>
                <span>{player.name}</span>
                <PlayerBadge player={player} game={props.game} />
              </div>
              <span>{player.score} pts</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

interface ContinueRoundButtonProps {
  game: Game;
  handleContinueRound: () => void;
}

function ContinueRoundButton(props: ContinueRoundButtonProps) {
  const isGameInProgress = props.game.status === "InProgress";
  const isGameRoundOver = props.game.round_status === "CardCzarChoseWinningCard";

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
  let { user } = useAuth();

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
