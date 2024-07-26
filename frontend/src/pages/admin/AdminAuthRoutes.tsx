import { Route, Routes } from "react-router-dom";
import AdminLogin from "./auth/AdminLogin";

const AdminAuthRoute: React.FC = () => (
    <Routes>
        <Route path="login" element={<AdminLogin />} />
    </Routes>
);

export { AdminAuthRoute };
