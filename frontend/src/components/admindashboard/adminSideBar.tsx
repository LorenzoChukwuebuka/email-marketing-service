import React, { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { Badge, Typography } from "antd";
import {
    DashboardOutlined,
    BarChartOutlined,
    TeamOutlined,
    UserOutlined,
    MailOutlined,
    CustomerServiceOutlined,
    LineChartOutlined,
    CreditCardOutlined,
    DownOutlined,
    UpOutlined,
    FlagOutlined
} from "@ant-design/icons";
import { useAdminPendingTicketsQuery } from "../../features/support/hooks/useAdminSupportTicketQuery";

const { Title, Text } = Typography;

interface SidebarProps {
    apiName: string;
}

const AdminSidebar: React.FC<SidebarProps> = ({ apiName }) => {
    const [manageUsersDropdownOpen, setManageUsersDropdownOpen] = useState(false);
    const location = useLocation();

    const toggleManageUsersDropdown = () => {
        setManageUsersDropdownOpen((prev) => !prev);
    };

    const isActiveRoute = (path: string): boolean => {
        if (path === "/zen") {
            return location.pathname === path;
        }
        return location.pathname.startsWith(path);
    };

    const firstFourLetters = apiName.slice(0, 4);
    const remainingLetters = apiName.slice(4);

    const menuItems = [
        {
            key: "/zen",
            icon: <DashboardOutlined />,
            label: "Dashboard",
            path: "/zen"
        },
        {
            key: "/zen/plan",
            icon: <BarChartOutlined />,
            label: "Plans",
            path: "/zen/plan"
        },
        {
            key: "manage-users",
            icon: <TeamOutlined />,
            label: "Manage Users",
            isDropdown: true
        },
        {
            key: "/zen/support",
            icon: <CustomerServiceOutlined />,
            label: "Support",
            path: "/zen/support"
        },
        {
            key: "/zen/analytics",
            icon: <LineChartOutlined />,
            label: "Analytics",
            path: "/zen/analytics"
        },
        {
            key: "/zen/billing",
            icon: <CreditCardOutlined />,
            label: "Billing",
            path: "/zen/billing"
        }
    ];


    const { data: apiResponse } = useAdminPendingTicketsQuery(
        undefined,
        undefined,
        undefined
    );


    return (
        <div className="w-64 h-full bg-gradient-to-b from-slate-900 to-slate-800 shadow-xl">
            <div className="p-6 border-b border-slate-700">
                {/* Logo/Brand */}
                <div className="flex items-center justify-center mb-2">
                    <div className="text-center">
                        <Title level={3} className="!text-white !mb-0 !font-bold">
                            <span className="text-white">{firstFourLetters}</span>
                            <span className="text-blue-400">{remainingLetters}</span>
                        </Title>
                        <FlagOutlined className="text-blue-400 text-xl ml-2" />
                    </div>
                </div>
                <Text className="text-slate-300 text-sm block text-center">
                    Admin Panel
                </Text>
            </div>

            {/* Navigation */}
            <nav className="p-4 flex-1 overflow-y-auto">
                <div className="space-y-2">
                    {menuItems.map((item) => {
                        if (item.isDropdown) {
                            return (
                                <div key={item.key} className="space-y-1">
                                    <button
                                        onClick={toggleManageUsersDropdown}
                                        className={`w-full flex items-center justify-between px-4 py-3 rounded-lg transition-all duration-200 ${isActiveRoute("/zen/users")
                                            ? "bg-blue-600 text-white shadow-lg"
                                            : "text-slate-300 hover:bg-slate-700 hover:text-white"
                                            }`}
                                    >
                                        <div className="flex items-center space-x-3">
                                            {item.icon}
                                            <span className="font-medium">{item.label}</span>
                                        </div>
                                        {manageUsersDropdownOpen ? (
                                            <UpOutlined className="text-xs" />
                                        ) : (
                                            <DownOutlined className="text-xs" />
                                        )}
                                    </button>

                                    {/* Dropdown items */}
                                    <div className={`ml-4 space-y-1 transition-all duration-200 ${manageUsersDropdownOpen ? "max-h-96 opacity-100" : "max-h-0 opacity-0 overflow-hidden"
                                        }`}>
                                        <Link
                                            to="/zen/users"
                                            className={`flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-200 ${isActiveRoute("/zen/users")
                                                ? "bg-blue-500 text-white shadow-md"
                                                : "text-slate-400 hover:bg-slate-700 hover:text-white"
                                                }`}
                                        >
                                            <UserOutlined className="text-sm" />
                                            <span className="text-sm font-medium">Users</span>
                                        </Link>
                                        <Link
                                            to="/zen/users/email"
                                            className={`flex items-center space-x-3 px-4 py-2 rounded-lg transition-all duration-200 ${isActiveRoute("/zen/users/email")
                                                ? "bg-blue-500 text-white shadow-md"
                                                : "text-slate-400 hover:bg-slate-700 hover:text-white"
                                                }`}
                                        >
                                            <MailOutlined className="text-sm" />
                                            <span className="text-sm font-medium">Email All Users</span>
                                        </Link>
                                    </div>
                                </div>
                            );
                        }

                        return (
                            <Link
                                key={item.key}
                                to={item.path!}
                                className={`flex items-center space-x-3 px-4 py-3 rounded-lg transition-all duration-200 ${isActiveRoute(item.path!)
                                    ? "bg-blue-600 text-white shadow-lg transform scale-[1.02]"
                                    : "text-slate-300 hover:bg-slate-700 hover:text-white hover:transform hover:scale-[1.01]"
                                    }`}
                            >
                                {item.icon}
                                <span className="font-medium">{item.label}</span>
                                {item.key === "/zen/support" && (
                                    <Badge
                                        count={apiResponse?.payload.total}
                                        size="small"
                                        className="ml-auto"
                                        style={{ backgroundColor: '#f59e0b' }}
                                    />
                                )}
                            </Link>
                        );
                    })}
                </div>

                {/* Footer section */}
                <div className="mt-8 pt-4 border-t border-slate-700">
                    <div className="px-4 py-2">
                        <Text className="text-slate-400 text-xs">
                            Version 2.1.0
                        </Text>
                    </div>
                </div>
            </nav>
        </div>
    );
};

export default AdminSidebar;