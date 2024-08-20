import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { errResponse } from '../../utils/isError';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { BaseEntity } from '../../interface/baseentity.interface';

type Campaign = {
    name: string;
    subject?: string;
    preview_text?: string;
    sender_id?: string;
    user_id: string;
    sender_from_name?: string;
    template_id?: string;
    sent_template_id?: string;
    recipient_info?: string;
    is_published: boolean;
    status: string;
    track_type: string;
    is_archived: boolean;
    sent_at?: Date;
    created_by: string;
    last_edited_by: string;
    template?: string;
    sender?: string;
};

type CreateCampaignValues = Pick<Campaign, "name">

type CampaignResponse = APIResponse<PaginatedResponse<Campaign & BaseEntity>>

type CampaignStore = {
    createCampaignValues: CreateCampaignValues
    setCreateCampaignValues: (newData: CreateCampaignValues) => void
    campaignData: (Campaign & BaseEntity)[]
    createCampaign: () => Promise<void>
    getAllCampaigns: () => Promise<void>
    setCampaignData: (newData: (Campaign & BaseEntity)[]) => void
}

const useCampaignStore = create<CampaignStore>((set, get) => ({
    createCampaignValues: { name: "" },
    campaignData: [],
    setCampaignData: (newData) => set({ campaignData: newData }),
    setCreateCampaignValues: (newData) => set({ createCampaignValues: newData }),

    createCampaign: async () => {
        try {
            const { createCampaignValues } = get()
            let response = await axiosInstance.post<ResponseT>("/campaigns/create-campaign", createCampaignValues)
            console.log(response.data)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },

    getAllCampaigns: async () => {
        try {
            const { campaignData } = get()
            let response = await axiosInstance.get<CampaignResponse>("/campaign")
            const { data, ...paginationInfo } = response.data.payload;
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    }
}))


export default useCampaignStore
