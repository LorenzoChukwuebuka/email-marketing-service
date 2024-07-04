import { Route, Routes } from "react-router-dom";
import UserDashLayout from "./userDashLayout";
import UserDashPage from "../pages/userDashboard/userHomeDashPage";
import { APISettingsDashTemplate } from "../templates";
import UserMgtSettingsDashTemplate from "../templates/user/SettingsTemplate/userManagementTemplate";

const UserDashRoute = () => (
  <Routes>
    <Route path="dash" element={<UserDashLayout />}>
      <Route index element={<UserDashPage />} />
      <Route path="settings/api" element={<APISettingsDashTemplate />} />
      <Route
        path="settings/user-management"
        element={<UserMgtSettingsDashTemplate />}
      />
    </Route>
  </Routes>
);

export { UserDashRoute };
