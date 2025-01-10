import { ContactAPIResponse } from './contact.interface';
import { BaseEntity } from '../../../../../frontend/src/interface/baseentity.interface';
export type ContactGroupFormValues = {
    group_name: string;
    description: string
}

export type AddToGroup = {
    group_id: string;
    contact_id: string
}

export type ContactGroupData = ContactGroupFormValues & BaseEntity & {
    userId: string;
    contacts: Omit<ContactAPIResponse, 'group'>[]
}

export type EditGroupValues = ContactGroupFormValues & { uuid: string }