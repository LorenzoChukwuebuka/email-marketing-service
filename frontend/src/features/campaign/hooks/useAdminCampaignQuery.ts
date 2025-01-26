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
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return useQuery<APIResponse<PaginatedResponse<CampaignData>>>(
        createCampaignQueryOptions(
            "user_campaigns",
            () => adminCampaignApi.getAllUserCampaign(userId, page, pageSize, query),
            [userId, page, pageSize, query]
        )
    );
};

// Get single campaign
export const useAdminUserSingleCampaignQuery = (campaignId: string) => {
    return useQuery<APIResponse<CampaignData>>(
        createCampaignQueryOptions(
            "single_campaign",
            () => adminCampaignApi.getSingleUserCampaign(campaignId),
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
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return queryClient.prefetchQuery(
        createCampaignQueryOptions(
            "user_campaigns",
            () => adminCampaignApi.getAllUserCampaign(userId, page, pageSize, query),
            [userId, page, pageSize, query]
        )
    );
};

export const prefetchSingleCampaignQuery = (
    queryClient: QueryClient,
    campaignId: string
) => {
    return queryClient.prefetchQuery(
        createCampaignQueryOptions(
            "single_campaign",
            () => adminCampaignApi.getSingleUserCampaign(campaignId),
            [campaignId]
        )
    );
};