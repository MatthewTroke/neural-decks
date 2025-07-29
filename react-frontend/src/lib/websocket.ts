// WebSocket utility function
export function getWebSocketUrl(path: string): string {
  // Replace http://localhost:8080 with ws://localhost:8080 for WebSocket connections
  const baseUrl = 'ws://localhost:8080';
  return `${baseUrl}${path}`;
} 