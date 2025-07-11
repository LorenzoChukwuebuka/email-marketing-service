import React from 'react';
import { Card, Row, Col, Typography, Badge } from 'antd';
import {
    UserOutlined,
    UserAddOutlined,
    UserDeleteOutlined,
    HeartOutlined,
    RiseOutlined,
    FallOutlined,
    TeamOutlined,
    BarChartOutlined
} from '@ant-design/icons';
import { useContactEngagementQuery } from "../../features/contacts/hooks/useContactQuery";

const { Title, Text } = Typography;

const OverviewStats: React.FC = () => {
    const { data: engagementCount } = useContactEngagementQuery();

    const statsData = [
        {
            title: "Total Active Subscribers",
            value: engagementCount?.payload?.total || 0,
            icon: <UserOutlined className="text-2xl" />,
            gradient: "from-blue-500 to-cyan-500",
            bgColor: "bg-blue-50",
            iconColor: "text-blue-600",
            trend: "up",
            description: "Your complete subscriber base"
        },
        {
            title: "New Subscribers",
            value: engagementCount?.payload?.new || 0,
            icon: <UserAddOutlined className="text-2xl" />,
            gradient: "from-green-500 to-emerald-500",
            bgColor: "bg-green-50",
            iconColor: "text-green-600",
            trend: "up",
            description: "Recent additions to your list"
        },
        {
            title: "Unsubscribed",
            value: engagementCount?.payload?.unsubscribed || 0,
            icon: <UserDeleteOutlined className="text-2xl" />,
            gradient: "from-red-500 to-pink-500",
            bgColor: "bg-red-50",
            iconColor: "text-red-600",
            trend: "down",
            description: "Users who opted out"
        },
        {
            title: "Engaged Subscribers",
            value: engagementCount?.payload?.engaged || 0,
            icon: <HeartOutlined className="text-2xl" />,
            gradient: "from-purple-500 to-pink-500",
            bgColor: "bg-purple-50",
            iconColor: "text-purple-600",
            trend: "up",
            description: "Active and interacting users"
        }
    ];

    return (
        <div className="w-full p-6 bg-gradient-to-br from-gray-50 to-blue-50/30 min-h-[400px]">
            {/* Header Section */}
            <div className="mb-8">
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-4">
                        <div className="p-3 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl shadow-lg">
                            <BarChartOutlined className="text-white text-2xl" />
                        </div>
                        <div>
                            <Title level={2} className="!mb-0 !text-gray-800">
                                Overview Stats
                            </Title>
                            <Text className="text-gray-600">
                                Track your subscription metrics and audience engagement
                            </Text>
                        </div>
                    </div>
                    <Badge
                        count="Live"
                        style={{
                            backgroundColor: '#52c41a',
                            boxShadow: '0 0 0 4px rgba(82, 196, 26, 0.1)'
                        }}
                    />
                </div>
            </div>

            {/* Section Header */}
            <div className="mb-6">
                <div className="flex items-center space-x-2 mb-2">
                    <TeamOutlined className="text-indigo-600" />
                    <Title level={4} className="!mb-0 !text-gray-700">
                        Subscription & Audience
                    </Title>
                </div>
                <Text className="text-gray-500">
                    Monitor your subscriber growth and engagement patterns
                </Text>
            </div>

            {/* Stats Cards */}
            <Row gutter={[24, 24]}>
                {statsData.map((stat, index) => (
                    <Col xs={24} sm={12} lg={6} key={index}>
                        <Card
                            className="relative overflow-hidden border-0 shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1"
                            style={{ borderRadius: '16px' }}
                        >
                            {/* Gradient Background */}
                            <div
                                className={`absolute inset-0 bg-gradient-to-br ${stat.gradient} opacity-5`}
                            />

                            <div className="relative z-10">
                                {/* Header with Icon */}
                                <div className="flex items-center justify-between mb-4">
                                    <div className={`p-3 ${stat.bgColor} rounded-xl`}>
                                        <span className={stat.iconColor}>
                                            {stat.icon}
                                        </span>
                                    </div>
                                    <div className="flex items-center space-x-1">
                                        {stat.trend === 'up' ? (
                                            <RiseOutlined className="text-green-500 text-sm" />
                                        ) : (
                                            <FallOutlined className="text-red-500 text-sm" />
                                        )}
                                        <span className={`text-xs font-medium ${stat.trend === 'up' ? 'text-green-600' : 'text-red-600'
                                            }`}>
                                            {stat.trend === 'up' ? 'Growth' : 'Decline'}
                                        </span>
                                    </div>
                                </div>

                                {/* Value */}
                                <div className="mb-3">
                                    <Text className="text-3xl font-bold text-gray-800">
                                        {typeof stat.value === 'number' ? stat.value.toLocaleString() : stat.value}
                                    </Text>
                                </div>

                                {/* Title */}
                                <div className="mb-2">
                                    <Text className="text-sm font-semibold text-gray-700">
                                        {stat.title}
                                    </Text>
                                </div>

                                {/* Description */}
                                <Text className="text-xs text-gray-500">
                                    {stat.description}
                                </Text>
                            </div>

                            {/* Decorative Elements */}
                            <div className="absolute -top-4 -right-4 w-16 h-16 bg-gradient-to-br from-white/20 to-transparent rounded-full blur-sm" />
                            <div className="absolute -bottom-2 -left-2 w-8 h-8 bg-gradient-to-br from-white/10 to-transparent rounded-full blur-sm" />
                        </Card>
                    </Col>
                ))}
            </Row>

            {/* Additional Insights */}
            <Row gutter={[24, 24]} className="mt-8">
                <Col xs={24} md={12}>
                    <Card
                        className="border-0 shadow-md hover:shadow-lg transition-shadow duration-300"
                        style={{ borderRadius: '12px' }}
                    >
                        <div className="flex items-center space-x-4">
                            <div className="p-3 bg-gradient-to-br from-blue-100 to-cyan-100 rounded-xl">
                                <BarChartOutlined className="text-blue-600 text-xl" />
                            </div>
                            <div>
                                <Title level={5} className="!mb-1 !text-gray-800">
                                    Engagement Rate
                                </Title>
                                <Text className="text-gray-600 text-sm">
                                    {engagementCount?.payload?.engaged && engagementCount?.payload?.total
                                        ? `${Math.round(((engagementCount.payload.engaged as any) / (engagementCount.payload.total as any)) * 100)}%`
                                        : '0%'
                                    } of subscribers are actively engaged
                                </Text>
                            </div>
                        </div>
                    </Card>
                </Col>

                <Col xs={24} md={12}>
                    <Card
                        className="border-0 shadow-md hover:shadow-lg transition-shadow duration-300"
                        style={{ borderRadius: '12px' }}
                    >
                        <div className="flex items-center space-x-4">
                            <div className="p-3 bg-gradient-to-br from-green-100 to-emerald-100 rounded-xl">
                                <TeamOutlined className="text-green-600 text-xl" />
                            </div>
                            <div>
                                <Title level={5} className="!mb-1 !text-gray-800">
                                    Growth Trend
                                </Title>
                                <Text className="text-gray-600 text-sm">
                                    {engagementCount?.payload?.new || 0} new subscribers this period
                                </Text>
                            </div>
                        </div>
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default OverviewStats;