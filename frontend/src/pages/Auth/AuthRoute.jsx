import { Route, Routes } from "react-router-dom";
import { LoginPage, ForgotPassword, ResetPassword } from "./components";
import SignUpPage from "./components/signupPage";
import OTPPage from "./components/OTP";

const AuthRoute = () => (
  <Routes>
    <Route path="login" element={<LoginPage />} />
    <Route path="forgot-password" element={<ForgotPassword />} />
    <Route path="reset-password" element={<ResetPassword />} />
    <Route path="sign-up" element={<SignUpPage />} />
    <Route path="account-verification" element={<OTPPage />} />
  </Routes>
);

export { AuthRoute };
