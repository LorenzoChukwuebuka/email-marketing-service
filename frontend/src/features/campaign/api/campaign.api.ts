import { APIResponse, ResponseT } from '../../../interface/api.interface';
import axiosInstance from "../../../utils/api";
import { Campaign, CampaignGroupValues, CampaignResponse, CreateCampaignValues, CampaignStats, CampaignEmailRecipientStats } from '../interface/campaign.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';

class CampaignApi {
    async createCampaign(createCampaignValues: Partial<CreateCampaignValues>): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>("/campaigns/create-campaign", createCampaignValues);
        return response.data;
    }

    async getAllCampaigns(page?: number, pageSize?: number, query?: string): Promise<APIResponse<PaginatedResponse<CampaignResponse[]>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<CampaignResponse[]>>>("/campaigns/get-all-campaigns", {
            params: { page, page_size: pageSize, search: query },
        });
        return response.data;
    }

    async getSingleCampaign(uuid: string): Promise<APIResponse<CampaignResponse>> {
        const response = await axiosInstance.get<APIResponse<CampaignResponse>>(`/campaigns/get-campaign/${uuid}`);
        return response.data;
    }

    async updateCampaign(uuid: string, updatePayload: Partial<Campaign>): Promise<ResponseT> {
        const response = await axiosInstance.put<ResponseT>(`/campaigns/update-campaign/${uuid}`, updatePayload);
        return response.data
    }

    async createCampaignGroup(values: CampaignGroupValues): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>("/campaigns/add-campaign-group", values);
        return response.data
    }

    async getScheduledCampaigns(page?: number, pageSize?: number, searchQuery?: string): Promise<APIResponse<PaginatedResponse<CampaignResponse>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<CampaignResponse>>>(`/campaigns/get-scheduled-campaigns`, {

            params: { page, page_size: pageSize, search: searchQuery || undefined },

        });
        return response.data;
    }

    async deleteCampaign(campaignId: string): Promise<ResponseT> {
        const response = await axiosInstance.delete<ResponseT>("/campaigns/delete-campaign/" + campaignId)
        return response.data
    }

    async sendCampaign(value: Record<string, string>): Promise<ResponseT> {
        const response = await axiosInstance.post<ResponseT>(`/campaigns/send-campaign`, value);
        return response.data

    }

    async getCampaignStats(id: string):Promise<APIResponse<CampaignStats>> {
        const response = await axiosInstance.get<APIResponse<CampaignStats>>("/campaigns/get-stats/" + id)
        return response.data
    }

    async getCampaignRecipients(id: string):Promise<APIResponse<CampaignEmailRecipientStats[]>> {
        const response = await axiosInstance.get<APIResponse<CampaignEmailRecipientStats[]>>("/campaigns/get-email-recipients/" + id)
        return response.data
    }

}



export default new CampaignApi();
