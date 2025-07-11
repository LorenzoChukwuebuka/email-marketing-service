import React, { useState, useEffect } from "react";
import { Link, Outlet, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import {
    Layout,
    Button,
    Dropdown,
    Badge,
    Avatar,
    Typography,
    Divider,
    Empty,
    Spin,
    Modal,
    Tag,
} from "antd";
import {
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    BellOutlined,
    UserOutlined,
    SettingOutlined,
    ApiOutlined,
    LogoutOutlined,
    NotificationOutlined,
    MailOutlined,
    CrownOutlined
} from "@ant-design/icons";
import { useNotificationQuery } from "./../../hooks/useNotificationQuery";
import eventBus from "../../utils/eventbus";
import { useMailCalcQuery } from "../../hooks/useMailDataQuery";
import useUserNotificationStore from "../../store/notification.store";

const { Header, Content } = Layout;
const { Text, Title } = Typography;

interface ContentProps {
    sidebarOpen: boolean;
    setSidebarOpen: (open: boolean) => void;
}

const DashContent: React.FC<ContentProps> = ({ sidebarOpen, setSidebarOpen }) => {
    const [userName, setUserName] = useState<string>("");
    const [isNotificationDropdownOpen, setIsNotificationDropdownOpen] = useState(false);
    const { data: notificationsData, isError: notificationError, isLoading: notificationLoader } = useNotificationQuery();
    const { data: mailData } = useMailCalcQuery();
    const navigate = useNavigate();
    const { updateReadStatus } = useUserNotificationStore();

    const hasNotifications = notificationsData && notificationsData.payload?.length > 0;
    const hasUnreadNotifications = hasNotifications && notificationsData.payload.some(notification => notification.read_status === false);

    if (notificationError) {
        eventBus.emit('error', 'Failed to fetch notifications');
    }

    const readNotifications = async () => {
        if (hasUnreadNotifications) {
            await updateReadStatus();
        }
    };

    const handleLogout = (): void => {
        Modal.confirm({
            title: "Confirm Logout",
            content: "Are you sure you want to logout?",
            okText: "Yes, Logout",
            cancelText: "Cancel",
            okButtonProps: { danger: true },
            onOk: () => {
                const cookies = Cookies.get("Cookies");
                if (cookies) {
                    Cookies.remove("Cookies");
                    navigate("/auth/login");
                }
            },
        });
    };

    useEffect(() => {
        const cookie = Cookies.get("Cookies");
        const user = cookie ? JSON.parse(cookie)?.details?.fullname : "";
        setUserName(user);
    }, []);

    // Notification dropdown content
    // const notificationContent = (
    //     <div className="w-80 max-h-96 overflow-hidden">
    //         <div className="p-4 border-b border-gray-100">
    //             <div className="flex items-center justify-between">
    //                 <Title level={5} className="!m-0 text-gray-800">
    //                     <NotificationOutlined className="mr-2 text-blue-500" />
    //                     Notifications
    //                 </Title>
    //                 {hasUnreadNotifications && (
    //                     <Badge
    //                         count={notificationsData?.payload?.filter(n => !n.read_status).length}
    //                         className="bg-red-500"
    //                     />
    //                 )}
    //             </div>
    //         </div>

    //         <div className="max-h-64 overflow-y-auto">
    //             {notificationLoader ? (
    //                 <div className="flex justify-center items-center py-8">
    //                     <Spin size="large" />
    //                 </div>
    //             ) : hasNotifications ? (
    //                 notificationsData?.payload?.map((notification, index) => (
    //                     <div key={index} className="p-3 border-b border-gray-50 last:border-b-0 hover:bg-gray-50 transition-colors">
    //                         <div className="flex items-start space-x-3">
    //                             <Avatar
    //                                 size="small"
    //                                 icon={<BellOutlined />}
    //                                 className="bg-blue-100 text-blue-600 mt-1"
    //                             />
    //                             <div className="flex-1 min-w-0">
    //                                 <Text className="font-medium text-gray-900 text-sm block">
    //                                     {notification.title}
    //                                 </Text>
    //                                 <Text className="text-xs text-gray-500 mt-1 block">
    //                                     {notification.created_at}
    //                                 </Text>
    //                             </div>
    //                             {!notification.read_status && (
    //                                 <div className="w-2 h-2 bg-blue-500 rounded-full mt-2"></div>
    //                             )}
    //                         </div>
    //                     </div>
    //                 ))
    //             ) : (
    //                 <div className="p-6">
    //                     <Empty
    //                         image={Empty.PRESENTED_IMAGE_SIMPLE}
    //                         description={
    //                             <Text className="text-gray-500">No new notifications</Text>
    //                         }
    //                     />
    //                 </div>
    //             )}
    //         </div>

    //         {hasNotifications && (
    //             <div className="p-3 border-t border-gray-100">
    //                 <Link to="/app/notifications">
    //                     <Button
    //                         type="link"
    //                         size="small"
    //                         className="w-full text-center text-blue-500 hover:bg-blue-50"
    //                     >
    //                         View all activity
    //                     </Button>
    //                 </Link>
    //             </div>
    //         )}
    //     </div>
    // );

    // User profile dropdown items
    const userMenuItems = [
        {
            key: 'profile-info',
            label: (
                <div className="px-3 py-2 border-b border-gray-100">
                    <div className="flex items-center space-x-3">
                        <Avatar
                            size="large"
                            icon={<UserOutlined />}
                            className="bg-gradient-to-br from-blue-500 to-purple-600"
                        />
                        <div>
                            <Text className="font-semibold text-gray-900 block">{userName}</Text>
                            <div className="flex items-center space-x-2 mt-1">
                                <MailOutlined className="text-gray-400 text-xs" />
                                <Text className="text-xs text-gray-500">
                                    {mailData?.payload?.remainingMails}/{mailData?.payload?.mailsPerDay} emails
                                </Text>
                            </div>
                            <Tag
                                icon={<CrownOutlined />}
                                color="gold"
                                className="mt-1 text-xs"
                            >
                                {mailData?.payload?.plan}
                            </Tag>
                        </div>
                    </div>
                </div>
            ),
        },
        {
            type: 'divider',
        },
        {
            key: 'profile',
            icon: <UserOutlined className="text-gray-600" />,
            label: <Link to="/app/settings/account-management">My Profile</Link>,
        },
        {
            key: 'api',
            icon: <ApiOutlined className="text-gray-600" />,
            label: <Link to="/app/settings/api">API & SMTP</Link>,
        },
        {
            key: 'settings',
            icon: <SettingOutlined className="text-gray-600" />,
            label: <Link to="/app/settings">Settings</Link>,
        },
        {
            type: 'divider',
        },
        {
            key: 'logout',
            icon: <LogoutOutlined className="text-red-500" />,
            label: (
                <Text className="text-red-500 font-medium">Logout</Text>
            ),
            onClick: handleLogout,
        },
    ];

    return (
        <Layout className="flex-1 bg-gray-50">
            <Header className="bg-white border-b border-gray-200 px-6 h-16 flex items-center justify-between shadow-sm">
                <div className="flex items-center space-x-4">
                    <Button
                        type="text"
                        icon={sidebarOpen ? <MenuFoldOutlined /> : <MenuUnfoldOutlined />}
                        onClick={() => setSidebarOpen(!sidebarOpen)}
                        className="text-gray-600 hover:text-gray-900 hover:bg-gray-100 border-0"
                        size="large"
                    />

                    <Divider type="vertical" className="h-6 bg-gray-300" />

                    <div className="hidden md:flex items-center space-x-2">
                        <Text className="text-gray-600 text-sm">Welcome back,</Text>
                        <Text className="font-semibold text-gray-900">{userName}</Text>
                    </div>
                </div>

                <div className="flex items-center space-x-4">
                
                    <Dropdown
                        open={isNotificationDropdownOpen}
                        onOpenChange={setIsNotificationDropdownOpen}
                        trigger={['click']}
                        placement="bottomRight"
                        dropdownRender={() => (
                            <div className="w-80 max-h-96 overflow-hidden bg-white shadow-lg rounded-md border">
                                <div className="p-4 border-b border-gray-100">
                                    <div className="flex items-center justify-between">
                                        <Title level={5} className="!m-0 text-gray-800">
                                            <NotificationOutlined className="mr-2 text-blue-500" />
                                            Notifications
                                        </Title>
                                        {hasUnreadNotifications && (
                                            <Badge
                                                count={notificationsData?.payload?.filter(n => !n.read_status).length}
                                                className="bg-red-500"
                                            />
                                        )}
                                    </div>
                                </div>

                                <div className="max-h-64 overflow-y-auto">
                                    {notificationLoader ? (
                                        <div className="flex justify-center items-center py-8">
                                            <Spin size="large" />
                                        </div>
                                    ) : hasNotifications ? (
                                        notificationsData?.payload?.map((notification, index) => (
                                            <div key={index} className="p-3 border-b border-gray-50 last:border-b-0 hover:bg-gray-50 transition-colors">
                                                <div className="flex items-start space-x-3">
                                                    <Avatar
                                                        size="small"
                                                        icon={<BellOutlined />}
                                                        className="bg-blue-100 text-blue-600 mt-1"
                                                    />
                                                    <div className="flex-1 min-w-0">
                                                        <Text className="font-medium text-gray-900 text-sm block">
                                                            {notification.title}
                                                        </Text>
                                                        <Text className="text-xs text-gray-500 mt-1 block">
                                                            {notification.created_at}
                                                        </Text>
                                                    </div>
                                                    {!notification.read_status && (
                                                        <div className="w-2 h-2 bg-blue-500 rounded-full mt-2"></div>
                                                    )}
                                                </div>
                                            </div>
                                        ))
                                    ) : (
                                        <div className="p-6">
                                            <Empty
                                                image={Empty.PRESENTED_IMAGE_SIMPLE}
                                                description={
                                                    <Text className="text-gray-500">No new notifications</Text>
                                                }
                                            />
                                        </div>
                                    )}
                                </div>

                                {hasNotifications && (
                                    <div className="p-3 border-t border-gray-100">
                                        <Link to="/app/notifications">
                                            <Button
                                                type="link"
                                                size="small"
                                                className="w-full text-center text-blue-500 hover:bg-blue-50"
                                            >
                                                View all activity
                                            </Button>
                                        </Link>
                                    </div>
                                )}
                            </div>
                        )}
                    >
                        <Badge
                            dot={hasUnreadNotifications}
                            offset={[-2, 2]}
                        >
                            <Button
                                type="text"
                                icon={<BellOutlined />}
                                className="text-gray-600 hover:text-gray-900 hover:bg-gray-100 border-0"
                                size="large"
                                onClick={readNotifications}
                            />
                        </Badge>
                    </Dropdown>

                    {/* User Profile Dropdown */}
                    <Dropdown
                        menu={{ items: userMenuItems as any}}
                        trigger={['click']}
                        placement="bottomRight"
                    >
                        <Button
                            type="text"
                            className="flex items-center space-x-2 text-gray-700 hover:text-gray-900 hover:bg-gray-100 border-0 h-10"
                        >
                            <Avatar
                                size="small"
                                icon={<UserOutlined />}
                                className="bg-gradient-to-br from-blue-500 to-purple-600"
                            />
                            <Text className="font-medium hidden sm:block">{userName}</Text>
                            <span className="text-xs text-gray-400">▼</span>
                        </Button>
                    </Dropdown>
                </div>
            </Header>

            <Content className="flex-1 overflow-auto bg-gradient-to-br from-gray-50 to-blue-50/20">
                <div className="p-6">
                    <Outlet />
                </div>
            </Content>

            <style dangerouslySetInnerHTML={{
                __html: `
                    .ant-layout-header {
                        padding: 0 24px !important;
                        line-height: 64px !important;
                    }
                    
                    .ant-dropdown-menu-item:hover {
                        background-color: rgba(59, 130, 246, 0.05) !important;
                    }
                    
                    .ant-dropdown-menu-item-selected {
                        background-color: rgba(59, 130, 246, 0.1) !important;
                    }
                    
                    .ant-badge-dot {
                        background-color: #ef4444 !important;
                        box-shadow: 0 0 0 2px #fff !important;
                    }
                    
                    .ant-empty-img-simple {
                        opacity: 0.3 !important;
                    }
                `
            }} />
        </Layout>
    );
};

export default DashContent;