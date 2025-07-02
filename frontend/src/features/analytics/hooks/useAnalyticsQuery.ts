import { QueryClient, useQuery } from "@tanstack/react-query";
import AnalyticsAPI from "../api/analytics.api";

export const CAMPAIGN_STATS_QUERY_KEY = ["CAMPAIGN_STATS_key"] as const;

export const campaignStatsQueryOptions = (page?: number, pageSize?: number, query?: string) => ({
    queryKey: [CAMPAIGN_STATS_QUERY_KEY, page, pageSize, query],
    queryFn: async () => {
        const response = await AnalyticsAPI.getAllCampaignStats(page, pageSize);
        return response ?? [];
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
    retry: 3,
    refetchInterval: 3 * 60 * 1000,
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useAllCampaignStatsQuery = (page?: number, pageSize?: number, query?: string) => {
    return useQuery(campaignStatsQueryOptions(page, pageSize, query));
};

export const invalidateAllCampaignStatsQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: CAMPAIGN_STATS_QUERY_KEY });
};

/************************************************ USER STATS *********************************************************** */
export const CAMPAIGN_USER_STATS_QUERY_KEY = ["CAMPAIGN_USER_STATS_key"] as const;

export const campaignUserStatsQueryOptions = {
    queryKey: [CAMPAIGN_USER_STATS_QUERY_KEY], // ✅ Fixed: Use correct query key
    queryFn: async () => {
        const response = await AnalyticsAPI.getCampaignUserStats();
        return response ?? [];
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
    retry: 3,
    refetchInterval: 3 * 60 * 1000,
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
};

export const useCampaignUserStatsQuery = () => {
    return useQuery(campaignUserStatsQueryOptions);
};

export const invalidateCampaignUserStatsQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: CAMPAIGN_USER_STATS_QUERY_KEY }); // ✅ Fixed: Use correct query key
};