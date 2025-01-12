import { BaseEntity } from '../../../interface/baseentity.interface';
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

