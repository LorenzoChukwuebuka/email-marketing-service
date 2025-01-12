import { QueryClient, useQuery } from "@tanstack/react-query";
import { APIResponse } from '../../../frontend/src/interface/api.interface';
import dailymailcalculationApi from "../api/dailymailcalculation.api";
import { MailData } from '../interface/mailcalc.interface';

export const MAILCALC_QUERY_KEY = ["mailCalc_key"] as const;

// Query options that will be used consistently across components
export const mailCalcQueryOptions = {
    // Add all the dynamic query to the key 
    queryKey: [MAILCALC_QUERY_KEY],
    queryFn: async () => {
        const response = await dailymailcalculationApi.getDailyMailCalculation();
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
};

// Custom hook for using the mailCalc query
export const useMailCalcQuery = () => {
    return useQuery<APIResponse<MailData>>(
        mailCalcQueryOptions
    );
};

// Prefetch function that can be used anywhere
export const prefetchMailCalcs = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(mailCalcQueryOptions);
};

export const invalidateMailCalc = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: MAILCALC_QUERY_KEY });
};
