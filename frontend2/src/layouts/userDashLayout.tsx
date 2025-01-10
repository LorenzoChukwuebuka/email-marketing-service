import React, { useState } from "react";


import Sidebar from "../components/dashboard/sidebar";
import DashContent from "../components/dashboard/dashcontent";

const UserDashLayout: React.FC = () => {
    const [sidebarOpen, setSidebarOpen] = useState<boolean>(true);
    const apiName = import.meta.env.VITE_API_NAME;

    return (
        <div className="flex h-screen bg-gray-100">
            <Sidebar sidebarOpen={sidebarOpen} apiName={apiName} />
            <DashContent sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />
        </div>
    );
};
export default UserDashLayout;

