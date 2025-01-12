import { BaseEntity } from '../../../interface/baseentity.interface';

export type ContactFormValues = {
    first_name: string
    last_name: string
    email: string
    from: string
    is_subscribed: boolean
}

export type Group = {
    group_name: string;
    user_id: string;
    description: string;
} & BaseEntity;



export type ContactBase = {
    user_id: string;
    groups: Group[] | null;
}

export type ContactAPIResponse = BaseEntity & ContactFormValues & ContactBase;

export type ContactCount = { recent: number; total: number }

export type EditContactValues = { uuid: string } & Partial<ContactFormValues>

export type ContactEngageCount = {
    engaged: number
    new: number
    total: number
    unsubscribed: number
}

export type FileCSVType = null | File;