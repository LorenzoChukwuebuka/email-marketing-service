
import { RouteObject } from "react-router-dom";
import ContactDashTemplate from "../templates";
import GroupContactList from "../components/contactgroup/groupContactListComponent";

const contactRoute: RouteObject[] = [

    {
        index: true,
        element: <ContactDashTemplate />
    },
    {
        path: "view-group",
        element: <GroupContactList/>

    }
]


export default contactRoute