import { BaseEntity } from '../../../interface/baseentity.interface';
import { Template } from '../../email-templates/interface/email-templates.interface';

 
export type Campaign = {
    name: string;
    subject?: string;
    preview_text?: string;
    sender_id?: string;
    user_id: string;
    sender_from_name?: string;
    template_id?: string;
    sent_template_id?: string;
    recipient_info?: string;
    is_published: boolean;
    status: string;
    track_type: string;
    is_archived: boolean;
    sent_at?: string;
    created_by: string;
    last_edited_by: string;
    template?: Template;
    scheduled_at?: string
    sender?: string;
    campaign_groups: CampaignGroup[]
};


export type CampaignResponse = Campaign  & BaseEntity

export type CampaignGroup = { campaign_id: number; group_id: number } & BaseEntity

export type CreateCampaignValues = Partial<Campaign>

export type CampaignGroupValues = { campaign_id: string; group_id: string }

export type CampaignData = BaseEntity & Campaign

export type CampaignStats = {
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
}

export type CampaignEmailRecipientStats = {
    campaign_id: string;
    recipient_email: string;
    version: string;
    sent_at: string;
    opened_at: string | null;
    open_count: number;
    clicked_at: string | null;
    click_count: number;
    conversion_at: string | null;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
} & BaseEntity;

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

