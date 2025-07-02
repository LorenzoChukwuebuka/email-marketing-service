import axiosInstance from "../../../utils/api"
import { APIResponse, ResponseT } from '../../../../../frontend/src/interface/api.interface';
import { PaginatedResponse } from '../../../../../frontend/src/interface/pagination.interface';
import { Sender, SenderFormValues, VerifySender } from '../interface/sender.interface';

class SenderAPI {

    private baseUrl = "/senders"

    async getAllSenders(page?: number, pageSize?: number, query?: string): Promise<APIResponse<PaginatedResponse<Sender>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<Sender>>>(`${this.baseUrl}/get`, {
            params: {
                page: page || undefined,
                page_size: pageSize || undefined,
                search: query || undefined
            }
        })
        return response.data
    }

    async createSender(formValues: SenderFormValues): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>(`${this.baseUrl}/create`, formValues)
        return response.data
    }

    async updateSender(id: string, data: Partial<Sender>): Promise<ResponseT> {
        const response = await axiosInstance.put<ResponseT>(`${this.baseUrl}/update/${id}`, data)
        return response.data
    }

    async deleteSender(id: string): Promise<ResponseT> {
        const response = await axiosInstance.delete<ResponseT>(`${this.baseUrl}/delete/${id}`)
        return response.data
    }

    async verifySender(values: VerifySender): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>(`${this.baseUrl}/verify`, values)
        return response.data
    }
}


export default new SenderAPI