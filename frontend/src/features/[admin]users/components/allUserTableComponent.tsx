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
import { useAllUsersQuery } from '../hooks/useAdminUsersQueryHook';
import useAdminUserStore from '../store/adminuser.store';
import LoadingSpinnerComponent from '../../../components/loadingSpinnerComponent';


const AllUsersTable: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState<string>("");
    const debouncedSearchQuery = useDebounce(searchQuery, 300);
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const navigate = useNavigate();

    const { blockUser, verifyUser, unBlockUser } = useAdminUserStore();
    const { data: userData, isLoading, refetch } = useAllUsersQuery(currentPage, pageSize, debouncedSearchQuery);

    const userdetailsData = useMemo(() => userData?.payload.data, [userData]);

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

                    setTimeout(() => refetch(), 1000);
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
                        className="bg-gradient-to-r from-blue-500 to-purple-600"
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
                        color={record.verified ? 'green' : 'orange'}
                        className="rounded-full"
                    >
                        {record.verified ? 'Verified' : 'Unverified'}
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
            title: 'Verified',
            key: 'verified',
            render: (_, record) => (
                <Switch
                    checked={record.verified}
                    onChange={() => handleToggle(record.id, 'verified', record.verified)}
                    checkedChildren="✓"
                    unCheckedChildren="✗"
                    className="bg-gray-300"
                />
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
            title: 'Created',
            dataIndex: 'created_at',
            key: 'created_at',
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
                        className="bg-blue-500 hover:bg-blue-600 border-none rounded-full"
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
                            placeholder="Search users by name, email, or company..."
                            prefix={<SearchOutlined className="text-gray-400" />}
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="rounded-lg border-gray-200 focus:border-blue-500"
                            size="large"
                        />
                    </div>
                    <div className="flex items-center space-x-4">
                        <Badge
                            count={userData?.payload?.total || 0}
                            showZero
                            className="bg-blue-100 text-blue-800"
                        />
                        <span className="text-sm text-gray-600">Total Users</span>
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
                        total: userData?.payload?.total || 0,
                        showSizeChanger: true,
                        showQuickJumper: true,
                        showTotal: (total, range) =>
                            `${range[0]}-${range[1]} of ${total} users`,
                        pageSizeOptions: ['10', '20', '50', '100'],
                        className: "px-4 pb-4"
                    }}
                    onChange={handleTableChange}
                    className="modern-table"
                    scroll={{ x: 1200 }}
                />
            </Card>

            <style
                dangerouslySetInnerHTML={{
                    __html: `
                        .modern-table .ant-table-thead > tr > th {
                            background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
                            border-bottom: 2px solid #e2e8f0;
                            font-weight: 600;
                            color: #475569;
                        }
                        
                        .modern-table .ant-table-tbody > tr:hover > td {
                            background: #f8fafc;
                        }
                        
                        .modern-table .ant-table-tbody > tr > td {
                            border-bottom: 1px solid #f1f5f9;
                            padding: 16px;
                        }
                        
                        .modern-table .ant-switch-checked {
                            background-color: #10b981;
                        }
                        
                        .modern-table .ant-tag {
                            border-radius: 20px;
                            padding: 4px 12px;
                            font-size: 12px;
                            font-weight: 500;
                        }
                    `
                }}
            />
        </div>
    );
};

export default AllUsersTable;