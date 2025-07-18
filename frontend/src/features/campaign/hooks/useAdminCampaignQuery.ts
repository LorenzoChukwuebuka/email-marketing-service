import { QueryClient, useQuery } from "@tanstack/react-query";
import adminCampaignApi from "../api/admincampaign.api";
import { APIResponse } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';
import { CampaignData } from '../interface/campaign.interface';

/***************** Campaign Query Keys *******************/
export const CAMPAIGN_QUERY_KEY = ["campaign"] as const;

// Reusable query options for all APIs
const createCampaignQueryOptions = (
    endpoint: string,
    fetchFn: () => Promise<any>,
    dependencies: any[] = []
) => ({
    queryKey: [CAMPAIGN_QUERY_KEY, endpoint, ...dependencies],
    queryFn: fetchFn,
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 15 * 60 * 1000, // 15 minutes
    retry: 3,
    refetchOnWindowFocus: true,
});

/***************** Hooks for APIs *******************/

// Get all user campaigns
export const useAdminUserCampaignsQuery = (
    userId: string,
    companyId: string,
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return useQuery<APIResponse<PaginatedResponse<CampaignData>>>(
        createCampaignQueryOptions(
            "user_campaigns",
            () => adminCampaignApi.getAllUserCampaign(userId, companyId, page, pageSize, query),
            [userId, page, pageSize, query]
        )
    );
};

// Get single campaign
export const useAdminUserSingleCampaignQuery = (campaignId: string, userId: string, companyId: string) => {
    return useQuery<APIResponse<CampaignData>>(
        createCampaignQueryOptions(
            "single_campaign",
            () => adminCampaignApi.getSingleUserCampaign(campaignId, userId, companyId),
            [campaignId]
        )
    );
};

/***************** Query Invalidations *******************/

export const invalidateCampaignQueries = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: CAMPAIGN_QUERY_KEY });
};

export const invalidateUserCampaignsQuery = (queryClient: QueryClient, userId: string) => {
    return queryClient.invalidateQueries({
        queryKey: [CAMPAIGN_QUERY_KEY, "user_campaigns", userId]
    });
};

export const invalidateSingleCampaignQuery = (queryClient: QueryClient, campaignId: string) => {
    return queryClient.invalidateQueries({
        queryKey: [CAMPAIGN_QUERY_KEY, "single_campaign", campaignId]
    });
};

/***************** Prefetch Functions *******************/

export const prefetchUserCampaignsQuery = (
    queryClient: QueryClient,
    userId: string,
    companyId: string,
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return queryClient.prefetchQuery(
        createCampaignQueryOptions(
            "user_campaigns",
            () => adminCampaignApi.getAllUserCampaign(userId, companyId, page, pageSize, query),
            [userId, page, pageSize, query]
        )
    );
};

export const prefetchSingleCampaignQuery = (
    queryClient: QueryClient,
    campaignId: string,
    userId: string,
    companyId: string
) => {
    return queryClient.prefetchQuery(
        createCampaignQueryOptions(
            "single_campaign",
            () => adminCampaignApi.getSingleUserCampaign(campaignId, userId, companyId),
            [campaignId]
        )
    );
};


/******************************** Get Campaign Recipients ********************************************/

export const CAMPAIGN_RECIPIENTS_QUERY_KEY = ["campaign_recipients_key"] as const;

// Query options for getting campaign recipients
export const campaignRecipientsQueryOptions = (id: string, companyId: string) => ({
    queryKey: [CAMPAIGN_RECIPIENTS_QUERY_KEY, id],
    queryFn: async () => {
        const response = await adminCampaignApi.getCampaignRecipients(id, companyId);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000, // Data considered fresh for 1 minute
    cacheTime: 10 * 60 * 1000, // Cache unused data for 10 minutes
    retry: 3, // Retry up to 3 times if query fails
    refetchOnWindowFocus: true, // Refetch on window focus
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minutes
    refetchIntervalInBackground: true, // Allow background refetching
});

// Custom hook for using the campaign recipients query
export const useAdminCampaignRecipientsQuery = (id: string, companyId: string) => {
    return useQuery(campaignRecipientsQueryOptions(id, companyId));
};

// Prefetch function for campaign recipients
export const prefetchCampaignRecipients = (queryClient: QueryClient, id: string, companyId: string) => {
    return queryClient.prefetchQuery(campaignRecipientsQueryOptions(id, companyId));
};

// Invalidate campaign recipients query
export const invalidateCampaignRecipients = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: CAMPAIGN_RECIPIENTS_QUERY_KEY });
};
