import { RouteObject } from "react-router-dom";
import CampaignDashTemplate from '../templates/campaignDashTemplate';
import EditCampaignTemplate from "../templates/editCampaignTemplate";
import CampaignReportTemplate from "../templates/campaignReportTemplate";

const campaignRoute: RouteObject[] = [
    {
        index: true,
        element: <CampaignDashTemplate />
    },

    {
        path: "edit/:id",
        element: <EditCampaignTemplate />
    },
    {
        path: "report/:id",
        element: <CampaignReportTemplate />
    }
]

export default campaignRoute

