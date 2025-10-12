import { useState, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import {
    Table,
    Button,
    Modal,
    Input,
    Tag,
    Dropdown,
    Card,
    Typography,
 
    Layout,
    Empty
} from 'antd';
import {
    MoreOutlined,
    EditOutlined,
    DeleteOutlined,
    EyeOutlined,
    PlusOutlined,
    SearchOutlined,
 
    MailOutlined
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import type { MenuProps } from 'antd';
import useCampaignStore from "../store/campaign.store";
import { useCampaignQuery } from "../hooks/useCampaignQuery";
import useDebounce from "../../../hooks/useDebounce";
import { parseDate } from '../../../../../frontend/src/utils/utils';
import CreateCampaignComponent from './createCampaignComponent';
import { CampaignResponse } from "../interface/campaign.interface";

const { Title } = Typography;
const { Content } = Layout;
const { Search } = Input;

const GetAllCampaignComponent: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { deleteCampaign } = useCampaignStore();
    const [searchQuery, setSearchQuery] = useState<string>("");
    const navigate = useNavigate();
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);

    const debouncedSearchQuery = useDebounce(searchQuery, 300);
    const { data: CampaignData, isLoading,refetch } = useCampaignQuery(currentPage, pageSize, debouncedSearchQuery);

    const deleteCamp = async (uuid: string) => {
        Modal.confirm({
            title: "Delete Campaign",
            content: "Are you sure you want to delete this campaign? This action cannot be undone.",
            okText: "Delete",
            cancelText: "Cancel",
            okType: "danger",
            centered: true,
            onOk: async () => {
                await deleteCampaign(uuid);
                await new Promise(resolve => setTimeout(resolve, 1000));
                location.reload();
            },
        });
    };

    const getStatusColor = (status: string): string => {
        switch (status.toLowerCase()) {
            case 'draft':
                return 'orange';
            case 'failed':
                return 'red';
            case 'sent':
                return 'green';
            case 'sending':
                return 'blue';
            case 'scheduled':
                return 'purple';
            default:
                return 'default';
        }
    };

    const getMenuItems = (campaign: CampaignResponse): MenuProps['items'] => [
        {
            key: 'edit',
            label: 'Edit Campaign',
            icon: <EditOutlined />,
            onClick: () => navigate(`/app/campaign/edit/${campaign.id}`),
        },
        {
            type: 'divider',
        },
        {
            key: 'delete',
            label: 'Delete Campaign',
            icon: <DeleteOutlined />,
            danger: true,
            onClick: () => deleteCamp(campaign.id),
        },
    ];

    const getActionButton = (campaign: CampaignResponse) => {
        const isSent = campaign.sent_at !== null;

        if (isSent) {
            return (
                <Button
                    type="primary"
                    icon={<EyeOutlined />}
                    size="small"
                    onClick={() => navigate(`/app/campaign/report/${campaign.id}`)}
                    className="bg-blue-600 hover:bg-blue-700"
                >
                    View Report
                </Button>
            );
        }

        return (
            <Dropdown
                menu={{ items: getMenuItems(campaign) }}
                trigger={['click']}
                placement="bottomRight"
            >
                <Button
                    type="text"
                    icon={<MoreOutlined />}
                    size="small"
                    className="hover:bg-gray-100"
                />
            </Dropdown>
        );
    };

    const columns: ColumnsType<CampaignResponse> = [
        {
            title: 'Campaign Name',
            dataIndex: 'name',
            key: 'name',
            sorter: true,
            ellipsis: true,
            render: (name: string) => (
                <div className="font-medium text-gray-900">{name}</div>
            ),
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            width: 120,
            render: (status: string) => (
                <Tag
                    color={getStatusColor(status)}
                    className="font-medium"
                >
                    {status.charAt(0).toUpperCase() + status.slice(1)}
                </Tag>
            ),
            filters: [
                { text: 'Draft', value: 'draft' },
                { text: 'Sent', value: 'Sent' },
                { text: 'Sending', value: 'Sending' },
                { text: 'Failed', value: 'Failed' },
                { text: 'Scheduled', value: 'scheduled' },
            ],
            onFilter: (value, record) => record.status.toLowerCase() === value,
        },
        {
            title: 'Created On',
            dataIndex: 'created_at',
            key: 'created_at',
            width: 200,
            render: (date: string) => (
                <div className="text-gray-600">
                    {parseDate(date).toLocaleString('en-US', {
                        timeZone: 'UTC',
                        year: 'numeric',
                        month: 'short',
                        day: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                    })}
                </div>
            ),
            sorter: true,
        },
        {
            title: 'Actions',
            key: 'action',
            width: 120,
            render: (_, campaign: CampaignResponse) => getActionButton(campaign),
        },
    ];

    const handleTableChange = (pagination: TablePaginationConfig) => {
        setCurrentPage(pagination.current || 1);
        setPageSize(pagination.pageSize || 20);
    };

    const handleSearch = (value: string) => {
        setSearchQuery(value);
        setCurrentPage(1); // Reset to first page when searching
    };

    const cData = useMemo(() => CampaignData?.payload.data || [], [CampaignData]);

    return (
        <Layout className="bg-gray-50 min-h-screen">
            <Content className="p-6">
                <div className="max-w-7xl mx-auto">
                    {/* Header Section */}
                    <div className="mb-8">
                        <div className="flex items-center justify-between">
                            <div className="flex items-center space-x-3">
                                <MailOutlined className="text-2xl text-blue-600" />
                                <Title level={2} className="!mb-0 !text-gray-800">
                                    All Campaigns
                                </Title>
                            </div>
                            <Button
                                type="primary"
                                icon={<PlusOutlined />}
                                size="large"
                                onClick={() => setIsModalOpen(true)}
                                className="bg-blue-600 hover:bg-blue-700 border-blue-600 shadow-lg"
                            >
                                Create Campaign
                            </Button>
                        </div>
                        <p className="text-gray-600 mt-2">
                            Manage all your email campaigns and track their performance
                        </p>
                    </div>

                    {/* Search Section */}
                    <Card className="mb-6 shadow-sm border-0">
                        <div className="flex items-center justify-between">
                            <Search
                                placeholder="Search campaigns by name..."
                                allowClear
                                size="large"
                                style={{ width: 400 }}
                                prefix={<SearchOutlined className="text-gray-400" />}
                                value={searchQuery}
                                onChange={(e) => handleSearch(e.target.value)}
                                onSearch={handleSearch}
                                className="rounded-lg"
                            />
                            <div className="text-sm text-gray-500">
                                {CampaignData?.payload?.total || 0} campaigns total
                            </div>
                        </div>
                    </Card>

                    {/* Table Section */}
                    <Card className="shadow-sm border-0 rounded-lg overflow-hidden">
                        <Table<CampaignResponse>
                            columns={columns}
                            dataSource={cData}
                            rowKey="id"
                            loading={isLoading}
                            pagination={{
                                current: currentPage,
                                pageSize: pageSize,
                                total: CampaignData?.payload?.total || 0,
                                showSizeChanger: true,
                                showQuickJumper: true,
                                pageSizeOptions: ['10', '20', '50', '100'],
                                showTotal: (total, range) =>
                                    `${range[0]}-${range[1]} of ${total} campaigns`,
                                onChange: (page, size) => {
                                    setCurrentPage(page);
                                    setPageSize(size || 20);
                                },
                                className: "px-6 pb-4",
                            }}
                            onChange={handleTableChange}
                            locale={{
                                emptyText: (
                                    <div className="py-16">
                                        <Empty
                                            image={Empty.PRESENTED_IMAGE_SIMPLE}
                                            description={
                                                <div className="space-y-2">
                                                    <div className="text-gray-500 text-lg font-medium">
                                                        No Campaigns Found
                                                    </div>
                                                    <div className="text-gray-400">
                                                        {searchQuery
                                                            ? `No campaigns match "${searchQuery}"`
                                                            : "Create your first campaign to get started"
                                                        }
                                                    </div>
                                                </div>
                                            }
                                        >
                                            {!searchQuery && (
                                                <Button
                                                    type="primary"
                                                    icon={<PlusOutlined />}
                                                    onClick={() => setIsModalOpen(true)}
                                                    className="bg-blue-600 hover:bg-blue-700"
                                                >
                                                    Create Your First Campaign
                                                </Button>
                                            )}
                                        </Empty>
                                    </div>
                                ),
                            }}
                            scroll={{ x: 800 }}
                            size="middle"
                            className="[&_.ant-table-thead>tr>th]:bg-gray-50 [&_.ant-table-thead>tr>th]:border-b [&_.ant-table-thead>tr>th]:border-gray-200 [&_.ant-table-thead>tr>th]:font-semibold [&_.ant-table-thead>tr>th]:text-gray-600 [&_.ant-table-thead>tr>th]:p-4 [&_.ant-table-tbody>tr>td]:p-4 [&_.ant-table-tbody>tr>td]:border-b [&_.ant-table-tbody>tr>td]:border-gray-100 [&_.ant-table-tbody>tr:last-child>td]:border-b-0"
                            rowClassName="hover:bg-gray-50 transition-colors duration-200"
                        />
                    </Card>

                    <CreateCampaignComponent
                        isOpen={isModalOpen}
                        onClose={() => setIsModalOpen(false)}
                        refetch={()=>refetch()}
                    />
                </div>
            </Content>
        </Layout>
    );
};

export default GetAllCampaignComponent;