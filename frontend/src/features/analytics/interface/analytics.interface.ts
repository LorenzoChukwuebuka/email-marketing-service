 export type EmailCampaignStats = {
    bounces: number;
    campaign_id: string;
    clicked: number;
    complaints: number;
    name: string;
    opened: number;
    recipients: number;
    sent_date: string | null;
    unsubscribed: number;
};
 

export type CampaignUserStats = {
    hard_bounces: number;
    open_rate: number;
    soft_bounces: number;
    total_bounces: number;
    total_clicks: number;
    total_deliveries: number;
    total_emails_sent: number;
    total_opens: number;
    unique_clicks: number;
    unique_opens: number;
};