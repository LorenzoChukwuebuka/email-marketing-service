import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Route, Routes } from "react-router-dom";
import { AuthRoute } from "./pages/Auth";
import IndexLandingPage from "./pages/landingPage/index";
import { useEffect } from "react";
import eventBus from "./utils/eventBus";
import { toast } from "react-toastify";
import { UserDashRoute } from "./layouts/userDashRoute";
import { ProtectedRoute } from "./utils/protectedRoute";
import { AdminAuthRoute } from "./pages/admin";
import { AdminDashRoute } from "./layouts/adminDashRoute";
import useTokenExpiration from './hooks/usetokenExpiration';
import PaymentSuccessPage from "./pages/userDashboard/paymentSuccesspage";

const App: React.FC = () => {
    const handleSuccess = (message: string) => {
        toast.success(message, {
            position: "top-right",
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
            position: "top-right",
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
            position: "top-right",
            autoClose: 5000,
            hideProgressBar: false,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
            progress: undefined,
            theme: "dark",
        });
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

    useTokenExpiration()

    return (
        <>


            <ToastContainer />
            <Routes>
                <Route index element={<IndexLandingPage />} />
                <Route path="/auth/*" element={<AuthRoute />} />
                <Route path="/next/*" element={<AdminAuthRoute />} />

                <Route element={<ProtectedRoute />}>
                    <Route path="/user/*" element={<UserDashRoute />} />
                    <Route path="/zen/*" element={<AdminDashRoute />} />
                    <Route path="/payment" element={<PaymentSuccessPage />} />
                </Route>
            </Routes>
        </>
    );
}

export default App;
