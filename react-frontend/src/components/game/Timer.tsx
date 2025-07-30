import { Clock } from "lucide-react";
import { useEffect, useState } from "react";

interface TimerProps {
  isActive?: boolean;
  onReset?: () => void;
}

export default function Timer({ isActive = false, onReset }: TimerProps) {
  const [secondsRemaining, setSecondsRemaining] = useState<number>(30);

  // Reset timer when game state changes
  useEffect(() => {
    if (isActive && onReset) {
      setSecondsRemaining(30);
    }
  }, [isActive, onReset]);

  // Countdown effect
  useEffect(() => {
    if (!isActive || secondsRemaining <= 0) {
      return;
    }

    const interval = setInterval(() => {
      setSecondsRemaining((prev) => {
        if (prev <= 1) return 0;
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [isActive, secondsRemaining]);

  if (!isActive || secondsRemaining <= 0) {
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
      <Clock className="h-4 w-4" />
      <span className="font-mono font-bold">
        {formatTime(secondsRemaining)}
      </span>
    </div>
  );
} 