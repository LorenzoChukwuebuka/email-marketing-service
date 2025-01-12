import { RouteObject } from "react-router-dom";
import HelpAndSupport from "../templates/supportTemplate";
import SupportTicketTemplate from "../templates/supportTicketTemplate";
import SupportTicketDetailsTemplate from "../templates/ticketDetailsTemplate";

const supportRoute: RouteObject[] = [
    {
        index: true,
        element: <HelpAndSupport />
    },
    {
        path: "ticket",
        element: <SupportTicketTemplate />
    },
    {
        path:"ticket/details/:id",
        element: <SupportTicketDetailsTemplate/>
    }
]

export default supportRoute