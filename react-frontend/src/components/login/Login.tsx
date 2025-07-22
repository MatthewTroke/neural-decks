import axios from "axios";
import { useMutation } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardContent } from "@/components/ui/card";
import cbhImage from "@/assets/cbh-white.png";
import { Separator } from "@/components/ui/separator";

import {
  ChromeIcon as Google,
  Menu,
  Github,
  Twitter,
  Brain,
  Users,
  Shuffle,
  X,
} from "lucide-react";
import { Link } from "react-router";
import { useState } from "react";

export default function LandingPage() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const mutation = useMutation({
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

  return (
    <div className="flex min-h-[100dvh] flex-col bg-[#0f1524] text-white">
      {/* Header */}
      <header className="sticky top-0 z-50 w-full border-b border-gray-800 bg-[#0f1524]/95 backdrop-blur supports-[backdrop-filter]:bg-[#0f1524]/60">
        <div
          className="container mx-auto px-4 md:px-6"
          style={{ maxWidth: "1280px" }}
        >
          <div className="flex h-16 items-center justify-between">
            <div className="flex items-center gap-2">
              <Link to="/" className="flex items-center gap-2">
                <span className="text-xl font-bold">Neural Decks</span>
                <div className="inline-block rounded-lg bg-purple-900/20 px-3 py-1 text-sm text-purple-400">
                  Beta
                </div>
              </Link>
            </div>

            {/* Mobile menu button */}
            <div className="md:hidden">
              <button
                className="flex items-center justify-center"
                onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                aria-label="Toggle menu"
              >
                {mobileMenuOpen ? (
                  <X className="h-6 w-6" />
                ) : (
                  <Menu className="h-6 w-6" />
                )}
              </button>
            </div>

            {/* Desktop navigation */}
            <nav className="hidden md:flex items-center gap-6">
              <Link
                href="#features"
                className="text-sm font-medium text-gray-300 hover:text-white"
              >
                Features
              </Link>
              <Link
                href="#how-it-works"
                className="text-sm font-medium text-gray-300 hover:text-white"
              >
                How It Works
              </Link>
              <Link
                href="#about"
                className="text-sm font-medium text-gray-300 hover:text-white"
              >
                About
              </Link>
              <Link href="#">
                <Button
                  onClick={() => mutation.mutate()}
                  asChild
                  variant="default"
                  className="bg-white text-[#0f1524] hover:bg-gray-200"
                >
                  <span>
                    <Google className="mr-2 h-4 w-4" />
                    Login with Google
                  </span>
                </Button>
              </Link>
            </nav>
          </div>
        </div>

        {/* Mobile menu */}
        {mobileMenuOpen && (
          <div className="md:hidden">
            <div className="space-y-1 px-2 pb-3 pt-2">
              <Link
                href="#features"
                className="block rounded-md px-3 py-2 text-base font-medium text-gray-300 hover:bg-gray-800 hover:text-white"
                onClick={() => setMobileMenuOpen(false)}
              >
                Features
              </Link>
              <Link
                href="#how-it-works"
                className="block rounded-md px-3 py-2 text-base font-medium text-gray-300 hover:bg-gray-800 hover:text-white"
                onClick={() => setMobileMenuOpen(false)}
              >
                How It Works
              </Link>
              <Link
                href="#about"
                className="block rounded-md px-3 py-2 text-base font-medium text-gray-300 hover:bg-gray-800 hover:text-white"
                onClick={() => setMobileMenuOpen(false)}
              >
                About
              </Link>
              <Link
                href="#login"
                className="block rounded-md px-3 py-2 text-base font-medium text-gray-300 hover:bg-gray-800 hover:text-white"
                onClick={() => setMobileMenuOpen(false)}
              >
                <div className="flex items-center">
                  <Google className="mr-2 h-4 w-4" />
                  Login with Google
                </div>
              </Link>
            </div>
          </div>
        )}
      </header>

      <main className="flex-1">
        {/* Hero Section */}
        <section className="w-full py-12 md:py-24 lg:py-32">
          <div
            className="flex container mx-auto px-4 md:px-6"
            style={{ maxWidth: "1280px" }}
          >
            <div className="grid gap-6 lg:gap-12 items-center mx-auto">
              <div className="flex flex-col justify-center space-y-4">
                <div className="space-y-2">
                  <h1 className="text-3xl font-bold tracking-tighter sm:text-5xl xl:text-6xl/none">
                    Play Neural Decks,{" "}
                    <span className="text-purple-400 block sm:inline">
                      Powered by AI
                    </span>
                  </h1>
                  <p className="max-w-[600px] text-gray-400 md:text-xl">
                    Play Neural Decks for free. Neural Decks is a real-time
                    online card game you can play with friends or random people.
                    Our AI generates unique, hilarious card decks that keep the
                    game fresh every time you play.
                  </p>
                </div>
                <div className="flex flex-col gap-2 sm:flex-row">
                  <Button
                    asChild
                    size="lg"
                    className="bg-white text-[#0f1524] hover:bg-gray-200"
                  >
                    <Link href="#login">
                      <Google className="mr-2 h-4 w-4" />
                      Login with Google
                    </Link>
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Card Preview Section */}
        <section className="w-full py-12 md:py-24 relative overflow-hidden">
          <div className="absolute inset-0 bg-gradient-to-b from-purple-900/20 to-transparent"></div>
          <div
            className="container mx-auto px-4 md:px-6 relative"
            style={{ maxWidth: "1280px" }}
          >
            <div className="mx-auto flex flex-col items-center space-y-4 text-center">
              <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">
                Endless Possibilities
              </h2>
              <p className="max-w-[700px] text-gray-400 md:text-xl">
                Our AI card generation feature creates card combinations that
                you won't find in other traditional card games, giving each game
                a feel of uniqueness.
              </p>
            </div>

            <div className="mt-16 flex flex-wrap justify-center gap-4 md:gap-8">
              {/* Black Card */}
              <div className="relative h-64 w-48 rounded-lg bg-black p-4 shadow-lg transform rotate-[-5deg]">
                <p className="font-medium">In a world overrun by zombies.</p>
                <div className="absolute bottom-4 left-4 text-xs text-gray-500">
                  Neural Decks
                </div>
              </div>

              {/* White Cards */}
              <div className="relative h-64 w-48 rounded-lg bg-white p-4 shadow-lg text-black transform rotate-[3deg]">
                <p className="font-medium">
                  A horde of kittens with laser eyes.
                </p>
                <div className="absolute bottom-4 left-4 text-xs text-gray-500">
                  Neural Decks
                </div>
              </div>

              <div className="relative h-64 w-48 rounded-lg bg-white p-4 shadow-lg text-black transform rotate-[8deg] hidden sm:block">
                <p className="font-medium">
                  The neighbor's garden gnomes coming to life.
                </p>
                <div className="absolute bottom-4 left-4 text-xs text-gray-500">
                  Neural Decks
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section id="features" className="w-full py-12 md:py-24 lg:py-32">
          <div
            className="container mx-auto px-4 md:px-6"
            style={{ maxWidth: "1280px" }}
          >
            <div className="mx-auto flex flex-col items-center space-y-4 text-center">
              <div className="inline-block rounded-lg bg-purple-900/20 px-3 py-1 text-sm text-purple-400">
                Features
              </div>
              <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">
                Why Play?
              </h2>
              <p className="max-w-[700px] text-gray-400 md:text-xl">
                Neural Decks takes a twist on traditional card games by using
                AI-powered features.
              </p>
            </div>

            <div
              className="mx-auto grid gap-8 py-12 sm:grid-cols-2 md:grid-cols-3"
              style={{ maxWidth: "1024px" }}
            >
              <Card className="bg-[#121a2f] border-gray-800 shadow-xl">
                <div className="flex flex-col items-center p-6 text-center">
                  <div className="flex h-12 w-12 items-center justify-center rounded-full bg-purple-900/20 text-purple-400">
                    <Brain className="h-6 w-6" />
                  </div>
                  <h3 className="mt-4 text-xl font-bold">AI-Generated Decks</h3>
                  <p className="mt-2 text-gray-400">
                    Our AI creates unique card combinations that keep the game
                    fresh and unpredictable.
                  </p>
                </div>
              </Card>

              <Card className="bg-[#121a2f] border-gray-800 shadow-xl">
                <div className="flex flex-col items-center p-6 text-center">
                  <div className="flex h-12 w-12 items-center justify-center rounded-full bg-purple-900/20 text-purple-400">
                    <Users className="h-6 w-6" />
                  </div>
                  <h3 className="mt-4 text-xl font-bold">Play With Friends</h3>
                  <p className="mt-2 text-gray-400">
                    Invite friends to join your game room and play together no
                    matter where you are.
                  </p>
                </div>
              </Card>

              <Card className="bg-[#121a2f] border-gray-800 shadow-xl sm:col-span-2 md:col-span-1">
                <div className="flex flex-col items-center p-6 text-center">
                  <div className="flex h-12 w-12 items-center justify-center rounded-full bg-purple-900/20 text-purple-400">
                    <Shuffle className="h-6 w-6" />
                  </div>
                  <h3 className="mt-4 text-xl font-bold">Endless Variety</h3>
                  <p className="mt-2 text-gray-400">
                    Never play the same game twice. Choose a subject to create a
                    new unique deck to keep games fresh and replayable.
                  </p>
                </div>
              </Card>
            </div>
          </div>
        </section>

        {/* How It Works Section */}
        <section
          id="how-it-works"
          className="w-full py-12 md:py-24 lg:py-32 bg-[#121a2f]"
        >
          <div
            className="container mx-auto px-4 md:px-6"
            style={{ maxWidth: "1280px" }}
          >
            <div className="mx-auto flex flex-col items-center space-y-4 text-center">
              <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">
                How It Works
              </h2>
              <p className="max-w-[700px] text-gray-400 md:text-xl">
                Playing is simple.
              </p>
            </div>

            <div
              className="mx-auto grid gap-8 py-12 sm:grid-cols-2 md:grid-cols-3"
              style={{ maxWidth: "1024px" }}
            >
              <div className="flex flex-col items-center text-center">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-white text-[#0f1524] font-bold text-xl">
                  1
                </div>
                <h3 className="mt-4 text-xl font-bold">Login with Google</h3>
                <p className="mt-2 text-gray-400">
                  Sign in with your Google account to access the game.
                </p>
              </div>

              <div className="flex flex-col items-center text-center">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-white text-[#0f1524] font-bold text-xl">
                  2
                </div>
                <h3 className="mt-4 text-xl font-bold">
                  Create or Join a Game
                </h3>
                <p className="mt-2 text-gray-400">
                  Start a new game room or join an existing one with friends or
                  randoms.
                </p>
              </div>

              <div className="flex flex-col items-center text-center sm:col-span-2 md:col-span-1">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-white text-[#0f1524] font-bold text-xl">
                  3
                </div>
                <h3 className="mt-4 text-xl font-bold">Play and Enjoy</h3>
                <p className="mt-2 text-gray-400">
                  Experience the uniqueness of AI-generated cards for maximum
                  replayability
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* Final CTA Section */}
        <section id="login" className="w-full py-12 md:py-24 lg:py-32">
          <div
            className="container mx-auto px-4 md:px-6"
            style={{ maxWidth: "1280px" }}
          >
            <div className="mx-auto flex flex-col items-center space-y-4 text-center">
              <div className="space-y-2">
                <h2 className="text-3xl font-bold tracking-tighter md:text-4xl/tight">
                  Ready to play?
                </h2>
                <p className="max-w-[600px] text-gray-400 md:text-xl/relaxed">
                  {/* Join many others already enjoying our AI-powered card game. */}
                </p>
              </div>
              <div className="flex flex-col gap-2 min-[400px]:flex-row">
                <Button
                  asChild
                  size="lg"
                  className="bg-white text-[#0f1524] hover:bg-gray-200"
                >
                  <Link href="#login">
                    <Google className="mr-2 h-4 w-4" />
                    Login with Google
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="w-full border-t border-gray-800 bg-[#0a0f1a] py-12">
        <div
          className="container mx-auto px-4 md:px-6"
          style={{ maxWidth: "1280px" }}
        >
          <div className="grid gap-8 lg:grid-cols-2">
            <div className="flex flex-col space-y-4">
              <Link href="/" className="flex items-center gap-2">
                {/* <Image
                  src="/placeholder.svg?height=32&width=32"
                  width={32}
                  height={32}
                  alt="Logo"
                  className="rounded"
                /> */}
                <span className="text-xl font-bold">Neural Decks </span>
                <div className="inline-block rounded-lg bg-purple-900/20 px-3 py-1 text-sm text-purple-400">
                  Beta
                </div>
              </Link>
              <p className="max-w-[400px] text-sm text-gray-400">
                Neural Decks is a real-time AI-powered card game you can play
                with friends or random people. Our platform uses artificial
                intelligence to create unique, hilarious card combinations to
                keep games fresh & unique.
              </p>
              <div className="flex gap-4">
                <Link href="#" className="text-gray-400 hover:text-white">
                  <Twitter className="h-5 w-5" />
                  <span className="sr-only">Twitter</span>
                </Link>
                <Link href="#" className="text-gray-400 hover:text-white">
                  <Github className="h-5 w-5" />
                  <span className="sr-only">GitHub</span>
                </Link>
              </div>
            </div>
            <div className="grid gap-8 sm:grid-cols-2">
              <div className="space-y-4">
                <h3 className="text-sm font-medium">Links</h3>
                <ul className="space-y-2 text-sm">
                  <li>
                    <Link
                      href="#features"
                      className="text-gray-400 hover:text-white"
                    >
                      Features
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="#how-it-works"
                      className="text-gray-400 hover:text-white"
                    >
                      How It Works
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="#about"
                      className="text-gray-400 hover:text-white"
                    >
                      About
                    </Link>
                  </li>
                </ul>
              </div>
              {/* <div className="space-y-4">
                <h3 className="text-sm font-medium">Legal</h3>
                <ul className="space-y-2 text-sm">
                  <li>
                    <Link href="#" className="text-gray-400 hover:text-white">
                      Privacy Policy
                    </Link>
                  </li>
                  <li>
                    <Link href="#" className="text-gray-400 hover:text-white">
                      Terms of Service
                    </Link>
                  </li>
                  <li>
                    <Link href="#" className="text-gray-400 hover:text-white">
                      Cookie Policy
                    </Link>
                  </li>
                </ul>
              </div> */}
            </div>
          </div>
          {/* <div className="flex flex-col items-center justify-between gap-4 border-t border-gray-800 py-6 md:h-24 md:flex-row md:py-0 mt-8"> */}
          {/* <p className="text-center text-sm text-gray-400 md:text-left">
              &copy; {new Date().getFullYear()} Neural Decks
            </p> */}
          {/* <p className="text-center text-sm text-gray-400 md:text-left">
              Not affiliated with Cards Against Humanity LLC.
            </p> */}
          {/* </div> */}
        </div>
      </footer>
    </div>
  );
}
