import axiosInstance from "../../../utils/api"
import { APIResponse } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';
import { CampaignData } from '../interface/campaign.interface';

class AdminCampaignAPI {
    async getAllUserCampaign(userId: string, page?: number, pageSize?: number, query?: string): Promise<APIResponse<PaginatedResponse<CampaignData>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<CampaignData>>>("/admin/campaign/user-campaigns/" + userId, {
            params: {
                page: page || undefined,
                page_size: pageSize || undefined,
                search: query || undefined
            }
        })

        return response.data
    }


    async getSingleUserCampaign(campaignId: string) {
        const response = await axiosInstance.get<APIResponse<CampaignData>>("/admin/campaign/campaign/" + campaignId)
        return response.data
    }
}



export default new AdminCampaignAPI()