import { Helmet, HelmetProvider } from "react-helmet-async";
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import EmptyState from '../../../../components/emptyStateComponent';
import { parseDate } from '../../../../utils/utils';
import { useMemo, useRef, useState } from "react";
import { useAdminUserCampaignsQuery } from "../../hooks/useAdminCampaignQuery";
import useDebounce from "../../../../hooks/useDebounce";
import {
    Pagination,
    Input,
    Button,
    Table,
    Tag,
    Dropdown,
    Space,
    Card,
    Typography,
    Breadcrumb,
    Row,
    Col,
    Empty
} from "antd";
import {
    ArrowLeftOutlined,
    SearchOutlined,
    EyeOutlined,
    PauseOutlined,
    DeleteOutlined,
    MoreOutlined,
    FileTextOutlined,
    UserOutlined,
    HomeOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;
const { Search } = Input;

const AdminUserCampaigns: React.FC = () => {
    const { id, companyId } = useParams();
    const [searchParams] = useSearchParams();
    const username = searchParams.get('username');
    const navigate = useNavigate();
    const [searchQuery, setSearchQuery] = useState<string>("");
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const debouncedSearchQuery = useDebounce(searchQuery, 300);

    const { data, isLoading } = useAdminUserCampaignsQuery(
        id as string,
        companyId as string,
        currentPage,
        pageSize,
        debouncedSearchQuery
    );

    const campaignData = useMemo(() => data?.payload.data, [data]);

    const handleSearch = (value: string) => {
        setSearchQuery(value);
        setCurrentPage(1); // Reset to first page when searching
    };

    const deleteCamp = (uuid: string) => {
        console.log('Delete campaign:', uuid);
    };

    const suspendCampaign = (uuid: string) => {
        console.log('Suspend campaign:', uuid);
    };

    const pageTitle = username ? `Campaigns for ${username}` : `User Campaigns`;

    const getStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case 'draft':
                return 'default';
            case 'failed':
                return 'error';
            case 'sent':
                return 'success';
            case 'sending':
                return 'processing';
            default:
                return 'default';
        }
    };

    const getDropdownItems = (campaign: any) => [
        {
            key: 'view',
            label: (
                <Space>
                    <EyeOutlined />
                    View
                </Space>
            ),
            onClick: () => navigate(`/zen/campaigns/specific/${campaign.id}/${id}/${companyId}`),
        },
        {
            key: 'suspend',
            label: (
                <Space>
                    <PauseOutlined />
                    Suspend
                </Space>
            ),
            onClick: () => suspendCampaign(campaign.id)
        },
        {
            key: 'delete',
            label: (
                <Space>
                    <DeleteOutlined />
                    Delete
                </Space>
            ),
            danger: true,
            onClick: () => deleteCamp(campaign.id)
        }
    ];

    const columns = [
        {
            title: 'Campaign Name',
            dataIndex: 'name',
            key: 'name',
            render: (text: string) => (
                <Text strong className="text-gray-900">{text}</Text>
            )
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status: string) => (
                <Tag
                    color={getStatusColor(status)}
                    className="font-medium"
                >
                    {status.charAt(0).toUpperCase() + status.slice(1)}
                </Tag>
            )
        },
        {
            title: 'Created On',
            dataIndex: 'created_at',
            key: 'created_at',
            render: (date: string) => (
                <Text className="text-gray-600">
                    {parseDate(date).toLocaleString('en-US', {
                        timeZone: 'UTC',
                        year: 'numeric',
                        month: 'short',
                        day: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                    })}
                </Text>
            )
        },
        {
            title: 'Actions',
            key: 'actions',
            width: 100,
            render: (record: any) => {
                const isSent = record.sent_at !== null;

                if (isSent) {
                    return (
                        <Button
                            type="link"
                            icon={<FileTextOutlined />}
                            onClick={() => navigate(`/user/dash/campaign/report/${record.id}`)}
                            className="text-blue-600 hover:text-blue-700"
                        >
                            View Report
                        </Button>
                    );
                }

                return (
                    <Dropdown
                        menu={{ items: getDropdownItems(record) }}
                        trigger={['click']}
                        placement="bottomRight"
                    >
                        <Button
                            type="text"
                            icon={<MoreOutlined />}
                            className="text-gray-500 hover:text-gray-700"
                        />
                    </Dropdown>
                );
            }
        }
    ];

    const breadcrumbItems = [
        {
            title: <HomeOutlined />,
            href: '/zen'
        },
        {
            title: (
                <Space>
                    <UserOutlined />
                    Users
                </Space>
            ),
            href: '/zen/users'
        },
        {
            title: username || 'User Campaigns'
        }
    ];

    return (
        <HelmetProvider>
            <Helmet title={pageTitle} />

            <div className="min-h-screen bg-gray-50 p-6">
                <div className="max-w-7xl mx-auto">
                    {/* Header Section */}
                    <div className="mb-6">
                        <div className="flex items-center gap-4 mb-4">
                            <Button
                                type="text"
                                icon={<ArrowLeftOutlined />}
                                onClick={() => window.history.back()}
                                className="text-gray-600 hover:text-gray-800"
                            />
                            <Breadcrumb items={breadcrumbItems} />
                        </div>

                        <Title level={2} className="!mb-2">
                            {username ? `Campaigns for ${username}` : 'User Campaigns'}
                        </Title>
                        <Text type="secondary" className="text-base">
                            Manage and monitor user campaign activities
                        </Text>
                    </div>

                    {/* Search and Filter Section */}
                    <Card className="mb-6 shadow-sm">
                        <Row gutter={[16, 16]} align="middle">
                            <Col xs={24} sm={12} md={8} lg={6}>
                                <Search
                                    placeholder="Search campaigns..."
                                    allowClear
                                    enterButton={<SearchOutlined />}
                                    size="large"
                                    onSearch={handleSearch}
                                    onChange={(e) => handleSearch(e.target.value)}
                                    className="w-full"
                                />
                            </Col>
                            <Col xs={24} sm={12} md={8} lg={6}>
                                <Text type="secondary">
                                    Total: {data?.payload?.total || 0} campaigns
                                </Text>
                            </Col>
                        </Row>
                    </Card>

                    {/* Table Section */}
                    <Card className="shadow-sm">
                        <Table
                            columns={columns}
                            dataSource={campaignData}
                            loading={isLoading}
                            pagination={false}
                            rowKey="uuid"
                            className="w-full"
                            locale={{
                                emptyText: (
                                    <Empty
                                        image={Empty.PRESENTED_IMAGE_SIMPLE}
                                        description={
                                            <div>
                                                <Text type="secondary" className="block mb-2">
                                                    No campaigns found
                                                </Text>
                                                <Text type="secondary" className="text-sm">
                                                    This user hasn't created any campaigns yet
                                                </Text>
                                            </div>
                                        }
                                    />
                                )
                            }}
                            scroll={{ x: 800 }}
                        />

                        {/* Pagination */}
                        {campaignData && campaignData.length > 0 && (
                            <div className="mt-6 flex justify-center">
                                <Pagination
                                    current={currentPage}
                                    pageSize={pageSize}
                                    total={data?.payload?.total || 0}
                                    onChange={(page, size) => {
                                        setCurrentPage(page);
                                        setPageSize(size);
                                    }}
                                    showSizeChanger
                                    pageSizeOptions={["10", "20", "50", "100"]}
                                    showTotal={(total, range) =>
                                        `${range[0]}-${range[1]} of ${total} campaigns`
                                    }
                                    className="mt-4"
                                />
                            </div>
                        )}
                    </Card>
                </div>
            </div>
        </HelmetProvider>
    );
};

export default AdminUserCampaigns;