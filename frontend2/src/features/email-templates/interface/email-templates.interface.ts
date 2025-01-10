export type Template = {
    user_id: string;
    template_name: string;
    campaignId: number | null;
    sender_name: string | null;
    from_email: string | null;
    subject: string | null;
    type: string;
    email_html: string;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    email_design: any;
    is_editable: boolean;
    is_published: boolean;
    is_public_template: boolean;
    is_gallery_template: boolean;
    tags: string;
    description: string | null;
    image_Url: string | null;
    is_active: boolean;
    editor_type: string | null;
};


export type SendTestMailValues = {
    template_id: string
    email_address: string
    subject: string
}