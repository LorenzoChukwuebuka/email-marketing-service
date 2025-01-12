import { QueryClient, useQuery } from "@tanstack/react-query";
import BillingApi from "../api/biiling.api";


export const BILLING_QUERY_KEY = ["billing_key"] as const;

/**************************** All Billing ********************************/

// Query options that will be used consistently across components
export const billingQueryOptions = (page?: number, pageSize?: number) => ({
    // Add all the dynamic query to the key 
    queryKey: [BILLING_QUERY_KEY, page, pageSize],
    queryFn: async () => {
        const response = await BillingApi.fetchBilling(page, pageSize);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

// Custom hook for using the billing query
export const useBillingQuery = (page?: number, pageSize?: number) => {
    return useQuery(
        billingQueryOptions(page, pageSize)
    );
};

// Prefetch function that can be used anywhere
export const prefetchBilling = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(billingQueryOptions());
};

export const invalidateBilling = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: BILLING_QUERY_KEY });
};
