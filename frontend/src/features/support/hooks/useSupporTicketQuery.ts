import { QueryClient, useQuery } from "@tanstack/react-query";
import { APIResponse } from '../../../interface/api.interface';
import { SupportTicket } from "../api/supportticket.api";
import { Ticket } from '../interface/support.interface';
export const SUPPORT_TICKET_QUERY_KEY = ["support_ticket_key"] as const;
export const TICKET_DETAILS_QUERY_KEY = ["ticket_details_key"] as const;

/**************************** All Tickets ********************************/

// Query options for all tickets
export const supportTicketQueryOptions = () => ({
    queryKey: [SUPPORT_TICKET_QUERY_KEY],
    queryFn: async () => {
        const response = await SupportTicket.getTickets();
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 10 * 60 * 1000, // Keep unused data in cache for 10 minutes
    retry: 3,
    refetchInterval: 3 * 60 * 1000, // Refetch every 3 minute
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
});

// Custom hook for using the support ticket query
export const useSupportTicketQuery = () => {
    return useQuery<APIResponse<Ticket[]>>(
        supportTicketQueryOptions()
    );
};

/**************************** Single Ticket Details ********************************/

// Query options for single ticket details
export const ticketDetailsQueryOptions = (uuid: string) => ({
    queryKey: [TICKET_DETAILS_QUERY_KEY, uuid],
    queryFn: async () => {
        const response = await SupportTicket.getTicketDetails(uuid);
        return response ?? null;
    },
    staleTime: 1 * 60 * 1000,
    cacheTime: 5 * 60 * 1000, // Keep unused data in cache for 5 minutes
    retry: 2,
    refetchInterval: 1 * 60 * 1000, // Refetch every minute to get latest updates
    refetchIntervalInBackground: true,
    refetchOnWindowFocus: true,
    enabled: !!uuid, // Only run query if uuid is provided
});

// Custom hook for using the ticket details query
export const useTicketDetailsQuery = (uuid: string) => {
    return useQuery<APIResponse<Ticket>>(
        ticketDetailsQueryOptions(uuid)
    );
};

/**************************** Prefetch & Invalidation ********************************/

// Prefetch function for all tickets
export const prefetchSupportTickets = (queryClient: QueryClient) => {
    return queryClient.prefetchQuery(supportTicketQueryOptions());
};

// Prefetch function for ticket details
export const prefetchTicketDetails = (queryClient: QueryClient, uuid: string) => {
    return queryClient.prefetchQuery(ticketDetailsQueryOptions(uuid));
};

// Invalidate all ticket queries
export const invalidateSupportTickets = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: SUPPORT_TICKET_QUERY_KEY });
};

// Invalidate specific ticket details
export const invalidateTicketDetails = async (queryClient: QueryClient, uuid: string) => {
    await queryClient.invalidateQueries({
        queryKey: [TICKET_DETAILS_QUERY_KEY, uuid]
    });
};

// Invalidate all support ticket related queries
export const invalidateAllTicketQueries = async (queryClient: QueryClient) => {
    await Promise.all([
        queryClient.invalidateQueries({ queryKey: SUPPORT_TICKET_QUERY_KEY }),
        queryClient.invalidateQueries({ queryKey: TICKET_DETAILS_QUERY_KEY })
    ]);
};