import { RouteObject } from "react-router-dom";
import AdminGetAllUserCampaign from "../templates/[admin]/getUserCampaignTemplate";

const adminCampaignRoute: RouteObject[] = [
    {
        path: "details/:id",
        element: <AdminGetAllUserCampaign />
    },

    // {
    //     path:""
    // }
]

export default adminCampaignRoute