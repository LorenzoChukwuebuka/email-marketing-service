import React from "react";
import { Route, Routes } from "react-router-dom";
import AdminDashLayout from "./adminDashLayout";
import AdminDashPage from "../pages/admin/pages/adminDash";
import AdminPlanPage from "../pages/admin/pages/adminPlan";
import AdminUsersPage from "../pages/admin/pages/adminUsers";
import AdminUserDetailComponent from "../templates/admin/templates/components/users/getUserDetailsComponent";
import SupportDashTemplate from "../templates/admin/templates/support/supportDashTemplate";

const AdminDashRoute: React.FC = () => (
    <Routes>
        <Route path="dash" element={<AdminDashLayout />}>
            <Route index element={<AdminDashPage />} />
            <Route path="plan" element={<AdminPlanPage />} />
            <Route path="users" element={<AdminUsersPage />} />
            <Route path="users/detail/:id" element={<AdminUserDetailComponent />} />
            <Route path="support" element={<SupportDashTemplate />} />
        </Route>
    </Routes>
);

export { AdminDashRoute };
