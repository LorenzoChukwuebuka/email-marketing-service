import { useState, useMemo } from "react";
import {
    Table,
    Button,
    Input,
    Dropdown,
    Tag,
    Card,
    Typography,
    Spin,
    Empty,
    MenuProps
} from "antd";
import {
    PlusOutlined,
    SearchOutlined,
    MoreOutlined,
    EditOutlined,
    DeleteOutlined,
    CalendarOutlined
} from "@ant-design/icons";
import CreateCampaignComponent from "./createCampaignComponent";
import { useNavigate } from "react-router-dom";
import { parseDate } from "../../../utils/utils";
import useDebounce from "../../../hooks/useDebounce";
import { useScheduledCampaignQuery } from "../hooks/useCampaignQuery";
import { CampaignResponse } from "../interface/campaign.interface";

const { Title } = Typography;
const { Search } = Input;

const GetScheduledCampaignComponent: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const navigate = useNavigate();
    const [searchQuery, setSearchQuery] = useState<string>("");
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const debouncedSearchQuery = useDebounce(searchQuery, 300);

    const { data: scheduleCampaignData, isLoading,refetch } = useScheduledCampaignQuery(
        currentPage,
        pageSize,
        debouncedSearchQuery
    );

    const schCdata = useMemo(() => scheduleCampaignData?.payload.data || [], [scheduleCampaignData]);

    const deleteCampaign = async (uuid: string) => {
        console.log(uuid);
        // Add your delete logic here
    };

    const getStatusColor = (status: string) => {
        const statusColors: { [key: string]: string } = {
            'active': 'green',
            'pending': 'orange',
            'paused': 'red',
            'completed': 'blue',
            'draft': 'default'
        };
        return statusColors[status.toLowerCase()] || 'default';
    };

    const getDropdownItems = (campaign: CampaignResponse): MenuProps['items'] => [
        {
            key: 'edit',
            icon: <EditOutlined />,
            label: 'Edit',
            onClick: () => navigate(`/user/dash/campaign/edit/${campaign.id}`)
        },
        {
            key: 'delete',
            icon: <DeleteOutlined />,
            label: 'Delete',
            danger: true,
            onClick: () => deleteCampaign(campaign.id)
        }
    ];

    const columns = [
        {
            title: 'Campaign Name',
            dataIndex: 'name',
            key: 'name',
            render: (text: string) => (
                <div className="font-medium text-gray-900">{text}</div>
            ),
            sorter: (a: CampaignResponse, b: CampaignResponse) => a.name.localeCompare(b.name),
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status: string) => (
                <Tag color={getStatusColor(status)} className="font-medium">
                    {status.charAt(0).toUpperCase() + status.slice(1)}
                </Tag>
            ),
            filters: [
                { text: 'Active', value: 'active' },
                { text: 'Pending', value: 'pending' },
                { text: 'Paused', value: 'paused' },
                { text: 'Completed', value: 'completed' },
                { text: 'Draft', value: 'draft' },
            ],
            onFilter: (value: any, record: CampaignResponse) => record.status === value,
        },
        {
            title: 'Created On',
            dataIndex: 'created_at',
            key: 'created_at',
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
            sorter: (a: CampaignResponse, b: CampaignResponse) =>
                new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
            defaultSortOrder: 'descend' as const,
        },
        {
            title: 'Actions',
            key: 'actions',
            width: 80,
            render: (_: any, record: CampaignResponse) => (
                <Dropdown
                    menu={{ items: getDropdownItems(record) }}
                    trigger={['click']}
                    placement="bottomRight"
                >
                    <Button
                        type="text"
                        icon={<MoreOutlined />}
                        className="hover:bg-gray-100"
                    />
                </Dropdown>
            ),
        },
    ];

    const handleTableChange = (pagination: any) => {
        setCurrentPage(pagination.current);
        setPageSize(pagination.pageSize);
    };

    const customEmptyState = (
        <div className="py-16">
            <Empty
                image={<CalendarOutlined className="text-6xl text-gray-300" />}
                description={
                    <div className="text-center">
                        <Title level={4} className="text-gray-500 mb-2">
                            No Scheduled Campaigns
                        </Title>
                        <p className="text-gray-400 mb-4">
                            Scheduled campaigns are campaigns that are sent on a later date
                        </p>
                        <Button
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={() => setIsModalOpen(true)}
                            className="bg-blue-600 hover:bg-blue-700"
                        >
                            Create Your First Campaign
                        </Button>
                    </div>
                }
            />
        </div>
    );

    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <Spin size="large" />
            </div>
        );
    }

    return (
        <div className="p-6 bg-gray-50 min-h-screen">
            <div className="max-w-7xl mx-auto">
                {/* Header Section */}
                <Card className="mb-6 shadow-sm border-0 rounded-xl">
                    <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                        <div>
                            <Title level={2} className="mb-2 text-gray-800">
                                Scheduled Campaigns
                            </Title>
                            <p className="text-gray-600 m-0">
                                Manage and monitor your scheduled email campaigns
                            </p>
                        </div>
                        <Button
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={() => setIsModalOpen(true)}
                            size="large"
                            className="bg-blue-600 hover:bg-blue-700 border-blue-600 shadow-sm rounded-lg"
                        >
                            Create Campaign
                        </Button>
                    </div>
                </Card>

                {/* Search and Filters */}
                <Card className="mb-6 shadow-sm border-0 rounded-xl">
                    <div className="flex flex-col sm:flex-row gap-4">
                        <Search
                            placeholder="Search campaigns..."
                            allowClear
                            onChange={(e) => setSearchQuery(e.target.value)}
                            prefix={<SearchOutlined className="text-gray-400" />}
                            className="sm:max-w-md rounded-lg"
                            size="large"
                        />
                    </div>
                </Card>

                {/* Table Section */}
                <Card className="shadow-sm border-0 rounded-xl">
                    <Table
                        columns={columns}
                        dataSource={schCdata}
                        rowKey="uuid"
                        loading={isLoading}
                        locale={{ emptyText: customEmptyState }}
                        pagination={{
                            current: currentPage,
                            pageSize: pageSize,
                            total: scheduleCampaignData?.payload?.total || 0,
                            showSizeChanger: true,
                            showQuickJumper: true,
                            showTotal: (total, range) =>
                                `${range[0]}-${range[1]} of ${total} campaigns`,
                            pageSizeOptions: ['10', '20', '50', '100'],
                            className: 'mt-6',
                        }}
                        onChange={handleTableChange}
                        className="campaign-table [&_.ant-table-thead>tr>th]:bg-gray-50 [&_.ant-table-thead>tr>th]:font-semibold [&_.ant-table-thead>tr>th]:text-gray-700 [&_.ant-table-tbody>tr>td]:border-b [&_.ant-table-tbody>tr>td]:border-gray-100 [&_.ant-table-tbody>tr>td]:py-4"
                        scroll={{ x: 800 }}
                        rowClassName="hover:bg-gray-50/50 transition-colors duration-200"
                    />
                </Card>

                <CreateCampaignComponent
                refetch={refetch}
                    isOpen={isModalOpen}
                    onClose={() => setIsModalOpen(false)}
                />
            </div>


        </div>
    );
};

export default GetScheduledCampaignComponent;