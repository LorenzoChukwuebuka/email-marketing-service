import axiosInstance from "../../../utils/api"
import { APIResponse, ResponseT } from '../../../../../frontend/src/interface/api.interface';
import { PaginatedResponse } from '../../../../../frontend/src/interface/pagination.interface';
import { Sender, SenderFormValues, VerifySender } from '../interface/sender.interface';

class SenderAPI {

    async getAllSenders(page?: number, pageSize?: number, query?: string): Promise<APIResponse<PaginatedResponse<Sender>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<Sender>>>("/sender/get-all-senders", {
            params: {
                page: page || undefined,
                page_size: pageSize || undefined,
                search: query || undefined
            }
        })
        return response.data
    }

    async createSender(formValues: SenderFormValues): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>("/sender/create-sender", formValues)
        return response.data
    }

    async updateSender(id: string, data: Partial<Sender>): Promise<ResponseT> {
        const response = await axiosInstance.put<ResponseT>("/sender/update-sender/" + id, data)
        return response.data
    }

    async deleteSender(id: string): Promise<ResponseT> {
        const response = await axiosInstance.delete<ResponseT>("/sender/delete-sender/" + id)
        return response.data
    }

    async verifySender(values: VerifySender): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>("/sender/verify-sender", values)
        return response.data
    }
}


export default new SenderAPI