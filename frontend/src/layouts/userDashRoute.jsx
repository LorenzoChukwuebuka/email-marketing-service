import { Route, Routes } from "react-router-dom";
import UserDashLayout from "./userDashLayout";
import UserDashPage from "../pages/userDashboard/userHomeDashPage";

const UserDashRoute = () => (
  <Routes>
    <Route path="dash" element={<UserDashLayout />}>
      <Route path="" element={<UserDashPage />} />
    </Route>
  </Routes>
);

export { UserDashRoute };
