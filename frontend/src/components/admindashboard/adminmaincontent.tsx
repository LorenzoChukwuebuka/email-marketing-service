import React from "react";
import { Outlet, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";

interface MainContentProps {
    sidebarOpen: boolean;
    toggleSidebar: () => void;
}

const AdminMainContent: React.FC<MainContentProps> = ({ sidebarOpen, toggleSidebar }) => {
    const navigate = useNavigate();

    const Logout = () => {
        if (window.confirm("Do you want to logout?")) {
            const cookies = Cookies.get("Cookies");
            if (cookies) {
                Cookies.remove("Cookies");
                navigate("/next/login");
            }
        }
    };

    return (
        <div className="flex-1 flex flex-col">
            <header className="bg-white h-12 p-4 flex justify-between items-center">
                <button
                    onClick={toggleSidebar}
                    className="text-gray-500 hover:text-gray-700"
                >
                    <span style={{ fontSize: "24px" }}>{sidebarOpen ? "≡" : "☰"}</span>
                </button>
                <h1 className="text-xl font-semibold">Home</h1>
                <button
                    className="hover:bg-blue-200 hover:rounded-btn hover:text-blue-500 font-semibold p-1"
                    onClick={Logout}
                >
                    Log Out
                </button>
            </header>
            <main className="flex-1 p-6 overflow-auto">
                <Outlet />
            </main>
        </div>
    );
};

export default AdminMainContent;
