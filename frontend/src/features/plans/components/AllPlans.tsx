import React, { useMemo, useState } from "react";
import {
    Table,
    Button,
    Tag,
    Space,
    Typography,
    Card,
    Tooltip,
    Badge,
    Statistic,
    Row,
    Col,
    Avatar,
    Divider
} from "antd";
import type { ColumnsType } from "antd/es/table";
import {
    EditOutlined,
    DollarOutlined,
    MailOutlined,
    CheckCircleOutlined,
    ThunderboltOutlined,
    UserOutlined,
    SendOutlined,
    CalendarOutlined
} from "@ant-design/icons";
import EditPlans from "./EditPlans";
import { usePlansQuery } from "../hooks/usePlanQuery";
import { PlanData } from '../interface/plan.interface';
import usePlanStore from "../store/plan.store";

const { Title, Text } = Typography;

const GetAllPlans: React.FC = () => {
    const { selectedId, setSelectedId } = usePlanStore();
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [selectedPlan, setSelectedPlan] = useState<PlanData | null>(null);

    const { data: planAPIData, isLoading } = usePlansQuery();

    const planData = useMemo(() => planAPIData?.payload || [], [planAPIData]);

    const handleSelectChange = (selectedRowKeys: React.Key[]) => {
        setSelectedId(selectedRowKeys as string[]);
    };

    const openEditModal = (plan: PlanData) => {
        setSelectedPlan(plan);
        setIsModalOpen(true);
    };

    const closeEditModal = () => {
        setIsModalOpen(false);
        setSelectedPlan(null);
    };

    const getStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case 'active': return 'green';
            case 'inactive': return 'red';
            case 'pending': return 'orange';
            default: return 'default';
        }
    };

    const formatPrice = (price: number) => {
        return price === 0 ? 'Free' : `$${price.toFixed(2)}`;
    };

    const formatLimit = (limit: number) => {
        return limit === 0 ? 'Unlimited' : limit.toLocaleString();
    };

    const expandedRowRender = (record: PlanData) => {
        return (
            <div className="p-6 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg">
                {/* Features Section */}
                <div className="mb-6">
                    <Title level={5} className="mb-4 text-gray-800">
                        <ThunderboltOutlined className="mr-2 text-blue-500" />
                        Plan Features
                    </Title>
                    <Row gutter={[16, 16]}>
                        {record.features.map((feature, index) => (
                            <Col xs={24} sm={12} md={8} lg={6} key={index}>
                                <Card
                                    size="small"
                                    className="shadow-sm hover:shadow-md transition-shadow duration-200"
                                    bodyStyle={{ padding: '12px' }}
                                >
                                    <div className="flex items-start justify-between mb-2">
                                        <Text strong className="text-gray-800 text-sm">
                                            {feature.name}
                                        </Text>
                                        <Tag 
                                            color={feature.value === 'true' ? 'green' : 
                                                   feature.value === 'false' ? 'red' : 'blue'}
                                            className="text-xs"
                                        >
                                            {feature.value === 'true' ? '✓' : 
                                             feature.value === 'false' ? '✗' : feature.value}
                                        </Tag>
                                    </div>
                                    {feature.description && (
                                        <div className="mt-2 pt-2 border-t border-gray-200">
                                            <Text type="secondary" className="text-xs">
                                                {feature.description}
                                            </Text>
                                        </div>
                                    )}
                                </Card>
                            </Col>
                        ))}
                    </Row>
                </div>

                {/* Mailing Limits Section */}
                <Divider className="my-4" />
                <div>
                    <Title level={5} className="mb-4 text-gray-800">
                        <MailOutlined className="mr-2 text-green-500" />
                        Mailing Limits
                    </Title>
                    <Row gutter={[16, 16]}>
                        <Col xs={24} sm={8}>
                            <Card
                                size="small"
                                className="shadow-sm hover:shadow-md transition-shadow duration-200 bg-green-50"
                                bodyStyle={{ padding: '16px' }}
                            >
                                <div className="text-center">
                                    <CalendarOutlined className="text-green-500 text-lg mb-2" />
                                    <div className="text-lg font-bold text-green-600">
                                        {formatLimit(record.mailing_limits.daily_limit)}
                                    </div>
                                    <Text type="secondary" className="text-xs">
                                        Daily Limit
                                    </Text>
                                </div>
                            </Card>
                        </Col>
                        <Col xs={24} sm={8}>
                            <Card
                                size="small"
                                className="shadow-sm hover:shadow-md transition-shadow duration-200 bg-blue-50"
                                bodyStyle={{ padding: '16px' }}
                            >
                                <div className="text-center">
                                    <SendOutlined className="text-blue-500 text-lg mb-2" />
                                    <div className="text-lg font-bold text-blue-600">
                                        {formatLimit(record.mailing_limits.monthly_limit)}
                                    </div>
                                    <Text type="secondary" className="text-xs">
                                        Monthly Limit
                                    </Text>
                                </div>
                            </Card>
                        </Col>
                        <Col xs={24} sm={8}>
                            <Card
                                size="small"
                                className="shadow-sm hover:shadow-md transition-shadow duration-200 bg-purple-50"
                                bodyStyle={{ padding: '16px' }}
                            >
                                <div className="text-center">
                                    <UserOutlined className="text-purple-500 text-lg mb-2" />
                                    <div className="text-lg font-bold text-purple-600">
                                        {formatLimit(record.mailing_limits.max_recipients_per_mail)}
                                    </div>
                                    <Text type="secondary" className="text-xs">
                                        Max Recipients
                                    </Text>
                                </div>
                            </Card>
                        </Col>
                    </Row>
                </div>
            </div>
        );
    };

    const columns: ColumnsType<PlanData> = [
        {
            title: 'Plan Details',
            key: 'planDetails',
            width: 300,
            render: (_, record) => (
                <div className="flex items-center space-x-3">
                    <Avatar
                        className="bg-gradient-to-r from-blue-500 to-purple-600"
                        size={40}
                    >
                        {record.name.charAt(0).toUpperCase()}
                    </Avatar>
                    <div>
                        <Text strong className="text-gray-800 text-base">
                            {record.name}
                        </Text>
                        <div className="text-sm text-gray-500 mt-1">
                            {record.description}
                        </div>
                    </div>
                </div>
            ),
        },
        {
            title: 'Pricing',
            key: 'pricing',
            width: 120,
            render: (_, record) => (
                <div className="text-center">
                    <div className="text-lg font-bold text-green-600">
                        {formatPrice(record.price)}
                    </div>
                    <Text type="secondary" className="text-xs">
                        per {record.billing_cycle}
                    </Text>
                </div>
            ),
        },
        {
            title: 'Email Limits',
            key: 'emailLimits',
            width: 150,
            render: (_, record) => (
                <div className="space-y-1">
                    <div className="flex items-center space-x-2">
                        <SendOutlined className="text-blue-500 text-sm" />
                        <Text className="text-sm font-medium">
                            {formatLimit(record.mailing_limits.monthly_limit)}
                        </Text>
                    </div>
                    <div className="text-xs text-gray-500">
                        Monthly emails
                    </div>
                    <div className="flex items-center space-x-2">
                        <CalendarOutlined className="text-green-500 text-sm" />
                        <Text className="text-sm">
                            {formatLimit(record.mailing_limits.daily_limit)}
                        </Text>
                    </div>
                    <div className="text-xs text-gray-500">
                        Daily emails
                    </div>
                </div>
            ),
        },
        {
            title: 'Status',
            key: 'status',
            width: 100,
            render: (_, record) => (
                <Tag
                    color={getStatusColor(record.status)}
                    className="rounded-full px-3 py-1 text-xs font-medium"
                >
                    {record.status.charAt(0).toUpperCase() + record.status.slice(1)}
                </Tag>
            ),
        },
        {
            title: 'Features',
            key: 'featureCount',
            width: 100,
            render: (_, record) => (
                <div className="text-center">
                    <Badge
                        count={record.features.length}
                        showZero
                        color="#1890ff"
                        className="text-sm"
                    />
                    <div className="text-xs text-gray-500 mt-1">
                        features
                    </div>
                </div>
            ),
        },
        {
            title: 'Actions',
            key: 'actions',
            width: 80,
            render: (_, record) => (
                <Space>
                    <Tooltip title="Edit plan">
                        <Button
                            type="text"
                            icon={<EditOutlined />}
                            onClick={() => openEditModal(record)}
                            className="text-gray-600 hover:text-blue-600 hover:bg-blue-50"
                        />
                    </Tooltip>
                </Space>
            ),
        },
    ];

    const rowSelection = {
        selectedRowKeys: selectedId,
        onChange: handleSelectChange,
        getCheckboxProps: (record: PlanData) => ({
            name: record.name,
        }),
    };

    return (
        <div className="p-6 bg-gray-50 min-h-screen">
            <div className="max-w-7xl mx-auto">
                {/* Header Section */}
                <div className="mb-8">
                    <div className="flex items-center justify-between">
                        <div>
                            <Title level={2} className="mb-2 text-gray-800">
                                📋 Plan Management
                            </Title>
                            <Text type="secondary" className="text-base">
                                Manage and configure your subscription plans
                            </Text>
                        </div>
                        <div className="flex items-center space-x-4">
                            <Statistic
                                title="Total Plans"
                                value={planData.length}
                                prefix={<DollarOutlined />}
                                valueStyle={{ color: '#3f8600' }}
                            />
                            {selectedId.length > 0 && (
                                <Statistic
                                    title="Selected"
                                    value={selectedId.length}
                                    valueStyle={{ color: '#1890ff' }}
                                />
                            )}
                        </div>
                    </div>
                </div>

                {/* Main Table */}
                <Card
                    className="shadow-lg rounded-xl border-0"
                    bodyStyle={{ padding: 0 }}
                >
                    <Table
                        columns={columns}
                        dataSource={planData}
                        rowKey="id"
                        loading={isLoading}
                        rowSelection={rowSelection}
                        expandable={{
                            expandedRowRender,
                            rowExpandable: (record) => record.features.length > 0,
                            expandRowByClick: false,
                        }}
                        pagination={{
                            total: planData.length,
                            showSizeChanger: true,
                            showQuickJumper: true,
                            showTotal: (total, range) =>
                                `${range[0]}-${range[1]} of ${total} plans`,
                            className: "px-6 py-4",
                        }}
                        className="rounded-xl overflow-hidden"
                        scroll={{ x: 800 }}
                        size="middle"
                    />
                </Card>

                {/* Selected Actions */}
                {selectedId.length > 0 && (
                    <Card className="mt-6 bg-blue-50 border-blue-200">
                        <div className="flex items-center justify-between">
                            <div className="flex items-center space-x-2">
                                <CheckCircleOutlined className="text-blue-500" />
                                <Text strong className="text-blue-700">
                                    {selectedId.length} plan(s) selected
                                </Text>
                            </div>
                            <Space>
                                <Button
                                    type="primary"
                                    size="small"
                                    className="bg-blue-600 hover:bg-blue-700"
                                >
                                    Bulk Actions
                                </Button>
                                <Button
                                    size="small"
                                    onClick={() => setSelectedId([])}
                                >
                                    Clear Selection
                                </Button>
                            </Space>
                        </div>
                    </Card>
                )}
            </div>

            <EditPlans
                isOpen={isModalOpen}
                onClose={closeEditModal}
                plan={selectedPlan}
            />
        </div>
    );
};

export default GetAllPlans;