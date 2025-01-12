import { RouteObject } from "react-router-dom";
import TemplateBuilderDashComponent from "../templates";
import CreateMarketingTemplateDashBoard from "../templates/marketing/createMarketingTemplate";
import CreateTransactionalTemplateDashBoard from "../templates/transactional/createTransactionalTemplate";

const emailTemplateRoute: RouteObject[] = [
    {
        index: true,
        element: <TemplateBuilderDashComponent />
    },
    {
        path:"marketing",
        element:<CreateMarketingTemplateDashBoard/>
    },
    {
        path:"transactional",
        element:<CreateTransactionalTemplateDashBoard/>
    }
]

export default emailTemplateRoute