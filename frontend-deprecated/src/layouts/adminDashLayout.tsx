import React, { useState } from "react";
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";

const AdminDashLayout: React.FC = () => {
    const [sidebarOpen, setSidebarOpen] = useState<boolean>(true);
    const [settingsDropdownOpen, setSettingsDropdownOpen] = useState<boolean>(false);
    const [manageUsersDropdownOpen, setManageUsersDropdownOpen] = useState<boolean>(false);
    const location = useLocation();
    const navigate = useNavigate();

    const getLinkClassName = (path: string): string => {
        const baseClass = "mb-2 text-center text-lg font-semibold";
        const activeClass = "text-white bg-[rgb(56,68,94)] p-2 px-2 rounded-md";
        const inactiveClass = "text-gray-300 hover:text-white hover:bg-[rgb(56,68,94)] px-2 p-2 rounded-md";

        if (path === "/zen/dash") {
            // Exact match for dashboard
            return `${baseClass} ${location.pathname === path ? activeClass : inactiveClass}`;
        } else {
            // Starts with for other routes
            return `${baseClass} ${location.pathname.startsWith(path) ? activeClass : inactiveClass}`;
        }
    };

    const toggleSettingsDropdown = () => {
        setSettingsDropdownOpen(prevState => !prevState);
    };

    const toggleManageUsersDropdown = () => {
        setManageUsersDropdownOpen(prevState => !prevState);
    };

    //for the name
    const apiName = import.meta.env.VITE_API_NAME;
    const firstFourLetters = apiName.slice(0, 4);
    const remainingLetters = apiName.slice(4);

    const Logout = () => {
        const confirmResult = window.confirm("Do you want to logout?");

        if (confirmResult) {
            const cookies = Cookies.get("Cookies");

            if (cookies) {
                Cookies.remove("Cookies");
                navigate("/next/login");
            }
        } else {
            console.log("Logout canceled");
        }

        const confirmDialog: any = document.querySelector("div[role='dialog']");
        if (confirmDialog) {
            confirmDialog.style.backgroundColor = "white";
            confirmDialog.style.color = "black";
            confirmDialog.style.padding = "20px";
            confirmDialog.style.borderRadius = "5px";
            confirmDialog.style.boxShadow = "0 0 10px rgba(0, 0, 0, 0.3)";
        }
    };

    return (
        <div className="flex h-screen bg-gray-100">
            {/* Sidebar */}
            <div
                className={`${sidebarOpen ? "w-64" : "w-0"
                    } transition-all duration-300 bg-[rgb(26,46,68)]`}
            >
                {sidebarOpen && (
                    <nav className="p-4 text-white h-full">
                        <h2 className="text-xl font-bold mt-4 text-center mb-4">
                            <span>{firstFourLetters}</span>
                            <span className="text-blue-500">{remainingLetters}</span> <i className="bi bi-mailbox2-flag text-blue-500"></i>
                        </h2>
                        <ul className="mt-12 w-full">
                            <li className={getLinkClassName("/zen/dash")}>
                                <Link
                                    to="/zen/dash"
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-house-fill mr-2"></i> DashBoard
                                </Link>
                            </li>
                            <li className={getLinkClassName("/zen/dash/plan")}>
                                <Link
                                    to="/zen/dash/plan"
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-bar-chart-fill mr-2"></i> Plans
                                </Link>
                            </li>
                            <li className={`${getLinkClassName("/zen/dash/users")} relative`}>
                                <button
                                    onClick={toggleManageUsersDropdown}
                                    className="flex font-semibold text-base items-center w-full"
                                >
                                    <i className="bi bi-people-fill"></i> &nbsp; Manage Users
                                    <i className={`bi ${manageUsersDropdownOpen ? 'bi-chevron-up' : 'bi-chevron-down'} ml-auto`}></i>
                                </button>
                                {manageUsersDropdownOpen && (
                                    <ul className="pl-4 mt-2">
                                        <li className={getLinkClassName("/zen/dash/users")}>
                                            <Link to="/zen/dash/users" className="flex font-semibold text-sm items-center">
                                                Users
                                            </Link>
                                        </li>
                                        <li className={getLinkClassName("/zen/dash/email-users")}>
                                            <Link to="/zen/dash/email-users" className="flex font-semibold text-sm items-center">
                                                Email All Users
                                            </Link>
                                        </li>
                                    </ul>
                                )}
                            </li>
                            <li className={getLinkClassName("/zen/dash/support")}>
                                <Link
                                    to="/zen/dash/support"
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-headset"></i> &nbsp; Support
                                </Link>
                            </li>
                            <li className={getLinkClassName("/zen/dash/analytics")}>
                                <Link
                                    to="/zen/dash/analytics"
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-people-fill"></i> &nbsp; Analytics
                                </Link>
                            </li>
                            <li className={getLinkClassName("/zen/dash/billing")}>
                                <Link
                                    to="/zen/dash/billing"
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-credit-card"></i> &nbsp; Billing
                                </Link>
                            </li>
                        </ul>
                    </nav>
                )}
            </div>

            {/* Main content */}
            <div className="flex-1 flex flex-col">
                {/* Header */}
                <header className="bg-white h-12 p-4 flex justify-between items-center">
                    <button
                        onClick={() => setSidebarOpen(!sidebarOpen)}
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

                {/* Content area */}
                <main className="flex-1 p-6 overflow-auto">
                    <Outlet />
                </main>
            </div>
        </div>
    );
};

export default AdminDashLayout;