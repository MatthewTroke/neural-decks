import { Clock } from "lucide-react";
import { useEffect, useState, useRef } from "react";

interface TimerProps {
  isActive?: boolean;
  onReset?: () => void;
  nextAutoProgressAt?: string | null; // ISO timestamp from the game object
  roundState?: string; // Current round state to detect changes
}

export default function Timer({ isActive = false, onReset, nextAutoProgressAt, roundState }: TimerProps) {
  const [secondsRemaining, setSecondsRemaining] = useState<number | null>(null);
  const previousRoundState = useRef<string | undefined>(undefined);

  // Calculate time remaining from nextAutoProgressAt
  useEffect(() => {
    if (!nextAutoProgressAt) {
      return;
    }

    const calculateTimeRemaining = () => {
      const now = new Date().getTime();
      const targetTime = new Date(nextAutoProgressAt).getTime();
      const diff = Math.max(0, Math.floor((targetTime - now) / 1000));
      setSecondsRemaining(diff);
    };

    calculateTimeRemaining();
  }, [nextAutoProgressAt]);

  // Reset timer to 30 seconds when round state changes (but not on initial render)
  useEffect(() => {
    if (roundState && previousRoundState.current !== undefined && previousRoundState.current !== roundState) {
      // Only reset if we have a previous round state and it's different
      setSecondsRemaining(30);
    }
    
    // Update the previous round state
    previousRoundState.current = roundState;
  }, [roundState]);

  // Countdown effect - decrement the timer every second
  useEffect(() => {
    if (!isActive || secondsRemaining === null || secondsRemaining <= 0) {
      return;
    }

    const interval = setInterval(() => {
      setSecondsRemaining((prev) => {
        if (prev === null || prev <= 1) return 0;
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [isActive, secondsRemaining]);

  if (!isActive || secondsRemaining === null || secondsRemaining <= 0) {
    return null;
  }

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  // Color based on time remaining
  const getTimerColor = () => {
    if (secondsRemaining <= 10) return "text-red-500";
    if (secondsRemaining <= 20) return "text-yellow-500";
    return "text-blue-500";
  };

  return (
    <div className={`flex items-center gap-2 text-sm ${getTimerColor()}`}>
      <Clock className="h-6 w-6" />
      <span className="font-mono font-bold text-xl">
        {formatTime(secondsRemaining)}
      </span>
    </div>
  );
} 