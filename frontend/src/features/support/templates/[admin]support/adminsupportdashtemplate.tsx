import { useState } from "react";
import { HelmetProvider, Helmet } from "react-helmet-async";
import { Tabs, Card, Badge } from "antd";
import {
    FileTextOutlined,
    ClockCircleOutlined,
    CheckCircleOutlined,
    InboxOutlined
} from "@ant-design/icons";
import AllSupportTicketComponentTable from "../../components/[admin]support/allSupportTicketsComponent";
import PendingSupportTicketComponentTable from "../../components/[admin]support/PendingSupportTicketsComponent";
import ClosedSupportTicketComponentTable from "../../components/[admin]support/closedSupportTicketsComponent";

type TabType = "all" | "pending" | "closed";

const SupportDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "all" || storedTab === "closed" || storedTab === "pending") ? storedTab : "all";
    });

    const handleTabChange = (key: string) => {
        const tab = key as TabType;
        setActiveTab(tab);
        localStorage.setItem("activeTab", tab);
    };

    const getTitle = () => {
        switch (activeTab) {
            case "all":
                return "All Support Tickets";
            case "pending":
                return "Pending Support Tickets";
            case "closed":
                return "Closed Support Tickets";
            default:
                return "Support Dashboard";
        }
    };

    const tabItems = [
        {
            key: "all",
            label: (
                <div className="flex items-center gap-2 px-2">
                    <FileTextOutlined className="text-blue-500" />
                    <span className="font-medium">All Tickets</span>
                </div>
            ),
            children: <AllSupportTicketComponentTable />,
        },
        {
            key: "pending",
            label: (
                <div className="flex items-center gap-2 px-2">
                    <ClockCircleOutlined className="text-orange-500" />
                    <span className="font-medium">Pending</span>
                    <Badge
                        count={0}
                        showZero={false}
                        className="ml-1"
                        style={{ backgroundColor: '#f59e0b' }}
                    />
                </div>
            ),
            children: <PendingSupportTicketComponentTable />,
        },
        {
            key: "closed",
            label: (
                <div className="flex items-center gap-2 px-2">
                    <CheckCircleOutlined className="text-green-500" />
                    <span className="font-medium">Closed</span>
                </div>
            ),
            children: <ClosedSupportTicketComponentTable />,
        },
    ];

    return (
        <HelmetProvider>
            <Helmet title={getTitle()} />
            <div className="min-h-screen bg-gray-50">
                <div className="max-w-7xl mx-auto p-6">
                    {/* Header */}
                    <div className="mb-8">
                        <div className="flex items-center gap-3 mb-2">
                            <div className="p-2 bg-blue-100 rounded-lg">
                                <InboxOutlined className="text-2xl text-blue-600" />
                            </div>
                            <div>
                                <h1 className="text-3xl font-bold text-gray-900">
                                    Support Dashboard
                                </h1>
                                <p className="text-gray-600 mt-1">
                                    Manage and track all support tickets
                                </p>
                            </div>
                        </div>
                    </div>

                    {/* Main Content Card */}
                    <Card
                        className="shadow-lg border-0 rounded-xl"
                        bodyStyle={{ padding: 0 }}
                    >
                        <Tabs
                            activeKey={activeTab}
                            onChange={handleTabChange}
                            items={tabItems}
                            className="modern-tabs"
                            size="large"
                            tabBarStyle={{
                                padding: '0 24px',
                                margin: 0,
                                borderBottom: '1px solid #f0f0f0',
                                backgroundColor: '#fafafa'
                            }}
                        />
                    </Card>
                </div>
            </div>


            <div
                dangerouslySetInnerHTML={{
                    __html: `
                        <style>
                            .modern-tabs .ant-tabs-tab {
                                padding: 16px 8px !important;
                                margin-right: 32px !important;
                                border: none !important;
                                background: transparent !important;
                            }
                            
                            .modern-tabs .ant-tabs-tab:hover {
                                background: rgba(59, 130, 246, 0.05) !important;
                                border-radius: 8px !important;
                            }
                            
                            .modern-tabs .ant-tabs-tab-active {
                                background: rgba(59, 130, 246, 0.1) !important;
                                border-radius: 8px !important;
                            }
                            
                            .modern-tabs .ant-tabs-tab-active .ant-tabs-tab-btn {
                                color: #3b82f6 !important;
                            }
                            
                            .modern-tabs .ant-tabs-ink-bar {
                                background: #3b82f6 !important;
                                height: 3px !important;
                                border-radius: 2px !important;
                            }
                            
                            .modern-tabs .ant-tabs-content-holder {
                                padding: 24px !important;
                                background: white !important;
                            }
                        </style>
                    `
                }}
            />
        </HelmetProvider>
    );
};

export default SupportDashTemplate;