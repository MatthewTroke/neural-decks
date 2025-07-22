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
    let hasRoundWinner = props.game.RoundWinner;

    if (!hasRoundWinner) {
      return;
    }

    const winner = props.game.Players.find(
      (p) => p.UserID === props.game.RoundWinner.UserID
    );

    if (winner) {
      setShowAnimation(true);
      const timer = setTimeout(() => setShowAnimation(false), 2000);
      return () => clearTimeout(timer);
    }
  }, [props.game.Players]);

  return (
    <div>
      <div className="flex justify-between items-start mb-4">
        <div className="space-y-1">
          <h3 className="font-semibold">Game #{props.game.ID}</h3>
          <Badge
          // variant={
          //   game.Status === "In Progress" ? "default" : "secondary"
          // }
          >
            {props.game.Status}
          </Badge>
          <Badge
          // variant={
          //   game.Status === "In Progress" ? "default" : "secondary"
          // }
          >
            {props.game.RoundStatus}
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
          {props.game.Players.map((player) => (
            <div
              key={player.UserID}
              className="flex justify-between items-center"
            >
              <div className="flex items-center gap-2">
                <div className="w-6 h-6 rounded-full bg-gray-700 flex items-center justify-center text-xs">
                  {player.Name.at(0)}
                </div>
                <span>{player.Name}</span>
                <PlayerBadge player={player} game={props.game} />
              </div>
              <span>{player.Score} pts</span>
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
  const isGameInProgress = props.game.Status === "InProgress";
  const isGameRoundOver = props.game.RoundStatus === "CardCzarChoseWinningCard";

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

  const isUserInGame = props.game.Players.some(
    (player: Player) => player.UserID === user?.user_id
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
  const isGameReadyToBegin = props.game.Players.length > 1;
  const isGameInSetupState = props.game.Status === "Setup";

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
