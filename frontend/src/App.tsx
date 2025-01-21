import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { toast } from "react-toastify";
import eventBus from "./utils/eventbus";
import { useEffect } from "react";
import { queryClient } from "./utils/queryclient";
import AppRouter from "./routes";
import { QueryClientProvider } from "@tanstack/react-query";
import { tokenRefreshOptions } from "./features/auth/hooks/useTokenRefreshQuery";
import { TokenRefreshProvider } from "./features/auth/hooks/useTokenRefreshProvider";
import Cookies from "js-cookie";
import { GoogleOAuthProvider } from "@react-oauth/google";


queryClient.setQueryDefaults(tokenRefreshOptions.queryKey, tokenRefreshOptions); // Pre-configure the token refresh query

function App() {
    const handleSuccess = (message: string) => {
        toast.success(message, {
            position: "bottom-right",
            autoClose: 5000,
            hideProgressBar: false,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
            progress: undefined,
            theme: "light",
        });
    };

    const handleError = (message: string) => {
        toast.error(message, {
            position: "bottom-right",
            autoClose: 5000,
            hideProgressBar: false,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
            progress: undefined,
            theme: "dark",
        });
    };

    const handleInfo = (message: string) => {
        toast.info(message, {
            position: "bottom-right",
            autoClose: 5000,
            hideProgressBar: false,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
            progress: undefined,
            theme: "dark",
        });
    };

    const handleRefreshError = (error: Error) => {
        console.error('Token refresh failed:', error);
        Cookies.remove('Cookies', { path: '/' });
        window.location.href = "/auth/login"
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
    }, [])


    return (
        <QueryClientProvider client={queryClient}>
            <ToastContainer />
            <TokenRefreshProvider onRefreshError={handleRefreshError}>
                <GoogleOAuthProvider clientId={import.meta.env.VITE_CLIENT_ID}>
                    <AppRouter />
                </GoogleOAuthProvider>
            </TokenRefreshProvider>
        </QueryClientProvider>
    )

}



export default App
