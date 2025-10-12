import { APIResponse, ResponseT } from '../../../interface/api.interface';
import axiosInstance from "../../../utils/api";
import { Campaign, CampaignGroupValues, CampaignResponse, CreateCampaignValues, CampaignStats, CampaignEmailRecipientStats } from '../interface/campaign.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';

class CampaignApi {
    private baseUrl: string = "/campaigns";
    async createCampaign(createCampaignValues: Partial<CreateCampaignValues>): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>(`${this.baseUrl}/create`, createCampaignValues);
        return response.data;
    }

    async getAllCampaigns(page?: number, pageSize?: number, query?: string): Promise<APIResponse<PaginatedResponse<CampaignResponse>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<CampaignResponse>>>(`${this.baseUrl}/get`, {
            params: { page, page_size: pageSize, search: query },
        });
        return response.data;
    }

    async getSingleCampaign(uuid: string): Promise<APIResponse<CampaignResponse>> {
        const response = await axiosInstance.get<APIResponse<CampaignResponse>>(`${this.baseUrl}/get/${uuid}`);
        return response.data;
    }

    async updateCampaign(uuid: string, updatePayload: Partial<Campaign>): Promise<ResponseT> {
        const response = await axiosInstance.put<ResponseT>(`${this.baseUrl}/update/${uuid}`, updatePayload);
        return response.data
    }

    async createCampaignGroup(values: CampaignGroupValues): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>(`${this.baseUrl}/add-campaign-group`, values);
        return response.data
    }

    async getScheduledCampaigns(page?: number, pageSize?: number, searchQuery?: string): Promise<APIResponse<PaginatedResponse<CampaignResponse>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<CampaignResponse>>>(`${this.baseUrl}/scheduled`, {
            params: { page, page_size: pageSize, search: searchQuery || undefined },

        });
        return response.data;
    }

    async deleteCampaign(campaignId: string): Promise<ResponseT> {
        const response = await axiosInstance.delete<ResponseT>(`${this.baseUrl}/delete/${campaignId}`)
        return response.data
    }

    async sendCampaign(value: Record<string, string>): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>(`${this.baseUrl}/send`, value);
        return response.data

    }

    async getCampaignStats(id: string): Promise<APIResponse<CampaignStats>> {
        const response = await axiosInstance.get<APIResponse<CampaignStats>>(`${this.baseUrl}/campaign-stat/${id}`)
        return response.data
    }

    async getCampaignRecipients(id: string): Promise<APIResponse<CampaignEmailRecipientStats[]>> {
        const response = await axiosInstance.get<APIResponse<CampaignEmailRecipientStats[]>>(`${this.baseUrl}/get-email-recipients/${id}`)
        return response.data
    }

}



export default new CampaignApi();
