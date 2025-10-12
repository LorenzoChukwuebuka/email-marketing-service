import React, { useState } from 'react';
import {
    Card,
    Typography,
    Tag,
    Button,
    Space,
    Modal,
    Avatar,
    Row,
    Col,
    Empty
} from 'antd';
import {
    EyeOutlined,
    CalendarOutlined,
    UserOutlined,
    MailOutlined,
    TeamOutlined,
    SendOutlined,
    LayoutOutlined
} from '@ant-design/icons';
import { CampaignData } from '../../interface/campaign.interface';
const { Title, Text } = Typography;

type CampaignInfoProps = {
    campaignData: CampaignData;
};

const CampaignInfo: React.FC<CampaignInfoProps> = ({ campaignData }) => {
    const [templatePreviewVisible, setTemplatePreviewVisible] = useState(false);

    // Get status color based on campaign status
    const getStatusColor = (status: string) => {
        switch (status?.toLowerCase()) {
            case 'sent': return 'success';
            case 'draft': return 'default';
            case 'sending': return 'processing';
            case 'scheduled': return 'warning';
            case 'failed': return 'error';
            default: return 'default';
        }
    };

    // Format date for display
    const formatDate = (dateString: string | undefined) => {
        if (!dateString || dateString === '0001-01-01T00:00:00Z') return 'Not set';
        return new Date(dateString).toLocaleDateString('en-GB', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    // Get group names from the API response
    const getGroupNames = () => {
        if (campaignData?.groups && campaignData?.groups?.length > 0) {
            return campaignData?.groups?.map(group => group?.group_name);
        }
        return [];
    };

    const groupNames = getGroupNames();

    return (
        <div style={{ padding: '24px' }}>
            <Card
                title={
                    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                            <MailOutlined style={{ fontSize: '20px', color: '#1890ff' }} />
                            <Title level={3} style={{ margin: 0 }}>
                                {campaignData?.name}
                            </Title>
                        </div>
                        <Tag color={getStatusColor(campaignData?.status)} style={{ fontSize: '12px' }}>
                            {campaignData?.status?.toUpperCase()}
                        </Tag>
                    </div>
                }
                bordered={false}
                style={{
                    borderRadius: '12px',
                    boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
                }}
            >
                <Row gutter={[24, 24]}>
                    {/* Basic Information */}
                    <Col xs={24} lg={12}>
                        <Card
                            type="inner"
                            title={
                                <Space>
                                    <SendOutlined />
                                    <span>Email Details</span>
                                </Space>
                            }
                            style={{ height: '100%' }}
                        >
                            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                                <div>
                                    <Text type="secondary">Subject</Text>
                                    <div>
                                        <Text strong>{campaignData?.subject || 'No subject'}</Text>
                                    </div>
                                </div>

                                <div>
                                    <Text type="secondary">Sender</Text>
                                    <div>
                                        <Text strong>{campaignData?.sender || 'Not specified'}</Text>
                                    </div>
                                </div>

                                <div>
                                    <Text type="secondary">From Name</Text>
                                    <div>
                                        <Text strong>{campaignData?.sender_from_name || 'Not specified'}</Text>
                                    </div>
                                </div>

                                <div>
                                    <Text type="secondary">Preview Text</Text>
                                    <div>
                                        <Text>{campaignData?.preview_text || 'No preview text'}</Text>
                                    </div>
                                </div>
                            </Space>
                        </Card>
                    </Col>

                    {/* Campaign Timeline */}
                    <Col xs={24} lg={12}>
                        <Card
                            type="inner"
                            title={
                                <Space>
                                    <CalendarOutlined />
                                    <span>Timeline</span>
                                </Space>
                            }
                            style={{ height: '100%' }}
                        >
                            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                                <div>
                                    <Text type="secondary">Created At</Text>
                                    <div>
                                        <Text strong>{formatDate(campaignData?.created_at)}</Text>
                                    </div>
                                </div>

                                <div>
                                    <Text type="secondary">Sent At</Text>
                                    <div>
                                        <Text strong>{formatDate(campaignData?.sent_at)}</Text>
                                    </div>
                                </div>

                                <div>
                                    <Text type="secondary">Scheduled At</Text>
                                    <div>
                                        <Text strong>{formatDate(campaignData?.scheduled_at)}</Text>
                                    </div>
                                </div>

                                <div>
                                    <Text type="secondary">Last Updated</Text>
                                    <div>
                                        <Text strong>{formatDate(campaignData?.updated_at)}</Text>
                                    </div>
                                </div>
                            </Space>
                        </Card>
                    </Col>

                    {/* Recipients */}
                    <Col xs={24}>
                        <Card
                            type="inner"
                            title={
                                <Space>
                                    <TeamOutlined />
                                    <span>Recipients</span>
                                </Space>
                            }
                        >
                            <div>
                                <Text type="secondary">Groups/Segments</Text>
                                <div style={{ marginTop: '8px' }}>
                                    {groupNames.length > 0 ? (
                                        <Space wrap>
                                            {groupNames.map((groupName, index) => (
                                                <Tag
                                                    key={index}
                                                    color="blue"
                                                    icon={<TeamOutlined />}
                                                    style={{ marginBottom: '4px' }}
                                                >
                                                    {groupName}
                                                </Tag>
                                            ))}
                                        </Space>
                                    ) : (
                                        <Empty
                                            description="No groups selected"
                                            image={Empty.PRESENTED_IMAGE_SIMPLE}
                                            style={{ margin: '16px 0' }}
                                        />
                                    )}
                                </div>
                            </div>
                        </Card>
                    </Col>

                    {/* Template Information */}
                    <Col xs={24}>
                        <Card
                            type="inner"
                            title={
                                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                    <Space>
                                        <LayoutOutlined />
                                        <span>Template</span>
                                    </Space>
                                    <Space>
                                        <Button
                                            type="primary"
                                            icon={<EyeOutlined />}
                                            onClick={() => setTemplatePreviewVisible(true)}
                                            disabled={!campaignData?.template?.email_html}
                                        >
                                            Preview Template
                                        </Button>
                                    </Space>
                                </div>
                            }
                        >
                            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                                <div>
                                    <Text type="secondary">Template Name</Text>
                                    <div>
                                        <Text strong>
                                            {campaignData?.template?.name || 'No template selected'}
                                        </Text>
                                    </div>
                                </div>

                                {campaignData?.template && (
                                    <>
                                        <div>
                                            <Text type="secondary">Template Type</Text>
                                            <div>
                                                <Tag color="geekblue">
                                                    {campaignData?.template?.type}
                                                </Tag>
                                            </div>
                                        </div>

                                        <div>
                                            <Text type="secondary">Template Description</Text>
                                            <div>
                                                <Text>{campaignData?.template?.description}</Text>
                                            </div>
                                        </div>
                                    </>
                                )}
                            </Space>
                        </Card>
                    </Col>

                    {/* User & Company Info */}
                    {(campaignData?.user || campaignData?.company) && (
                        <Col xs={24}>
                            <Card
                                type="inner"
                                title={
                                    <Space>
                                        <UserOutlined />
                                        <span>Created By</span>
                                    </Space>
                                }
                            >
                                <Row gutter={[16, 16]}>
                                    {campaignData?.user && (
                                        <Col xs={24} sm={12}>
                                            <div>
                                                <Text type="secondary">User</Text>
                                                <div style={{ marginTop: '8px' }}>
                                                    <Space>
                                                        <Avatar icon={<UserOutlined />} />
                                                        <div>
                                                            <div>
                                                                <Text strong>{campaignData?.user?.user_fullname}</Text>
                                                            </div>
                                                            <div>
                                                                <Text type="secondary" style={{ fontSize: '12px' }}>
                                                                    {campaignData?.user?.user_email}
                                                                </Text>
                                                            </div>
                                                        </div>
                                                    </Space>
                                                </div>
                                            </div>
                                        </Col>
                                    )}

                                    {campaignData?.company && (
                                        <Col xs={24} sm={12}>
                                            <div>
                                                <Text type="secondary">Company</Text>
                                                <div style={{ marginTop: '8px' }}>
                                                    <Text strong>{campaignData?.company?.company_name}</Text>
                                                </div>
                                            </div>
                                        </Col>
                                    )}
                                </Row>
                            </Card>
                        </Col>
                    )}
                </Row>
            </Card>

            {/* Template Preview Modal */}
            <Modal
                title={
                    <Space>
                        <EyeOutlined />
                        <span>Template Preview</span>
                    </Space>
                }
                open={templatePreviewVisible}
                onCancel={() => setTemplatePreviewVisible(false)}
                footer={[
                    <Button key="close" onClick={() => setTemplatePreviewVisible(false)}>
                        Close
                    </Button>
                ]}
                width="80%"
                style={{ top: 20 }}
            >
                {campaignData?.template?.email_html ? (
                    <div
                        style={{
                            border: '1px solid #d9d9d9',
                            borderRadius: '6px',
                            padding: '16px',
                            backgroundColor: '#fff',
                            maxHeight: '60vh',
                            overflow: 'auto'
                        }}
                        dangerouslySetInnerHTML={{
                            __html: campaignData.template.email_html
                        }}
                    />
                ) : (
                    <Empty description="No template content available" />
                )}
            </Modal>
        </div>
    );
};

export default CampaignInfo;