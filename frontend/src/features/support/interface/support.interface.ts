import { AdminData } from '../../auth/interface/admin.auth.interface';
import { UserDetails } from '../../../../../frontend/src/store/userstore/AuthStore';
import { BaseEntity } from '../../../../../frontend/src/interface/baseentity.interface';

export type TicketFile = {
    file_name: string
    file_path: string
} & BaseEntity


export type TicketMessage = {
    user_id: string
    message: string
    is_admin: boolean
    user: Partial<UserDetails>
    admin: Partial<AdminData>
    files: TicketFile[]
} & BaseEntity

export type Ticket = {
    user_id: string;
    name: string;
    email: string;
    subject: string;
    description: string;
    status: string;
    ticket_number: string
    priority: string;
    last_reply: string;
    messages: TicketMessage[];
} & BaseEntity;

export type SupportRequestValues = Partial<Ticket>