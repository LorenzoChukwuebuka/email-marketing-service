import { QueryClient, useQuery } from "@tanstack/react-query";
import PlanAPI from "../api/plan.api";


export const PLANS_QUERY_KEY = ["plans_key"] as const;

/**************************** All Plans ********************************/

// Query options that will be used consistently across components
export const plansQueryOptions = () => ({
    // Add all the dynamic query to the key 
    queryKey: [PLANS_QUERY_KEY],
    queryFn: async () => {
        const response = await PlanAPI.getPlans();
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

// Custom hook for using the plans query
export const usePlansQuery = () => {
    return useQuery(
        plansQueryOptions()
    );
};

// Prefetch function that can be used anywhere
export const prefetchPlans = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(plansQueryOptions());
};

export const invalidatePlans = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: PLANS_QUERY_KEY });
};
