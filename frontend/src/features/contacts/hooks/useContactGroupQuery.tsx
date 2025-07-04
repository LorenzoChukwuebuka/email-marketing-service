import { QueryClient, useQuery } from "@tanstack/react-query";
import { contactGroupAPI } from "../api/contactgroup.api";


/***************** All Contact Groups ******************* */

export const CONTACT_GROUP_QUERY_KEY = ["contact_group_key"] as const;

export const contactGroupQueryOptions = (page?: number, pageSize?: number, query?: string) => ({
    queryKey: [CONTACT_GROUP_QUERY_KEY, page, pageSize, query],
    queryFn: async () => {
        const response = await contactGroupAPI.getAllGroups(page, pageSize, query);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useContactGroupQuery = (page?: number, pageSize?: number, query?: string) => {
    return useQuery(contactGroupQueryOptions(page, pageSize, query));
};

export const invalidateContactGroupQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: CONTACT_GROUP_QUERY_KEY });
};

export const prefetchContactGroupQuery = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(contactGroupQueryOptions());
};

/***************** Single Contact Group ******************* */

export const SINGLE_CONTACT_GROUP_QUERY_KEY = ["single_contact_group_key"] as const;

export const singleContactGroupQueryOptions = (uuid: string) => ({
    queryKey: [SINGLE_CONTACT_GROUP_QUERY_KEY, uuid],
    queryFn: async () => {
        const response = await contactGroupAPI.getSingleGroup(uuid);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
    retry: 3,
    refetchInterval: 3 * 60 * 1000,
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useSingleContactGroupQuery = (uuid: string) => {
    return useQuery(singleContactGroupQueryOptions(uuid));
};

export const invalidateSingleContactGroupQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: SINGLE_CONTACT_GROUP_QUERY_KEY });
};

export const prefetchSingleContactGroupQuery = (queryClient: QueryClient, uuid: string) => {
    return queryClient.prefetchQuery(singleContactGroupQueryOptions(uuid));
};
