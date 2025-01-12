import { QueryClient, useQuery } from "@tanstack/react-query";

import { TemplateAPI } from "../api/email-template.api";


/***************** All transaction templates ******************* */
export const TRANSACTIONAL_TEMPLATE_QUERY_KEY = ["TRANSACTIONAL_template_key"] as const;

export const TransactionalTemplateQueryOptions = (page?: number, pageSize?: number, query?: string) => ({
    queryKey: [TRANSACTIONAL_TEMPLATE_QUERY_KEY, page, pageSize, query],
    queryFn: async () => {
        const response = await TemplateAPI.getAllTransactionalTemplates(page, pageSize, query);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
    retry: 3,
    refetchInterval: 3 * 60 * 1000,
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useTransactionalTemplateQuery = (page?: number, pageSize?: number, query?: string) => {
    return useQuery(TransactionalTemplateQueryOptions(page, pageSize, query));
};

export const invalidateTransactionalTemplateQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: TRANSACTIONAL_TEMPLATE_QUERY_KEY });
}

export const prefetchTransactionalTemplateQuery = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(TransactionalTemplateQueryOptions());
}

/***************** Single transaction template ******************* */
export const SINGLE_TRANSACTIONAL_TEMPLATE_QUERY_KEY = ["single_TRANSACTIONAL_template_key"] as const;

export const singleTransactionalTemplateQueryOptions = (uuid: string) => ({
    queryKey: [SINGLE_TRANSACTIONAL_TEMPLATE_QUERY_KEY, uuid],
    queryFn: async () => {
        const response = await TemplateAPI.getSingleTransactionalTemplate(uuid);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
    retry: 3,
    refetchInterval: 3 * 60 * 1000,
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useSingleTransactionalTemplateQuery = (uuid: string) => {
    return useQuery(singleTransactionalTemplateQueryOptions(uuid));
};

export const invalidateSingleTransactionalTemplateQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: SINGLE_TRANSACTIONAL_TEMPLATE_QUERY_KEY });
}

export const prefetchSingleMarketingTemplateQuery = (queryClient: QueryClient, uuid: string) => {
    return queryClient.prefetchQuery(singleTransactionalTemplateQueryOptions(uuid));
}

