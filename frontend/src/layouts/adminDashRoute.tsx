import React from "react";
import { Route, Routes } from "react-router-dom";
import AdminDashLayout from "./adminDashLayout";
import AdminDashPage from "../pages/admin/pages/adminDash";
import AdminPlanPage from "../pages/admin/pages/adminPlan";
import AdminUsersPage from "../pages/admin/pages/adminUsers";
import AdminUserDetailComponent from "../templates/admin/templates/components/users/getUserDetailsComponent";
import SupportDashTemplate from "../templates/admin/templates/support/supportDashTemplate";
import AdminTicketDetails from "../templates/admin/templates/components/support/ticketDetailscomponent";
import AdminSendUsersMail from "../templates/admin/templates/user/sendMailToUsersTemplate";
import AdminUserCampaigns from "../templates/admin/templates/campaigns/getAllUserCampaigns";

const AdminDashRoute: React.FC = () => (
    <Routes>
        <Route path="dash" element={<AdminDashLayout />}>
            <Route index element={<AdminDashPage />} />
            <Route path="plan" element={<AdminPlanPage />} />
            <Route path="users" element={<AdminUsersPage />} />
            <Route path="users/detail/:id" element={<AdminUserDetailComponent />} />
            <Route path="email-users" element={<AdminSendUsersMail />} />
            <Route path="support" element={<SupportDashTemplate />} />
            <Route path="support/details/:id" element={<AdminTicketDetails />} />
            <Route path="users/campaigns/:userid" element={<AdminUserCampaigns />} />
        </Route>
    </Routes>
);

export { AdminDashRoute };
