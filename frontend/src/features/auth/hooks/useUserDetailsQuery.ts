import { QueryClient, useQuery } from "@tanstack/react-query";
import authApi from "../api/auth.api";

export const USER_DETAILS_QUERY_KEY = ["user_details_key"] as const;

/**************************** All Campaigns ********************************/

// Query options that will be used consistently across components
export const userDetailsQueryOptions = () => ({
    // Add all the dynamic query to the key 
    queryKey: [USER_DETAILS_QUERY_KEY,],
    queryFn: async () => {
        const response = await authApi.getUserDetails();
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
export const useUserDetailsQuery = () => {
    return useQuery(
        userDetailsQueryOptions()
    );
};

// Prefetch function that can be used anywhere
export const prefetchUserDetails = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(userDetailsQueryOptions());
};

export const invalidateUserDetails = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: USER_DETAILS_QUERY_KEY });
};
