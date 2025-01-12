import { BaseEntity } from '../../../interface/baseentity.interface';
export type DomainRecord = {
    user_id: string;
    domain: string;
    txt_record: string;
    dmarc_record: string;
    dkim_selector: string;
    dkim_public_key: string;
    spf_record: string;
    mx_record: string
    verified: boolean;
} & BaseEntity;


export type DomainFormValues = Pick<DomainRecord, "domain">;
