import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";

interface ProtectedRouteProps {
  children: React.ReactNode;
  requireAuth?: boolean; // true for protected routes, false for login page
}

function ProtectedRoute({ children, requireAuth = true }: ProtectedRouteProps) {
  const { isLoggedIn } = useAuth();
  const location = useLocation();

  // If this is the login page and user is already logged in, redirect to /games
  if (!requireAuth && isLoggedIn) {
    return <Navigate to="/games" replace />;
  }

  // If this is a protected route and user is not logged in, redirect to /login
  if (requireAuth && !isLoggedIn) {
    return <Navigate to="/login" replace />;
  }

  return children;
}

export default ProtectedRoute;
