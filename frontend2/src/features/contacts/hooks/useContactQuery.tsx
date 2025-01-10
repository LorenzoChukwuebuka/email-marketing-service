import { QueryClient, useQuery } from "@tanstack/react-query";
import contactApi from "../api/contact.api";


/***************** All contacts ******************* */

export const CONTACT_QUERY_KEY = ["contact_key"] as const;

// Query options that will be used consistently across components
export const contactQueryOptions = (page?: number, pageSize?: number, query?: string) => ({
    // Add all the dynamic query to the key 
    queryKey: [CONTACT_QUERY_KEY, page, pageSize, query],
    queryFn: async () => {
        const response = await contactApi.getAllContacts(page, pageSize, query);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useContactQuery = (page?: number, pageSize?: number, query?: string) => {
    return useQuery(contactQueryOptions(page, pageSize, query));
};

export const invalidateContactQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: CONTACT_QUERY_KEY });
}

export const prefechContactQuery = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(contactQueryOptions());
}


/***************** Single contact ******************* */

export const SINGLE_CONTACT_QUERY_KEY = ["single_contact_key"] as const;

// Query options that will be used consistently across components
export const singleContactQueryOptions = (uuid: string) => ({
    // Add all the dynamic query to the key 
    queryKey: [SINGLE_CONTACT_QUERY_KEY, uuid],
    queryFn: async () => {
        const response = await contactApi.getContact(uuid);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useSingleContactQuery = (uuid: string) => {
    return useQuery(singleContactQueryOptions(uuid));
};

export const invalidateSingleContactQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: SINGLE_CONTACT_QUERY_KEY });
}

export const prefechSingleContactQuery = (queryClient: QueryClient, uuid: string) => {
    return queryClient.prefetchQuery(singleContactQueryOptions(uuid));
}

/**************************** Get contact count ***************************************** */

export const CONTACT_COUNT_QUERY_KEY = ["contact_count_key"] as const;

// Query options that will be used consistently across components
export const contactCountQueryOptions = () => ({
    // Add all the dynamic query to the key 
    queryKey: CONTACT_COUNT_QUERY_KEY,
    queryFn: async () => {
        const response = await contactApi.getContactCount();
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useContactCountQuery = () => {
    return useQuery(contactCountQueryOptions());
};

/**************************************** Get contact engagement **********************************************************/

export const CONTACT_ENGAGEMENT_QUERY_KEY = ["contact_engagement_key"] as const;

// Query options that will be used consistently across components

export const contactEngagementQueryOptions = () => ({
    // Add all the dynamic query to the key 
    queryKey: CONTACT_ENGAGEMENT_QUERY_KEY,
    queryFn: async () => {
        const response = await contactApi.getContactEngagement();
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useContactEngagementQuery = () => {
    return useQuery(contactEngagementQueryOptions());
};

export const invalidateContactEngagementQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: CONTACT_ENGAGEMENT_QUERY_KEY });
}

export const prefechContactEngagementQuery = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(contactEngagementQueryOptions());
}



