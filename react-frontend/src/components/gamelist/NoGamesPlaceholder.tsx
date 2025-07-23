import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { PlusCircle, Gamepad2 } from "lucide-react"

export function NoGamesPlaceholder() {
  return (
    <Card className="bg-card text-card-foreground border-border">
      <CardHeader className="flex flex-col items-center text-center">
        <Gamepad2 className="h-12 w-12 text-muted-foreground mb-4" />
        <CardTitle className="text-2xl font-bold">No Active Game Rooms</CardTitle>
        <CardDescription className="text-muted-foreground mt-2">
          It looks like there are no active game rooms right now.
          <br />
          Be the first to create one and invite your friends!
        </CardDescription>
      </CardHeader>
    </Card>
  )
}
