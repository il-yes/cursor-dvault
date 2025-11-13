import { useState, useEffect } from 'react';

interface UseConnectivityReturn {
  isOnline: boolean;
  wasOffline: boolean;
  onReconnect: (callback: () => void) => void;
}

export function useConnectivity(): UseConnectivityReturn {
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [wasOffline, setWasOffline] = useState(false);
  const [reconnectCallbacks, setReconnectCallbacks] = useState<(() => void)[]>([]);

  useEffect(() => {
    const handleOnline = () => {
      setIsOnline(true);
      
      // Execute all reconnect callbacks
      reconnectCallbacks.forEach(callback => {
        try {
          callback();
        } catch (error) {
          console.error('Reconnect callback error:', error);
        }
      });
      
      // Clear callbacks after execution
      setReconnectCallbacks([]);
      
      // Track that we were offline
      if (!isOnline) {
        setWasOffline(true);
        setTimeout(() => setWasOffline(false), 5000); // Reset after 5s
      }
    };

    const handleOffline = () => {
      setIsOnline(false);
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, [isOnline, reconnectCallbacks]);

  const onReconnect = (callback: () => void) => {
    setReconnectCallbacks(prev => [...prev, callback]);
  };

  return { isOnline, wasOffline, onReconnect };
}
