import { jwtDecode } from "jwt-decode";
import React, { createContext, useContext, useState, useEffect } from "react";
import Cookies from 'js-cookie';
import api from '@/lib/axios';

interface AuthContextType {
  user: User | null;
  logout: () => void;
  isLoggedIn: boolean;
  loading: boolean;
}

const defaultAuthContext: AuthContextType = {
  user: null,
  logout: () => {},
  isLoggedIn: false,
  loading: true,
};

const AuthContext = createContext(defaultAuthContext);

export const useAuth = () => {
  return useContext(AuthContext);
};

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  
  const logout = async () => {
    try {
      await api.post("/auth/logout");
    } catch (error) {
      console.error("❌ [FRONTEND] Logout error:", error);
    } finally {
      setUser(null);
      Cookies.remove("neural_decks_jwt");
    }
  };

  useEffect(() => {
    const initializeAuth = async () => {
      const token = Cookies.get("neural_decks_jwt");

      if (!token) {
        setUser(null);
        setLoading(false);
        return;
      }

      try {
        // Decode the token to check expiration
        const decodedToken: any = jwtDecode(token);
        const currentTime = Date.now() / 1000;
        const timeUntilExpiry = decodedToken.exp - currentTime;
        
        // If token is expired, clear user state (backend will handle refresh)
        if (decodedToken.exp && (decodedToken.exp - currentTime) <= 0) {
          setUser(null);
          Cookies.remove("neural_decks_jwt");
        } else {
          // Token is still valid
          setUser(decodedToken as User);
        }
      } catch (error) {
        console.error("❌ [FRONTEND] Token decode error:", error);
        setUser(null);
        Cookies.remove("neural_decks_jwt");
      }
      
      setLoading(false);
    };

    initializeAuth();
  }, []);

  const value = {
    user,
    logout,
    isLoggedIn: !!user,
    loading,
  };

  if (loading) {
    return null;
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
