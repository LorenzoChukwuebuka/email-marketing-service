import { useParams } from 'react-router-dom';
import {
    Button,
    Card,
    Statistic,
    Row,
    Col,
    Typography,
    Space,
    Spin,
    Divider,
    Layout
} from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import CampaignRecipientComponent from './campaignRecipientComponent';
import CampaignInfo from './campaignInfoComponent';
import { useCampaignStatsQuery, useSingleCampaignQuery } from '../../hooks/useCampaignQuery';
import { CampaignData } from '../../interface/campaign.interface';

const { Title } = Typography;
const { Content } = Layout;

interface StatItemProps {
    value: string | number;
    label: string;
    suffix?: string;
}

const StatItem: React.FC<StatItemProps> = ({ value, label, suffix }) => (
    <Card size="small" className="h-full">
        <Statistic
            value={value}
            title={label}
            suffix={suffix}
            valueStyle={{ fontSize: '1.5rem', fontWeight: 600 }}
        />
    </Card>
);

const CampaignReport: React.FC = () => {
    const { id } = useParams<{ id: string }>() as { id: string };
    const { data: campaignStatData, isLoading } = useCampaignStatsQuery(id);
    const { data: campaignData } = useSingleCampaignQuery(id);

    const stats = [
        {
            value: campaignStatData?.payload.total_emails_sent ?? 0,
            label: 'Emails Sent'
        },
        {
            value: campaignStatData?.payload.total_deliveries ?? 0,
            label: 'Delivered'
        },
        {
            value: campaignStatData?.payload.total_bounces ?? 0,
            label: 'Bounced'
        },
        {
            value: 0,
            label: 'Complaints'
        },
        {
            value: campaignStatData?.payload.hard_bounces ?? 0,
            label: 'Rejected'
        },
        {
            value: campaignStatData?.payload.total_opens ?? 0,
            label: 'Opens'
        },
        {
            value: campaignStatData?.payload.unique_opens ?? 0,
            label: 'Unique Opens'
        },
        {
            value: campaignStatData?.payload.open_rate ?? 0,
            label: 'Open Rate',
            suffix: '%'
        },
        {
            value: campaignStatData?.payload.total_clicks ?? 0,
            label: 'Total Clicks'
        },
        {
            value: campaignStatData?.payload.unique_clicks ?? 0,
            label: 'Unique Clicks'
        },
    ];

    const handleGoBack = () => {
        window.history.back();
    };

    return (
        <Layout>
            <Content style={{ padding: '24px', minHeight: '100vh', backgroundColor: '#f5f5f5' }}>
                {isLoading ? (
                    <div style={{
                        display: 'flex',
                        justifyContent: 'center',
                        alignItems: 'center',
                        minHeight: '400px'
                    }}>
                        <Spin size="large" />
                    </div>
                ) : (
                    <Space direction="vertical" size="large" style={{ width: '100%' }}>
                        {/* Header Section */}
                        <Space align="center" size="middle">
                            <Button
                                type="text"
                                icon={<ArrowLeftOutlined />}
                                onClick={handleGoBack}
                                size="large"
                            />
                            <Title level={2} style={{ margin: 0 }}>
                                {campaignData?.payload?.name}
                            </Title>
                        </Space>

                        {/* Statistics Grid */}
                        <Card title="Campaign Statistics" style={{ borderRadius: '8px' }}>
                            <Row gutter={[16, 16]}>
                                {stats.map((stat, index) => (
                                    <Col
                                        key={index}
                                        xs={12}
                                        sm={8}
                                        md={6}
                                        lg={4}
                                        xl={4}
                                    >
                                        <StatItem
                                            value={stat.value}
                                            label={stat.label}
                                            suffix={stat.suffix}
                                        />
                                    </Col>
                                ))}
                            </Row>
                        </Card>

                        <Divider />

                        {/* Campaign Info Section */}
                        <CampaignInfo campaignData={campaignData?.payload as CampaignData} />

                        {/* Campaign Recipients Section */}
                        <CampaignRecipientComponent
                            campaignId={campaignData?.payload?.id as string}
                        />
                    </Space>
                )}
            </Content>
        </Layout>
    );
};

export default CampaignReport;