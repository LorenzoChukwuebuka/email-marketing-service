import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { errResponse } from '../../utils/isError';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { BaseEntity } from '../../interface/baseentity.interface';
import Cookies from 'js-cookie'

export type Campaign = {
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
    campaign_groups: CampaignGroup[] | CampaignGroup
};

type CampaignGroup = { campaign_id: number; group_id: number } & BaseEntity

type CreateCampaignValues = Partial<Campaign>

type CampaignResponse = APIResponse<PaginatedResponse<Campaign & BaseEntity>>

type CampaignStore = {
    createCampaignValues: CreateCampaignValues
    paginationInfo: Omit<PaginatedResponse<Campaign>, 'data'>;
    setCreateCampaignValues: (newData: CreateCampaignValues) => void
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<Campaign>, 'data'>) => void;
    campaignData: (Campaign & BaseEntity)[] | Campaign & BaseEntity
    createCampaign: () => Promise<void>
    getAllCampaigns: (page?: number, pageSize?: number) => Promise<void>
    setCampaignData: (newData: (Campaign & BaseEntity)[] | Campaign & BaseEntity) => void
    getSingleCampaign: (uuid: string) => Promise<void>
    updateCampaign: (uuid: string) => Promise<void>
}

const useCampaignStore = create<CampaignStore>((set, get) => ({
    createCampaignValues: {
        name: "",
        subject: "",
        preview_text: "",
        sender_id: "",
        sender_from_name: "",
        template_id: "",
        sent_template_id: "",
        recipient_info: "",
        is_published: false,
        status: "draft",
        track_type: "",
        is_archived: false,
        created_by: "",
        last_edited_by: "",
    },
    campaignData: [],
    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },

    setCampaignData: (newData) => set({ campaignData: newData }),
    setCreateCampaignValues: (newData) => set({ createCampaignValues: newData }),
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),

    createCampaign: async () => {
        try {
            const { createCampaignValues } = get()

            let cookie: any = Cookies.get("Cookies");
            let user = JSON.parse(cookie)?.details?.company;
            const updatedCampaignValues = { ...createCampaignValues, sender_from_name: user };
            let response = await axiosInstance.post<ResponseT>("/campaigns/create-campaign", updatedCampaignValues)
            if (response.data.status == true) {
                eventBus.emit("success", "Campaign created successfully")

                window.location.href = "/user/dash/campaign/edit/" + response.data.payload.templateId
            }
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            get().setCreateCampaignValues({ name: "", sender_from_name: "" })
        }
    },

    getAllCampaigns: async (page = 1, pageSize = 10) => {
        try {
            const { setCampaignData, setPaginationInfo } = get()
            let response = await axiosInstance.get<CampaignResponse>(`/campaigns/get-all-campaigns?page=${page}&page_size=${pageSize}`)
            const { data, ...paginationInfo } = response.data.payload;
            setCampaignData(data)
            setPaginationInfo(paginationInfo)
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
    getSingleCampaign: async (uuid: string) => {
        try {
            const { setCampaignData } = get()
            let response = await axiosInstance.get<APIResponse<(BaseEntity & Campaign)>>("/campaigns/get-campaign/" + uuid)
            if (response.data.status == true) {
                const campaign: Campaign & BaseEntity = response.data.payload as Campaign & BaseEntity;
                setCampaignData(campaign);
            }
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

    updateCampaign: async (uuid: string) => {

        try {
            const { createCampaignValues } = get()
            let response = await axiosInstance.put<ResponseT>("/campaigns/update-campaign/" + uuid, createCampaignValues)
            if (response.data.status === true) {
                eventBus.emit('success', "campaign updated successfully")
            }
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
