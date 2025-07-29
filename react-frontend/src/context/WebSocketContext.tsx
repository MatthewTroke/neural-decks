import { createContext, useContext, ReactNode } from 'react';
import useWebSocket from 'react-use-websocket';
import { getWebSocketUrl } from '@/lib/websocket';

interface WebSocketContextType {
  sendJsonMessage: (message: any) => void;
  lastJsonMessage: any;
  readyState: number;
  lastMessage: any;
  getWebSocket: () => any;
  sendMessage: (message: string) => void;
}

const WebSocketContext = createContext<WebSocketContextType | null>(null);

interface WebSocketProviderProps {
  children: ReactNode;
  gameId: string;
  onMessage: (message: any) => void;
}

export function WebSocketProvider({ children, gameId, onMessage }: WebSocketProviderProps) {
  console.log('WebSocketProvider: Connecting to game', gameId);
  
  const ws = useWebSocket(getWebSocketUrl(`/ws/game/${gameId}`), {
    onMessage: (event) => {
      console.log('WebSocket received message:', event.data);
      console.log('Event data type:', typeof event.data);
      console.log('Event data length:', event.data?.length);
      
      if (!event.data) {
        console.warn('WebSocket received empty data');
        return;
      }
      
      try {
        const message = JSON.parse(event.data);
        console.log('Parsed message:', message);
        onMessage(message);
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error, event.data);
      }
    },
    onError: (error) => console.error("WebSocket error:", error),
    onOpen: () => console.log('WebSocket connected'),
    onClose: () => console.log('WebSocket disconnected'),
    shouldReconnect: (closeEvent) => true,
  });

  console.log('WebSocket state:', ws.readyState);

  return (
    <WebSocketContext.Provider value={ws}>
      {children}
    </WebSocketContext.Provider>
  );
}

export const useGameWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('useGameWebSocket must be used within a WebSocketProvider');
  }
  return context;
}; 