import React from 'react';
import { Card,  Row, Col, Typography, Avatar,  Spin } from 'antd';
import { UserOutlined, UsergroupAddOutlined, TeamOutlined } from '@ant-design/icons';
import { useContactCountQuery } from "../../features/contacts/hooks/useContactQuery";

const { Title } = Typography;

const ContactsDashboard: React.FC = () => {
    const { data: contactCount, isLoading } = useContactCountQuery();

    if (isLoading) {
        return (
            <div className="flex justify-center items-center min-h-[400px]">
                <Spin size="large" />
            </div>
        );
    }

    return (
        <div className="w-full mx-auto p-6 bg-gradient-to-br from-blue-50/50 to-indigo-50/30 min-h-screen">
            {/* Header Section */}
            <div className="mb-8">
                <div className="flex items-center space-x-4 mb-2">
                    <div className="p-3 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl shadow-lg">
                        <TeamOutlined className="text-white text-2xl" />
                    </div>
                    <div>
                        <Title level={2} className="!mb-0 !text-gray-800">
                            Contacts Dashboard
                        </Title>
                        <p className="text-gray-600 mt-1">Manage and track your contact network</p>
                    </div>
                </div>
            </div>

            {/* Stats Cards */}
            <Row gutter={[24, 24]} className="mb-8">
                <Col xs={24} md={12}>
                    <Card
                        className="relative overflow-hidden border-0 shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1"
                        style={{
                            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                            borderRadius: '16px'
                        }}
                    >
                        <div className="relative z-10">
                            <div className="flex items-center justify-between mb-4">
                                <Avatar 
                                    size={56} 
                                    className="bg-white/20 backdrop-blur-sm border-2 border-white/30"
                                    icon={<UserOutlined className="text-white text-xl" />}
                                />
                                <div className="text-right">
                                    <div className="text-white/80 text-sm font-medium">Total Contacts</div>
                                    <div className="text-white text-3xl font-bold">
                                        {contactCount?.payload?.total || 0}
                                    </div>
                                </div>
                            </div>
                            
                            <div className="flex items-center justify-between">
                                <div className="text-white/90 text-sm">
                                    Your complete network
                                </div>
                                <div className="flex items-center space-x-1 text-white/80">
                                    <TeamOutlined className="text-sm" />
                                    <span className="text-xs">Active</span>
                                </div>
                            </div>
                        </div>
                        
                        {/* Background decorative elements */}
                        <div className="absolute -top-4 -right-4 w-24 h-24 bg-white/10 rounded-full blur-xl"></div>
                        <div className="absolute -bottom-8 -left-8 w-32 h-32 bg-white/5 rounded-full blur-2xl"></div>
                    </Card>
                </Col>

                <Col xs={24} md={12}>
                    <Card
                        className="relative overflow-hidden border-0 shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1"
                        style={{
                            background: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
                            borderRadius: '16px'
                        }}
                    >
                        <div className="relative z-10">
                            <div className="flex items-center justify-between mb-4">
                                <Avatar 
                                    size={56} 
                                    className="bg-white/20 backdrop-blur-sm border-2 border-white/30"
                                    icon={<UsergroupAddOutlined className="text-white text-xl" />}
                                />
                                <div className="text-right">
                                    <div className="text-white/80 text-sm font-medium">New This Month</div>
                                    <div className="text-white text-3xl font-bold">
                                        {contactCount?.payload?.recent || 0}
                                    </div>
                                </div>
                            </div>
                            
                            <div className="flex items-center justify-between">
                                <div className="text-white/90 text-sm">
                                    Last 30 days growth
                                </div>
                                <div className="flex items-center space-x-1 text-white/80">
                                    <TeamOutlined className="text-sm" />
                                    <span className="text-xs">Growing</span>
                                </div>
                            </div>
                        </div>
                        
                        {/* Background decorative elements */}
                        <div className="absolute -top-4 -right-4 w-24 h-24 bg-white/10 rounded-full blur-xl"></div>
                        <div className="absolute -bottom-8 -left-8 w-32 h-32 bg-white/5 rounded-full blur-2xl"></div>
                    </Card>
                </Col>
            </Row>

            {/* Additional Info Cards */}
            <Row gutter={[24, 24]}>
                <Col xs={24} lg={8}>
                    <Card 
                        className="border-0 shadow-md hover:shadow-lg transition-shadow duration-300"
                        style={{ borderRadius: '12px' }}
                    >
                        <div className="text-center p-4">
                            <div className="mb-4">
                                <div className="inline-flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full">
                                    <UserOutlined className="text-blue-600 text-2xl" />
                                </div>
                            </div>
                            <Title level={4} className="!mb-2 !text-gray-800">
                                Contact Management
                            </Title>
                            <p className="text-gray-600 text-sm">
                                Organize and maintain your professional network efficiently
                            </p>
                        </div>
                    </Card>
                </Col>

                <Col xs={24} lg={8}>
                    <Card 
                        className="border-0 shadow-md hover:shadow-lg transition-shadow duration-300"
                        style={{ borderRadius: '12px' }}
                    >
                        <div className="text-center p-4">
                            <div className="mb-4">
                                <div className="inline-flex items-center justify-center w-16 h-16 bg-green-100 rounded-full">
                                    <UsergroupAddOutlined className="text-green-600 text-2xl" />
                                </div>
                            </div>
                            <Title level={4} className="!mb-2 !text-gray-800">
                                Network Growth
                            </Title>
                            <p className="text-gray-600 text-sm">
                                Track your networking progress and expansion over time
                            </p>
                        </div>
                    </Card>
                </Col>

                <Col xs={24} lg={8}>
                    <Card 
                        className="border-0 shadow-md hover:shadow-lg transition-shadow duration-300"
                        style={{ borderRadius: '12px' }}
                    >
                        <div className="text-center p-4">
                            <div className="mb-4">
                                <div className="inline-flex items-center justify-center w-16 h-16 bg-purple-100 rounded-full">
                                    <TeamOutlined className="text-purple-600 text-2xl" />
                                </div>
                            </div>
                            <Title level={4} className="!mb-2 !text-gray-800">
                                Team Collaboration
                            </Title>
                            <p className="text-gray-600 text-sm">
                                Share and collaborate on contacts with your team members
                            </p>
                        </div>
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default ContactsDashboard;