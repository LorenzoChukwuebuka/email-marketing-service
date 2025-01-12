import { QueryClient, useQuery } from "@tanstack/react-query";
import { APIResponse } from '../../../interface/api.interface';
import domainApi from "../api/domain.api";
import { PaginatedResponse } from "../../../interface/pagination.interface";
import { DomainRecord } from "../interface/domain.interface";

export const DOMAIN_QUERY_KEY = ["domain_key"] as const;

/**************************** All Domains ********************************/

// Query options that will be used consistently across components
export const domainQueryOptions = (page?: number, pageSize?: number, query?: string) => ({
    // Add all the dynamic query to the key 
    queryKey: [DOMAIN_QUERY_KEY, page, pageSize, query],
    queryFn: async () => {
        const response = await domainApi.getAllDomains(page, pageSize, query);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minutes
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

// Custom hook for using the domain query
export const useDomainQuery = (page?: number, pageSize?: number, query?: string) => {
    return useQuery<APIResponse<PaginatedResponse<DomainRecord[]>>>(
        domainQueryOptions(page, pageSize, query)
    );
};

// Prefetch function that can be used anywhere
export const prefetchDomain = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(domainQueryOptions());
};

export const invalidateDomain = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: DOMAIN_QUERY_KEY });
};

/**************************** Single Domain ********************************/

export const singleDomainQueryOptions = (domainId: string) => ({
    queryKey: [DOMAIN_QUERY_KEY, domainId],
    queryFn: async () => {
        const response = await domainApi.getDomain(domainId);
        return response ?? null;
    },
    staleTime: 5 * 60 * 1000, // Keep single domain data fresh for 5 minutes
    cacheTime: 15 * 60 * 1000, // Cache single domain data for 15 minutes
    retry: 2,
});

// Custom hook for using the single domain query
export const useSingleDomainQuery = (domainId: string) => {
    return useQuery<APIResponse<DomainRecord>>(
        singleDomainQueryOptions(domainId)
    );
};

// Prefetch function for a single domain
export const prefetchSingleDomain = (queryClient: QueryClient, domainId: string) => {
    return queryClient.prefetchQuery(singleDomainQueryOptions(domainId));
};

export const invalidateSingleDomain = async (queryClient: QueryClient, domainId: string) => {
    await queryClient.invalidateQueries({ queryKey: [DOMAIN_QUERY_KEY, domainId] });
};
