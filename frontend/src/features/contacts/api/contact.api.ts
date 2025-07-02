
import { ContactFormValues, ContactAPIResponse } from '../interface/contact.interface';
import { APIResponse, ResponseT } from '../../../interface/api.interface';
import axiosInstance from '../../../utils/api';
import { PaginatedResponse } from '../../../interface/pagination.interface';

class ContactAPI {
    private baseurl = "/contacts";
    async createContact(contactFormValues: ContactFormValues): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>(`${this.baseurl}/create`, contactFormValues)
        return response.data;
    }

    async deleteContact(uuid: string): Promise<ResponseT> {
        const response = await axiosInstance.delete<ResponseT>(`${this.baseurl}/delete/${uuid}`)
        return response.data;
    }

    async getContact(uuid: string): Promise<ResponseT> {
        const response = await axiosInstance.get<ResponseT>(`${this.baseurl}/get-contact/${uuid}`)
        return response.data;
    }

    async editContact(uuid: string, contactFormValues: Partial<ContactFormValues>): Promise<ResponseT> {
        const response = await axiosInstance.put<ResponseT>(`${this.baseurl}/update/${uuid}`, contactFormValues)
        return response.data;
    }

    async getContactCount(): Promise<APIResponse<Record<string, string>>> {
        const response = await axiosInstance.get<APIResponse<Record<string, string>>>(`${this.baseurl}/count`)
        return response.data
    }

    async getAllContacts(page?: number, pageSize?: number, query?: string): Promise<APIResponse<PaginatedResponse<ContactAPIResponse>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<ContactAPIResponse>>>(`${this.baseurl}/get`, {
            params: {
                page: page || undefined,
                page_size: pageSize || undefined,
                search: query || undefined
            }
        });

        return response.data;
    }

    async getContactEngagement(): Promise<APIResponse<Record<string, string>>> {
        const response = await axiosInstance.get<APIResponse<Record<string, string>>>(`${this.baseurl}/engagement`)
        return response.data
    }

    async batchUploadContact(file: FormData): Promise<APIResponse<Record<string, string>>> {
        const response = await axiosInstance.post('${this.baseurl}/upload-csv', file)
        return response.data
    }
}

export default new ContactAPI();