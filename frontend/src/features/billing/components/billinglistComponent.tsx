import React, { useMemo } from "react";
import { Table, Tag, Space, Typography, Card, Badge } from "antd";
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table';
import { useBillingQuery } from '../hooks/useBillingQuery';
import { BillingData } from '../interface/billing.interface';

const { Text } = Typography;

const BillingList: React.FC = () => {
    const [currentPage, setCurrentPage] = React.useState(1);
    const [pageSize, setPageSize] = React.useState(20);
    const { data: billingData, isLoading } = useBillingQuery(currentPage, pageSize);

    const billingRecords = useMemo(() => billingData?.payload?.data || [], [billingData]);

    // Status color mapping
    const getStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case 'completed':
                return 'success';
            case 'active':
                return 'processing';
            case 'expired':
                return 'error';
            case 'pending':
                return 'warning';
            default:
                return 'default';
        }
    };

    // Subscription status color mapping
    const getSubscriptionStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case 'active':
                return 'success';
            case 'expired':
                return 'error';
            case 'trial':
                return 'processing';
            default:
                return 'default';
        }
    };

    // Format currency
    const formatCurrency = (amount: string, currency: string) => {
        const numAmount = parseFloat(amount);
        if (numAmount === 0) return 'Free';
        return `${currency} ${numAmount.toLocaleString()}`;
    };

    // Format date
    const formatDate = (dateString: string | null) => {
        if (!dateString) return 'N/A';
        return new Date(dateString).toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: '2-digit'
        });
    };

    const columns: ColumnsType<BillingData> = [
        {
            title: 'Paid By',
            key: 'user',
            width: 180,
            fixed: 'left',
            render: (_, record) => (
                <Space direction="vertical" size={1}>
                    <Text strong style={{ fontSize: '13px' }}>
                        {record.user.userfullname}
                    </Text>
                    <Text type="secondary" style={{ fontSize: '11px' }} ellipsis>
                        {record.user.useremail}
                    </Text>
                </Space>
            ),
        },
        {
            title: 'Plan',
            dataIndex: ['subscription', 'plan', 'plan_name'],
            key: 'plan',
            width: 90,
            render: (planName: string) => (
                <Tag color="blue" style={{ fontSize: '11px' }}>
                    {planName}
                </Tag>
            ),
        },
        {
            title: 'Amount',
            key: 'amount',
            width: 90,
            render: (_, record) => (
                <Text strong style={{ fontSize: '12px' }}>
                    {formatCurrency(record.amount, record.currency)}
                </Text>
            ),
            sorter: (a, b) => parseFloat(a.amount) - parseFloat(b.amount),
        },
        {
            title: 'Method',
            dataIndex: 'payment_method',
            key: 'payment_method',
            width: 85,
            render: (method: string) => (
                <Tag 
                    color={method === 'None' ? 'default' : 'geekblue'} 
                    style={{ fontSize: '10px' }}
                >
                    {method}
                </Tag>
            ),
        },
        {
            title: 'Status',
            key: 'statuses',
            width: 110,
            render: (_, record) => (
                <Space direction="vertical" size={2}>
                    <Badge 
                        status={getStatusColor(record.status) as any} 
                        text={<span style={{ fontSize: '11px' }}>{record.status}</span>}
                    />
                    <Tag 
                        color={getSubscriptionStatusColor(record.subscription.subscriptionstatus)}
                        style={{ fontSize: '10px', margin: 0 }}
                    >
                        Sub: {record.subscription.subscriptionstatus}
                    </Tag>
                </Space>
            ),
        },
        {
            title: 'Cycle',
            dataIndex: ['subscription', 'subscriptionbillingcycle'],
            key: 'billing_cycle',
            width: 70,
            render: (cycle: string) => (
                <Text style={{ fontSize: '11px' }}>{cycle}</Text>
            ),
        },
        {
            title: 'Start Date',
            key: 'start_date',
            width: 100,
            render: (_, record) => (
                <Text style={{ fontSize: '11px' }}>
                    {formatDate(record.subscription.subscriptionstartsat)}
                </Text>
            ),
            sorter: (a, b) => {
                const dateA = a.subscription.subscriptionstartsat ? new Date(a.subscription.subscriptionstartsat).getTime() : 0;
                const dateB = b.subscription.subscriptionstartsat ? new Date(b.subscription.subscriptionstartsat).getTime() : 0;
                return dateA - dateB;
            },
        },
        {
            title: 'End Date',
            key: 'end_date',
            width: 100,
            render: (_, record) => (
                <Text style={{ fontSize: '11px' }}>
                    {formatDate(record.subscription.subscriptionendsat)}
                </Text>
            ),
            sorter: (a, b) => {
                const dateA = a.subscription.subscriptionendsat ? new Date(a.subscription.subscriptionendsat).getTime() : 0;
                const dateB = b.subscription.subscriptionendsat ? new Date(b.subscription.subscriptionendsat).getTime() : 0;
                return dateA - dateB;
            },
        },
        {
            title: 'Transaction',
            dataIndex: 'payment_id',
            key: 'payment_id',
            width: 100,
            render: (paymentId: string | null) => (
                <Text 
                    code 
                    style={{ fontSize: '10px' }}
                    copyable={paymentId ? { text: paymentId } : false}
                >
                    {paymentId ? paymentId.substring(0, 6) + '...' : 'N/A'}
                </Text>
            ),
        },
    ];

    const handleTableChange = (pagination: TablePaginationConfig) => {
        setCurrentPage(pagination.current || 1);
        setPageSize(pagination.pageSize || 20);
    };

    const paginationConfig: TablePaginationConfig = {
        current: currentPage,
        pageSize: pageSize,
        total: billingData?.payload?.total || 0,
        showSizeChanger: true,
        showQuickJumper: true,
        showTotal: (total, range) => 
            `${range[0]}-${range[1]} of ${total} billing records`,
        pageSizeOptions: ['10', '20', '50', '100'],
    };

    return (
        <Card 
            title="Billing Records" 
            className="mt-4"
            styles={{
                header: { 
                    borderBottom: '1px solid #f0f0f0',
                    marginBottom: '16px'
                },
                body: {
                    padding: 0,
                    overflow: 'hidden'
                }
            }}
        >
            <div style={{ 
                width: '100%', 
                overflowX: 'auto',
                overflowY: 'hidden'
            }}>
                <Table<BillingData>
                    columns={columns}
                    dataSource={billingRecords}
                    rowKey="id"
                    loading={isLoading}
                    pagination={{
                        ...paginationConfig,
                        style: { padding: '16px' }
                    }}
                    onChange={handleTableChange}
                    scroll={{ x: 'max-content' }}
                    size="middle"
                    locale={{
                        emptyText: 'No billing records found'
                    }}
                    className="billing-table"
                    sticky
                />
            </div>
        </Card>
    );
};

export default BillingList;