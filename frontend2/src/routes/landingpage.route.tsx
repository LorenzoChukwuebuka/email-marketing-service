import { RouteObject } from "react-router-dom";
import IndexLandingPage from "../pages/landingPage";
import PrivacyPolicy from "../pages/landingPage/privacypage";
import TermsOfService from "../pages/landingPage/tos";
import GDPRExplanation from "../pages/landingPage/gdpr";

const landingPageRoute: RouteObject[] = [
    {
        index: true,
        element: <IndexLandingPage />
    },
    {
        path: "privacy",
        element: <PrivacyPolicy />
    },
    {
        path: "tos",
        element: <TermsOfService />
    },
    {
        path: "gdpr",
        element: <GDPRExplanation />
    }
]

export default landingPageRoute