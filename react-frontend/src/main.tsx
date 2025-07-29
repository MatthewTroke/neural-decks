import Game from "@/components/game/Game";
import GameList from "@/components/gamelist/GameList.tsx";
import Login from "@/components/login/Login.tsx";
import ProtectedRoute from "@/components/shared/ProtectedRoute.tsx";
import { AuthProvider, useAuth } from "@/context/AuthContext";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createRoot } from "react-dom/client";
import {
  Navigate,
  Route,
  BrowserRouter as Router,
  Routes,
} from "react-router-dom";
import { Navigation } from "./components/shared/Navigation";
import "./index.css";
import LandingPage from "@/components/landing/LandingPage";

const queryClient = new QueryClient();

function RedirectToAppropriatePage() {
  const { isLoggedIn } = useAuth();

  return isLoggedIn ? <Navigate to="/games" /> : <Navigate to="/login" />;
}

createRoot(document.getElementById("root")!).render(
  <AuthProvider>
    <QueryClientProvider client={queryClient}>
      <Router>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route 
            path="/login" 
            element={
              <ProtectedRoute requireAuth={false}>
                <Login />
              </ProtectedRoute>
            } 
          />
          <Route
            path="/games"
            element={
              <ProtectedRoute>
                <Navigation>
                  <GameList />
                </Navigation>
              </ProtectedRoute>
            }
          />
          <Route
            path="/games/:gameId"
            element={
              <ProtectedRoute>
                <Navigation>
                  <Game />
                </Navigation>
              </ProtectedRoute>
            }
          />
          <Route path="*" element={<RedirectToAppropriatePage />} />
        </Routes>
      </Router>
    </QueryClientProvider>
  </AuthProvider>
);
