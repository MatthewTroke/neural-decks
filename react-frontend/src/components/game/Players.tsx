import PlayerBadge from "@/components/game/players/Badge";
import { Button } from "@/components/ui/button";
import { useIsMobile } from "@/hooks/use-mobile";
import { ChevronDown, ChevronUp, Users } from "lucide-react";
import { useEffect, useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import { Card, CardContent, CardHeader } from "../ui/card";
import { GameButtons } from "./GameBoardButtons";

export default function Players(props: {
  game: Game;
  handleJoinGame: () => void;
  handleBeginGame: () => void;
  handleContinueRound: () => void;
}) {
  const isMobile = useIsMobile();
  const [expanded, setExpanded] = useState(true);
  const [showAnimation, setShowAnimation] = useState(false);

  useEffect(() => {
    setExpanded(!isMobile);
  }, [isMobile]);

  // Reset animation when a new winner is chosen
  useEffect(() => {
    const hasRoundWinner = props.game.round_winner;

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
    <div className="flex flex-col w-full">
      <GameButtons
        game={props.game}
        handleBeginGame={props.handleBeginGame}
        handleContinueRound={props.handleContinueRound}
        handleJoinGame={props.handleJoinGame}
      />

      <Card className="w-full">
        <CardHeader className="flex flex-row justify-between">
          <h2 className="text-lg font-medium flex items-center gap-2">
            <Users className="h-5 w-5" /> Players
          </h2>

          <Button
            variant="ghost"
            title={expanded ? "Expand Players" : "Collapse Players"}
            onClick={() => setExpanded((expanded) => !expanded)}
          >
            {expanded ? <ChevronDown /> : <ChevronUp />}
          </Button>
        </CardHeader>
        {expanded && (
          <CardContent>
            <div className="space-y-2">
              {props.game.players.map((player) => (
                <div
                  key={player.user_id}
                  className="flex justify-between items-center"
                >
                  <div className="flex items-center gap-2">
                    <Avatar>
                      <AvatarImage alt={player.name} />
                      <AvatarFallback>{player.name.at(0)}</AvatarFallback>
                    </Avatar>
                    <span>{player.name}</span>
                    <PlayerBadge player={player} game={props.game} />
                  </div>
                  <span>{player.score} pts</span>
                </div>
              ))}
            </div>
          </CardContent>
        )}
      </Card>
    </div>
  );
}
