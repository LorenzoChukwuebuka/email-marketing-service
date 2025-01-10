import { APIResponse, ResponseT } from '../../../interface/api.interface';
import axiosInstance from "../../../utils/api";
import { Campaign, CampaignGroupValues, CampaignResponse } from "../interface/campaign.interface";
import eventBus from '../../../utils/eventbus';
import { errResponse } from '../../../utils/isError';
import { PaginatedResponse } from '../../../interface/pagination.interface';

class CampaignApi {
    async createCampaign(createCampaignValues: Partial<Campaign>): Promise<ResponseT> {
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

    async createCampaignGroup(uuid: string, groupIds: string[]): Promise<void> {
        try {
            const promises = groupIds.map(groupId => {
                const data: CampaignGroupValues = { group_id: groupId, campaign_id: uuid };
                return axiosInstance.post<ResponseT>("/campaigns/add-campaign-group", data);
            });
            const results = await Promise.all(promises);
            const allSuccessful = results.every(response => response.data.status);
            if (allSuccessful) {
                eventBus.emit("success", "Recipients added successfully");
            } else {
                eventBus.emit("error", "Some recipients could not be created");
            }
        } catch (error) {
            this.handleError(error);
        }
    }

    async getScheduledCampaigns(page?: number, pageSize?: number): Promise<CampaignResponse> {
            const response = await axiosInstance.get<CampaignResponse>(`/campaigns/get-scheduled-campaigns?page=${page}&page_size=${pageSize}`);
            return response.data;
    }

    async sendCampaign(campaignId: string): Promise<void> {
        try {
            const campaign = await this.getSingleCampaign(campaignId);
            if (!campaign || !campaign?.payload?.template?.email_html || !campaign?.payload?.subject || !campaign?.payload?.campaign_groups?.length) {
                throw new Error("Incomplete campaign details");
            }
            const response = await axiosInstance.post<ResponseT>(`/campaigns/send/${campaignId}`);
            if (response.data.status) {
                eventBus.emit("success", "Campaign sent successfully");
            }
        } catch (error) {
            this.handleError(error);
        }
    }

    private handleError(error: unknown): void {
        if (errResponse(error)) {
            eventBus.emit("error", error?.response?.data.message);
        } else if (error instanceof Error) {
            eventBus.emit("error", error.message);
        } else {
            console.error("Unknown error:", error);
        }
    }
}



export default new CampaignApi();
