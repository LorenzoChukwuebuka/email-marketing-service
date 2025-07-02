export type ContactGroupFormValues = {
    group_name: string;
    description: string;
}

export type AddToGroup = {
    group_id: string;
    contact_id: string;
}

export type ContactInGroup = {
    contact_id: string;
    contact_first_name: string;
    contact_last_name: string;
    contact_email: string;
    contact_from_origin: string;
    contact_is_subscribed: boolean;
    contact_created_at: string;
}


export type ContactGroupData = {
    group_id: string;
    group_name: string;
    description: string;
    group_created_at: string;
    contacts: ContactInGroup[];
}

export type EditGroupValues = ContactGroupFormValues & { group_id: string; }