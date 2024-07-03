import { Route, Routes } from "react-router-dom";
import UserDashLayout from "./userDashLayout";
import UserDashPage from "../pages/userDashboard/userHomeDashPage";
import UserSettingPage from "../pages/userDashboard/userSettingsPage";

const UserDashRoute = () => (
  <Routes>
    <Route path="dash" element={<UserDashLayout />}>
      <Route index element={<UserDashPage />} />
      <Route path="setting" element={<UserSettingPage />} />
    </Route>
  </Routes>
);

export { UserDashRoute };
