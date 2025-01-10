import { QueryClient, useQuery } from "@tanstack/react-query";
import { APIResponse } from '../../../interface/api.interface';
import campaignApi from "../api/campaign.api";
import { PaginatedResponse } from "../../../interface/pagination.interface";
import { CampaignResponse } from "../interface/campaign.interface";

export const CAMPAIGN_QUERY_KEY = ["campaign_key"] as const;

/**************************** All Campaigns ********************************/

// Query options that will be used consistently across components
export const campaignQueryOptions = (page?: number, pageSize?: number, query?: string) => ({
    // Add all the dynamic query to the key 
    queryKey: [CAMPAIGN_QUERY_KEY, page, pageSize, query],
    queryFn: async () => {
        const response = await campaignApi.getAllCampaigns(page, pageSize, query);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

// Custom hook for using the mailCalc query
export const useCampaignQuery = (page?: number, pageSize?: number, query?: string) => {
    return useQuery<APIResponse<PaginatedResponse<CampaignResponse[]>>>(
        campaignQueryOptions(page, pageSize, query)
    );
};

// Prefetch function that can be used anywhere
export const prefetchCampaign = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(campaignQueryOptions());
};

export const invalidateCampaign = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: CAMPAIGN_QUERY_KEY });
};


/******************************** Get single campaign ********************************************/

export const SINGLE_CAMPAIGN_QUERY_KEY = ["campaign_key"] as const;

// Query options that will be used consistently across components
export const singleCampaignQueryOptions = (uuid: string) => ({
    // Add all the dynamic query to the key 
    queryKey: [SINGLE_CAMPAIGN_QUERY_KEY, uuid],
    queryFn: async () => {
        const response = await campaignApi.getSingleCampaign(uuid);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

// Custom hook for using the mailCalc query
export const useSingleCampaignQuery = (uuid: string) => {
    return useQuery<APIResponse<CampaignResponse>>(
        singleCampaignQueryOptions(uuid)
    );
};

// Prefetch function that can be used anywhere
export const prefetchSingleCampaign = (queryClient: QueryClient,uuid:string) => {
    return queryClient.prefetchQuery(singleCampaignQueryOptions(uuid));
};

export const invalidateSingleCampaign = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: SINGLE_CAMPAIGN_QUERY_KEY });
};

/*********************************** Get Scheduled campaign **********************************/

export const SCHEDULED_CAMPAIGN_QUERY_KEY = ["scheduled_campaign_key"] as const;

// Query options that will be used consistently across components
export const scheduledCampaignQueryOptions = (page?: number, pageSize?: number) => ({
    // Add all the dynamic query to the key 
    queryKey: [SCHEDULED_CAMPAIGN_QUERY_KEY, page, pageSize],
    queryFn: async () => {
        const response = await campaignApi.getScheduledCampaigns(page, pageSize);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

// Custom hook for using the mailCalc query
export const useScheduledCampaignQuery = (page?: number, pageSize?: number) => {
    return useQuery<CampaignResponse>(
        scheduledCampaignQueryOptions(page, pageSize)
    );
};

// Prefetch function that can be used anywhere
export const prefetchScheduledCampaign = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(scheduledCampaignQueryOptions());
};

export const invalidateScheduledCampaign = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: SCHEDULED_CAMPAIGN_QUERY_KEY });
};

