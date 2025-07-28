import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"
import { useMutation } from "@tanstack/react-query"
import axios from "axios"
import { Link } from "react-router-dom"

export default function LoginPage() {
    const googleAuthMutation = useMutation({
        mutationFn: async () => {
          const response = await axios.get("http://localhost:8080/auth/google", {
            withCredentials: true,
          });
    
          return response.data;
        },
        onSuccess: (data: any) => {
          window.location.href = data.redirect_url;
        },
        onError: (error: any) => {
          console.error("Error generating OAuth link:", error);
        },
      });

      const discordAuthMutation = useMutation({
        mutationFn: async () => {
          const response = await axios.get("http://localhost:8080/auth/discord", {
            withCredentials: true,
          });
    
          return response.data;
        },
        onSuccess: (data: any) => {
          window.location.href = data.redirect_url;
        },
        onError: (error: any) => {
          console.error("Error generating OAuth link:", error);
        },
      });

    return (
        <div className="min-h-screen bg-background flex items-center justify-center p-4">
            {/* Login Card */}
            <Card className="w-full max-w-md bg-card/50 border-border backdrop-blur-sm">
                <CardHeader className="text-center space-y-2">
                    <CardTitle className="text-2xl font-bold text-foreground flex justify-center">
                        <div className="flex items-center gap-3">
                            <h1 className="text-xl font-bold text-foreground">Neural Decks</h1>
                            <span className="bg-primary text-primary-foreground text-xs px-2 py-1 rounded-full">Beta</span>
                        </div></CardTitle>
                    <CardDescription className="text-muted-foreground">Sign in to a provider to access the game.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                    {/* Google Login Button */}
                    <Button
                        onClick={() => googleAuthMutation.mutate()}
                        variant="outline"
                        className="w-full bg-muted hover:bg-muted text-foreground border-border font-medium"
                    >
                        <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24">
                            <path
                                fill="#4285F4"
                                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                            />
                            <path
                                fill="#34A853"
                                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                            />
                            <path
                                fill="#FBBC05"
                                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                            />
                            <path
                                fill="#EA4335"
                                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                            />
                        </svg>
                        Login with Google
                    </Button>

                    {/* Login with Discord Button */}
                    <Button
                        onClick={() => discordAuthMutation.mutate()}
                        variant="outline"
                        className="w-full bg-muted hover:bg-muted text-foreground border-border font-medium"
                    >
                        <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M20.317 4.3698C18.9998 3.8066 17.6336 3.3378 16.2376 3.0178C16.2152 3.0134 16.1922 3.0222 16.1772 3.0402C15.9632 3.323 15.7362 3.6618 15.5932 3.9126C13.8072 3.635 12.0262 3.635 10.2522 3.9126C10.1092 3.6618 9.8822 3.323 9.6672 3.0402C9.6522 3.0222 9.6292 3.0134 9.6068 3.0178C8.2108 3.3378 6.8446 3.8066 5.5274 4.3698C5.513 4.3762 5.5012 4.3878 5.4932 4.4022C2.5332 9.0366 1.4192 13.5222 1.6692 17.9658C1.6696 17.9738 1.6742 17.981 1.6812 17.9846C3.5042 19.3238 5.2762 20.0938 7.0122 20.4846C7.0352 20.4898 7.0582 20.4782 7.0702 20.457C7.2772 20.1042 7.4612 19.7438 7.6242 19.377C7.6512 19.3162 7.6202 19.2478 7.5562 19.233C7.1532 19.1402 6.7652 19.0218 6.3892 18.8802C6.3242 18.855 6.3192 18.7614 6.3802 18.7326C6.4712 18.6894 6.5622 18.6442 6.6512 18.597C10.4322 20.2938 13.5682 20.2938 17.3482 18.597C17.4372 18.6442 17.5282 18.6894 17.6192 18.7326C17.6802 18.7614 17.6752 18.855 17.6102 18.8802C17.2342 19.0218 16.8462 19.1402 16.4432 19.233C16.3792 19.2478 16.3482 19.3162 16.3752 19.377C16.5382 19.7438 16.7222 20.1042 16.9292 20.457C16.9412 20.4782 16.9642 20.4898 16.9872 20.4846C18.7232 20.0938 20.4952 19.3238 22.3182 17.9846C22.3252 17.981 22.3298 17.9738 22.3302 17.9658C22.5802 13.5222 21.4662 9.0366 18.5062 4.4022C18.4982 4.3878 18.4864 4.3762 18.472 4.3698H20.317ZM8.0202 15.3318C7.2952 15.3318 6.7032 14.6934 6.7032 13.9314C6.7032 13.1694 7.2812 12.531 8.0202 12.531C8.7652 12.531 9.3532 13.1754 9.3402 13.9314C9.3402 14.6934 8.7652 15.3318 8.0202 15.3318ZM15.9792 15.3318C15.2542 15.3318 14.6622 14.6934 14.6622 13.9314C14.6622 13.1694 15.2402 12.531 15.9792 12.531C16.7242 12.531 17.3122 13.1754 17.2992 13.9314C17.2992 14.6934 16.7242 15.3318 15.9792 15.3318Z" fill="#5865F2" />
                        </svg>
                        Login with Discord
                    </Button>
                </CardContent>
            </Card>
        </div>
    )
}
