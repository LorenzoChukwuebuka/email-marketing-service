import { create } from "zustand";
import { persist } from 'zustand/middleware';
import { CreateCampaignValues, CampaignGroupValues } from '../interface/campaign.interface';
import { Campaign } from '../../../../../frontend/src/store/userstore/campaignStore';
import campaignApi from "../api/campaign.api";
import eventBus from "../../../utils/eventbus";
import { invalidateCampaign } from "../hooks/useCampaignQuery";
import { queryClient } from "../../../utils/queryclient";
import { handleError } from "../../../utils/isError";

type CampaignState = {
    createCampaignValues: Partial<CreateCampaignValues>;
    selectedCampaign: Campaign[];
    selectedGroupIds: string[];
    currentCampaignId: string | null;
}

type CampaignActions = {
    setCreateCampaignValues: (newValues: Partial<CreateCampaignValues>) => void;
    setSelectedGroupIds: (groupIds: string[]) => void;
    setCurrentCampaignId: (id: string | null) => void;
    clearCurrentCampaignId: () => void;
}

type CampaignAsyncActions = {
    createCampaignGroup: (uuid: string, groupIds: string[]) => Promise<void>;
    deleteCampaign: (campaignId: string) => Promise<void>;
    updateCampaign: (uuid: string) => Promise<void>;
    createCampaign: () => Promise<void>;
    sendCampaign: (campaignId: string) => Promise<void>
}

type CampaignStore = CampaignState & CampaignActions & CampaignAsyncActions;

const initialState: CampaignState = {
    createCampaignValues: {},
    selectedCampaign: [],
    selectedGroupIds: [],
    currentCampaignId: null
};

const useCampaignStore = create<CampaignStore>()(
    persist(
        (set, get) => ({
            ...initialState,
            setCurrentCampaignId: (id) => set({ currentCampaignId: id }),
            clearCurrentCampaignId: () => set({ currentCampaignId: null }),
            setCreateCampaignValues: (newData) => set({ createCampaignValues: newData }),
            setSelectedGroupIds: (groupIds) => set({ selectedGroupIds: groupIds }),

            createCampaign: async () => {
                const { createCampaignValues } = get()
                try {
                    const response = await campaignApi.createCampaign(createCampaignValues)
                    if (response) {
                        eventBus.emit('success', 'campaign created successfully')
                    }
                    invalidateCampaign(queryClient)
                } catch (error) {
                    console.error('Failed to create campaign:', error);
                    throw error;
                }
            },
            createCampaignGroup: async (uuid, groupIds) => {
                try {
                    if (groupIds.length > 0) {
                        const promises = groupIds.map(groupId => {
                            const data = {
                                group_id: groupId,
                                campaign_id: uuid
                            } satisfies CampaignGroupValues
                            return campaignApi.createCampaignGroup(data)
                        })

                        const results = await Promise.all(promises)
                        const allSuccessful = results.every(response => response.payload.status === true)

                        if (allSuccessful) {
                            eventBus.emit('success', "Recipients added successfully")
                        }
                    } else {
                        eventBus.emit('error', "No groups selected")
                    }
                    invalidateCampaign(queryClient)
                } catch (error) {
                    console.error('Failed to create campaign group:', error);
                    throw error;
                }
            },
            updateCampaign: async (uuid) => {
                const { createCampaignValues, setCreateCampaignValues } = get();
                try {
                    // fetch the current campaign data
                    const currentCampaign = await campaignApi.getSingleCampaign(uuid)

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
                        if (!currentCampaign.payload.template?.email_html) {
                            eventBus.emit('error', 'You haven\'t created a template yet');
                            return;
                        }
                        if (!currentCampaign.payload.subject) {
                            eventBus.emit('error', 'You have not created a subject yet');
                            return;
                        }
                        if (!currentCampaign.payload.campaign_groups?.length) {
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
                    const response = await campaignApi.updateCampaign(uuid, updatePayload)

                    if (response.status) {
                        eventBus.emit('success', successMessage);
                        // Update the local state
                        setCreateCampaignValues({
                            ...createCampaignValues,
                            ...updatePayload
                        });
                    }
                    invalidateCampaign(queryClient)

                } catch (error) {
                    console.error('Failed to update campaign:', error);
                    throw error;
                }
            },
            deleteCampaign: async (campaignId) => {
                try {
                    const response = await campaignApi.deleteCampaign(campaignId)

                    if (response) {
                        eventBus.emit('success', 'campaign deleted successfully')
                    }
                } catch (error) {
                    console.error('Failed to delete campaign:', error);
                    throw error;
                }
            },
            sendCampaign: async (campaignid) => {
                try {
                    const campaign = await campaignApi.getSingleCampaign(campaignid)

                    if (!campaign) {
                        eventBus.emit('error', 'Campaign not found');
                        return;
                    }

                    // Check if campaign is defined and has the necessary properties
                    if (!campaign) {
                        eventBus.emit('error', 'You must fill all the neccessary fields')
                        return
                    }

                    // Check if template exists and has content
                    if (!campaign?.payload.template || !campaign?.payload.template.email_html) {
                        eventBus.emit('error', 'You haven`t created a template yet')
                        return
                    }

                    // Check if subject exists
                    if (!campaign?.payload.subject) {
                        eventBus.emit('error', 'You have not created a subject yet')
                        return
                    }

                    // Check if campaign_groups exist and have recipients
                    if (!campaign?.payload.campaign_groups || campaign?.payload.campaign_groups.length === 0) {
                        eventBus.emit('error', 'You have not created a recipient yet')
                        return
                    }

                    if (!campaign?.payload.sender || !campaign?.payload.sender_from_name) {
                        eventBus.emit('error', 'You have not created a sender ')
                    }

                    const value = { campaign_id: campaignid }

                    const response = await campaignApi.sendCampaign(value)

                    if (response) {
                        eventBus.emit('success', 'campaign sent successfully')
                    }
                } catch (error) {
                    handleError(error)
                }
            },
            // Add a method to manually persist the store
            // persistStore: () => {
            //     const state = get();
            //     customStorage.setItem('campaign-store', JSON.stringify(state));
            // },

            //  resetCampaignData: () => set({ campaignData: null }),
        }),
        {
            name: 'campaign-store',
            partialize: (state) => ({ selectedCampaign: state.selectedCampaign,currentCampaignId:state.currentCampaignId })
        }
    )
);

export default useCampaignStore;