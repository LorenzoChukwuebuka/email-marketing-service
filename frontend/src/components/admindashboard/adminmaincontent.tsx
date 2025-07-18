import React from "react";
import { Outlet, useNavigate } from "react-router-dom";
import { Button, Dropdown, Avatar, Space, Typography, Modal } from "antd";
import { MenuOutlined, LogoutOutlined, UserOutlined, SettingOutlined } from "@ant-design/icons";
import Cookies from "js-cookie";

const { Title } = Typography;

interface MainContentProps {
    sidebarOpen: boolean;
    toggleSidebar: () => void;
}

const AdminMainContent: React.FC<MainContentProps> = ({ sidebarOpen, toggleSidebar }) => {
    const navigate = useNavigate();

    const handleLogout = () => {
        Modal.confirm({
            title: 'Confirm Logout',
            content: 'Are you sure you want to logout?',
            okText: 'Yes, Logout',
            okType: 'danger',
            cancelText: 'Cancel',
            centered: true,
            onOk: () => {
                const cookies = Cookies.get("Cookies");
                if (cookies) {
                    Cookies.remove("Cookies");
                    navigate("/next/login");
                }
            },
        });
    };

    const userMenuItems = [
        {
            key: 'profile',
            icon: <UserOutlined />,
            label: 'Profile',
            onClick: () => {
                // Navigate to profile page
                console.log('Navigate to profile');
            },
        },
        {
            key: 'settings',
            icon: <SettingOutlined />,
            label: 'Settings',
            onClick: () => {
                // Navigate to settings page
                console.log('Navigate to settings');
            },
        },
        {
            type: 'divider',
        },
        {
            key: 'logout',
            icon: <LogoutOutlined />,
            label: 'Logout',
            onClick: handleLogout,
            danger: true,
        },
    ];

    return (
        <div className="flex-1 flex flex-col bg-gray-50">
            {/* Header */}
            <header className="bg-white shadow-sm border-b border-gray-200 h-16 px-6 flex justify-between items-center">
                {/* Left side - Menu toggle and title */}
                <div className="flex items-center space-x-4">
                    <Button
                        type="text"
                        icon={<MenuOutlined />}
                        onClick={toggleSidebar}
                        className="text-gray-600 hover:text-gray-800 hover:bg-gray-100 flex items-center justify-center w-10 h-10"
                        size="large"
                    />
                    <Title level={4} className="!mb-0 text-gray-800">
                        Dashboard
                    </Title>
                </div>

                {/* Right side - User menu */}
                <div className="flex items-center space-x-4">
                    {/* Notification bell (optional) */}
                    <Button
                        type="text"
                        icon={<i className="bi bi-bell text-lg"></i>}
                        className="text-gray-600 hover:text-gray-800 hover:bg-gray-100 flex items-center justify-center w-10 h-10"
                        size="large"
                    />
                    
                    {/* User dropdown */}
                    <Dropdown
                        menu={{ items: userMenuItems as any }}
                        trigger={['click']}
                        placement="bottomRight"
                        arrow
                    >
                        <Space className="cursor-pointer hover:bg-gray-100 px-3 py-2 rounded-lg transition-colors">
                            <Avatar 
                                icon={<UserOutlined />} 
                                className="bg-blue-500"
                                size={32}
                            />
                            <span className="text-gray-700 font-medium hidden sm:inline">
                                Admin
                            </span>
                        </Space>
                    </Dropdown>
                </div>
            </header>

            {/* Main content area */}
            <main className="flex-1 p-6 overflow-auto bg-gray-50">
                <div className="max-w-7xl mx-auto">
                    <Outlet />
                </div>
            </main>
        </div>
    );
};

export default AdminMainContent;