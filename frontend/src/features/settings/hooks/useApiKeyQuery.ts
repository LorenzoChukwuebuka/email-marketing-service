import { QueryClient, useQuery } from "@tanstack/react-query";
import apikeyApi from "../api/apikey.api";


export const APIKEY_QUERY_KEY = ["apikey_key"] as const;

/**************************** All API Keys ********************************/

// Query options that will be used consistently across components
export const apikeyQueryOptions = () => ({
    // Add all the dynamic query to the key 
    queryKey: [APIKEY_QUERY_KEY],
    queryFn: async () => {
        const response = await apikeyApi.getAPIKey();
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

// Custom hook for using the apikey query
export const useAPIKeyQuery = () => {
    return useQuery(
        apikeyQueryOptions()
    );
};

// Prefetch function that can be used anywhere
export const prefetchAPIKey = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(apikeyQueryOptions());
};

export const invalidateAPIKey = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: APIKEY_QUERY_KEY });
};
