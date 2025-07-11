import { notification } from 'antd';
import {
    CheckCircleOutlined,
    CloseCircleOutlined,
    InfoCircleOutlined
} from '@ant-design/icons';
import eventBus from "./utils/eventbus";
import { useEffect } from "react";
import { queryClient } from "./utils/queryclient";
import AppRouter from "./routes";
import { QueryClientProvider } from "@tanstack/react-query";
import { tokenRefreshOptions } from "./features/auth/hooks/useTokenRefreshQuery";
import { TokenRefreshProvider } from "./features/auth/hooks/useTokenRefreshProvider";
import Cookies from "js-cookie";
import { GoogleOAuthProvider } from "@react-oauth/google";

// Configure notification defaults
notification.config({
    placement: 'bottomRight',
    duration: 4.5,
    maxCount: 3,
    rtl: false,
});

queryClient.setQueryDefaults(tokenRefreshOptions.queryKey, tokenRefreshOptions); // Pre-configure the token refresh query

function App() {
    const handleSuccess = (message: string) => {
        notification.success({
            message: 'Success',
            description: message,
            icon: <CheckCircleOutlined style={{ color: '#52c41a' }} />,
            duration: 4.5,
            style: {
                borderRadius: '8px',
                boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
            },
        });
    };

    const handleError = (message: string) => {
        notification.error({
            message: 'Error',
            description: message,
            icon: <CloseCircleOutlined style={{ color: '#ff4d4f' }} />,
            duration: 5,
            style: {
                borderRadius: '8px',
                boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
            },
        });
    };

    const handleInfo = (message: string) => {
        notification.info({
            message: 'Information',
            description: message,
            icon: <InfoCircleOutlined style={{ color: '#1890ff' }} />,
            duration: 4.5,
            style: {
                borderRadius: '8px',
                boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
            },
        });
    };

    const handleRefreshError = (error: Error) => {
        console.error('Token refresh failed:', error);

        // Show error notification before redirect
        notification.error({
            message: 'Session Expired',
            description: 'Your session has expired. Please log in again.',
            icon: <CloseCircleOutlined style={{ color: '#ff4d4f' }} />,
            duration: 3,
            style: {
                borderRadius: '8px',
                boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
            },
        });

        // Clean up and redirect after a short delay
        setTimeout(() => {
            Cookies.remove('Cookies', { path: '/' });
            window.location.href = "/auth/login";
        }, 1000);
    };

    useEffect(() => {
        const successListener = (message: string) => handleSuccess(message);
        const errorListener = (message: string) => handleError(message);
        const infoListener = (message: string) => handleInfo(message);

        eventBus.on("success", successListener);
        eventBus.on("error", errorListener);
        eventBus.on("message", infoListener);

        // Clean up the event listeners on unmount
        return () => {
            eventBus.off("success", successListener);
            eventBus.off("error", errorListener);
            eventBus.off("message", infoListener);
        };
    }, []);

    return (
        <QueryClientProvider client={queryClient}>
            <TokenRefreshProvider onRefreshError={handleRefreshError}>
                <GoogleOAuthProvider clientId={import.meta.env.VITE_CLIENT_ID}>
                    <AppRouter />
                </GoogleOAuthProvider>
            </TokenRefreshProvider>
        </QueryClientProvider>
    );
}

export default App;