import { RouteObject } from "react-router-dom";
import SignupTemplate from "../templates/signuptemplate";
import LoginTemplate from "../templates/loginTemplate";
import ForgetPasswordTemplate from "../templates/forgetpasswordtemplate";
import VerifyAccountTemplate from "../templates/verifyaccountTemplate";
import ResetPasswordTemplate from "../templates/resetpasswordtemplate";
import GoogleAuthCallbackTemplate from "../templates/googleauthcallbacktemplate";

const authRoute: RouteObject[] = [
    {
        path: "sign-up",
        element: <SignupTemplate />
    },
    {
        path: "login",
        element: <LoginTemplate />
    },
    {
        path: "forgot-password",
        element: <ForgetPasswordTemplate />
    },
    {
        path: "account-verification",
        element: <VerifyAccountTemplate />
    },
    {
        path: "reset-password",
        element: <ResetPasswordTemplate />
    },
    {
        path: "callback",
        element: <GoogleAuthCallbackTemplate />
    },
    {
        path:"signup/callback"
    }
]

export default authRoute