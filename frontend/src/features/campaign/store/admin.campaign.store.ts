import { create } from 'zustand'
import { handleError } from '../../../utils/isError';


type AdminUserSpecificCampaignStore = {
    suspendCampaign: (campaignId: string) => Promise<void>
    deleteCampaign: (campaignId: string) => Promise<void>
}

const useAdminUserCamapaignStore = create<AdminUserSpecificCampaignStore>(() => ({
    // Suspend a campaign by ID
    suspendCampaign: async (campaignId: string) => {
        try {
            console.log(campaignId)
        } catch (error) {
            handleError(error)
        }
    },

    // Delete a campaign by ID
    deleteCampaign: async (campaignId: string) => {
        try {
            console.log(campaignId)
        } catch (error) {
            handleError(error)
        }
    }
}));

export default useAdminUserCamapaignStore;
