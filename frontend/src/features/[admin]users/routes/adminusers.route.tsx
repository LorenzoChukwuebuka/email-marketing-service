import { RouteObject } from "react-router-dom";
import AdminUserDashTemplate from '../templates/user/userDashTemplate';
import AdminUserDetailTemplate from "../templates/user/userDetailsTemplate";
import AdminSendUsersMail from "../templates/user/sendMailToUsersTemplate";

const adminUsersRoute: RouteObject[] = [
    {
        index: true,
        element: <AdminUserDashTemplate />
    },
    {
        path: "detail/:id",
        element: <AdminUserDetailTemplate />
    },
    {
        path: "email",
        element: <AdminSendUsersMail />
    },
    {
        path:"campaigns",
        element:<> hello world </>
    }
]


export default adminUsersRoute