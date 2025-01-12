import { BaseEntity } from '../../../interface/baseentity.interface';
import { PlanData } from '../../plans/interface/plan.interface';
import { APIResponse } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';

export type PaymentMethod = "Paystack" | "FlutterWave"

export type PaymentValue = {
    plan_id: string
    amount_to_pay: number
    duration: string
    payment_method: PaymentMethod
}

export interface User extends BaseEntity {
    fullname: string;
    email: string;
    company: string;
    phonenumber: string;
    verified: boolean;
    blocked: boolean;
    verified_at: string | null;
}

export interface BillingData extends BaseEntity {
    user_id: number;
    amount_paid: number;
    plan_id: number;
    duration: string;
    expiry_date: string;
    reference: string;
    transaction_id: string;
    payment_method: string;
    status: string;
    user: User;
    plan: PlanData;
}



export type InitializeData = {
    data: {
        access_code: string
        authorization_url: string
        reference: string
    }
}



export type BillingAPIResponse = APIResponse<PaginatedResponse<BillingData>>;



