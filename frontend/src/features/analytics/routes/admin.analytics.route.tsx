import { RouteObject } from "react-router-dom";
import AdminAnalyticsDashTemplate from "../templates/admin/adminAnalyticsDashTemplate";

const adminAnalyticsRoute: RouteObject[] = [
    {
        index: true,
        element: <AdminAnalyticsDashTemplate />
    }
]

export default adminAnalyticsRoute