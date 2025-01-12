import { QueryClient, useQuery } from "@tanstack/react-query";
 
import { TemplateAPI } from "../api/email-template.api";


/***************** All marketing templates ******************* */
export const MARKETING_TEMPLATE_QUERY_KEY = ["marketing_template_key"] as const;

export const marketingTemplateQueryOptions = (page?: number, pageSize?: number, query?: string) => ({
    queryKey: [MARKETING_TEMPLATE_QUERY_KEY, page, pageSize, query],
    queryFn: async () => {
        const response = await TemplateAPI.getAllMarketingTemplates(page, pageSize, query);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
    retry: 3,
    refetchInterval: 3 * 60 * 1000,
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useMarketingTemplateQuery = (page?: number, pageSize?: number, query?: string) => {
    return useQuery(marketingTemplateQueryOptions(page, pageSize, query));
};

export const invalidateMarketingTemplateQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: MARKETING_TEMPLATE_QUERY_KEY });
}

export const prefetchMarketingTemplateQuery = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(marketingTemplateQueryOptions());
}

/***************** Single marketing template ******************* */
export const SINGLE_MARKETING_TEMPLATE_QUERY_KEY = ["single_marketing_template_key"] as const;

export const singleMarketingTemplateQueryOptions = (uuid: string) => ({
    queryKey: [SINGLE_MARKETING_TEMPLATE_QUERY_KEY, uuid],
    queryFn: async () => {
        const response = await TemplateAPI.getSingleMarketingTemplate(uuid);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
    retry: 3,
    refetchInterval: 3 * 60 * 1000,
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

export const useSingleMarketingTemplateQuery = (uuid: string) => {
    return useQuery(singleMarketingTemplateQueryOptions(uuid));
};

export const invalidateSingleMarketingTemplateQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: SINGLE_MARKETING_TEMPLATE_QUERY_KEY });
}

export const prefetchSingleMarketingTemplateQuery = (queryClient: QueryClient, uuid: string) => {
    return queryClient.prefetchQuery(singleMarketingTemplateQueryOptions(uuid));
}

