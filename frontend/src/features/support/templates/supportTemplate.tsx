import React from 'react';
import { Card, Typography, Row, Col } from 'antd';
import { Mail, Book, HelpCircle, PlayCircle } from 'lucide-react';
import { Link } from 'react-router-dom';
import { HelmetProvider, Helmet } from 'react-helmet-async';

const { Title, Paragraph } = Typography;

interface SupportCardProps {
    icon: React.ReactNode;
    title: string;
    description: string;
    link: string;
}

const SupportCard: React.FC<SupportCardProps> = ({ icon, title, description, link }) => {
    return (
        <Link to={link}>
            <Card
                hoverable
                className="h-full"
                style={{
                    padding: 24,
                    height: '100%',
                    display: 'flex',
                    flexDirection: 'column'
                }}
            >
                <div className="flex items-center mb-4">
                    {icon}
                    <Title level={4} className="mb-0 ml-3">
                        {title}
                    </Title>
                </div>
                <Paragraph className="text-gray-600 m-0">
                    {description}
                </Paragraph>
            </Card>
        </Link>
    );
};

const HelpAndSupport: React.FC = () => {
    const supportCards = [
        {
            icon: <HelpCircle className="w-8 h-8 text-blue-600" />,
            title: "FAQs",
            description: "Read our frequently asked questions here, this is a quick starting point to answers common questions",
            link: "#"
        },
        {
            icon: <Book className="w-8 h-8 text-blue-600" />,
            title: "Help Articles",
            description: "Easy short advice, answers & best practices from the crabmailer team & contributors",
            link: "#"
        },
        {
            icon: <PlayCircle className="w-8 h-8 text-blue-600" />,
            title: "Video Tutorials",
            description: "Watch step-by-step guides on how to use crabmailer effectively",
            link: "#"
        },
        {
            icon: <Mail className="w-8 h-8 text-blue-600" />,
            title: "Contact Support",
            description: "Get in touch with our support team for personalized assistance",
            link: "/app/support/ticket"
        }
    ];

    return (
        <HelmetProvider>
            <Helmet title="Help and support - CrabMailer" />
            <div className="max-w-6xl mx-auto p-6 mt-10">
                <Title
                    level={1}
                    className="text-center mb-8"
                    style={{ color: '#1890ff' }}
                >
                    Help & Support
                </Title>

                <Row gutter={[24, 24]}>
                    {supportCards.map((card, index) => (
                        <Col xs={24} md={12} key={index}>
                            <SupportCard {...card} />
                        </Col>
                    ))}
                </Row>
            </div>
        </HelmetProvider>
    );
};

export default HelpAndSupport;