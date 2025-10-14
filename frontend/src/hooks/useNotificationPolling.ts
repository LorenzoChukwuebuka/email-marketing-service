import { useState, useEffect, useRef, useCallback } from 'react';

interface Notification {
    id: string;
    message: string;
    created_at?: string;
    [key: string]: any;
}

interface NotificationResponse {
    data: Notification[];
}

interface UseNotificationPollingOptions {
    enabled?: boolean;
    retryDelay?: number;
    pollInterval?: number;
}

interface UseNotificationPollingReturn {
    notifications: Notification[];
    error: string | null;
    isPolling: boolean;
    clearNotifications: () => void;
}

/**
 * Custom hook for long polling notifications
 * @param baseUrl - The base URL of the API
 * @param token - The authentication token
 * @param options - Configuration options
 * @returns Object containing notifications, error, isPolling status, and clearNotifications function
 */
export function useNotificationPolling(
    baseUrl: string,
    token: string,
    options: UseNotificationPollingOptions = {}
): UseNotificationPollingReturn {
    const {
        enabled = true,
        retryDelay = 5000,
        pollInterval = 100
    } = options;

    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [isPolling, setIsPolling] = useState(false);

    const lastNotificationIdRef = useRef<string | null>(null);
    const pollingRef = useRef(false);
    const abortControllerRef = useRef<AbortController | null>(null);

    const clearNotifications = useCallback(() => {
        setNotifications([]);
    }, []);

    useEffect(() => {
        // Don't start polling if disabled, no baseUrl, or no token
        if (!enabled || !baseUrl || !token) {
            return;
        }

        pollingRef.current = true;
        setIsPolling(true);

        const poll = async () => {
            while (pollingRef.current) {
                // Create a new AbortController for this request
                abortControllerRef.current = new AbortController();

                try {
                    const url = new URL(`${baseUrl}/notifications/poll`);

                    if (lastNotificationIdRef.current) {
                        url.searchParams.append('sinceId', lastNotificationIdRef.current);
                    }

                    const response = await fetch(url, {
                        method: 'GET',
                        headers: {
                            'Authorization': `Bearer ${token}`,
                            'Content-Type': 'application/json'
                        },
                        signal: abortControllerRef.current.signal
                    });

                    if (!response.ok) {
                        throw new Error(`HTTP error! status: ${response.status}`);
                    }

                    const data: NotificationResponse = await response.json();

                    if (data.data && data.data.length > 0) {
                        // Update lastNotificationId to the newest notification's ID
                        lastNotificationIdRef.current = data.data[0].id;

                        // Add new notifications to state
                        setNotifications(prev => [...data.data, ...prev]);

                        // Clear any previous errors
                        setError(null);
                    }

                    // Wait before next poll
                    await new Promise(resolve => setTimeout(resolve, pollInterval));
                } catch (err) {
                    // Ignore abort errors (happens during cleanup)
                    if (err instanceof Error && err.name === 'AbortError') {
                        break;
                    }

                    const errorMessage = err instanceof Error ? err.message : 'Unknown error';
                    console.error('Polling error:', errorMessage);
                    setError(errorMessage);

                    // Wait before retrying on error
                    await new Promise(resolve => setTimeout(resolve, retryDelay));
                }
            }
        };

        poll();

        // Cleanup function
        return () => {
            pollingRef.current = false;
            setIsPolling(false);

            // Abort any ongoing fetch request
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        };
    }, [baseUrl, token, enabled, retryDelay, pollInterval]);

    return {
        notifications,
        error,
        isPolling,
        clearNotifications
    };
}