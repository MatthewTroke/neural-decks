import Board from "@/components/game/Board";
import PlayerHand from "@/components/game/Hand";
import PlayerList from "@/components/game/PlayerList";
import { useAuth } from "@/context/AuthContext";
import { AvatarImage } from "@radix-ui/react-avatar";
import { Crown, Users } from "lucide-react";
import { useState } from "react";
import { useParams } from "react-router";
import useWebSocket from "react-use-websocket";
import { SiteHeader } from "../site-header";
import { Avatar, AvatarFallback } from "../ui/avatar";
import { Card, CardContent } from "../ui/card";

export default function GameComponent() {
  const { gameId } = useParams<{ gameId: string }>();
  const [game, setGame] = useState<Game | null>(null);
  const { user } = useAuth();

  // Function to handle incoming WebSocket messages
  const handleIncomingWebSocketMessage = (message: {
    type: string;
    payload: any;
  }) => {
    switch (message.type) {
      case "GAME_UPDATE":
        setGame(message.payload);
        break;
      case "CHAT_MESSAGE":
        break;
      default:
    }
  };

  const {
    sendJsonMessage,
    lastJsonMessage,
    readyState,
    lastMessage,
    getWebSocket,
    sendMessage,
  } = useWebSocket(`ws://localhost:8080/ws/game/${gameId}`, {
    onMessage: (event) => {
      handleIncomingWebSocketMessage(JSON.parse(event.data));
    },
    onError: (error) => console.error("WebSocket error:", error),
    shouldReconnect: (closeEvent) => true,
  });

  const handleJoinGame = () => {
    sendMessage(
      JSON.stringify({
        type: "JoinedGame",
        payload: {
          //TODO remove user_id and use the claim on the backend
          game_id: gameId,
          user_id: user?.user_id,
        },
      })
    );
  };

  const handleBeginGame = () => {
    sendMessage(
      JSON.stringify({
        type: "GameBegins",
        payload: {
          //TODO remove user_id and use the claim on the backend
          game_id: gameId,
          user_id: user?.user_id,
        },
      })
    );
  };

  const handlePlayCard = (card: Card) => {
    sendMessage(
      JSON.stringify({
        type: "CardPlayed",
        payload: {
          card_id: card.id,
          game_id: gameId,
        },
      })
    );
  };

  const handlePickWinningCard = (cardId: string) => {
    sendMessage(
      JSON.stringify({
        type: "CardCzarChoseWinningCard",
        payload: {
          card_id: cardId,
          game_id: gameId,
        },
      })
    );
  };

  const handleContinueRound = () => {
    sendMessage(
      JSON.stringify({
        type: "RoundContinued",
        payload: {
          game_id: gameId,
        },
      })
    );
  };

  if (!game) {
    return null;
  }

  if (readyState !== 1) {
    return <div>Loading...</div>;
  }

  return (
    <>
      <SiteHeader title={`Game: ${game.Name}`} />

      <div className="flex flex-col">
        <div className="container mx-auto px-4 py-6">
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
            {/* Sidebar - stacks on mobile (top), side by side on desktop */}
            <div className="lg:col-span-4 order-1 lg:order-2">
              {/* Player List and Chat in 2 columns on small screens, but still within the sidebar */}
              <div className="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-1 gap-4 mb-4 lg:mb-0">
                {/* Player List */}
                <div className="sm:col-span-2">
                  <PlayerList
                    handleJoinGame={handleJoinGame}
                    handleBeginGame={handleBeginGame}
                    handleContinueRound={handleContinueRound}
                    game={game}
                  />
                </div>

                {/* Chat */}
                <div className="sm:col-span-1">
                  <div>Chatroom</div>
                  <div>{game.status}</div>
                  <div>{game.round_status}</div>
                </div>
              </div>
            </div>

            {/* Main game area - takes full width on mobile, 8/12 on desktop */}
            <div className="lg:col-span-8 order-2 lg:order-1">
              {/* Board - make it larger and responsive */}
              <div className="flex flex-col gap-4">
                {renderGameBoard(game, handlePickWinningCard)}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="w-full absolute bottom-0 top-0 flex justify-center">
        {/* Player's hand */}
        <PlayerHand game={game} handlePlayCard={handlePlayCard} />
      </div>
    </>
  );
}

function renderGameBoard(
  game: Game,
  handlePickWinningCard: (cardId: string) => void
) {
  if (game.status === "Setup") {
    return <JoinGameGrid game={game} />;
  }

  return <Board game={game} handlePickWinningCard={handlePickWinningCard} />;
}

function JoinGameGrid(props: { game: Game }) {
  const players = props.game.players;

  let joinableSlots = props.game.max_player_count - props.game.players.length;
  let emptySlots = new Array(joinableSlots).fill(null);

  const playerGrid = [...players, ...emptySlots];

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {playerGrid.map((player, index) => (
        <Card key={index} className={`${!player && "border-dashed"}`}>
          <CardContent className="p-4">
            {player ? (
              <div className="flex flex-col items-center text-center gap-2">
                <Avatar className="h-16 w-16">
                  <AvatarImage src={player.image} />
                  <AvatarFallback>{player.name.substring(0, 2)}</AvatarFallback>
                </Avatar>
                <div>
                  <div className="font-medium">{player.name}</div>
                  <Crown className="h-3 w-3 mr-1" />
                </div>
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-[116px] text-center">
                <Users className="h-8 w-8 text-zinc-600 mb-2" />
                <span className="text-sm text-zinc-500">
                  Waiting for player...
                </span>
              </div>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
