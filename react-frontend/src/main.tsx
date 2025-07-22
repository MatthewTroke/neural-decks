import { createRoot } from "react-dom/client";
import {
  BrowserRouter as Router,
  Route,
  Routes,
  Link,
  Navigate,
} from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { AuthProvider, useAuth } from "@/context/AuthContext";
import Login from "@/components/login/Login.tsx";
import ProtectedRoute from "@/components/shared/ProtectedRoute.tsx";
import "./index.css";
import GameList from "@/components/gamelist/GameList.tsx";
import Game from "@/components/game/Game";

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
          <Route path="/login" element={<Login />} />
          <Route
            path="/games"
            element={
              <ProtectedRoute>
                <GameList />
              </ProtectedRoute>
            }
          />
          <Route
            path="/games/:gameId"
            element={
              <ProtectedRoute>
                <Game />
              </ProtectedRoute>
            }
          />
          <Route path="*" element={<RedirectToAppropriatePage />} />
        </Routes>
      </Router>
    </QueryClientProvider>
  </AuthProvider>
);
