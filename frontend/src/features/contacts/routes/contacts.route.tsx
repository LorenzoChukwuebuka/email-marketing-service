
import { RouteObject } from "react-router-dom";
import ContactDashTemplate from "../templates";
import ContactGroupListTemplate from "../templates/contactgrouplist";

const contactRoute: RouteObject[] = [

    {
        index: true,
        element: <ContactDashTemplate />
    },
    {
        path: "view-group",
        element: <ContactGroupListTemplate/>

    }
]


export default contactRoute