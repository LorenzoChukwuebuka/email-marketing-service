import React from 'react';
import { Card, Typography, Row, Col, Space, Divider } from 'antd';
import { Mail, Book, HelpCircle, PlayCircle, ArrowRight } from 'lucide-react';
import { Link } from 'react-router-dom';
import { HelmetProvider, Helmet } from 'react-helmet-async';

const { Title, Paragraph, Text } = Typography;

interface SupportCardProps {
    icon: React.ReactNode;
    title: string;
    description: string;
    link: string;
    gradient: string;
}

const SupportCard: React.FC<SupportCardProps> = ({ icon, title, description, link, gradient }) => {
    return (
        <Link to={link} className="block group">
            <Card
                hoverable
                className="h-full border-0 shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 overflow-hidden"
                bodyStyle={{ padding: 0, height: '100%' }}
            >
                <div className="relative h-full">
                    {/* Gradient header */}
                    <div className={`h-2 w-full ${gradient}`} />
                    
                    <div className="p-6 h-full flex flex-col">
                        {/* Icon and title section */}
                        <div className="flex items-center mb-4">
                            <div className="p-3 bg-gradient-to-br from-blue-50 to-indigo-50 rounded-xl group-hover:from-blue-100 group-hover:to-indigo-100 transition-all duration-300">
                                {icon}
                            </div>
                            <div className="ml-4 flex-1">
                                <Title level={4} className="mb-0 text-gray-800 group-hover:text-blue-600 transition-colors duration-300">
                                    {title}
                                </Title>
                            </div>
                            <ArrowRight className="w-5 h-5 text-gray-400 group-hover:text-blue-600 transform group-hover:translate-x-1 transition-all duration-300" />
                        </div>

                        {/* Description */}
                        <Paragraph className="text-gray-600 mb-0 flex-1 leading-relaxed">
                            {description}
                        </Paragraph>

                        {/* Bottom accent */}
                        <div className="mt-4 pt-4 border-t border-gray-100">
                            <Text className="text-sm text-blue-600 font-medium group-hover:text-blue-700 transition-colors duration-300">
                                Learn more →
                            </Text>
                        </div>
                    </div>
                </div>
            </Card>
        </Link>
    );
};

const HelpAndSupport: React.FC = () => {
    const supportCards = [
        {
            icon: <HelpCircle className="w-6 h-6 text-emerald-600" />,
            title: "FAQs",
            description: "Find quick answers to common questions. Our comprehensive FAQ section covers everything from getting started to advanced features.",
            link: "#",
            gradient: "bg-gradient-to-r from-emerald-400 to-teal-500"
        },
        {
            icon: <Book className="w-6 h-6 text-blue-600" />,
            title: "Help Articles",
            description: "Explore detailed guides and best practices. Get expert advice from our team and community contributors.",
            link: "#",
            gradient: "bg-gradient-to-r from-blue-400 to-indigo-500"
        },
        {
            icon: <PlayCircle className="w-6 h-6 text-purple-600" />,
            title: "Video Tutorials",
            description: "Watch step-by-step video guides that show you exactly how to use CrabMailer's features effectively.",
            link: "#",
            gradient: "bg-gradient-to-r from-purple-400 to-pink-500"
        },
        {
            icon: <Mail className="w-6 h-6 text-orange-600" />,
            title: "Contact Support",
            description: "Need personalized help? Our support team is here to assist you with any questions or issues you may have.",
            link: "/app/support/ticket",
            gradient: "bg-gradient-to-r from-orange-400 to-red-500"
        }
    ];

    return (
        <HelmetProvider>
            <Helmet title="Help and support - CrabMailer" />
            <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50">
                <div className="max-w-7xl mx-auto px-6 py-12">
                    {/* Header Section */}
                    <div className="text-center mb-12">
                        <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-2xl mb-6 shadow-lg">
                            <HelpCircle className="w-8 h-8 text-white" />
                        </div>
                        
                        <Title
                            level={1}
                            className="mb-4 bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent"
                            style={{ fontSize: '3rem', fontWeight: 'bold' }}
                        >
                            Help & Support
                        </Title>
                        
                        <Paragraph className="text-xl text-gray-600 max-w-2xl mx-auto leading-relaxed">
                            Get the help you need to make the most of CrabMailer. From quick answers to detailed guides, we're here to support your success.
                        </Paragraph>
                    </div>

                    <Divider className="mb-12 border-gray-200" />

                    {/* Support Cards Grid */}
                    <Row gutter={[32, 32]} className="mb-12">
                        {supportCards.map((card, index) => (
                            <Col xs={24} lg={12} key={index}>
                                <SupportCard {...card} />
                            </Col>
                        ))}
                    </Row>

                    {/* Bottom CTA Section */}
                    <div className="text-center mt-16">
                        <Card className="bg-gradient-to-r from-blue-600 to-indigo-600 border-0 shadow-xl">
                            <div className="py-8 px-6">
                                <Title level={3} className="text-white mb-4">
                                    Still need help?
                                </Title>
                                <Paragraph className="text-blue-100 mb-6 text-lg">
                                    Can't find what you're looking for? Our support team is always ready to help.
                                </Paragraph>
                                <Space size="large">
                                    <Link to="/app/support/ticket">
                                        <Card
                                            hoverable
                                            className="bg-white/10 border-white/20 backdrop-blur-sm hover:bg-white/20 transition-all duration-300"
                                            bodyStyle={{ padding: '12px 24px' }}
                                        >
                                            <div className="flex items-center space-x-3">
                                                <Mail className="w-5 h-5 text-white" />
                                                <Text className="text-white font-medium">Create Support Ticket</Text>
                                            </div>
                                        </Card>
                                    </Link>
                                    <Card
                                        hoverable
                                        className="bg-white/10 border-white/20 backdrop-blur-sm hover:bg-white/20 transition-all duration-300"
                                        bodyStyle={{ padding: '12px 24px' }}
                                    >
                                        <div className="flex items-center space-x-3">
                                            <Book className="w-5 h-5 text-white" />
                                            <Text className="text-white font-medium">Browse Documentation</Text>
                                        </div>
                                    </Card>
                                </Space>
                            </div>
                        </Card>
                    </div>
                </div>
            </div>
        </HelmetProvider>
    );
};

export default HelpAndSupport;