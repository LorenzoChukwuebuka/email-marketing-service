import axiosInstance from "../../../utils/api";
import { DomainFormValues, DomainRecord } from "../interface/domain.interface";
import { APIResponse, ResponseT } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';

class DomainAPI {
    private static baseUrl = "/domains"
    static async createDomain(domainValues: DomainFormValues): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>(`${this.baseUrl}/create`, domainValues)
        return response.data
    }

    static async deleteDomain(uuid: string): Promise<ResponseT> {
        const response = await axiosInstance.delete<ResponseT>(`${this.baseUrl}/delete/${uuid}`)
        return response.data
    }

    static async authenticateDomain(uuid: string): Promise<ResponseT> {
        const response = await axiosInstance.put<ResponseT>(`${this.baseUrl}/verify/${uuid}`)
        return response.data
    }

    static async getDomain(uuid: string): Promise<APIResponse<DomainRecord>> {
        const response = await axiosInstance.get<APIResponse<DomainRecord>>(`${this.baseUrl}/get/${uuid}`)
        return response.data
    }

    static async getAllDomains(page?: number, pageSize?: number, query?: string): Promise<APIResponse<PaginatedResponse<DomainRecord[]>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<DomainRecord[]>>>(`${this.baseUrl}/get`, {
            params: {
                page: page,
                page_size: pageSize,
                search: query
            }
        })

        return response.data
    }
}

export default DomainAPI