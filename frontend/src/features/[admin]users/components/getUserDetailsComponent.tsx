import { Link, useParams } from 'react-router-dom';
import { Helmet, HelmetProvider } from 'react-helmet-async';
import { useSingleUserQuery, useUserStatsQuery } from '../hooks/useAdminUsersQueryHook';
import { useMemo } from 'react';
import {
    Card,
    Statistic,
    Button,
    Descriptions,
    Tag,
    Space,
    Typography,
    Row,
    Col,
    Avatar
} from 'antd';
import {
    UserOutlined,
    ThunderboltOutlined,
    FileTextOutlined,
    TeamOutlined,
    ContactsOutlined,
    EyeOutlined,
    CheckCircleOutlined,
    CloseCircleOutlined,
    CalendarOutlined,
    MailOutlined,
    PhoneOutlined,
    BankOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;

const AdminUserDetailComponent = () => {
    const { id } = useParams<{ id: string }>() as { id: string };
    const { data: userData } = useSingleUserQuery(id);
    const { data: statsData } = useUserStatsQuery(id);

    const userdetailsData = useMemo(() => userData?.payload, [userData]);
    const userStatsData = useMemo(() => statsData?.payload, [statsData]);

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString('en-US', {
            timeZone: 'UTC',
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: 'numeric',
            minute: 'numeric',
            second: 'numeric'
        });
    };

    const statsCards = [
        {
            title: "Total Campaigns Created",
            value: userStatsData?.total_campaigns || 0,
            icon: <ThunderboltOutlined className="text-blue-500" />,
            color: "from-blue-500 to-blue-600",
            action: (
                <Link to={{ pathname: "/zen/campaigns/details/" + id, search: `?username=${userdetailsData?.fullname}` }}>
                    <Button type="primary" size="small" icon={<EyeOutlined />}>
                        View Campaigns
                    </Button>
                </Link>
            )
        },
        {
            title: "Total Templates Created",
            value: userStatsData?.total_templates || 0,
            icon: <FileTextOutlined className="text-purple-500" />,
            color: "from-purple-500 to-purple-600"
        },
        {
            title: "Total Groups Created",
            value: userStatsData?.total_groups || 0,
            icon: <TeamOutlined className="text-indigo-500" />,
            color: "from-indigo-500 to-indigo-600"
        },
        {
            title: "Total Contacts",
            value: userStatsData?.total_contacts || 0,
            icon: <ContactsOutlined className="text-green-500" />,
            color: "from-green-500 to-green-600"
        }
    ];

    return (
        <HelmetProvider>
            <Helmet title={`Details for ${userdetailsData?.fullname}`} />
            <div className="min-h-screen bg-gray-50 p-6">
                <div className="max-w-7xl mx-auto">
                    {/* Header */}
                    <div className="mb-8">
                        <Title level={2} className="!mb-2">
                            <UserOutlined className="mr-3" />
                            User Details
                        </Title>
                        <Text type="secondary" className="text-lg">
                            Comprehensive overview of user information and statistics
                        </Text>
                    </div>

                    {/* User Profile Card */}
                    <Card className="mb-6 shadow-lg border-0">
                        <div className="flex items-center space-x-4 mb-6">
                            <Avatar
                                size={64}
                                icon={<UserOutlined />}
                                className="bg-gradient-to-r from-blue-500 to-purple-600"
                            />
                            <div>
                                <Title level={3} className="!mb-1">
                                    {userdetailsData?.fullname || "Unknown User"}
                                </Title>
                                <Text type="secondary">{userdetailsData?.email}</Text>
                                <div className="flex items-center space-x-2 mt-2">
                                    {userdetailsData?.verified ? (
                                        <Tag icon={<CheckCircleOutlined />} color="success">
                                            Verified
                                        </Tag>
                                    ) : (
                                        <Tag icon={<CloseCircleOutlined />} color="default">
                                            Not Verified
                                        </Tag>
                                    )}
                                    {userdetailsData?.blocked ? (
                                        <Tag color="error">Blocked</Tag>
                                    ) : (
                                        <Tag color="success">Active</Tag>
                                    )}
                                </div>
                            </div>
                        </div>
                    </Card>

                    {/* Statistics Cards */}
                    <Row gutter={[24, 24]} className="mb-8">
                        {statsCards.map((stat, index) => (
                            <Col xs={24} sm={12} lg={6} key={index}>
                                <Card
                                    className={`shadow-lg border-0 bg-gradient-to-r ${stat.color} text-white hover:shadow-xl transition-all duration-300`}
                                    bodyStyle={{ padding: '24px' }}
                                >
                                    <div className="flex items-center justify-between mb-4">
                                        <div className="text-white/80 text-4xl">
                                            {stat.icon}
                                        </div>
                                        <Statistic
                                            value={stat.value}
                                            valueStyle={{
                                                color: 'white',
                                                fontSize: '2rem',
                                                fontWeight: 'bold'
                                            }}
                                        />
                                    </div>
                                    <Text className="text-white/90 text-sm font-medium">
                                        {stat.title}
                                    </Text>
                                    {stat.action && (
                                        <div className="mt-4">
                                            {stat.action}
                                        </div>
                                    )}
                                </Card>
                            </Col>
                        ))}
                    </Row>

                    {/* User Information */}
                    <Card
                        title={
                            <Title level={4} className="!mb-0">
                                <UserOutlined className="mr-2" />
                                User Information
                            </Title>
                        }
                        className="shadow-lg border-0"
                    >
                        <Descriptions
                            bordered
                            column={{ xxl: 3, xl: 3, lg: 2, md: 2, sm: 1, xs: 1 }}
                            size="middle"
                            labelStyle={{
                                backgroundColor: '#f8fafc',
                                fontWeight: 'bold',
                                width: '200px'
                            }}
                        >
                            <Descriptions.Item
                                label={
                                    <Space>
                                        <UserOutlined />
                                        Full Name
                                    </Space>
                                }
                            >
                                <Text strong>{userdetailsData?.fullname || "N/A"}</Text>
                            </Descriptions.Item>

                            <Descriptions.Item
                                label={
                                    <Space>
                                        <MailOutlined />
                                        Email
                                    </Space>
                                }
                            >
                                <Text copyable>{userdetailsData?.email || "N/A"}</Text>
                            </Descriptions.Item>

                            <Descriptions.Item
                                label={
                                    <Space>
                                        <PhoneOutlined />
                                        Mobile Number
                                    </Space>
                                }
                            >
                                <Text copyable>
                                    {userdetailsData?.phonenumber || "N/A"}
                                </Text>
                            </Descriptions.Item>

                            <Descriptions.Item
                                label={
                                    <Space>
                                        <BankOutlined />
                                        Company
                                    </Space>
                                }
                                span={2}
                            >
                                <Text strong>
                                    {userdetailsData?.company?.companyname || "N/A"}
                                </Text>
                            </Descriptions.Item>

                            <Descriptions.Item
                                label={
                                    <Space>
                                        <CheckCircleOutlined />
                                        Verified Status
                                    </Space>
                                }
                            >
                                {userdetailsData?.verified ? (
                                    <Tag icon={<CheckCircleOutlined />} color="success">
                                        Verified
                                    </Tag>
                                ) : (
                                    <Tag icon={<CloseCircleOutlined />} color="default">
                                        Not Verified
                                    </Tag>
                                )}
                            </Descriptions.Item>

                            <Descriptions.Item
                                label={
                                    <Space>
                                        <CalendarOutlined />
                                        Verified At
                                    </Space>
                                }
                            >
                                <Text>
                                    {userdetailsData?.verified_at
                                        ? formatDate(userdetailsData.verified_at)
                                        : "Not Verified Yet"
                                    }
                                </Text>
                            </Descriptions.Item>

                            <Descriptions.Item
                                label={
                                    <Space>
                                        <CloseCircleOutlined />
                                        Blocked Status
                                    </Space>
                                }
                            >
                                {userdetailsData?.blocked ? (
                                    <Tag color="error">Blocked</Tag>
                                ) : (
                                    <Tag color="success">Active</Tag>
                                )}
                            </Descriptions.Item>

                            <Descriptions.Item
                                label={
                                    <Space>
                                        <CalendarOutlined />
                                        Joined On
                                    </Space>
                                }
                                span={2}
                            >
                                <Text>
                                    {userdetailsData?.created_at
                                        ? formatDate(userdetailsData.created_at)
                                        : "N/A"
                                    }
                                </Text>
                            </Descriptions.Item>
                        </Descriptions>
                    </Card>
                </div>
            </div>
        </HelmetProvider>
    );
};

export default AdminUserDetailComponent;