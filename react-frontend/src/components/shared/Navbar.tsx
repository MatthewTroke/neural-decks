import { Button } from "@/components/ui/button";
import { useAuth } from "@/context/AuthContext";
import { useState } from "react";
import { Link } from "react-router";
import { Menu, X } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { Badge } from "../ui/badge";

export function Navbar() {
  const { user } = useAuth();
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
    <header className="sticky top-0 z-50 w-full border-b border-gray-800 bg-[#0f1524]/95 backdrop-blur supports-[backdrop-filter]:bg-[#0f1524]/60">
      <div
        className="container mx-auto px-4 md:px-6"
        style={{ maxWidth: "1280px" }}
      >
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center gap-2">
            <Link to="/" className="flex items-center gap-2">
              <span className="text-xl font-bold">Neural Decks</span>
              <Badge>Beta</Badge>
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
            {user?.email}
            {/* <Link
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
              </Link> */}
          </nav>
        </div>
      </div>

      {/* Mobile menu */}
      {mobileMenuOpen && (
        <div className="md:hidden">
          <div className="space-y-1 px-2 pb-3 pt-2">
            <Link
              to="/login"
              className="block rounded-md px-3 py-2 text-base font-medium text-gray-300 hover:bg-gray-800 hover:text-white"
              onClick={() => setMobileMenuOpen(false)}
            >
              Report a bug
              <div className="flex items-center">
                Login
              </div>
            </Link>
          </div>
        </div>
      )}
    </header>
  );
}
