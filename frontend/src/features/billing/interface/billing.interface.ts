import { APIResponse } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';

export type PaymentMethod = "Paystack" | "FlutterWave" | "None"

export type PaymentValue = {
    plan_id: string
    amount_to_pay: number
    duration: string
    payment_method: PaymentMethod
}

export interface Company {
    companyname: string;
    companycreatedat: string;
    companyupdatedat: string;
}

export interface User {
    userfullname: string;
    useremail: string;
    userphonenumber: string | null;
    userpicture: string | null;
    userverified: boolean;
    userblocked: boolean;
    userstatus: string;
    userlastloginat: string | null;
    usercreatedat: string;
}

export interface Plan {
    plan_name: string;
}

export interface Subscription {
    subscriptionplanid: string;
    subscriptionamount: string;
    subscriptionbillingcycle: string;
    subscriptiontrialstartsat: string | null;
    subscriptiontrialendsat: string | null;
    subscriptionstartsat: string | null;
    subscriptionendsat: string | null;
    subscriptionstatus: string;
    subscriptioncreatedat: string;
    plan: Plan;
}

export interface BillingData {
    id: string;
    company_id: string;
    user_id: string;
    subscription_id: string;
    payment_id: string | null;
    amount: string;
    currency: string;
    payment_method: string;
    status: string;
    notes: string;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
    company: Company;
    subscription: Subscription;
    user: User;
}

export type InitializeData = {
    data: {
        access_code: string
        authorization_url: string
        reference: string
    }
}

export type BillingAPIResponse = APIResponse<PaginatedResponse<BillingData>>;