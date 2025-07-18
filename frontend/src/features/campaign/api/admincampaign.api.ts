import axiosInstance from "../../../utils/api"
import { APIResponse } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';
import { CampaignData, CampaignEmailRecipientStats } from '../interface/campaign.interface';

class AdminCampaignAPI {
    async getAllUserCampaign(userId: string, companyId: string, page?: number, pageSize?: number, query?: string): Promise<APIResponse<PaginatedResponse<CampaignData>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<CampaignData>>>("/admin/campaigns/get/" + userId + "/" + companyId, {
            params: {
                page: page || undefined,
                page_size: pageSize || undefined,
                search: query || undefined
            }
        })

        return response.data
    }

    async getSingleUserCampaign(campaignId: string, userId: string, companyId: string) {
        const response = await axiosInstance.get<APIResponse<CampaignData>>("/admin/campaigns/get/single/" + userId + "/" + companyId + "/" + campaignId)
        return response.data
    }

    async getCampaignRecipients(id: string, companyId: string): Promise<APIResponse<CampaignEmailRecipientStats[]>> {
        const response = await axiosInstance.get<APIResponse<CampaignEmailRecipientStats[]>>("admin/campaigns/get-campaign-recipients/" + id + "/" + companyId)
        return response.data
    }
}



export default new AdminCampaignAPI()