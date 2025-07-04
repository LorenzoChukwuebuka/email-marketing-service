import { useState } from "react";
import { Table, Button, Modal, Card, Typography, Space, Tag, Tooltip, message } from 'antd';
import { CopyOutlined, DeleteOutlined, CheckOutlined, ApiOutlined } from '@ant-design/icons';
import { convertToNormalTime, maskAPIKey, copyToClipboard } from "../../../../utils/utils";
import EmptyState from "../../../../components/emptyStateComponent";
import useAPIKeyStore from "../../store/apikey.store";
import { useAPIKeyQuery } from "../../hooks/useApiKeyQuery";

const { Text } = Typography;

interface APIKey {
    id: string;
    user_id: string;
    company_id: string;
    name: string;
    api_key: string;
    created_at: string;
    updated_at: string;
}

interface APIResponse {
    message: string;
    payload: APIKey[];
    status: boolean;
}

const APIKeysTableComponent: React.FC = () => {
    const { deleteAPIKey } = useAPIKeyStore();
    const [deletingId, setDeletingId] = useState<string | null>(null);
    const [copyingKey, setCopyingKey] = useState<string | null>(null);

    const { data: apiKeyData } = useAPIKeyQuery() as { data: APIResponse };

    // Safe date formatting function
    const formatDate = (dateString: string | undefined) => {
        if (!dateString) return 'N/A';
        try {
            return convertToNormalTime(dateString);
        } catch (error) {
            return new Date(dateString).toLocaleDateString('en-US', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: 'numeric',
                minute: 'numeric'
            });
        }
    };

    const shouldRenderNoKey = () => {
        return (
            !apiKeyData ||
            !apiKeyData.payload ||
            !Array.isArray(apiKeyData.payload) ||
            apiKeyData.payload.length === 0
        );
    };

    const handleDelete = async (id: string, name: string) => {
        Modal.confirm({
            title: "Delete API Key",
            content: `Are you sure you want to delete the API key "${name}"? This action cannot be undone.`,
            okText: "Yes, Delete",
            cancelText: "Cancel",
            okType: "danger",
            onOk: async () => {
                setDeletingId(id);
                try {
                    await deleteAPIKey(id);
                    message.success('API key deleted successfully');
                    await new Promise(resolve => setTimeout(resolve, 1000));
                    location.reload();
                } catch (error) {
                    message.error('Failed to delete API key');
                } finally {
                    setDeletingId(null);
                }
            },
        });
    };

    const handleCopy = async (key: string, name: string) => {
        try {
            await copyToClipboard(key);
            setCopyingKey(key);
            message.success(`API key "${name}" copied to clipboard`);
            setTimeout(() => {
                setCopyingKey(null);
            }, 2000);
        } catch (error) {
            message.error('Failed to copy to clipboard');
        }
    };

    const columns = [
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
            render: (text: string) => <Text strong>{text}</Text>,
        },
        {
            title: 'API Key',
            dataIndex: 'api_key',
            key: 'api_key',
            render: (apiKey: string, record: APIKey) => (
                <Space>
                    <Text code style={{ fontSize: '12px' }}>{maskAPIKey(apiKey)}</Text>
                    <Tooltip title="Copy API key">
                        <Button
                            type="text"
                            size="small"
                            icon={copyingKey === apiKey ? <CheckOutlined /> : <CopyOutlined />}
                            onClick={() => handleCopy(apiKey, record.name)}
                            style={{ color: copyingKey === apiKey ? '#52c41a' : '#1890ff' }}
                        />
                    </Tooltip>
                </Space>
            ),
        },
        {
            title: 'Status',
            key: 'status',
            render: () => (
                <Tag color="green">ACTIVE</Tag>
            ),
        },
        {
            title: 'Created On',
            dataIndex: 'created_at',
            key: 'created_at',
            render: (date: string) => <Text>{formatDate(date)}</Text>,
        },
        {
            title: 'Actions',
            key: 'actions',
            render: (record: APIKey) => (
                <Button
                    type="text"
                    danger
                    icon={<DeleteOutlined />}
                    onClick={() => handleDelete(record.id, record.name)}
                    loading={deletingId === record.id}
                    disabled={deletingId === record.id}
                    size="small"
                >
                    Delete
                </Button>
            ),
        },
    ];

    if (shouldRenderNoKey()) {
        return (
            <Card>
                <EmptyState
                    title="You have not generated any API Key"
                    description="Kindly Generate an API Key to enjoy our services"
                    icon={<ApiOutlined style={{ fontSize: '48px', color: '#1890ff' }} />}
                />
            </Card>
        );
    }

    return (
        <Card
            title={
                <Space>
                    <ApiOutlined />
                    <span>API Keys</span>
                </Space>
            }
            extra={
                <Text type="secondary">
                    {apiKeyData.payload.length} key{apiKeyData.payload.length !== 1 ? 's' : ''}
                </Text>
            }
        >
            <Table
                columns={columns}
                dataSource={apiKeyData.payload}
                rowKey="id"
                pagination={{
                    pageSize: 10,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: (total, range) =>
                        `${range[0]}-${range[1]} of ${total} items`,
                }}
                scroll={{ x: 800 }}
                size="middle"
            />
        </Card>
    );
};

export default APIKeysTableComponent