import React, { useState } from "react";
import { Link, useLocation } from "react-router-dom";

interface SidebarProps {
    apiName: string;
}

const AdminSidebar: React.FC<SidebarProps> = ({ apiName }) => {
    const [manageUsersDropdownOpen, setManageUsersDropdownOpen] = useState(false);
    const location = useLocation();

    const getLinkClassName = (path: string): string => {
        const baseClass = "mb-2 text-center text-lg font-semibold";
        const activeClass = "text-white bg-[rgb(56,68,94)] p-2 px-2 rounded-md";
        const inactiveClass = "text-gray-300 hover:text-white hover:bg-[rgb(56,68,94)] px-2 p-2 rounded-md";

        if (path === "/zen") {
            return `${baseClass} ${location.pathname === path ? activeClass : inactiveClass}`;
        } else {
            return `${baseClass} ${location.pathname.startsWith(path) ? activeClass : inactiveClass}`;
        }
    };

    const toggleManageUsersDropdown = () => {
        setManageUsersDropdownOpen((prev) => !prev);
    };

    const firstFourLetters = apiName.slice(0, 4);
    const remainingLetters = apiName.slice(4);

    return (
        <div className="w-64 transition-all duration-300 bg-[rgb(26,46,68)]">
            <nav className="p-4 text-white h-full">
                <h2 className="text-xl font-bold mt-4 text-center mb-4">
                    <span>{firstFourLetters}</span>
                    <span className="text-blue-500">{remainingLetters}</span> 
                    <i className="bi bi-mailbox2-flag text-blue-500"></i>
                </h2>
                <ul className="mt-12 w-full">
                    <li className={getLinkClassName("/zen")}>
                        <Link to="/zen" className="flex font-semibold text-base items-center">
                            <i className="bi bi-house-fill mr-2"></i> DashBoard
                        </Link>
                    </li>
                    <li className={getLinkClassName("/zen/plan")}>
                        <Link to="/zen/plan" className="flex font-semibold text-base items-center">
                            <i className="bi bi-bar-chart-fill mr-2"></i> Plans
                        </Link>
                    </li>
                    <li className={`${getLinkClassName("/zen/users")} relative`}>
                        <button
                            onClick={toggleManageUsersDropdown}
                            className="flex font-semibold text-base items-center w-full"
                        >
                            <i className="bi bi-people-fill"></i> &nbsp; Manage Users
                            <i className={`bi ${manageUsersDropdownOpen ? 'bi-chevron-up' : 'bi-chevron-down'} ml-auto`}></i>
                        </button>
                        {manageUsersDropdownOpen && (
                            <ul className="pl-4 mt-2">
                                <li className={getLinkClassName("/zen/users")}>
                                    <Link to="/zen/users" className="flex font-semibold text-sm items-center">
                                        Users
                                    </Link>
                                </li>
                                <li className={getLinkClassName("/zen/email-users")}>
                                    <Link to="/zen/users/email" className="flex font-semibold text-sm items-center">
                                        Email All Users
                                    </Link>
                                </li>
                            </ul>
                        )}
                    </li>
                    <li className={getLinkClassName("/zen/support")}>
                        <Link to="/zen/support" className="flex font-semibold text-base items-center">
                            <i className="bi bi-headset"></i> &nbsp; Support
                        </Link>
                    </li>
                    <li className={getLinkClassName("/zen/analytics")}>
                        <Link to="/zen/analytics" className="flex font-semibold text-base items-center">
                            <i className="bi bi-people-fill"></i> &nbsp; Analytics
                        </Link>
                    </li>
                    <li className={getLinkClassName("/zen/billing")}>
                        <Link to="/zen/billing" className="flex font-semibold text-base items-center">
                            <i className="bi bi-credit-card"></i> &nbsp; Billing
                        </Link>
                    </li>
                </ul>
            </nav>
        </div>
    );
};

export default AdminSidebar;
