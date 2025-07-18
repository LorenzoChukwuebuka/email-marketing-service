// AllSupportTicketComponentTable.tsx
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAdminAllTicketsQuery } from "../../hooks/useAdminSupportTicketQuery";
import useDebounce from "../../../../hooks/useDebounce";
import {
    Table,
    Input,
    Space,
    Tag,
    Button,
    Card,
    Typography,
    Tooltip,
    Badge
} from "antd";
import {
    SearchOutlined,
    EyeOutlined,
    UserOutlined,
    MailOutlined,
    NumberOutlined,
    ClockCircleOutlined
} from "@ant-design/icons";
import type { ColumnsType } from 'antd/es/table';
import { AdminTicketData } from "../../interface/support.interface";

const { Title } = Typography;


const AllSupportTicketComponentTable: React.FC = () => {
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [searchQuery, setSearchQuery] = useState<string>("");
    const debouncedSearchQuery = useDebounce(searchQuery, 300);
    const navigate = useNavigate();

    const { data: apiResponse, isLoading } = useAdminAllTicketsQuery(
        currentPage,
        pageSize,
        debouncedSearchQuery
    );

    const supportData = apiResponse?.payload?.data || [];
    const totalRecords = apiResponse?.payload?.total || 0;

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'open': return 'success';
            case 'pending': return 'warning';
            case 'closed': return 'default';
            case 'resolved': return 'success';
            default: return 'default';
        }
    };

    const getPriorityColor = (priority: string) => {
        switch (priority) {
            case 'high': return 'red';
            case 'medium': return 'orange';
            case 'low': return 'blue';
            default: return 'default';
        }
    };

    const formatDate = (dateString: string | null) => {
        if (!dateString) return "Not available";
        return new Date(dateString).toLocaleString('en-US', {
            timeZone: 'UTC',
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };



    const columns: ColumnsType<AdminTicketData> = [
        {
            title: <Space><UserOutlined /> Name</Space>,
            dataIndex: 'name',
            key: 'name',
            width: 150,
            ellipsis: true,
            render: (text: string) => (
                <Tooltip title={text}>
                    <span className="font-medium">{text}</span>
                </Tooltip>
            ),
        },
        {
            title: <Space><MailOutlined /> Email</Space>,
            dataIndex: 'email',
            key: 'email',
            width: 200,
            ellipsis: true,
            render: (text: string) => (
                <Tooltip title={text}>
                    <span className="text-blue-600">{text}</span>
                </Tooltip>
            ),
        },
        {
            title: 'Subject',
            dataIndex: 'subject',
            key: 'subject',
            width: 200,
            ellipsis: true,
            render: (text: string) => (
                <Tooltip title={text}>
                    <span className="font-medium">{text}</span>
                </Tooltip>
            ),
        },
        {
            title: 'Description',
            dataIndex: 'description',
            key: 'description',
            width: 250,
            ellipsis: true,
            render: (text: string) => (
                <Tooltip title={text}>
                    <span className="text-gray-600">{text || "N/A"}</span>
                </Tooltip>
            ),
        },
        {
            title: <Space><NumberOutlined /> Ticket #</Space>,
            dataIndex: 'ticket_number',
            key: 'ticket_number',
            width: 120,
            render: (text: string) => (
                <Badge
                    count={`#${text}`}
                    style={{ backgroundColor: '#52c41a' }}
                    className="font-mono"
                />
            ),
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status: string) => (
                <Tag color={getStatusColor(status)} className="font-medium capitalize">
                    {status}
                </Tag>
            ),
            filters: [
                { text: 'Open', value: 'open' },
                { text: 'Pending', value: 'pending' },
                { text: 'Closed', value: 'closed' },
                { text: 'Resolved', value: 'resolved' },
            ],
            onFilter: (value, record) => record.status === value,
        },
        {
            title: 'Priority',
            dataIndex: 'priority',
            key: 'priority',
            width: 100,
            render: (priority: string) => (
                <Tag color={getPriorityColor(priority)} className="font-medium capitalize">
                    {priority}
                </Tag>
            ),
            filters: [
                { text: 'High', value: 'high' },
                { text: 'Medium', value: 'medium' },
                { text: 'Low', value: 'low' },
            ],
            onFilter: (value, record) => record.priority === value,
        },
        {
            title: <Space><ClockCircleOutlined /> Last Reply</Space>,
            dataIndex: 'last_reply',
            key: 'last_reply',
            width: 160,
            render: (date: string | null) => (
                <span className="text-sm text-gray-500">
                    {formatDate(date)}
                </span>
            ),
            sorter: (a, b) => {
                const dateA = a.last_reply ? new Date(a.last_reply).getTime() : 0;
                const dateB = b.last_reply ? new Date(b.last_reply).getTime() : 0;
                return dateB - dateA;
            },
        },
        {
            title: 'Action',
            key: 'action',
            width: 80,
            render: (_, record) => (
                <Button
                    type="primary"
                    icon={<EyeOutlined />}
                    size="small"
                    onClick={() => navigate(`/zen/support/details/${record.id}`)}
                    className="bg-blue-500 hover:bg-blue-600"
                >
                    View
                </Button>
            ),
        },
    ];

    return (
        <div className="p-6 bg-gray-50 min-h-screen">
            <Card
                className="shadow-lg border-0"
                title={
                    <div className="flex justify-between items-center">
                        <Title level={3} className="mb-0 text-gray-800">
                            📋 All Support Tickets
                        </Title>
                        <Badge
                            count={totalRecords}
                            style={{ backgroundColor: '#1890ff' }}
                            className="ml-2"
                        />
                    </div>
                }
                extra={
                    <Input
                        placeholder="Search tickets..."
                        prefix={<SearchOutlined className="text-gray-400" />}
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="w-64"
                        allowClear
                    />
                }
            >
                <Table<AdminTicketData>
                    columns={columns}
                    dataSource={supportData as any}
                    loading={isLoading}
                    rowKey="id"
                    pagination={{
                        current: currentPage,
                        pageSize: pageSize,
                        total: totalRecords,
                        onChange: (page, size) => {
                            setCurrentPage(page);
                            setPageSize(size || 20);
                        },
                        showSizeChanger: true,
                        showQuickJumper: true,
                        showTotal: (total, range) =>
                            `${range[0]}-${range[1]} of ${total} tickets`,
                        pageSizeOptions: ['10', '20', '50', '100'],
                    }}
                    scroll={{ x: 1200 }}
                    className="custom-table"
                    rowClassName="hover:bg-blue-50 transition-colors duration-200"
                />
            </Card>
        </div>
    );
};

export default AllSupportTicketComponentTable;