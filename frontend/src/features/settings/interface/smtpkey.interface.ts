import { BaseEntity } from '../../../../../frontend/src/interface/baseentity.interface';
export interface SMTPKeyFormValues {
    key_name: string;
}



export type SMTPKey = BaseEntity & {
    user_id: string;
    key_name: string;
    smtp_login: string;
    password: string;
    status: string;
}

export interface SMTPKeyDATA {
    keys: SMTPKey[];
    smtp_login: string;
    smtp_master: string;
    smtp_master_password: string;
    smtp_master_status: string;
    smtp_port: string;
    smtp_server: string;
    smtp_created_at: string
}