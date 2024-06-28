import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Route, Routes } from "react-router-dom";
import { AuthRoute } from "./pages/Auth";
import IndexLandingPage from "./pages/landingPage/index";
import { useEffect } from "react";
import eventBus from "./utils/eventBus";
import { toast } from "react-toastify";
import { UserDashRoute } from "./layouts/userDashRoute";
function App() {
  const handleSuccess = (message) => {
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

  const handleError = (message) => {
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

  const handleInfo = (message) => {
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
    const successListener = (message) => handleSuccess(message);
    const errorListener = (message) => handleError(message);
    const infoListener = (message) => handleInfo(message);

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
    <>
      <ToastContainer />
      <Routes>
        <Route index element={<IndexLandingPage />} />
        <Route path="/auth/*" element={<AuthRoute />} />
        <Route path="/user/*" element={<UserDashRoute />} />
      </Routes>
    </>
  );
}

export default App;
