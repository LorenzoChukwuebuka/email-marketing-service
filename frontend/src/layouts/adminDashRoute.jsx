import { Route, Routes } from "react-router-dom";
import AdminDashLayout from "./adminDashLayout";

const AdminDashRoute = () => (
  <Routes>
    <Route path="dash" element={<AdminDashLayout />} />
  </Routes>
);

export { AdminDashRoute };
