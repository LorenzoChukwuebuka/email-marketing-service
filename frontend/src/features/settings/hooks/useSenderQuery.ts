import { QueryClient, useQuery } from "@tanstack/react-query";
import { APIResponse } from '../../../interface/api.interface';
import senderApi from "../api/sender.api";
import { PaginatedResponse } from "../../../interface/pagination.interface";
import { Sender } from "../interface/sender.interface";

export const SENDER_QUERY_KEY = ["sender_key"] as const;

/**************************** All Senders ********************************/

// Query options that will be used consistently across components
export const senderQueryOptions = (page?: number, pageSize?: number, query?: string) => ({
    // Add all the dynamic query to the key 
    queryKey: [SENDER_QUERY_KEY, page, pageSize, query],
    queryFn: async () => {
        const response = await senderApi.getAllSenders(page, pageSize, query);
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
export const useSenderQuery = (page?: number, pageSize?: number, query?: string) => {
    return useQuery<APIResponse<PaginatedResponse<Sender>>>(
        senderQueryOptions(page, pageSize, query)
    );
};

// Prefetch function that can be used anywhere
export const prefetchSender = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(senderQueryOptions());
};

export const invalidateSender = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: SENDER_QUERY_KEY });
};
