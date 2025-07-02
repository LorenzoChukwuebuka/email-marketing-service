import { BaseEntity } from '../../../interface/baseentity.interface';


type APIFeature = {
    name: string;
    description: string;
    value: string;
}


type MailingLimits = {
    daily_limit: number;
    monthly_limit: number;
    max_recipients_per_mail: number;
}

// API Plan structure (matches your API response)
export type APIPlanData = {
    id: string;
    name: string;
    description: string;
    price: number;
    billing_cycle: string;
    status: string;
    features: APIFeature[];
    mailing_limits: MailingLimits;
    created_at: string;
    updated_at: string;
}


type Feature = BaseEntity & {
    plan_Id: number;
    name: string;
    identifier: string;
    count_limit: number;
    size_limit: number;
    is_active: boolean;
    description: string;
}

type BasePlan = {
    planname: string;
    duration: string;
    price: number;
    details: string;
    number_of_mails_per_day: string;
    status: string;
}

export type PlanValues = Omit<BasePlan, "ID"> & {
    features: Omit<Feature, keyof BaseEntity>[];
}

export type EditPlanValues = Partial<PlanValues> & {
    uuid: string;
}

export interface PlanData extends BaseEntity, BasePlan {
    features: Feature[];
}