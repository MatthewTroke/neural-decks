import { cn } from "@/lib/utils";
import { Crown, Eye } from "lucide-react";
import { Badge } from "@/components/ui/badge";

interface PlayerBadgeProps {
  game: Game;
  player: Player;
}

export default function PlayerBadge(props: PlayerBadgeProps) {
  if (props.game.status === "Setup") {
    return null;
  }

  if (props.player?.user_id === props.game?.round_winner?.user_id) {
    return (
      <div>
        <Badge variant="outline" className="text-xs">
          <Crown className="h-4 w-4 text-primary" />
          Winner!
        </Badge>
      </div>
    );
  }

  if (!props.player.is_card_czar && !props.player.placed_card) {
    return (
      <div>
        <Badge
          variant="secondary"
          className={cn(
            "text-xs animate-pulse",
            "bg-primary/10 text-primary hover:bg-primary/10"
          )}
        >
          <Eye className="h-4 w-4 text-primary" />
          Choosing...
        </Badge>
      </div>
    );
  }

  if (props.player.is_card_czar) {
    return (
      <div>
        <Badge variant="default" className={cn("text-xs animate-pulse")}>
          Picking...
        </Badge>
      </div>
    );
  }

  return null;
}
