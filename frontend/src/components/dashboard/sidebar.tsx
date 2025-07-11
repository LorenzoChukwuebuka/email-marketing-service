import { Link, useLocation } from "react-router-dom";
import { Menu, Layout, Typography, Button } from "antd";
import {
    HomeOutlined,
    SoundOutlined,
    ContactsOutlined,
    FileTextOutlined,
    BarChartOutlined,
    WalletOutlined,
    CustomerServiceOutlined,
    SettingOutlined,
    ApiOutlined,
    UserOutlined,
    CloudOutlined,
    MailOutlined,
    MenuFoldOutlined,
    MenuUnfoldOutlined
} from "@ant-design/icons";

const { Sider } = Layout;
const { Title, Text } = Typography;

interface SidebarProps {
    sidebarOpen: boolean;
    apiName: string;
    onToggle?: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ sidebarOpen, apiName, onToggle }) => {
    const location = useLocation();

    const firstFourLetters = apiName.slice(0, 4);
    const remainingLetters = apiName.slice(4);

    const getSelectedKeys = (): string[] => {
        const path = location.pathname;
        if (path === "/app") return ["dashboard"];
        if (path.startsWith("/app/campaign")) return ["campaigns"];
        if (path.startsWith("/app/contacts")) return ["contacts"];
        if (path.startsWith("/app/templates")) return ["templates"];
        if (path.startsWith("/app/analytics")) return ["analytics"];
        if (path.startsWith("/app/billing")) return ["billing"];
        if (path.startsWith("/app/support")) return ["support"];
        if (path.startsWith("/app/settings/api")) return ["api-tokens"];
        if (path.startsWith("/app/settings/account-management")) return ["account-settings"];
        if (path.startsWith("/app/settings/domain")) return ["domain-settings"];
        if (path.startsWith("/app/settings")) return ["settings"];
        return [];
    };

    const getOpenKeys = (): string[] => {
        const path = location.pathname;
        if (path.startsWith("/app/settings")) return ["settings"];
        return [];
    };

    const menuItems = [
        {
            key: 'dashboard',
            icon: <HomeOutlined />,
            label: <Link to="/app">Dashboard</Link>,
        },
        {
            key: 'campaigns',
            icon: <SoundOutlined />,
            label: <Link to="/app/campaign">Campaigns</Link>,
        },
        {
            key: 'contacts',
            icon: <ContactsOutlined />,
            label: <Link to="/app/contacts">Contacts</Link>,
        },
        {
            key: 'templates',
            icon: <FileTextOutlined />,
            label: <Link to="/app/templates">Templates</Link>,
        },
        {
            key: 'analytics',
            icon: <BarChartOutlined />,
            label: <Link to="/app/analytics">Analytics</Link>,
        },
        {
            type: 'divider',
        },
        {
            key: 'billing',
            icon: <WalletOutlined />,
            label: <Link to="/app/billing">Billing</Link>,
        },
        {
            key: 'support',
            icon: <CustomerServiceOutlined />,
            label: <Link to="/app/support">Help & Support</Link>,
        },
        {
            key: 'settings',
            icon: <SettingOutlined />,
            label: 'Settings',
            children: [
                {
                    key: 'api-tokens',
                    icon: <ApiOutlined />,
                    label: <Link to="/app/settings/api">API Tokens</Link>,
                },
                {
                    key: 'account-settings',
                    icon: <UserOutlined />,
                    label: <Link to="/app/settings/account-management">Account Settings</Link>,
                },
                {
                    key: 'domain-settings',
                    icon: <CloudOutlined />,
                    label: <Link to="/app/settings/domain">Senders and Domain</Link>,
                },
            ],
        },
    ];

    return (
        <>
            <Sider
                trigger={null}
                collapsible
                collapsed={!sidebarOpen}
                width={280}
                collapsedWidth={64}
                className="relative shadow-xl transition-all duration-300 ease-in-out"
                style={{
                    background: 'linear-gradient(180deg, #1e293b 0%, #0f172a 100%)',
                    borderRight: '1px solid rgba(255, 255, 255, 0.1)',
                }}
            >
                <div className="h-full flex flex-col">
                    {/* Header */}
                    <div className={`${sidebarOpen ? 'p-6' : 'p-4'} border-b border-white/10`}>
                        <div className="flex items-center justify-between">
                            <Link to="/app" className="flex items-center space-x-2">
                                <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
                                    <MailOutlined className="text-white text-lg" />
                                </div>
                                {sidebarOpen && (
                                    <div>
                                        <Title level={4} className="!m-0 !text-white">
                                            <span>{firstFourLetters}</span>
                                            <span className="text-blue-400">{remainingLetters}</span>
                                        </Title>
                                    </div>
                                )}
                            </Link>

                            {onToggle && sidebarOpen && (
                                <Button
                                    type="text"
                                    icon={<MenuFoldOutlined />}
                                    onClick={onToggle}
                                    className="text-white/70 hover:text-white hover:bg-white/10 border-0"
                                />
                            )}
                        </div>

                        {/* Collapsed state toggle button */}
                        {onToggle && !sidebarOpen && (
                            <div className="flex justify-center mt-4">
                                <Button
                                    type="text"
                                    icon={<MenuUnfoldOutlined />}
                                    onClick={onToggle}
                                    className="text-white/70 hover:text-white hover:bg-white/10 border-0"
                                />
                            </div>
                        )}
                    </div>

                    {/* Navigation Menu */}
                    <div className="flex-1 py-4">
                        <Menu
                            mode="inline"
                            selectedKeys={getSelectedKeys()}
                            defaultOpenKeys={getOpenKeys()}
                            items={menuItems as any}
                            className="bg-transparent border-0"
                            inlineCollapsed={!sidebarOpen}
                            style={{
                                backgroundColor: 'transparent',
                            }}
                        />
                    </div>

                    {/* Footer */}
                    {sidebarOpen && (
                        <div className="p-4 border-t border-white/10">
                            <div className="flex items-center space-x-3 p-3 rounded-lg bg-white/5 backdrop-blur-sm">
                                <div className="w-8 h-8 bg-gradient-to-br from-green-400 to-blue-500 rounded-full flex items-center justify-center">
                                    <UserOutlined className="text-white text-sm" />
                                </div>
                                <div className="flex-1 min-w-0">
                                    <Text className="text-white font-medium text-sm block truncate">
                                        Welcome back
                                    </Text>
                                    <Text className="text-white/60 text-xs block truncate">
                                        {apiName}
                                    </Text>
                                </div>
                            </div>
                        </div>
                    )}
                </div>

                <style dangerouslySetInnerHTML={{
                    __html: `
                        .ant-menu-inline {
                            background: transparent !important;
                            border-right: none !important;
                        }
                        
                        .ant-menu-item,
                        .ant-menu-submenu-title {
                            margin: 4px 12px !important;
                            width: calc(100% - 24px) !important;
                            border-radius: 8px !important;
                            color: rgba(255, 255, 255, 0.8) !important;
                            font-weight: 500 !important;
                            transition: all 0.3s ease !important;
                        }
                        
                        .ant-menu-inline-collapsed .ant-menu-item,
                        .ant-menu-inline-collapsed .ant-menu-submenu-title {
                            margin: 4px 8px !important;
                            width: calc(100% - 16px) !important;
                            text-align: center !important;
                            padding: 0 !important;
                        }
                        
                        .ant-menu-item:hover,
                        .ant-menu-submenu-title:hover {
                            background: rgba(255, 255, 255, 0.1) !important;
                            color: white !important;
                        }
                        
                        .ant-menu-item-selected {
                            background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%) !important;
                            color: white !important;
                            box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3) !important;
                        }
                        
                        .ant-menu-item-selected:hover {
                            background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%) !important;
                            color: white !important;
                        }
                        
                        .ant-menu-submenu-selected .ant-menu-submenu-title {
                            color: #3b82f6 !important;
                        }
                        
                        .ant-menu-sub {
                            background: rgba(0, 0, 0, 0.2) !important;
                            border-radius: 8px !important;
                            margin: 4px 12px !important;
                            width: calc(100% - 24px) !important;
                        }
                        
                        .ant-menu-sub .ant-menu-item {
                            margin: 2px 8px !important;
                            width: calc(100% - 16px) !important;
                            padding-left: 24px !important;
                        }
                        
                        .ant-menu-submenu-arrow {
                            color: rgba(255, 255, 255, 0.6) !important;
                        }
                        
                        .ant-menu-item a,
                        .ant-menu-submenu-title a {
                            color: inherit !important;
                            text-decoration: none !important;
                        }
                        
                        .ant-menu-item-icon {
                            font-size: 16px !important;
                        }
                        
                        .ant-divider {
                            background: rgba(255, 255, 255, 0.1) !important;
                            margin: 8px 16px !important;
                        }
                        
                        .ant-menu-inline-collapsed .ant-menu-submenu-title {
                            padding: 0 !important;
                        }
                        
                        .ant-menu-inline-collapsed .ant-menu-item-icon {
                            margin: 0 !important;
                        }
                    `
                }} />
            </Sider>
        </>
    );
};

export default Sidebar;