import React, { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { 
    Table, 
    Input, 
    Switch, 
    Tag, 
    Space, 
    Button, 
    Modal, 
    Avatar, 
    Tooltip,
    Card,
    Badge
} from 'antd';
import { 
    SearchOutlined, 
    EyeOutlined, 
    UserOutlined, 
    CheckCircleOutlined, 
    StopOutlined,
    MailOutlined,
    PhoneOutlined,
    BankOutlined
} from '@ant-design/icons';
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import useDebounce from '../../../hooks/useDebounce';
import { useVerifiedUsersQuery } from '../hooks/useAdminUsersQueryHook';
import useAdminUserStore from '../store/adminuser.store';
import LoadingSpinnerComponent from '../../../components/loadingSpinnerComponent';

const VerifiedUsersTable: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState<string>("");
    const debouncedSearchQuery = useDebounce(searchQuery, 300);
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const navigate = useNavigate();

    const { blockUser, verifyUser, unBlockUser } = useAdminUserStore();
    const { data: userdata, isLoading } = useVerifiedUsersQuery(currentPage, pageSize, debouncedSearchQuery);

    const userdetailsData = useMemo(() => userdata?.payload.data, [userdata]);

    const handleToggle = async (userId: string, field: "verified" | "blocked", isChecked: boolean) => {
        const actionText = isChecked ?
            (field === "blocked" ? "unblock" : "unverify") :
            (field === "blocked" ? "block" : "verify");

        try {
            Modal.confirm({
                title: "Confirm Action",
                content: `Are you sure you want to ${actionText} this user?`,
                okText: "Yes",
                cancelText: "No",
                okButtonProps: { 
                    className: field === "blocked" ? "bg-red-500 hover:bg-red-600" : "bg-blue-500 hover:bg-blue-600"
                },
                onOk: async () => {
                    if (field === "blocked") {
                        if (isChecked) {
                            await unBlockUser(userId);
                        } else {
                            await blockUser(userId);
                        }
                    } else if (field === "verified") {
                        await verifyUser(userId);
                    }
                    setTimeout(() => location.reload(), 1000);
                },
            });
        } catch (error) {
            console.error(error);
        }
    };

    const columns: ColumnsType<any> = [
        {
            title: 'User',
            key: 'user',
            render: (_, record) => (
                <div className="flex items-center space-x-3">
                    <Avatar 
                        size={40} 
                        icon={<UserOutlined />} 
                        className="bg-gradient-to-r from-green-500 to-blue-600"
                    />
                    <div>
                        <div className="font-medium text-gray-900">{record.fullname}</div>
                        <div className="flex items-center text-sm text-gray-500">
                            <MailOutlined className="mr-1" />
                            {record.email}
                        </div>
                    </div>
                </div>
            ),
            width: 280,
        },
        {
            title: 'Company',
            dataIndex: 'company',
            key: 'company',
            render: (company) => (
                <div className="flex items-center">
                    <BankOutlined className="mr-2 text-gray-400" />
                    <span>{company || 'N/A'}</span>
                </div>
            ),
        },
        {
            title: 'Phone',
            dataIndex: 'phonenumber',
            key: 'phonenumber',
            render: (phone) => (
                <div className="flex items-center">
                    <PhoneOutlined className="mr-2 text-gray-400" />
                    <span>{phone || 'N/A'}</span>
                </div>
            ),
        },
        {
            title: 'Status',
            key: 'status',
            render: (_, record) => (
                <Space direction="vertical" size={4}>
                    <Tag 
                        icon={<CheckCircleOutlined />} 
                        color="green"
                        className="rounded-full"
                    >
                        Verified
                    </Tag>
                    {record.blocked && (
                        <Tag 
                            icon={<StopOutlined />} 
                            color="red" 
                            className="rounded-full"
                        >
                            Blocked
                        </Tag>
                    )}
                </Space>
            ),
        },
        {
            title: 'Blocked',
            key: 'blocked',
            render: (_, record) => (
                <Switch
                    checked={record.blocked}
                    onChange={() => handleToggle(record.id, 'blocked', record.blocked)}
                    checkedChildren="✓"
                    unCheckedChildren="✗"
                    className={record.blocked ? "bg-red-500" : "bg-gray-300"}
                />
            ),
        },
        {
            title: 'Verified Date',
            dataIndex: 'verified_at',
            key: 'verified_at',
            render: (date) => (
                <div className="text-sm">
                    {date ? new Date(date).toLocaleDateString('en-US', {
                        year: 'numeric',
                        month: 'short',
                        day: 'numeric'
                    }) : 'N/A'}
                </div>
            ),
        },
        {
            title: 'Actions',
            key: 'actions',
            render: (_, record) => (
                <Tooltip title="View Details">
                    <Button
                        type="primary"
                        icon={<EyeOutlined />}
                        size="small"
                        className="bg-green-500 hover:bg-green-600 border-none rounded-full"
                        onClick={() => navigate("/zen/users/detail/" + record.id)}
                    />
                </Tooltip>
            ),
        },
    ];

    const handleTableChange = (pagination: TablePaginationConfig) => {
        setCurrentPage(pagination.current || 1);
        setPageSize(pagination.pageSize || 20);
    };

    if (isLoading) {
        return <LoadingSpinnerComponent />;
    }

    return (
        <div className="space-y-6">
            {/* Search Section */}
            <Card className="shadow-sm border-0">
                <div className="flex items-center justify-between">
                    <div className="flex-1 max-w-md">
                        <Input
                            placeholder="Search verified users..."
                            prefix={<SearchOutlined className="text-gray-400" />}
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="rounded-lg border-gray-200 focus:border-green-500"
                            size="large"
                        />
                    </div>
                    <div className="flex items-center space-x-4">
                        <Badge 
                            count={userdata?.payload?.total || 0} 
                            showZero 
                            className="bg-green-100 text-green-800"
                        />
                        <span className="text-sm text-gray-600">Verified Users</span>
                    </div>
                </div>
            </Card>

            {/* Table Section */}
            <Card className="shadow-sm border-0 overflow-hidden">
                <Table
                    columns={columns}
                    dataSource={userdetailsData}
                    rowKey="id"
                    pagination={{
                        current: currentPage,
                        pageSize: pageSize,
                        total: userdata?.payload?.total || 0,
                        showSizeChanger: true,
                        showQuickJumper: true,
                        showTotal: (total, range) => 
                            `${range[0]}-${range[1]} of ${total} verified users`,
                        pageSizeOptions: ['10', '20', '50', '100'],
                        className: "px-4 pb-4"
                    }}
                    onChange={handleTableChange}
                    className="modern-table"
                    scroll={{ x: 1200 }}
                />
            </Card>
        </div>
    );
};

export default VerifiedUsersTable;