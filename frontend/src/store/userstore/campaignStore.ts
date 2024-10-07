import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { errResponse } from '../../utils/isError';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { BaseEntity } from '../../interface/baseentity.interface';
import Cookies from 'js-cookie'
import { persist, createJSONStorage } from 'zustand/middleware';
import { Template } from './templateStore';
import { StateStorage } from 'zustand/middleware';

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
    sent_at?: string;
    created_by: string;
    last_edited_by: string;
    template?: Template;
    scheduled_at?: string
    sender?: string;
    campaign_groups: CampaignGroup[]
};

export type CampaignGroup = { campaign_id: number; group_id: number } & BaseEntity

type CreateCampaignValues = Partial<Campaign>

type CampaignResponse = APIResponse<PaginatedResponse<Campaign & BaseEntity>>

type CampaignGroupValues = { campaign_id: string; group_id: string }

export type CampaignData = BaseEntity & Campaign

type CampaignStats = {
    hard_bounces: number;
    open_rate: number;
    soft_bounces: number;
    total_bounces: number;
    total_clicks: number;
    total_deliveries: number;
    total_emails_sent: number;
    total_opens: number;
    unique_clicks: number;
    unique_opens: number;
}

type CampaignEmailRecipientStats = {
    campaign_id: string;
    recipient_email: string;
    version: string;
    sent_at: string;
    opened_at: string | null;
    open_count: number;
    clicked_at: string | null;
    click_count: number;
    conversion_at: string | null;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
} & BaseEntity;

type CampaignUserStats = {
    hard_bounces: number;
    open_rate: number;
    soft_bounces: number;
    total_bounces: number;
    total_clicks: number;
    total_deliveries: number;
    total_emails_sent: number;
    total_opens: number;
    unique_clicks: number;
    unique_opens: number;
};

type EmailCampaignStats = {
    bounces: number;
    campaign_id: string;
    clicked: number;
    complaints: number;
    name: string;
    opened: number;
    recipients: number;
    sent_date: string | null;
    unsubscribed: number;
};


type CampaignStore = {
    createCampaignValues: CreateCampaignValues
    selectedCampaign: Campaign[]
    selectedGroupIds: string[]
    campaignStatData: CampaignStats
    campaignUserStatsData: CampaignUserStats
    allCampaignStatsData: EmailCampaignStats[]
    campaignRecipientData: CampaignEmailRecipientStats[]
    paginationInfo: Omit<PaginatedResponse<Campaign>, 'data'>;
    setCreateCampaignValues: (newData: CreateCampaignValues) => void
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<Campaign>, 'data'>) => void;
    campaignData: CampaignData[] | CampaignData | null
    scheduledCampaignData: (Campaign & BaseEntity)[] | Campaign & BaseEntity | null
    currentCampaignId: string | null;
    setCurrentCampaignId: (id: string | null) => void;
    setCampaignStats: (newData: CampaignStats) => void
    setAllCampaignStats: (newData: EmailCampaignStats[]) => void
    clearCurrentCampaignId: () => void;
    setCampaignEmailRecipients: (newData: CampaignEmailRecipientStats[]) => void
    setCampaignUserStats: (newData: CampaignUserStats) => void
    createCampaign: () => Promise<void>
    getAllCampaigns: (page?: number, pageSize?: number, search?: string) => Promise<void>
    setCampaignData: (newData: (Campaign & BaseEntity)[] | Campaign & BaseEntity) => void
    setScheduledCampaignData: (newData: (Campaign & BaseEntity)[] | Campaign & BaseEntity) => void
    getSingleCampaign: (uuid: string) => Promise<CampaignData | null>
    updateCampaign: (uuid: string) => Promise<void>
    setSelectedGroupIds: (newData: string[]) => void
    createCampaignGroup: (uuid: string, groupIds: string[]) => Promise<void>
    resetCampaignData: () => void;
    getScheduledCampaign: (page?: number, pageSize?: number) => Promise<void>
    sendCampaign: (campaignId: string) => Promise<void>
    deleteCampaign: (campaignId: string) => Promise<void>
    getCampaignStats: (campaignId: string) => Promise<void>
    getCampaignRecipients: (campaignId: string) => Promise<void>
    getCampaignUserStats: () => Promise<void>
    getAllCampaignStats: () => Promise<void>
    searchCampaign: (query?: string) => void
    persistStore: () => void
}



// Create a custom storage object
const customStorage: StateStorage = {
    getItem: (key) => {
        const value = localStorage.getItem(key);
        return value ? JSON.parse(value) : null;
    },
    setItem: (key, value) => {
        localStorage.setItem(key, JSON.stringify(value));
    },
    removeItem: (key) => {
        localStorage.removeItem(key);
    },
};


const useCampaignStore = create(persist<CampaignStore>((set, get) => ({

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
    scheduledCampaignData: [],
    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },
    allCampaignStatsData: [],
    campaignUserStatsData: {
        hard_bounces: 0,
        open_rate: 0,
        soft_bounces: 0,
        total_bounces: 0,
        total_clicks: 0,
        total_deliveries: 0,
        total_emails_sent: 0,
        total_opens: 0,
        unique_clicks: 0,
        unique_opens: 0,
    },
    campaignStatData: {
        hard_bounces: 0,
        open_rate: 0,
        soft_bounces: 0,
        total_bounces: 0,
        total_clicks: 0,
        total_deliveries: 0,
        total_emails_sent: 0,
        total_opens: 0,
        unique_clicks: 0,
        unique_opens: 0,
    },
    campaignRecipientData: [],
    selectedCampaign: [],
    selectedGroupIds: [] as string[],
    currentCampaignId: null,
    setCurrentCampaignId: (id) => set({ currentCampaignId: id }),
    clearCurrentCampaignId: () => set({ currentCampaignId: null }),
    setCampaignStats: (newData) => set({ campaignStatData: newData }),
    setCampaignData: (newData) => set({ campaignData: newData }),
    setCreateCampaignValues: (newData) => set({ createCampaignValues: newData }),
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),
    setSelectedGroupIds: (groupIds: string[]) => set({ selectedGroupIds: groupIds }),
    setScheduledCampaignData: (newData) => set({ scheduledCampaignData: newData }),
    setCampaignEmailRecipients: (newData) => set({ campaignRecipientData: newData }),
    setCampaignUserStats: (newData) => set({ campaignUserStatsData: newData }),
    setAllCampaignStats: (newData) => set({ allCampaignStatsData: newData }),

    createCampaign: async () => {
        try {
            const { createCampaignValues } = get()

            let response = await axiosInstance.post<ResponseT>("/campaigns/create-campaign", createCampaignValues)

            if (response.data.status == true) {
                eventBus.emit("success", "Campaign created successfully")
                await new Promise(resolve => setTimeout(resolve, 3000))
                window.location.href = "/user/dash/campaign/edit/" + response.data.payload.campaignId
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

    getAllCampaigns: async (page = 1, pageSize = 10, query = "") => {
        try {
            const { setCampaignData, setPaginationInfo } = get()
            let response = await axiosInstance.get<CampaignResponse>("/campaigns/get-all-campaigns", {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            })
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
    getSingleCampaign: async (uuid: string): Promise<CampaignData | null> => {
        try {
            const { setCampaignData } = get()
            let response = await axiosInstance.get<APIResponse<(BaseEntity & Campaign)>>("/campaigns/get-campaign/" + uuid)
            if (response.data.status == true) {
                const campaign: Campaign & BaseEntity = response.data.payload as Campaign & BaseEntity;
                setCampaignData(campaign);
            }
            return response.data.payload
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }

        return get().campaignData as CampaignData
    },

    updateCampaign: async (uuid: string) => {
        try {
            const { createCampaignValues, getSingleCampaign, setCreateCampaignValues } = get();

            // Fetch the current campaign data
            const currentCampaign = await getSingleCampaign(uuid) as CampaignData;
            if (!currentCampaign) {
                eventBus.emit('error', 'Campaign not found');
                return;
            }

            // Prepare the update payload
            const updatePayload: Partial<Campaign> = {};
            let successMessage = "Campaign updated successfully";

            // Check and update each field
            const fieldsToUpdate = [
                'name', 'subject', 'preview_text', 'sender', 'sender_from_name',
                'template_id', 'scheduled_at', 'status'
            ] as const;

            fieldsToUpdate.forEach(field => {
                if (createCampaignValues[field] !== undefined && createCampaignValues[field] !== currentCampaign[field]) {
                    updatePayload[field] = createCampaignValues[field];
                }
            });

            // Special handling for scheduled campaigns
            if (createCampaignValues.scheduled_at) {
                if (!currentCampaign.template?.email_html) {
                    eventBus.emit('error', 'You haven\'t created a template yet');
                    return;
                }
                if (!currentCampaign.subject) {
                    eventBus.emit('error', 'You have not created a subject yet');
                    return;
                }
                if (!currentCampaign.campaign_groups?.length) {
                    eventBus.emit('error', 'You have not created a recipient yet');
                    return;
                }
                updatePayload.status = 'scheduled';
                successMessage = "Your campaign has been scheduled successfully";
            }

            // Only proceed if there are changes to update
            if (Object.keys(updatePayload).length === 0) {
                eventBus.emit('info', 'No changes detected');
                return;
            }

            // Send the update request
            const response = await axiosInstance.put<ResponseT>(
                `/campaigns/update-campaign/${uuid}`,
                updatePayload
            );

            if (response.data.status) {
                eventBus.emit('success', successMessage);
                // Update the local state
                setCreateCampaignValues({
                    ...createCampaignValues,
                    ...updatePayload
                });

                // Manually persist the updated store
                get().persistStore();
            }

        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
                eventBus.emit('error', "An unknown error occurred");
            }
        }
    },
    createCampaignGroup: async (uuid: string, groupIds: string[]) => {
        try {
            console.log("uuid", uuid)
            console.log("groupIds", groupIds)

            if (groupIds.length > 0) {
                const promises = groupIds.map(groupId => {
                    const data = {
                        group_id: groupId,
                        campaign_id: uuid
                    } satisfies CampaignGroupValues

                    return axiosInstance.post<ResponseT>("/campaigns/add-campaign-group", data)
                })

                const results = await Promise.all(promises)

                const allSuccessful = results.every(response => response.data.status === true)

                if (allSuccessful) {
                    eventBus.emit('success', "Recipients added successfully")
                } else {
                    eventBus.emit('error', "Some recipients could not be created")
                }
            } else {
                eventBus.emit('error', "No groups selected")
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
    getScheduledCampaign: async (page = 1, pageSize = 10) => {
        try {
            const { setScheduledCampaignData, setPaginationInfo } = get()
            let response = await axiosInstance.get<CampaignResponse>(`/campaigns/get-scheduled-campaigns?page=${page}&page_size=${pageSize}`)
            const { data, ...paginationInfo } = response.data.payload;
            setScheduledCampaignData(data)
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
    sendCampaign: async (campaignId: string) => {
        try {
            const { getSingleCampaign } = get();
            const campaign = await getSingleCampaign(campaignId) as CampaignData || null

            // Check if campaign is defined and has the necessary properties
            if (!campaign) {
                eventBus.emit('error', 'You must fill all the neccessary fields')
                return
            }

            // Check if template exists and has content
            if (!campaign?.template || !campaign?.template.email_html) {
                eventBus.emit('error', 'You haven`t created a template yet')
                return
            }

            // Check if subject exists
            if (!campaign?.subject) {
                eventBus.emit('error', 'You have not created a subject yet')
                return
            }

            // Check if campaign_groups exist and have recipients
            if (!campaign?.campaign_groups || campaign?.campaign_groups.length === 0) {
                eventBus.emit('error', 'You have not created a recipient yet')
                return
            }

            if (!campaign?.sender || !campaign?.sender_from_name) {
                eventBus.emit('error', 'You have not created a sender ')
            }

            let value = { campaign_id: campaignId }

            //  If all checks pass, proceed with sending the campaign
            let response = await axiosInstance.post<ResponseT>(`/campaigns/send-campaign`, value);

            if (response.data.status === true) {
                eventBus.emit('success', "Campaign sent successfully");
            } else {
                throw new Error("Failed to send campaign");
            }

        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }

        }
    },

    deleteCampaign: async (campaignId: string) => {
        try {
            let response = await axiosInstance.delete<ResponseT>("/campaigns/delete-campaign/" + campaignId)

            if (response.data.status == true) {
                eventBus.emit("success", "Campaign deleted successfully")
            }
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },

    getCampaignStats: async (campaignId: string) => {
        try {
            let response = await axiosInstance.get<APIResponse<CampaignStats>>("/campaigns/get-stats/" + campaignId)
            get().setCampaignStats(response.data.payload)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },

    getCampaignRecipients: async (campaignId: string) => {
        try {
            let response = await axiosInstance.get<APIResponse<CampaignEmailRecipientStats[]>>("/campaigns/get-email-recipients/" + campaignId)
            get().setCampaignEmailRecipients(response.data.payload)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },
    getCampaignUserStats: async () => {
        try {
            let response = await axiosInstance.get<APIResponse<CampaignUserStats>>("/campaigns/user-campaign-stats")
            get().setCampaignUserStats(response.data.payload)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },
    getAllCampaignStats: async () => {
        try {
            let response = await axiosInstance.get<APIResponse<EmailCampaignStats[]>>("/campaigns/user-campaigns-stats")
            get().setAllCampaignStats(response.data.payload)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload);
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },
    searchCampaign: async (query?: string) => {

        const { getAllCampaigns } = get();

        if (!query) {
            await getAllCampaigns();
            return;
        }

        await getAllCampaigns(1, 10, query)
    },

    // Add a method to manually persist the store
    persistStore: () => {
        const state = get();
        customStorage.setItem('campaign-store', JSON.stringify(state));
    },

    resetCampaignData: () => set({ campaignData: null }),
}),

    {
        name: 'campaign-store',
        storage: createJSONStorage(() => localStorage),
    }
))


export default useCampaignStore
