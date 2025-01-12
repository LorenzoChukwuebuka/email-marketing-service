import { QueryClient, useQuery } from "@tanstack/react-query";
import AnalyticsAPI from "../api/analytics.api";

export const CAMPAIGN_STATS_QUERY_KEY = ["CAMPAIGN_STATS_key"] as const;

export const campaignStatsQueryOptions = {
    queryKey: [CAMPAIGN_STATS_QUERY_KEY],
    queryFn: async () => {
        const response = await AnalyticsAPI.getAllCampaignStats();
        return response ?? [];
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
    retry: 3,
    refetchInterval: 3 * 60 * 1000,
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
};

export const useAllCampaignStatsQuery = () => {
    return useQuery(campaignStatsQueryOptions);
};

export const invalidateAllCampaignStatsQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: CAMPAIGN_STATS_QUERY_KEY });
}


/************************************************ USER STATS *********************************************************** */

export const CAMPAIGN_USER_STATS_QUERY_KEY = ["CAMPAIGN_STATS_key"] as const;

export const campaignUserStatsQueryOptions = {
    queryKey: [CAMPAIGN_STATS_QUERY_KEY],
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
    return queryClient.invalidateQueries({ queryKey: CAMPAIGN_STATS_QUERY_KEY });
}




