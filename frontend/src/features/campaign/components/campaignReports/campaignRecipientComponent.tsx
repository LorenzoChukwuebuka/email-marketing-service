import { useMemo } from "react"
import { Table, Typography } from 'antd';
import { parseDate } from "../../../../utils/utils"
import { useCampaignRecipientsQuery } from "../../hooks/useCampaignQuery"
import { CampaignEmailRecipientStats } from "../../interface/campaign.interface";

const { Title } = Typography;

type Props = { campaignId: string }

const CampaignRecipientComponent: React.FC<Props> = ({ campaignId }) => {
    const { data: campaignRecipientData, isLoading } = useCampaignRecipientsQuery(campaignId)

    const cRdata = useMemo(() => campaignRecipientData?.payload || [], [campaignRecipientData])

    const columns = [
        {
            title: 'Recipient Email',
            dataIndex: 'recipient_email',
            key: 'recipient_email',
            width: 200,
        },
        {
            title: 'Sent At',
            dataIndex: 'sent_at',
            key: 'sent_at',
            width: 180,
            render: (date: string) => {
                return new Date(date).toLocaleDateString('en-GB', {
                    day: '2-digit',
                    month: '2-digit',
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                    second: '2-digit',
                });
            },
        },
        {
            title: 'Opened At',
            dataIndex: 'opened_at',
            key: 'opened_at',
            width: 180,
            render: (date: string) => {
                if (!date) return 'N/A';
                return new Date(date).toLocaleDateString('en-GB', {
                    day: '2-digit',
                    month: '2-digit',
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                    second: '2-digit',
                });
            },
        },
        {
            title: 'Open Count',
            dataIndex: 'open_count',
            key: 'open_count',
            width: 120,
            align: 'center' as const,
        },
        {
            title: 'Clicked At',
            dataIndex: 'clicked_at',
            key: 'clicked_at',
            width: 200,
            render: (date: string | null) => {
                if (!date) return 'N/A';
                return parseDate(date).toLocaleString('en-US', {
                    timeZone: 'UTC',
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                    hour: 'numeric',
                    minute: 'numeric',
                    second: 'numeric'
                });
            },
        },
        {
            title: 'Click Count',
            dataIndex: 'click_count',
            key: 'click_count',
            width: 120,
            align: 'center' as const,
        },
    ];

    return (
        <div className="mt-8">
            <Title level={4} className="mb-4">
                Campaign Recipients
            </Title>
            
            <Table<CampaignEmailRecipientStats>
                columns={columns}
                dataSource={cRdata}
                rowKey="id"
                loading={isLoading}
                pagination={{
                    pageSize: 10,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: (total, range) =>
                        `${range[0]}-${range[1]} of ${total} recipients`,
                    pageSizeOptions: ['10', '20', '50', '100'],
                }}
                scroll={{ x: 1000 }}
                size="middle"
                className="bg-white rounded-lg shadow-sm"
                locale={{
                    emptyText: 'No data available'
                }}
            />
        </div>
    );
}

export default CampaignRecipientComponent