import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useIsMobile } from "@/hooks/use-mobile";
import { ChevronDown, ChevronUp, MessageSquare } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { Button } from "../ui/button";

export default function SystemMessages(props: { chatMessages: any[] }) {
  const isMobile = useIsMobile();
  const [expanded, setExpanded] = useState(true);
  const scrollAreaRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setExpanded(!isMobile);
  }, [isMobile]);

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    if (scrollAreaRef.current) {
      const scrollElement = scrollAreaRef.current.querySelector(
        "[data-radix-scroll-area-viewport]"
      );
      if (scrollElement) {
        scrollElement.scrollTop = scrollElement.scrollHeight;
      }
    }
  }, [props.chatMessages]); // Re-run when game state changes (new messages)

  return (
    <Card className="h-full flex flex-col">
      <CardHeader className=" flex flex-shrink-0 flex-row justify-between">
        <h2 className="text-lg font-medium flex items-center gap-2">
          <MessageSquare className="h-5 w-5" /> System Messages
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
        <CardContent className="flex-1 p-0">
          <ScrollArea
            ref={scrollAreaRef}
            className="h-full max-h-32 min-h-32 p-4"
          >
            <div className="space-y-4 flex flex-col">
              {props.chatMessages.map((chat, i) => (
                <div key={i} className="space-y-1 gap-3 flex">
                  <Badge>Event</Badge>
                  <p className="text-sm text-muted-foreground">{chat}</p>
                </div>
              ))}
            </div>
          </ScrollArea>
        </CardContent>
      )}
    </Card>
  );
}
