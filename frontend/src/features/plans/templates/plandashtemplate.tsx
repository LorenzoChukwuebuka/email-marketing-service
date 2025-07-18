import { FormEvent, useState } from "react";
import { Helmet, HelmetProvider } from "react-helmet-async";
import { Button, Typography, Card, Popconfirm, Badge, message } from "antd";
import { PlusOutlined, DeleteOutlined, SettingOutlined } from "@ant-design/icons";
import usePlanStore from "../store/plan.store";
import GetAllPlans from "../components/AllPlans";
import CreatePlan from "../components/createPlan";

const { Title, Text } = Typography;

const PlansDashTemplate: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [deleteLoading, setDeleteLoading] = useState<boolean>(false);
    const { selectedId, deletePlan } = usePlanStore();

    const handleDelete = async (e: FormEvent<HTMLButtonElement>) => {
        try {
            e.preventDefault();
            setDeleteLoading(true);

            await deletePlan();
            message.success(`Successfully deleted ${selectedId.length} plan(s)`);
        } catch (error) {
            console.log(error);
            message.error("Failed to delete plan(s)");
        } finally {
            setDeleteLoading(false);
        }
    };

    const handleCreatePlan = () => {
        setIsModalOpen(true);
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
    };

    return (
        <HelmetProvider>
            <Helmet title="Plans Dashboard" />

            <div className="min-h-screen bg-gray-50 p-6">
                <div className="max-w-7xl mx-auto">
                    {/* Header Section */}
                    <div className="mb-8">
                        <div className="flex items-center justify-between">
                            <div className="flex items-center space-x-3">
                                <div className="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center">
                                    <SettingOutlined className="text-white text-lg" />
                                </div>
                                <div>
                                    <Title level={2} className="mb-0 text-gray-800">
                                        Plans Management
                                    </Title>
                                    <Text className="text-gray-600">
                                        Create and manage your subscription plans
                                    </Text>
                                </div>
                            </div>

                            {selectedId.length > 0 && (
                                <Badge
                                    count={selectedId.length}
                                    className="mr-4"
                                    color="blue"
                                >
                                    <Text className="text-sm text-gray-600">
                                        Selected Plans
                                    </Text>
                                </Badge>
                            )}
                        </div>
                    </div>

                    {/* Action Bar */}
                    <Card className="mb-6 shadow-sm border-gray-200">
                        <div className="flex justify-between items-center">
                            <div className="flex items-center space-x-4">
                                <Button
                                    type="primary"
                                    icon={<PlusOutlined />}
                                    onClick={handleCreatePlan}
                                    className="h-10 px-6 rounded-lg bg-blue-600 hover:bg-blue-700 border-blue-600 hover:border-blue-700"
                                >
                                    Create New Plan
                                </Button>

                                {selectedId.length > 0 && (
                                    <Popconfirm
                                        title="Delete Plans"
                                        description={`Are you sure you want to delete ${selectedId.length} plan(s)? This action cannot be undone.`}
                                        onConfirm={handleDelete as any}
                                        okText="Yes, Delete"
                                        cancelText="Cancel"
                                        okButtonProps={{
                                            danger: true,
                                            loading: deleteLoading
                                        }}
                                        placement="topRight"
                                    >
                                        <Button
                                            danger
                                            icon={<DeleteOutlined />}
                                            loading={deleteLoading}
                                            className="h-10 px-6 rounded-lg"
                                        >
                                            Delete Selected ({selectedId.length})
                                        </Button>
                                    </Popconfirm>
                                )}
                            </div>

                            <div className="flex items-center space-x-2 text-sm text-gray-500">
                                <span>Last updated: {new Date().toLocaleDateString()}</span>
                            </div>
                        </div>
                    </Card>

                    {/* Plans Content */}
                    <Card
                        className="shadow-sm border-gray-200 min-h-[400px]"
                        bodyStyle={{ padding: 0 }}
                    >
                        <div className="p-6">
                            <GetAllPlans />
                        </div>
                    </Card>

                    {/* Create Plan Modal */}
                    <CreatePlan
                        isOpen={isModalOpen}
                        onClose={handleCloseModal}
                    />
                </div>
            </div>
        </HelmetProvider>
    );
};

export default PlansDashTemplate;