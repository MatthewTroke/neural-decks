import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Users,
  Timer,
  Trophy,
  PlayCircle,
  RefreshCw,
  Sparkles,
  Lock,
  Unlock,
} from "lucide-react";
import { Navbar } from "@/components/shared/Navbar";
import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { CreateGameDialog } from "./CreateGame";
import { Navigate, useNavigate } from "react-router";
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
      <Navbar />
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

        <div className="flex flex-col md:flex-row gap-8 mb-16 items-start">
          <div className="w-full md:w-2/3">
            <div className="lg:col-span-2">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-2xl font-bold">Active Game Rooms</h2>
                <div className="flex gap-2">
                  <Button variant="outline" size="icon" onClick={() => {}}>
                    <RefreshCw className="h-4 w-4" />
                  </Button>
                  <CreateGameDialog />
                </div>
              </div>
              <div className="space-y-4">
                <Games games={games} />
              </div>
            </div>
          </div>
          {/* Online Players Section */}
          <div className="w-full md:w-1/3">
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
          </div>
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
    return <NoGamesPlaceholder />
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-6">
      {props.games.map((game: Game) => (
        <Card key={game.id} className="p-6">
          <div className="flex justify-between items-start mb-4">
            <div className="space-y-1">
              <h3 className="font-semibold">Game #{game.id}</h3>
              <Badge
              // variant={
              //   game.status === "In Progress" ? "default" : "secondary"
              // }
              >
                {game.status}
              </Badge>
            </div>

            <Button
              onClick={() => handleEnterGameRoom(game.id)}
              variant="outline"
              size="sm"
            >
              Enter room
            </Button>
          </div>

          <div className="space-y-4">
            {/* Game details */}
            <div className="flex items-center gap-4 text-sm">
              <div className="flex items-center">
                <Timer className="h-4 w-4 mr-1" />
                0:00
              </div>
                              <div className="flex items-center">
                  <Trophy className="h-4 w-4 mr-1" />
                  Round {game.current_game_round}
                </div>
            </div>

            {/* Players */}
            <div>
              <div className="flex items-center gap-2 mb-2">
                <Users className="h-4 w-4" />
                <span className="text-sm font-medium">Players</span>
              </div>
              <ScrollArea className="h-24">
                <div className="space-y-2">
                  {game.players.map((player: Player) => (
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
  );
}
