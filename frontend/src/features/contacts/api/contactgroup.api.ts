import axiosInstance from '../../../utils/api';
import { APIResponse, ResponseT } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';
import { AddToGroup, ContactGroupData, ContactGroupFormValues } from '../interface/contactgroup.interface';


export class ContactGroupAPI {
    private baseurl = "/contacts";
    async getAllGroups(page?: number, pageSize?: number, query?: string): Promise<APIResponse<PaginatedResponse<ContactGroupData>>> {
        const response = await axiosInstance
            .get<APIResponse<PaginatedResponse<ContactGroupData>>>
            (`${this.baseurl}/getgroupwithcontacts`, {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            });

        return response.data;
    }

    async addContactToGroup(data: AddToGroup) {
        const response = await axiosInstance.post<ResponseT>(`${this.baseurl}/addcontacttogroup`, data);
        return response.data;
    }

    async getSingleGroup(uuid: string) {
        const response = await axiosInstance.get<APIResponse<ContactGroupData>>(`${this.baseurl}/getgroupwithcontacts/` + uuid);
        return response.data.payload;
    }

    async createGroup(formValues: ContactGroupFormValues) {
        const response = await axiosInstance.post<ResponseT>(`${this.baseurl}/creategroup`, formValues);
        return response.data;
    }

    async updateGroup(uuid: string, values: ContactGroupFormValues) {
        const response = await axiosInstance.put<ResponseT>(`${this.baseurl}/updatecontactgroup/` + uuid, values);
        return response.data;
    }

    async deleteGroup(groupId: string): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.delete<APIResponse<ResponseT>>(`${this.baseurl}/deletecontactgroup/` + groupId);
        return response.data;
    }

    async removeContactFromGroup(groupId: string, contactId: string) {
        const data = {
            group_id: groupId,
            contact_id: contactId
        };

        const response = await axiosInstance.post<ResponseT>(`${this.baseurl}/removecontactfromgroup`, data);
        return response.data;
    } 
}

export const contactGroupAPI = new ContactGroupAPI();