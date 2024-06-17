import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Route, Routes } from "react-router-dom";
import { AuthRoute } from "./pages/Auth";
import IndexLandingPage from "./pages/landingPage/index";

function App() {
  return (
    <>
      <ToastContainer />
      
      <Routes>
        <Route index element={<IndexLandingPage />} />
        <Route path="/auth/*" element={<AuthRoute />} />
      </Routes>
    </>
  );
}

export default App;
