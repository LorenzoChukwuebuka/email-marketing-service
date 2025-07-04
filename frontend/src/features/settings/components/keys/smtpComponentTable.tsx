import { useState, useMemo } from 'react';
import { Table, Button, Modal, Card, Row, Col, Typography, Space, Tag, Tooltip, message } from 'antd';
import { CopyOutlined, DeleteOutlined, ReloadOutlined, CheckOutlined } from '@ant-design/icons';
import { convertToNormalTime, copyToClipboard, maskAPIKey } from "../../../../utils/utils";
import useSMTPKeyStore from "../../store/smtpkey.store";
import { useSMTPKeyQuery } from "../../hooks/useSmtpkeyQuery";

const { Title, Text } = Typography;

interface SMTPKey {
    id: string;
    key_name: string;
    password: string;
    status: string;
    created_at: string;
    smtp_login: string;
    company_id: string;
    user_id: string;
    updated_at: string;
    deleted_at: string | null;
}

const SMTPKeysTableComponent: React.FC = () => {
    const { deleteSMTPKey, generateSMTPKey } = useSMTPKeyStore();
    const [deletingId, setDeletingId] = useState<string | null>(null);
    const [copyingKey, setCopyingKey] = useState<string | null>(null);
    const [regenerating, setRegenerating] = useState(false);

    const { data: smtpKeyData,refetch } = useSMTPKeyQuery();

    const smkdata = useMemo(() => smtpKeyData?.payload, [smtpKeyData]);

    const handleDelete = async (id: string) => {
        Modal.confirm({
            title: "Delete SMTP Key",
            content: "Are you sure you want to delete this SMTP key? This action cannot be undone.",
            okText: "Yes, Delete",
            cancelText: "Cancel",
            okType: "danger",
            onOk: async () => {
                setDeletingId(id);
                try {
                    await deleteSMTPKey(id);
                    message.success('SMTP key deleted successfully');
                    await new Promise(resolve => setTimeout(resolve, 1000));
                    refetch()
                } catch (error) {
                    message.error('Failed to delete SMTP key');
                } finally {
                    setDeletingId(null);
                }
            },
        });
    };

    const handleCopy = async (key: string, keyName: string) => {
        try {
            await copyToClipboard(key);
            setCopyingKey(key);
            message.success(`${keyName} copied to clipboard`);
            setTimeout(() => {
                setCopyingKey(null);
            }, 2000);
        } catch (error) {
            message.error('Failed to copy to clipboard');
        }
    };

    const handleRegenerate = async () => {
        setRegenerating(true);
        try {
            await generateSMTPKey();
            message.success('SMTP credentials regenerated successfully');
            await new Promise(resolve => setTimeout(resolve, 1000));
        } catch (error) {
            message.error('Failed to regenerate SMTP credentials');
        } finally {
            setRegenerating(false);
        }
    };

    // Safe date formatting function
    const formatDate = (dateString: string | undefined) => {
        if (!dateString) return 'N/A';
        try {
            return convertToNormalTime(dateString);
        } catch (error) {
            // Fallback to basic date formatting if convertToNormalTime fails
            return new Date(dateString).toLocaleDateString('en-US', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: 'numeric',
                minute: 'numeric'
            });
        }
    };

    const columns = [
        {
            title: 'Key Name',
            dataIndex: 'key_name',
            key: 'key_name',
            render: (text: string) => <Text strong>{text}</Text>,
        },
        {
            title: 'Key Value',
            dataIndex: 'password',
            key: 'password',
            render: (password: string) => (
                <Space>
                    <Text code>{maskAPIKey(password)}</Text>
                    <Tooltip title="Copy to clipboard">
                        <Button
                            type="text"
                            size="small"
                            icon={copyingKey === password ? <CheckOutlined /> : <CopyOutlined />}
                            onClick={() => handleCopy(password, 'Key')}
                            style={{ color: copyingKey === password ? '#52c41a' : '#1890ff' }}
                        />
                    </Tooltip>
                </Space>
            ),
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: (status: string) => (
                <Tag color={status === 'active' ? 'green' : 'red'}>
                    {status?.toUpperCase()}
                </Tag>
            ),
        },
        {
            title: 'Created On',
            dataIndex: 'created_at',
            key: 'created_at',
            render: (date: string) => formatDate(date),
        },
        {
            title: 'Actions',
            key: 'actions',
            render: (record: SMTPKey) => (
                <Button
                    type="text"
                    danger
                    icon={<DeleteOutlined />}
                    onClick={() => handleDelete(record.id)}
                    loading={deletingId === record.id}
                    disabled={deletingId === record.id}
                >
                    Delete
                </Button>
            ),
        },
    ];

    // Prepare data for the table
    const tableData = useMemo(() => {
        const data: any[] = [];

        // Add master key as first row
        if (smkdata?.smtp_master) {
            data.push({
                key: 'master',
                id: 'master',
                key_name: smkdata.smtp_master,
                password: smkdata.smtp_master_password,
                status: smkdata.smtp_master_status,
                created_at: smkdata.smtp_created_at,
                isMaster: true,
            });
        }

        // Add regular keys
        if (Array.isArray(smkdata?.keys) && smkdata.keys.length > 0) {
            smkdata.keys.forEach((key: any) => {
                data.push({
                    key: key.id,
                    ...key,
                    isMaster: false,
                });
            });
        }

        return data;
    }, [smkdata]);

    // Custom columns for master key (no delete action)
    const masterColumns = columns.filter(col => col.key !== 'actions');

    return (
        <div className="max-w-6xl mx-auto p-6">
            <Title level={2}>SMTP Configuration</Title>

            {/* SMTP Settings Card */}
            <Card
                title="SMTP Settings"
                style={{ marginBottom: 24 }}
                extra={
                    <Button
                        type="primary"
                        icon={<ReloadOutlined />}
                        onClick={handleRegenerate}
                        loading={regenerating}
                    >
                        Regenerate Credentials
                    </Button>
                }
            >
                <Row gutter={[16, 16]}>
                    <Col xs={24} sm={8}>
                        <Space direction="vertical" size={0}>
                            <Text type="secondary">SMTP Server</Text>
                            <Text strong>{smkdata?.smtp_server || 'N/A'}</Text>
                        </Space>
                    </Col>
                    <Col xs={24} sm={8}>
                        <Space direction="vertical" size={0}>
                            <Text type="secondary">Port</Text>
                            <Text strong>{smkdata?.smtp_port || 'N/A'}</Text>
                        </Space>
                    </Col>
                    <Col xs={24} sm={8}>
                        <Space direction="vertical" size={0}>
                            <Text type="secondary">Login</Text>
                            <Text strong>{smkdata?.smtp_login || 'N/A'}</Text>
                        </Space>
                    </Col>
                </Row>
            </Card>

            {/* Master Key Table */}
            {smkdata?.smtp_master && (
                <Card title="Master SMTP Key" style={{ marginBottom: 24 }}>
                    <Table
                        columns={[
                            ...masterColumns,
                            {
                                title: 'Key Value',
                                dataIndex: 'password',
                                key: 'password',
                                render: (password: string) => (
                                    <Space>
                                        <Text code>{maskAPIKey(password)}</Text>
                                        <Tooltip title="Copy master password">
                                            <Button
                                                type="text"
                                                size="small"
                                                icon={copyingKey === password ? <CheckOutlined /> : <CopyOutlined />}
                                                onClick={() => handleCopy(password, 'Master password')}
                                                style={{ color: copyingKey === password ? '#52c41a' : '#1890ff' }}
                                            />
                                        </Tooltip>
                                    </Space>
                                ),
                            }
                        ]}
                        dataSource={tableData.filter(item => item.isMaster)}
                        pagination={false}
                        size="small"
                    />
                </Card>
            )}

            {/* Regular Keys Table */}
            <Card title="SMTP Keys">
                <Table
                    columns={columns}
                    dataSource={tableData.filter(item => !item.isMaster)}
                    pagination={{
                        pageSize: 10,
                        showSizeChanger: true,
                        showQuickJumper: true,
                        showTotal: (total, range) =>
                            `${range[0]}-${range[1]} of ${total} items`,
                    }}
                    loading={!smkdata}
                    locale={{
                        emptyText: 'No SMTP keys found. Generate your first key to get started.',
                    }}
                    scroll={{ x: 800 }}
                />
            </Card>
        </div>
    );
};

export default SMTPKeysTableComponent;