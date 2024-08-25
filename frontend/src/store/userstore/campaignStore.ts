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


type CampaignStore = {
    createCampaignValues: CreateCampaignValues
    selectedCampaign: Campaign[]
    selectedGroupIds: string[]
    paginationInfo: Omit<PaginatedResponse<Campaign>, 'data'>;
    setCreateCampaignValues: (newData: CreateCampaignValues) => void
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<Campaign>, 'data'>) => void;
    campaignData: CampaignData[] | CampaignData | null
    scheduledCampaignData: (Campaign & BaseEntity)[] | Campaign & BaseEntity | null
    currentCampaignId: string | null;
    setCurrentCampaignId: (id: string | null) => void;
    clearCurrentCampaignId: () => void;
    createCampaign: () => Promise<void>
    getAllCampaigns: (page?: number, pageSize?: number) => Promise<void>
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
}

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
    selectedCampaign: [],
    selectedGroupIds: [] as string[],
    currentCampaignId: null,
    setCurrentCampaignId: (id) => set({ currentCampaignId: id }),
    clearCurrentCampaignId: () => set({ currentCampaignId: null }),
    setCampaignData: (newData) => set({ campaignData: newData }),
    setCreateCampaignValues: (newData) => set({ createCampaignValues: newData }),
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),
    setSelectedGroupIds: (groupIds: string[]) => set({ selectedGroupIds: groupIds }),
    setScheduledCampaignData: (newData) => set({ scheduledCampaignData: newData }),

    createCampaign: async () => {
        try {
            const { createCampaignValues } = get()

            let cookie: any = Cookies.get("Cookies");
            let user = JSON.parse(cookie)?.details?.company;
            const updatedCampaignValues = { ...createCampaignValues, sender_from_name: user };
            let response = await axiosInstance.post<ResponseT>("/campaigns/create-campaign", updatedCampaignValues)

            if (response.data.status == true) {
                eventBus.emit("success", "Campaign created successfully")

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
            const { createCampaignValues, getSingleCampaign } = get();

            let successMessage = ""

            if (createCampaignValues.scheduled_at) {
                successMessage = "Your campaign has been scheduled successfully";
            } else if (createCampaignValues.subject) {
                successMessage = "Subject was added successfully";
            } else if (createCampaignValues.template_id) {
                successMessage = "Template added successfully";
            } else {
                console.error("No recognizable field was found");
            }


            let updatedValues = {}
            if (createCampaignValues.scheduled_at) {
                const campaign = await getSingleCampaign(uuid) as CampaignData || null;

                if (!campaign) {
                    eventBus.emit('success', 'You must fill all the necessary fields');
                    return;
                }

                const { template, subject, campaign_groups } = campaign;

                if (!template?.email_html) {
                    eventBus.emit('success', 'You haven`t created a template yet');
                    return;
                }

                if (!subject) {
                    eventBus.emit('success', 'You have not created a subject yet');
                    return;
                }

                if (!campaign_groups?.length) {
                    eventBus.emit('success', 'You have not created a recipient yet');
                    return;
                }

                updatedValues = { ...createCampaignValues, status: "scheduled" }
            }

            const response = await axiosInstance.put<ResponseT>(
                `/campaigns/update-campaign/${uuid}`,
                updatedValues
            );

            if (response.data.status) {
                eventBus.emit('success', successMessage)
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
                eventBus.emit('success', 'You must fill all the neccessary fields')
                return
            }

            // Check if template exists and has content
            if (!campaign?.template || !campaign?.template.email_html) {
                eventBus.emit('success', 'You haven`t created a template yet')
                return
            }

            // Check if subject exists
            if (!campaign?.subject) {
                eventBus.emit('success', 'You have not created a subject yet')
                return
            }

            // Check if campaign_groups exist and have recipients
            if (!campaign?.campaign_groups || campaign?.campaign_groups.length === 0) {
                eventBus.emit('success', 'You have not created a recipient yet')
                return
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

    resetCampaignData: () => set({ campaignData: null }),
}),

    {
        name: 'campaign-store',
        storage: createJSONStorage(() => localStorage),
    }
))


export default useCampaignStore
