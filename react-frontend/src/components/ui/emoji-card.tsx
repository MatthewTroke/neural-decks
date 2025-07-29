import * as React from "react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { MessageSquare } from "lucide-react";
import { ScrollArea } from "@radix-ui/react-scroll-area";

interface EmojiCardProps {
  className?: string;
  emojis?: string[];
  title?: string;
}

const DEFAULT_EMOJIS = [
  "😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣", "😊", "😇",
];

export function EmojiCard({
  className,
  emojis = DEFAULT_EMOJIS,
  title = "Click an emoji!"
}: EmojiCardProps) {
  const [scrollingEmojis, setScrollingEmojis] = React.useState<Array<{id: string, emoji: string, timestamp: number, rightOffset: number}>>([]);

  const handleEmojiClick = (emoji: string) => {
    const newEmoji = {
      id: `${Date.now()}-${Math.random()}`,
      emoji,
      timestamp: Date.now(),
      rightOffset: Math.random() * 50 + 10 // Random offset between 10-60px from right
    };
    
    setScrollingEmojis(prev => [...prev, newEmoji]);

    // Remove emoji after animation completes
    setTimeout(() => {
      setScrollingEmojis(prev => prev.filter(e => e.id !== newEmoji.id));
    }, 3000);
  };

  return (
    <>
      <Card variant="ghost" className="h-full flex flex-col p-0 flex-0">
        <CardContent className="flex-1 p-0">
          {emojis.slice(0, 32).map((emoji, index) => (
            <button
              key={index}
              onClick={() => handleEmojiClick(emoji)}
              className="text-2xl hover:scale-110 transition-transform duration-200 p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
              aria-label={`Click ${emoji}`}
            >
              {emoji}
            </button>
          ))}
        </CardContent>
      </Card>


      {/* Scrolling emojis overlay */}
      {scrollingEmojis.length > 0 && (
        <div className="fixed inset-0 pointer-events-none z-50 overflow-hidden">
          {scrollingEmojis.map((emojiData) => (
            <div
              key={emojiData.id}
              className="absolute text-4xl animate-scroll-up"
              style={{
                right: `${emojiData.rightOffset}px`,
                animationDelay: '0s',
                animationDuration: '3s',
              }}
            >
              {emojiData.emoji}
            </div>
          ))}
        </div>
      )}
    </>
  );
}

// Add the scroll-up animation to your global CSS or create a style tag
const scrollUpStyles = `
  @keyframes scroll-up {
    0% {
      transform: translateY(100vh) rotate(0deg);
      opacity: 1;
    }
    100% {
      transform: translateY(-100px) rotate(360deg);
      opacity: 0;
    }
  }
  
  .animate-scroll-up {
    animation: scroll-up 3s ease-out forwards;
  }
`;

// Inject styles
if (typeof document !== 'undefined') {
  const styleId = 'emoji-scroll-styles';
  if (!document.getElementById(styleId)) {
    const style = document.createElement('style');
    style.id = styleId;
    style.textContent = scrollUpStyles;
    document.head.appendChild(style);
  }
} 