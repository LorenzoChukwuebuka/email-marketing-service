 import { BaseEntity } from '../../../interface/baseentity.interface';

 
export type Template = {
   id: string;
   user_id: string;
   company_id: string;
   name: string;
   sender_name: string;
   from_email: string;
   subject: string;
   type: string;
   email_html: string;
   email_design: string | null;
   is_editable: boolean;
   is_published: boolean;
   is_public_template: boolean;
   is_gallery_template: boolean;
   tags: string;
   description: string;
   image_url: string;
   is_active: boolean;
   editor_type: string;
   created_at: string;
   updated_at: string;
   deleted_at: string | null;
};

 
export type CampaignUser = {
    user_id_2: string;
    user_fullname: string;
    user_email: string;
    user_phonenumber: string | null;
    user_picture: string | null;
    user_verified: boolean;
    user_blocked: boolean;
    user_verified_at: string;
    user_status: string;
    user_last_login_at: string | null;
    user_created_at: string;
    user_updated_at: string;
};

 
export type CampaignCompany = {
    company_id_ref: string;
    company_name: string;
    company_created_at: string;
    company_updated_at: string;
};

 
export type CampaignGroupInfo = {
    id: string;
    group_name: string;
    description: string;
    created_at: string;
};

export type Campaign = {
    name: string;
    subject?: string;
    preview_text?: string;
    sender_id?: string;
    user_id: string;
    sender_from_name?: string;
    template_id?: string;
    sent_template_id?: string | null;
    recipient_info?: string;
    is_published: boolean;
    status: string;
    track_type: string;
    is_archived: boolean;
    sent_at?: string;
    created_by?: string;
    last_edited_by?: string;
    template?: Template;
    scheduled_at?: string;
    sender?: string;
    has_custom_logo?: boolean;
    user?: CampaignUser;
    company?: CampaignCompany;
    groups?: CampaignGroupInfo[];
};

export type CampaignResponse = Campaign & BaseEntity;

export type CampaignGroup = { 
    campaign_id: string; 
    group_id: string; 
} & BaseEntity;

export type CreateCampaignValues = Partial<Campaign>;

export type CampaignGroupValues = { 
    campaign_id: string; 
    group_id: string; 
};

export type CampaignData = BaseEntity & Campaign;

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
};

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