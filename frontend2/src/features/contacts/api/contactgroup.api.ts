import axiosInstance from '../../../utils/api';
import { APIResponse, ResponseT } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';
import { AddToGroup, ContactGroupData, ContactGroupFormValues } from '../interface/contactgroup.interface';


export class ContactGroupAPI {
    async getAllGroups(page?: number, pageSize?: number, query?: string): Promise<APIResponse<PaginatedResponse<ContactGroupData>>> {
        const response = await axiosInstance
            .get<APIResponse<PaginatedResponse<ContactGroupData>>>
            (`/contact/get-all-contact-groups`, {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            });

        return response.data;
    }

    async addContactToGroup(data: AddToGroup) {

        const response = await axiosInstance.post<ResponseT>("/contact/add-contact-to-group", data);
        return response.data;
    }

    async getSingleGroup(uuid: string) {
        const response = await axiosInstance.get<APIResponse<ContactGroupData>>('/contact/get-single-group/' + uuid);
        return response.data.payload;
    }

    async createGroup(formValues: ContactGroupFormValues) {
        const response = await axiosInstance.post<ResponseT>("/contact/create-contact-group", formValues);
        return response.data;
    }

    async updateGroup(uuid: string, values: ContactGroupFormValues) {
        const response = await axiosInstance.put<ResponseT>("/contact/edit-group/" + uuid, values);
        return response.data;
    }

    async deleteGroup(groupId: string): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.delete<APIResponse<ResponseT>>("/contact/delete-group/" + groupId);
        return response.data;
    }

    async removeContactFromGroup(groupId: string, contactId: string) {
        const data = {
            group_id: groupId,
            contact_id: contactId
        };

        const response = await axiosInstance.post<ResponseT>("/contact/remove-contact-from-group", data);
        return response.data;
    }


}

export const contactGroupAPI = new ContactGroupAPI();