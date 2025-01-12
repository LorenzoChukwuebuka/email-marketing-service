import React, { useState } from "react";
import { Link, useLocation } from "react-router-dom";

interface SidebarProps {
    sidebarOpen: boolean;
    apiName: string;
}

const Sidebar: React.FC<SidebarProps> = ({ sidebarOpen, apiName }) => {
    const [settingsDropdownOpen, setSettingsDropdownOpen] = useState<boolean>(false);
    const location = useLocation();

    const getLinkClassName = (path: string): string => {
        const baseClass = "mb-2 text-center text-lg font-semibold";
        const activeClass = "text-white bg-[rgb(56,68,94)] p-2 px-2 rounded-md";
        const inactiveClass = "text-gray-300 hover:text-white hover:bg-[rgb(56,68,94)] px-2 p-2 rounded-md";

        if (path === "/app") {
            return `${baseClass} ${location.pathname === path ? activeClass : inactiveClass}`;
        } else {
            return `${baseClass} ${location.pathname.startsWith(path) ? activeClass : inactiveClass}`;
        }
    };

    const toggleSettingsDropdown = (): void => {
        setSettingsDropdownOpen(!settingsDropdownOpen);
    };

    const firstFourLetters = apiName.slice(0, 4);
    const remainingLetters = apiName.slice(4);

    return (
        <div className={`${sidebarOpen ? "w-64" : "w-0"} transition-all duration-300 bg-[rgb(26,46,68)]`}>
            {sidebarOpen && (
                <nav className="p-4 text-white h-full">
                    <h2 className="text-xl p-4 font-bold mt-0 cursor-pointer mb-2">
                        <Link to="/app">
                            <span>{firstFourLetters}</span>
                            <span className="text-blue-500">{remainingLetters}</span>
                            <i className="bi bi-mailbox2-flag text-blue-500"></i>
                        </Link>
                    </h2>
                    <ul className="mt-5 w-full">
                        <li className={getLinkClassName("/app")}>
                            <Link to="/app" className="flex font-semibold text-base items-center">
                                <i className="bi bi-house-fill mr-2"></i> Dashboard
                            </Link>
                        </li>

                        <li className={getLinkClassName("/app/campaign")}>
                            <Link to="/app/campaign" className="flex font-semibold text-base items-center">
                                <i className="bi bi-megaphone-fill"></i> &nbsp; Campaigns
                            </Link>
                        </li>

                        <li className={getLinkClassName("/app/contacts")}>
                            <Link to="/app/contacts" className="flex font-semibold text-base items-center">
                                <i className="bi bi-person-lines-fill"></i> &nbsp; Contacts
                            </Link>
                        </li>

                        <li className={getLinkClassName("/app/templates")}>
                            <Link to="/app/templates" className="flex font-semibold text-base items-center">
                                <i className="bi bi-stack"></i> &nbsp; Templates
                            </Link>
                        </li>

                        <li className={getLinkClassName("/app/analytics")}>
                            <Link to="/app/analytics" className="flex font-semibold text-base items-center">
                                <i className="bi bi-bar-chart-fill mr-2"></i> Analytics
                            </Link>
                        </li>

                        <li className="bg-white h-[1px] mt-3 mb-3"></li>

                        <li className={getLinkClassName("/app/billing")}>
                            <Link to="/app/billing" className="flex font-semibold text-base items-center">
                                <i className="bi bi-wallet-fill"></i> &nbsp; Billing
                            </Link>
                        </li>

                        <li className={getLinkClassName("/app/support")}>
                            <Link to="/app/support" className="flex font-semibold text-base items-center">
                                <i className="bi bi-headset"></i> &nbsp; Help & Support
                            </Link>
                        </li>

                        <li className={`${getLinkClassName("/user/setting")} relative`}>
                            <button onClick={toggleSettingsDropdown} className="flex items-center w-full justify-between">
                                <span className="flex font-semibold text-base items-center">
                                    <i className="bi bi-gear-fill mr-2"></i> Settings
                                </span>
                                <i className={`bi ${settingsDropdownOpen ? "bi-chevron-up" : "bi-chevron-down"} ml-2`}></i>
                            </button>
                            {settingsDropdownOpen && (
                                <ul className="mt-2 bg-[rgb(36,56,78)] rounded-md p-2">
                                    <li className={`py-1 ${getLinkClassName("/app/settings/api")}`}>
                                        <Link to="/app/settings/api" className="block text-sm hover:bg-[rgb(56,68,94)] rounded">
                                            API Tokens
                                        </Link>
                                    </li>
                                    <li className={`py-1 ${getLinkClassName("/app/settings/account-management")}`}>
                                        <Link to="/app/settings/account-management" className="block text-sm hover:bg-[rgb(56,68,94)] rounded">
                                            Account Settings
                                        </Link>
                                    </li>
                                    <li className={`py-1 ${getLinkClassName("/app/settings/domain")}`}>
                                        <Link to="/app/settings/domain" className="block text-sm hover:bg-[rgb(56,68,94)] rounded">
                                            Senders and Domain
                                        </Link>
                                    </li>
                                </ul>
                            )}
                        </li>
                    </ul>
                </nav>
            )}
        </div>
    );
};

export default Sidebar;