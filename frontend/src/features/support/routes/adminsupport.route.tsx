import { RouteObject } from "react-router-dom";
import SupportDashTemplate from "../templates/[admin]support/adminsupportdashtemplate";
import AdminTicketDetails from "../components/[admin]support/ticketDetailscomponent";

const adminSupportRoute: RouteObject[] = [
    {
        index: true,
        element: <SupportDashTemplate />
    },
    {
        path: "details/:id",
        element: <AdminTicketDetails />
    }
]

export default adminSupportRoute