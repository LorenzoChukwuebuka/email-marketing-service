import { RouteObject } from "react-router-dom";
import AdminGetAllUserCampaign from "../templates/[admin]/getUserCampaignTemplate";
import AdminGetSpecificUserCampaignTemplate from "../templates/[admin]/getSpecificUserCampaign";

const adminCampaignRoute: RouteObject[] = [
    {
        path: "details/:id/:companyId",
        element: <AdminGetAllUserCampaign />
    },

    {
        path: "specific/:campaignId/:userId/:companyId",
        element: <AdminGetSpecificUserCampaignTemplate />
    }
]

export default adminCampaignRoute