import React, { useEffect, useState } from "react";
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import reactSVG from "./../assets/0832855c4b75f4d5a1dd822d6fb0d38d.jpg";
import useDailyUserMailSentCalc from "../store/userstore/userDashStore";


const UserDashLayout: React.FC = () => {
    const [sidebarOpen, setSidebarOpen] = useState<boolean>(true);
    const [userName, setUserName] = useState<string>("");
    const [settingsDropdownOpen, setSettingsDropdownOpen] = useState<boolean>(false);
    const location = useLocation();
    const navigate = useNavigate();

    const { mailData, getUserMailData } = useDailyUserMailSentCalc();

    const getLinkClassName = (path: string): string => {
        const baseClass = "mb-4 text-center text-lg font-semibold";
        const activeClass = "text-white bg-[rgb(56,68,94)] p-2 px-2 rounded-md";
        const inactiveClass =
            "text-gray-300 hover:text-white hover:bg-[rgb(56,68,94)] px-2 p-2 rounded-md";

        if (path === "/user/dash") {
            return `${baseClass} ${location.pathname === path ? activeClass : inactiveClass}`;
        } else {
            return `${baseClass} ${location.pathname.startsWith(path) ? activeClass : inactiveClass}`;
        }
    };

    const toggleSettingsDropdown = (): void => {
        setSettingsDropdownOpen(!settingsDropdownOpen);
    };

    const Logout = (): void => {
        const confirmResult = confirm("Do you want to logout?");

        if (confirmResult) {
            let cookies = Cookies.get("Cookies");

            if (cookies) {
                Cookies.remove("Cookies");
                navigate("/auth/login");
            }
        } else {
            console.log("Logout canceled");
        }

        const confirmDialog = document.querySelector("div[role='dialog']");
        if (confirmDialog) {
            (confirmDialog as HTMLElement).style.backgroundColor = "white";
            (confirmDialog as HTMLElement).style.color = "black";
            (confirmDialog as HTMLElement).style.padding = "20px";
            (confirmDialog as HTMLElement).style.borderRadius = "5px";
            (confirmDialog as HTMLElement).style.boxShadow = "0 0 10px rgba(0, 0, 0, 0.3)";
        }
    };

    useEffect(() => {
        let cookie = Cookies.get("Cookies");
        let user = cookie ? JSON.parse(cookie)?.details?.fullname : "";
        setUserName(user);
    }, []);

    useEffect(() => {
        getUserMailData();
    }, [getUserMailData]);

    return (
        <div className="flex h-screen bg-gray-100">
            {/* Sidebar */}
            <div
                className={`${sidebarOpen ? "w-64" : "w-0"
                    } transition-all duration-300 bg-[rgb(26,46,68)]`}
            >
                {sidebarOpen && (
                    <nav className="p-4 text-white h-full">
                        <h2 className="text-xl font-bold mt-4 text-center cursor-pointer mb-4">
                            <Link to="/user/dash"> {import.meta.env.VITE_API_NAME} </Link>
                        </h2>
                        <ul className="mt-8 w-full">
                            <li className={getLinkClassName("/user/dash")}>
                                <Link
                                    to="/user/dash"
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-house-fill mr-2"></i> Dashboard
                                </Link>
                            </li>

                            <li className={getLinkClassName("/campaigns")}>
                                <Link
                                    to=""
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-megaphone-fill"></i> &nbsp; Campaigns
                                </Link>
                            </li>

                            <li className={getLinkClassName("/user/dash/contacts")}>
                                <Link
                                    to="/user/dash/contacts"
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-person-lines-fill"></i> &nbsp; Contacts
                                </Link>
                            </li>

                            <li className={getLinkClassName("/user/dash/templates")}>
                                <Link
                                    to="/user/dash/templates"
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-stack"></i> &nbsp; Templates
                                </Link>
                            </li>

                            <li className={getLinkClassName("/analytics")}>
                                <Link
                                    to=""
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-bar-chart-fill mr-2"></i> Analytics
                                </Link>
                            </li>
                            <li className={getLinkClassName("/user/dash/billing")}>
                                <Link
                                    to="/user/dash/billing"
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-wallet-fill"></i> &nbsp; Billing
                                </Link>
                            </li>
                            <li className={getLinkClassName("/support")}>
                                <Link
                                    to=""
                                    className="flex font-semibold text-base items-center"
                                >
                                    <i className="bi bi-headset"></i> &nbsp; Help & Support
                                </Link>
                            </li>
                            <li className={`${getLinkClassName("/user/setting")} relative`}>
                                <button
                                    onClick={toggleSettingsDropdown}
                                    className="flex items-center w-full justify-between"
                                >
                                    <span className="flex font-semibold text-base items-center">
                                        <i className="bi bi-gear-fill mr-2"></i> Settings
                                    </span>
                                    <i
                                        className={`bi ${settingsDropdownOpen ? "bi-chevron-up" : "bi-chevron-down"
                                            } ml-2`}
                                    ></i>
                                </button>
                                {settingsDropdownOpen && (
                                    <ul className="mt-2 bg-[rgb(36,56,78)] rounded-md p-2">
                                        {/* <li
                                            className={`py-1 ${getLinkClassName(
                                                "/user/dash/settings/user-management"
                                            )}`}
                                        >
                                            <Link
                                                to="/user/dash/settings/user-management"
                                                className="block  text-sm hover:bg-[rgb(56,68,94)] rounded"
                                            >
                                                User Management
                                            </Link>
                                        </li> */}
                                        <li
                                            className={`py-1 ${getLinkClassName(
                                                "/user/dash/settings/api"
                                            )}`}
                                        >
                                            <Link
                                                to="/user/dash/settings/api"
                                                className="block  text-sm hover:bg-[rgb(56,68,94)] rounded"
                                            >
                                                API Tokens
                                            </Link>
                                        </li>
                                        <li
                                            className={`py-1 ${getLinkClassName(
                                                "/user/dash/settings/account-management"
                                            )}`}
                                        >
                                            <Link
                                                to="/user/dash/settings/account-management"
                                                className="block  text-sm hover:bg-[rgb(56,68,94)] rounded"
                                            >
                                                Account Settings
                                            </Link>
                                        </li>
                                    </ul>
                                )}
                            </li>
                        </ul>
                    </nav>
                )}
            </div>

            {/* Main content */}
            <div className="flex-1 flex flex-col">
                {/* Header */}
                <header className="bg-white h-14 p-4 flex justify-between items-center">
                    <button
                        onClick={() => setSidebarOpen(!sidebarOpen)}
                        className="text-gray-500 hover:text-gray-700"
                    >
                        <span style={{ fontSize: "24px" }}>{sidebarOpen ? "≡" : "☰"}</span>
                    </button>
                    {/* <h1 className="text-xl font-semibold"> Dashboard </h1> */}

                    <div className="dropdown dropdown-end">
                        <div tabIndex={0} role="button" className="m-1">
                            {userName}
                        </div>
                        <ul
                            tabIndex={0}
                            className="dropdown-content menu bg-white rounded-box z-[50] mt-4 w-52 p-2 shadow"
                        >
                            <li>
                                <div className="flex flex-col items-center">
                                    <span className="text-base-300 border-b-2 border-black mb-4">
                                        Emails sent: {mailData?.remainingMails}/
                                        {mailData?.mailsPerDay}
                                    </span>
                                    <span className="text-black bg-gray-300 rounded-md">
                                        Plan: {mailData?.plan}
                                    </span>
                                    <img
                                        className="h-8 w-8 rounded-full"
                                        src={reactSVG}
                                        alt="User avatar"
                                    />
                                    {userName}
                                    <span className="text-blue-500">
                                        <Link to="/user/setting/profile"> My Profile </Link>
                                    </span>
                                    <a onClick={Logout}>Logout</a>
                                </div>
                            </li>
                        </ul>
                    </div>
                </header>

                {/* Content area */}
                <main className="flex-1 p-2 w-full overflow-auto">
                    <Outlet />
                </main>
            </div>
        </div>
    );
};

export default UserDashLayout;