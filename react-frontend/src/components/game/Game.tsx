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
import { getWebSocketUrl } from "@/lib/websocket";
import { Avatar, AvatarFallback } from "../ui/avatar";
import { Card, CardContent } from "../ui/card";
import { Spinner } from "../ui/spinner";
import SystemMessages from "./SystemMessages";
import { EmojiCard } from "../ui/emoji-card";
import Timer from "./Timer";

export default function GameComponent() {
  const { gameId } = useParams<{ gameId: string }>();
  const [game, setGame] = useState<Game | null>(null);
  const [chatMessages, setChatMessages] = useState<any[]>([]);
  const [scrollingEmojis, setScrollingEmojis] = useState<Array<{ id: string, emoji: string, timestamp: number, rightOffset: number }>>([]);
  const { user } = useAuth();

  const handleEmojiClick = (emoji: string) => {
    const newEmoji = {
      id: `${Date.now()}-${Math.random()}`,
      emoji,
      timestamp: Date.now(),
      rightOffset: Math.random() * 50 + 10 // Random offset between 10-60px from right
    };

    // setScrollingEmojis(prev => [...prev, newEmoji]);
    handleEmojiClicked(emoji);

    // Remove emoji after animation completes
    setTimeout(() => {
      setScrollingEmojis(prev => prev.filter(e => e.id !== newEmoji.id));
    }, 3000);
  };

  // Function to handle incoming WebSocket messages
  const handleIncomingWebSocketMessage = (message: {
    type: string;
    payload: any;
  }) => {

    
    switch (message.type) {
      case "GAME_UPDATE":
        debugger;
        setGame(message.payload);
        break;
      case "CHAT_MESSAGE":
        setChatMessages([...chatMessages, message.payload]);
        break;
      case "EMOJI_CLICKED":

        const newEmoji = {
          id: `${Date.now()}-${Math.random()}`,
          emoji: message.payload.emoji,
          timestamp: Date.now(),
          rightOffset: Math.random() * 50 + 10 // Random offset between 10-60px from right
        };

        setScrollingEmojis(prev => [...prev, newEmoji]);
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
  } = useWebSocket(getWebSocketUrl(`/ws/game/${gameId}`), {
    onMessage: (event) => {
      handleIncomingWebSocketMessage(JSON.parse(event.data));
    },
    onError: (error) => console.error("WebSocket error:", error),
    shouldReconnect: (closeEvent) => true,
  });

  const handleEmojiClicked = (emoji: string) => {
    sendMessage(
      JSON.stringify({
        type: "EmojiClicked",
        payload: {
          emoji,
          game_id: gameId,
          user_id: user?.user_id,
        },
      })
    );
  };

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
        type: "JudgeChoseWinningCard",
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

  if (readyState !== 1 || !game) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold mb-4">Connecting to game...</h2>
          <Spinner />
        </div>
      </div>
    );
  }

  return (
    <>
      <SiteHeader title={`Game: ${game.name}`} />
      <div className="flex flex-col">
        {/* Timer display */}
        <div className="container mx-auto px-6 py-2 flex justify-center">
          <Timer
            isActive={game.status === "InProgress"}
            nextAutoProgressAt={game.next_auto_progress_at}
            roundState={game.round_status}
          />
        </div>
        <div className="container mx-auto p-6">
          <div className="grid grid-cols-3 gap-4">
            <div className="col-span-3 order-1 lg:order-2 lg:col-span-1">
              <PlayerList
                handleJoinGame={handleJoinGame}
                handleBeginGame={handleBeginGame}
                handleContinueRound={handleContinueRound}
                game={game}
              />
            </div>

            {/* Main game area - takes full width on mobile, 8/12 on desktop */}
            <div className="col-span-3 lg:col-span-2 order-2 lg:order-1">
              {/* Board - make it larger and responsive */}
              <div className="flex flex-col gap-4">
                {renderGameBoard(game, handlePickWinningCard)}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Player's hand */}
      <PlayerHand game={game} handlePlayCard={handlePlayCard} />

      {/* Emoji List */}
      <div className="mx-auto container p-6">
        <EmojiCard scrollingEmojis={scrollingEmojis} handleEmojiClick={handleEmojiClick} />
      </div>

      {/* Chat */}
      <div className="mx-auto container p-6">
        <SystemMessages chatMessages={chatMessages} />
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

  const joinableSlots = props.game.max_player_count - props.game.players.length;
  const emptySlots = new Array(joinableSlots).fill(null);

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
                  Join seat
                </span>
              </div>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
