import { useState } from "react";
import { Helmet, HelmetProvider } from "react-helmet-async";
import { Tabs, Card } from "antd";
import { UserOutlined, CheckCircleOutlined, ClockCircleOutlined, TeamOutlined } from "@ant-design/icons";
import AllUsersTable from "../../components/allUserTableComponent";
import UnVerifiedUsersTable from "../../components/unverifiedUsersComponent";
import VerifiedUsersTable from "../../components/verifiedUsersTableComponent";

type TabType = "All" | "Verified" | "Unverified";

const AdminUserDashTemplate = () => {
    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "All" || storedTab === "Verified" || storedTab === "Unverified") ? storedTab : "All";
    });

    const handleTabChange = (key: string) => {
        const tab = key as TabType;
        setActiveTab(tab);
        localStorage.setItem("activeTab", tab);
    };

    const getTitle = () => {
        switch (activeTab) {
            case "All":
                return "All Users";
            case "Verified":
                return "Verified Users";
            case "Unverified":
                return "Unverified Users";
            default:
                return "";
        }
    };

    const tabItems = [
        {
            key: "All",
            label: (
                <div className="flex items-center gap-2 px-2">
                    <TeamOutlined className="text-lg" />
                    <span className="font-medium">All Users</span>
                </div>
            ),
            children: (
                <div className="animate-fade-in">
                    <AllUsersTable />
                </div>
            ),
        },
        {
            key: "Verified",
            label: (
                <div className="flex items-center gap-2 px-2">
                    <CheckCircleOutlined className="text-lg text-green-600" />
                    <span className="font-medium">Verified Users</span>
                </div>
            ),
            children: (
                <div className="animate-fade-in">
                    <VerifiedUsersTable />
                </div>
            ),
        },
        {
            key: "Unverified",
            label: (
                <div className="flex items-center gap-2 px-2">
                    <ClockCircleOutlined className="text-lg text-orange-500" />
                    <span className="font-medium">Unverified Users</span>
                </div>
            ),
            children: (
                <div className="animate-fade-in">
                    <UnVerifiedUsersTable />
                </div>
            ),
        },
    ];

    return (
        <HelmetProvider>
            <Helmet title={getTitle()} />
            
            <div className="min-h-screen bg-gradient-to-br from-slate-50 to-gray-100">
                <div className="p-6 max-w-7xl mx-auto">
                    {/* Header Section */}
                    <div className="mb-8">
                        <div className="flex items-center gap-4 mb-4">
                            <div className="flex items-center justify-center w-12 h-12 bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl shadow-lg">
                                <UserOutlined className="text-white text-xl" />
                            </div>
                            <div>
                                <h1 className="text-3xl font-bold text-gray-800 mb-1">
                                    User Management
                                </h1>
                                <p className="text-gray-600">
                                    Manage and monitor user accounts across your platform
                                </p>
                            </div>
                        </div>
                    </div>

                    {/* Main Content Card */}
                    <Card 
                        className="shadow-xl border-0 rounded-2xl overflow-hidden"
                        bodyStyle={{ padding: 0 }}
                    >
                        <div className="bg-white">
                            <Tabs
                                activeKey={activeTab}
                                onChange={handleTabChange}
                                size="large"
                                className="modern-tabs"
                                tabBarStyle={{
                                    margin: 0,
                                    padding: '24px 24px 0 24px',
                                    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                                    borderRadius: '16px 16px 0 0',
                                }}
                                items={tabItems}
                            />
                        </div>
                    </Card>
                </div>
            </div>

            <style
                dangerouslySetInnerHTML={{
                    __html: `
                        .modern-tabs .ant-tabs-tab {
                            background: rgba(255, 255, 255, 0.1);
                            border: 1px solid rgba(255, 255, 255, 0.2);
                            border-radius: 12px;
                            margin-right: 8px;
                            padding: 12px 20px;
                            transition: all 0.3s ease;
                            backdrop-filter: blur(10px);
                        }

                        .modern-tabs .ant-tabs-tab:hover {
                            background: rgba(255, 255, 255, 0.15);
                            transform: translateY(-2px);
                        }

                        .modern-tabs .ant-tabs-tab-active {
                            background: rgba(255, 255, 255, 0.95);
                            border-color: rgba(255, 255, 255, 0.3);
                            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
                            transform: translateY(-2px);
                        }

                        .modern-tabs .ant-tabs-tab-active .ant-tabs-tab-btn {
                            color: #4f46e5;
                        }

                        .modern-tabs .ant-tabs-tab .ant-tabs-tab-btn {
                            color: rgba(255, 255, 255, 0.9);
                            font-weight: 500;
                        }

                        .modern-tabs .ant-tabs-ink-bar {
                            display: none;
                        }

                        .modern-tabs .ant-tabs-content-holder {
                            padding: 32px;
                            background: white;
                        }

                        .animate-fade-in {
                            animation: fadeIn 0.5s ease-in-out;
                        }

                        @keyframes fadeIn {
                            from {
                                opacity: 0;
                                transform: translateY(20px);
                            }
                            to {
                                opacity: 1;
                                transform: translateY(0);
                            }
                        }
                    `
                }}
            />
        </HelmetProvider>
    );
};

export default AdminUserDashTemplate;