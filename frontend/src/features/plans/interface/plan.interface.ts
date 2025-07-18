import { BaseEntity } from '../../../interface/baseentity.interface';

type MailingLimits = {
    daily_limit: number;
    monthly_limit: number;
    max_recipients_per_mail: number;
}

type Feature = {
    name: string;
    value: string
    description: string;
}

type BasePlan = {
    name: string
    description: string
    price: number
    billing_cycle: string
    status: string
}



export type EditPlanValues = Partial<PlanValues> & {
    id: string;
}

export interface PlanData extends BaseEntity, BasePlan {
    features: Feature[];
    mailing_limits: MailingLimits
}



export interface PlanValues {
    plan_name: string;
    description: string;
    price: number;
    billing_cycle: string;
    status: string;
    features: Feature[];
    mailing_limits: MailingLimits;
}