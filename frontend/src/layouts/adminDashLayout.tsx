import React, { useState } from "react";
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";

const AdminDashLayout: React.FC = () => {
    const [sidebarOpen, setSidebarOpen] = useState<boolean>(true);
    const [settingsDropdownOpen, setSettingsDropdownOpen] = useState<boolean>(false);
    const location = useLocation();
    const navigate = useNavigate();

    const getLinkClassName = (path: string): string => {
        const baseClass = "mb-4 text-center text-lg font-semibold";
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
                            CrabMailer
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
                            <li className={getLinkClassName("/zen/dash/users")}>
                                <Link
                                    to="/zen/dash/users"
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-people-fill"></i> &nbsp; Users
                                </Link>
                            </li>
                            {/* <li className={`${getLinkClassName("")} relative`}>
                <button
                  onClick={toggleSettingsDropdown}
                  className="flex items-center w-full justify-between"
                >
                  <span className="flex font-semibold text-base items-center">
                    <i className="bi bi-gear mr-2"></i> Settings
                  </span>
                  <i
                    className={`bi ${
                      settingsDropdownOpen ? "bi-chevron-up" : "bi-chevron-down"
                    } ml-2`}
                  ></i>
                </button>
                {settingsDropdownOpen && (
                  <ul className="mt-2 bg-[rgb(36,56,78)] rounded-md p-2">
                    <li className={`py-1 ${getLinkClassName("")}`}>
                      <Link
                        to=""
                        className="block  text-sm hover:bg-[rgb(56,68,94)] rounded"
                      >
                        User Management
                      </Link>
                    </li>
                    <li className={`py-1 ${getLinkClassName("")}`}>
                      <Link
                        to=""
                        className="block  text-sm hover:bg-[rgb(56,68,94)] rounded"
                      >
                        API Tokens
                      </Link>
                    </li>
                    <li className={`py-1 ${getLinkClassName("")}`}>
                      <Link
                        to=""
                        className="block  text-sm hover:bg-[rgb(56,68,94)] rounded"
                      >
                        Account Settings
                      </Link>
                    </li>
                  </ul>
                )}
              </li> */}
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
