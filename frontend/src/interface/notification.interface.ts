import { BaseEntity } from '../../../frontend/src/interface/baseentity.interface';
export type UserNotification = {
    title: string
    read_status: boolean
} & BaseEntity

