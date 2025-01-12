import axiosInstance from "../../../utils/api"
import { APIResponse } from '../../../../../frontend/src/interface/api.interface';
import { CampaignUserStats, EmailCampaignStats } from '../interface/analytics.interface';

class AnalyticsAPI {
    static async getAllCampaignStats(): Promise<APIResponse<EmailCampaignStats[]>> {
        const response = await axiosInstance.get<APIResponse<EmailCampaignStats[]>>("/campaigns/user-campaigns-stats")
        return response.data
    }

    static async getCampaignUserStats(): Promise<APIResponse<CampaignUserStats>> {
        const response = await axiosInstance.get<APIResponse<CampaignUserStats>>("/campaigns/user-campaign-stats")
        return response.data
    }
}

export default AnalyticsAPI

