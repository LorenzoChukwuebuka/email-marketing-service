import { useEffect, useState } from "react";
import { Card, Tabs, Typography, Empty, Spin } from "antd";
import { 
    AppstoreOutlined, 
    FileAddOutlined, 
    CodeOutlined, 
    EditOutlined,
    SmileOutlined 
} from "@ant-design/icons";
import CreateTransactionalTemplateComponent from '../../components/transactional/createTransactionalComponent';

const { Title, Text } = Typography;

type templateTypes = "Templates Gallery" | "Blank Template" | "Custom HTML" | "Rich Text";

const CreateTransactionalTemplateDashBoard: React.FC = () => {
    const [activeTab, setActiveTab] = useState<templateTypes>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "Templates Gallery" || storedTab === "Blank Template" || storedTab === "Custom HTML" || storedTab === "Rich Text") ? storedTab : "Templates Gallery";
    });
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(false);

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as templateTypes);
        }
    }, []);

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setIsLoading(true);
        setTimeout(() => {
            setIsLoading(false);
        }, 30000);
    };

    const handleTabChange = (newTab: templateTypes) => {
        if (isLoading) {
            setIsLoading(false);
        }
        setActiveTab(newTab);
        if (newTab !== "Templates Gallery") {
            setIsModalOpen(true);
        }
    };

    const tabItems = [
        {
            key: 'Templates Gallery',
            label: (
                <div className="flex items-center gap-2">
                    <AppstoreOutlined />
                    <span>Templates Gallery</span>
                </div>
            ),
            children: (
                <div className="py-8">
                    <Empty
                        image={<SmileOutlined className="text-6xl text-blue-400" />}
                        imageStyle={{
                            height: 80,
                            display: 'flex',
                            justifyContent: 'center',
                            alignItems: 'center',
                        }}
                        description={
                            <div className="text-center">
                                <Text className="text-lg font-medium text-gray-600 block mb-2">
                                    Template Gallery
                                </Text>
                                <Text className="text-gray-500">
                                    Template Gallery... coming soon
                                </Text>
                            </div>
                        }
                    />
                </div>
            ),
        },
        {
            key: 'Blank Template',
            label: (
                <div className="flex items-center gap-2">
                    <FileAddOutlined />
                    <span>Blank Template</span>
                </div>
            ),
            children: null,
        },
        {
            key: 'Custom HTML',
            label: (
                <div className="flex items-center gap-2">
                    <CodeOutlined />
                    <span>Custom HTML</span>
                </div>
            ),
            children: null,
        },
        {
            key: 'Rich Text',
            label: (
                <div className="flex items-center gap-2">
                    <EditOutlined />
                    <span>Text Editor</span>
                </div>
            ),
            children: null,
        },
    ];

    return (
        <div className="min-h-screen bg-gradient-to-br from-gray-50 to-blue-50/20 p-6">
            <div className="max-w-7xl mx-auto">
                {/* Header */}
                <div className="mb-8">
                    <Title level={2} className="!mb-2 text-gray-800">
                        Create Transactional Templates
                    </Title>
                    <Text className="text-gray-600">
                        Choose your preferred template creation method
                    </Text>
                </div>

                {/* Main Content Card */}
                <Card className="shadow-lg border-0 rounded-xl overflow-hidden">
                    <Tabs
                        activeKey={activeTab}
                        onChange={(key) => handleTabChange(key as templateTypes)}
                        items={tabItems}
                        size="large"
                        className="custom-tabs"
                    />

                    {/* Loading State */}
                    {isLoading && (
                        <div className="flex items-center justify-center py-20">
                            <Spin size="large" />
                        </div>
                    )}
                </Card>

                {/* Modal */}
                <CreateTransactionalTemplateComponent 
                    isOpen={isModalOpen} 
                    onClose={handleCloseModal} 
                    editorType={
                        activeTab === "Blank Template"
                            ? "drag-and-drop"
                            : activeTab === "Custom HTML"
                                ? "html-editor"
                                : "rich-text"
                    } 
                />
            </div>

            <style dangerouslySetInnerHTML={{
                __html: `
                    .custom-tabs .ant-tabs-tab {
                        padding: 12px 24px !important;
                        font-weight: 500 !important;
                        border-radius: 8px 8px 0 0 !important;
                        margin-right: 4px !important;
                        transition: all 0.3s ease !important;
                        background: transparent !important;
                        border: none !important;
                    }

                    .custom-tabs .ant-tabs-tab:hover {
                        background: rgba(59, 130, 246, 0.05) !important;
                        color: #3b82f6 !important;
                    }

                    .custom-tabs .ant-tabs-tab-active {
                        background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%) !important;
                        color: white !important;
                    }

                    .custom-tabs .ant-tabs-tab-active:hover {
                        background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%) !important;
                        color: white !important;
                    }

                    .custom-tabs .ant-tabs-ink-bar {
                        display: none !important;
                    }

                    .custom-tabs .ant-tabs-content-holder {
                        padding: 0 !important;
                    }

                    .custom-tabs .ant-tabs-tabpane {
                        padding: 0 !important;
                    }
                `
            }} />
        </div>
    );
};

export default CreateTransactionalTemplateDashBoard;