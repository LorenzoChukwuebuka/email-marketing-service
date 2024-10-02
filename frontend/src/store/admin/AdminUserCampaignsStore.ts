import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'
import Cookies from 'js-cookie'
import { errResponse } from '../../utils/isError';
import { Ticket } from '../userstore/support.store';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse } from '../../interface/api.interface';
import AdminUserSpecificCampaigns from '../../templates/admin/templates/campaigns/getSpecificUserCampaign';
import { CampaignData } from '../userstore/campaignStore';


type AdminUserSpecificCampaignStore = {
    campaignData: CampaignData[] | CampaignData
    paginationInfo: Omit<PaginatedResponse<CampaignData>, 'data'>;
    setCampaignData: (newData: CampaignData | CampaignData[]) => void
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<CampaignData>, 'data'>) => void;
    getAllUserCampaign: (userId: string, page?: number, pageSize?: number, search?: string) => Promise<void>
    getSingleUserCampaign: (campaignId: string) => Promise<void>
    suspendCampaign: (campaignId: string) => Promise<void>
    deleteCampaign: (campaignId: string) => Promise<void>
}


const useAdminUserCamapaignStore = create<AdminUserSpecificCampaignStore>((set, get) => ({
    campaignData: [],
    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },

    // Set campaign data in the store, can accept an array or single campaign data
    setCampaignData: (newData: CampaignData | CampaignData[]) => {
        set({ campaignData: newData });
    },
    setPaginationInfo: (newData) => set({ paginationInfo: newData }),

    // Fetch all campaigns for a specific user
    getAllUserCampaign: async (userId: string, page = 1, pageSize = 10, query = "") => {
        try {
            let response = await axiosInstance.get("/admin/campaign/user-campaigns/" + userId, {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            })
            if (response.data.status === true) {
                const { data, ...paginationInfo } = response.data.payload;
                get().setCampaignData(data);
                get().setPaginationInfo(paginationInfo);
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

    // Fetch details of a single campaign by ID
    getSingleUserCampaign: async (campaignId: string) => {
        try {

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

    // Suspend a campaign by ID
    suspendCampaign: async (campaignId: string) => {
        try {

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

    // Delete a campaign by ID
    deleteCampaign: async (campaignId: string) => {
        try {

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
}));

export default useAdminUserCamapaignStore;
