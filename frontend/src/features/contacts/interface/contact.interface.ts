import { BaseEntity, WithPrefix } from '../../../interface/baseentity.interface';
type GroupBase = WithPrefix<BaseEntity, 'group_'>;

export type ContactFormValues = {
    first_name: string
    last_name: string
    email: string
    from: string
    is_subscribed: boolean
}

export type Group = {
    group_name: string;
    group_id: string;
    group_creator_id: string;
    description: string;
} & GroupBase;


export type ContactBase = {
    user_id: string;
    groups: Group[] | null;
}

export type ContactAPIResponse = BaseEntity & Omit<ContactFormValues, 'from'> & {
    from_origin: string;
    contact_id: string
} & ContactBase;

export type ContactCount = { recent: number; total: number }

export type EditContactValues = { id: string } & Partial<ContactFormValues>

export type ContactEngageCount = {
    engaged: number
    new: number
    total: number
    unsubscribed: number
}

export type FileCSVType = null | File;