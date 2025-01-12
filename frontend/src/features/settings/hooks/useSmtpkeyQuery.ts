import { QueryClient, useQuery } from "@tanstack/react-query";
import smtpkeyApi from "../api/smtpkey.api";
export const SMTPKEY_QUERY_KEY = ["smtpkey_key"] as const;

/**************************** All SMTP Keys ********************************/

// Query options that will be used consistently across components
export const smtpkeyQueryOptions = () => ({
    // Add all the dynamic query to the key 
    queryKey: [SMTPKEY_QUERY_KEY],
    queryFn: async () => {
        const response = await smtpkeyApi.getSmtpKeys();
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

// Custom hook for using the smtpkey query
export const useSMTPKeyQuery = () => {
    return useQuery (
        smtpkeyQueryOptions()
    );
};

// Prefetch function that can be used anywhere
export const prefetchSMTPKey = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(smtpkeyQueryOptions());
};

export const invalidateSMTPKey = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: SMTPKEY_QUERY_KEY });
};
