import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn } from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import {
  CircleArrowRight,
  LoaderIcon,
  RefreshCw,
  Sparkles,
  Timer,
  Trophy,
  Users,
} from "lucide-react";
import { useNavigate } from "react-router";
import { SiteHeader } from "../site-header";
import { Separator } from "../ui/separator";
import { CreateGameDialog } from "./CreateGame";
import { NoGamesPlaceholder } from "./NoGamesPlaceholder";

// Mock data for demonstration
const onlinePlayers = [
  {
    name: "Emily Chen",
    avatar: "https://i.pravatar.cc/150?u=emily",
    status: "Looking for game",
  },
  {
    name: "David Kim",
    avatar: "https://i.pravatar.cc/150?u=david",
    status: "Looking for game",
  },
  {
    name: "Maria Garcia",
    avatar: "https://i.pravatar.cc/150?u=maria",
    status: "Looking for game",
  },
  {
    name: "Tom Wilson",
    avatar: "https://i.pravatar.cc/150?u=tom",
    status: "Looking for game",
  },
  {
    name: "Lisa Park",
    avatar: "https://i.pravatar.cc/150?u=lisa",
    status: "Looking for game",
  },
  {
    name: "James Lee",
    avatar: "https://i.pravatar.cc/150?u=james",
    status: "Looking for game",
  },
];

export default function GameLobby() {
  const { data: games, isLoading } = useQuery({
    queryKey: ["games"],
    queryFn: async () => {
      return axios
        .get("http://localhost:8080/games", {
          withCredentials: true,
        })
        .then((data) => data.data);
    },
  });

  if (isLoading) {
    return null;
  }

  return (
    <>
      <SiteHeader title="Games" />
      <div className="container mx-auto px-4 py-6">
        <section className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
          <FeatureCard
            icon={<Sparkles className="h-8 w-8 text-primary" />}
            title="AI-Generated cards"
            description="Our AI creates unique, hilarious cards based on your preferences and themes."
          />
          <FeatureCard
            icon={<Users className="h-8 w-8 text-primary" />}
            title="Play with anyone"
            description="Join public rooms or create private games with friends. Play with 3-10 players per room."
          />
          <FeatureCard
            icon={<Trophy className="h-8 w-8 text-primary" />}
            title="Custom game rules"
            description="Create your own rules, time limits, and scoring systems for unique gameplay."
          />
        </section>

        <Card>
          <CardHeader className="flex flex-row justify-between">
            <CardTitle className="text-2xl">Active Game Rooms</CardTitle>
            <div className="flex gap-2">
              <Button variant="outline" size="icon" onClick={() => {}}>
                <RefreshCw className="h-4 w-4" />
              </Button>
              <CreateGameDialog />
            </div>
          </CardHeader>
          <CardContent>
            <Games games={games} />
          </CardContent>
        </Card>

        <div className="flex flex-col md:flex-row gap-8 mb-16 items-start">
          {/* Online Players Section */}
          {/* <div className="w-full md:w-1/3">
            <div className="grid gap-4">
              <Card>
                <CardHeader>
                  <CardTitle>Live Stats</CardTitle>
                  <CardDescription>
                    Real-time platform statistics
                  </CardDescription>
                </CardHeader>
                <CardContent>Coming soon...</CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle>Live Stats</CardTitle>
                  <CardDescription>
                    Real-time platform statistics
                  </CardDescription>
                </CardHeader>
                <CardContent>Coming soon...</CardContent>
              </Card>
            </div>
          </div> */}
        </div>
      </div>
    </>
  );
}

function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <Card>
      <CardHeader>
        <div className="mb-2">{icon}</div>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <CardDescription className="text-zinc-400">
          {description}
        </CardDescription>
      </CardContent>
    </Card>
  );
}

function Games(props: { games: Game[] }) {
  const navigate = useNavigate();

  const handleEnterGameRoom = (gameId: string) => {
    navigate("/games/" + gameId);
  };

  if (props.games.length === 0) {
    return <NoGamesPlaceholder />;
  }

  return (
    <>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {props.games.map((game: Game) => (
          <Card key={game.id} className="p-6" variant="ghost">
            <div className="flex justify-between items-center">
              <div className="flex gap-2">
                <h3 className="font-semibold">Game: {game.name}</h3>
              </div>

              <Button
                variant="secondary"
                onClick={() => handleEnterGameRoom(game.id)}
                size="sm"
              >
                <CircleArrowRight />
                Enter room
              </Button>
            </div>

            <div className="space-y-4">
              {/* Game details */}
              <div className="flex flex-row-reverse gap-4 text-sm">
                <Badge
                  className={cn({
                    "animate-pulse": game.status === "InProgress",
                  })}
                >
                  <LoaderIcon className="animate-spin" />
                  {game.status === "InProgress" ? "In Progress" : "Waiting"}
                </Badge>
                <div className="flex items-center">
                  <Timer className="h-4 w-4 mr-1" />
                  0:00
                </div>
                <div className="flex items-center">
                  <Trophy className="h-4 w-4 mr-1" />
                  Round {game.current_game_round}
                </div>
              </div>

              <Separator />

              {/* Players */}
              <div>
                <div className="flex items-center gap-2 mb-4">
                  <Users className="h-4 w-4" />
                  <span className="text-sm font-medium">Players</span>
                </div>
                <ScrollArea className="min-h-24">
                  <div className="space-y-2">
                    {game.players.map((player) => (
                      <div
                        key={player.name}
                        className="flex items-center justify-between"
                      >
                        <div className="flex items-center gap-2">
                          <Avatar className="h-6 w-6">
                            <AvatarImage src={player.image} />
                            <AvatarFallback>{player.name[0]}</AvatarFallback>
                          </Avatar>
                          <span className="text-sm">{player.name}</span>
                        </div>
                        <span className="text-sm font-medium">0 pts</span>
                      </div>
                    ))}
                  </div>
                </ScrollArea>
              </div>
              
            </div>
          </Card>
        ))}
      </div>
    </>
  );
}
