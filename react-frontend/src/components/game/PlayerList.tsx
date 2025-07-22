import Players from "@/components/game/Players";

export default function PlayerList(props: {
    game: Game;
    handleJoinGame: () => void;
    handleBeginGame: () => void;
    handleContinueRound: () => void;
  }) {
    return (
        <Players
          handleJoinGame={props.handleJoinGame}
          handleBeginGame={props.handleBeginGame}
          handleContinueRound={props.handleContinueRound}
          game={props.game}
        />
    );
  }
  