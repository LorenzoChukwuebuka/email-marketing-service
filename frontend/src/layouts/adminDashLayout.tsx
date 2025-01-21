import React, { useState } from "react";

import AdminSidebar from "../components/admindashboard/adminSideBar";
import AdminMainContent from "../components/admindashboard/adminmaincontent";

const AdminDashLayout: React.FC = () => {
    const [sidebarOpen, setSidebarOpen] = useState(true);
    const apiName = import.meta.env.VITE_API_NAME;

    return (
        <div className="flex h-screen bg-gray-100">
            {sidebarOpen && <AdminSidebar apiName={apiName} />}
            <AdminMainContent
                sidebarOpen={sidebarOpen}
                toggleSidebar={() => setSidebarOpen((prev) => !prev)}
            />
        </div>
    );
};

export default AdminDashLayout;
