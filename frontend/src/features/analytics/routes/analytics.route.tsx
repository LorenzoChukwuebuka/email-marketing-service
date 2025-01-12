import { RouteObject } from "react-router-dom";
import AnalyticsTemplateDash from "../templates/analyticsTemplate";

const analyticsRoute: RouteObject[] = [
    {
        index: true,
        element: <AnalyticsTemplateDash />
    }
]

export default analyticsRoute