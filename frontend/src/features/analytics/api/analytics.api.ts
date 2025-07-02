import axiosInstance from "../../../utils/api"
import { APIResponse } from '../../../../../frontend/src/interface/api.interface';
import { CampaignUserStats, EmailCampaignStats } from '../interface/analytics.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';

class AnalyticsAPI {
    static async getAllCampaignStats(page?: number, pageSize?: number): Promise<APIResponse<PaginatedResponse<EmailCampaignStats[]>>> {
        const response = await axiosInstance.get<APIResponse<PaginatedResponse<EmailCampaignStats[]>>>("/campaigns/all-campaign-stats", {
            params: {
                page,
                page_size: pageSize
            }
        })
        return response.data
    }

    static async getCampaignUserStats(): Promise<APIResponse<CampaignUserStats>> {
        const response = await axiosInstance.get<APIResponse<CampaignUserStats>>("/campaigns/user-campaign-stat")
        return response.data
    }
}

export default AnalyticsAPI

