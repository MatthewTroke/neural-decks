import { Send } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export default function GameBoardChat(props: { game: Game }) {
    return (
      <Card className="flex-1 flex flex-col">
        <div className="p-4 border-b">
          <h3 className="font-semibold">Chat</h3>
        </div>
        <ScrollArea className="flex-1 p-4">
          <div className="space-y-4">
            {[
              { user: "John", message: "Good game everyone!" },
              { user: "Jane", message: "Nice move!" },
              { user: "Alex", message: "Thanks for playing" },
            ].map((msg, i) => (
              <div key={i} className="space-y-1">
                <span className="text-sm font-medium">{msg.user}</span>
                <p className="text-sm text-muted-foreground">{msg.message}</p>
              </div>
            ))}
          </div>
        </ScrollArea>
        <div className="p-4 border-t">
          <form className="flex gap-2">
            <Input placeholder="Type a message..." className="flex-1" />
            <Button size="icon">
              <Send className="h-4 w-4" />
            </Button>
          </form>
        </div>
      </Card>
    );
  }
  