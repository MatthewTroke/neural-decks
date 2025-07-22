import { jwtDecode } from "jwt-decode";
import React, { createContext, useContext, useState, useEffect } from "react";
import Cookies from 'js-cookie';

interface AuthContextType {
  user: User | null;
  logout: () => void;
  isLoggedIn: boolean;
}
const defaultAuthContext: AuthContextType = {
  user: null,
  logout: () => {},
  isLoggedIn: false,
};

const AuthContext = createContext(defaultAuthContext);

export const useAuth = () => {
  return useContext(AuthContext);
};

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  
  useEffect(() => {
    const token = Cookies.get("neural_decks_jwt");

    if (!token) {
      setUser(null);
      setLoading(false);
      return;
    }

    let user: User = jwtDecode(token);

    if (user) {
      setUser(user);
    }
    setLoading(false);
  }, []);

  const logout = () => {
    setUser(null);
    localStorage.removeItem("neural_decks_jwt");
  };

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
