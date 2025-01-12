import { BaseEntity } from '../../../../../frontend/src/interface/baseentity.interface';

export interface Sender extends BaseEntity {
    name: string;
    email: string;
    verified: boolean
    is_signed: boolean
}


export interface VerifySender {
    email: string,
    token: string,
    user_id: string
}

export interface SenderFormValues {
    name: string;
    email: string;
}

