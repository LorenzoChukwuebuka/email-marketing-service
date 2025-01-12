import { RouteObject } from "react-router-dom";
import BillingDashTemplate from "../templates/billingTemplate";

const billingRoutes: RouteObject[] = [
    {
        index: true,
        element: <BillingDashTemplate />

    }
]

export default billingRoutes