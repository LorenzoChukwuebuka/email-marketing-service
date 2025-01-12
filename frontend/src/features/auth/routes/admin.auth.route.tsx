import { RouteObject } from "react-router-dom";
import AdminLoginTemplate from "../templates/admin/admin.logintemplate";

const adminAuthRoute: RouteObject[] = [
    {
        index: true,
        element: <AdminLoginTemplate />
    }
]

export default adminAuthRoute