import { QueryClient, useQuery } from "@tanstack/react-query";
import adminSupportApi from "../api/adminsupport.api";
import { APIResponse } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';
import { Ticket } from '../interface/support.interface';

/***************** Support Tickets Query Keys *******************/
export const ADMIN_SUPPORT_TICKETS_QUERY_KEY = ["support_tickets"] as const;

// Reusable query options for all APIs
const createSupportTicketsQueryOptions = (
    endpoint: string,
    fetchFn: () => Promise<any>,
    dependencies: any[] = []
) => ({
    queryKey: [ADMIN_SUPPORT_TICKETS_QUERY_KEY, endpoint, ...dependencies],
    queryFn: fetchFn,
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 15 * 60 * 1000, // 15 minutes
    retry: 3,
    refetchOnWindowFocus: true,
});

/***************** Hooks for APIs *******************/

// Get all closed tickets
export const useAdminClosedTicketsQuery = (
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return useQuery<APIResponse<PaginatedResponse<Ticket>>>(
        createSupportTicketsQueryOptions(
            "closed_tickets",
            () => adminSupportApi.getAllClosedTickets(page, pageSize, query),
            [page, pageSize, query]
        )
    );
};

// Get all tickets
export const useAdminAllTicketsQuery = (
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return useQuery<APIResponse<PaginatedResponse<Ticket>>>(
        createSupportTicketsQueryOptions(
            "all_tickets",
            () => adminSupportApi.getAllTickets(page, pageSize, query),
            [page, pageSize, query]
        )
    );
};

// Get all pending tickets
export const useAdminPendingTicketsQuery = (
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return useQuery<APIResponse<PaginatedResponse<Ticket>>>(
        createSupportTicketsQueryOptions(
            "pending_tickets",
            () => adminSupportApi.getAllPendingTickets(page, pageSize, query),
            [page, pageSize, query]
        )
    );
};

/***************** Query Invalidations *******************/

export const invalidateSupportTicketsQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: ADMIN_SUPPORT_TICKETS_QUERY_KEY });
};

/***************** Prefetch Functions *******************/

export const prefetchClosedTicketsQuery = (
    queryClient: QueryClient,
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return queryClient.prefetchQuery(
        createSupportTicketsQueryOptions(
            "closed_tickets",
            () => adminSupportApi.getAllClosedTickets(page, pageSize, query),
            [page, pageSize, query]
        )
    );
};

export const prefetchAllTicketsQuery = (
    queryClient: QueryClient,
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return queryClient.prefetchQuery(
        createSupportTicketsQueryOptions(
            "all_tickets",
            () => adminSupportApi.getAllTickets(page, pageSize, query),
            [page, pageSize, query]
        )
    );
};

export const prefetchPendingTicketsQuery = (
    queryClient: QueryClient,
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return queryClient.prefetchQuery(
        createSupportTicketsQueryOptions(
            "pending_tickets",
            () => adminSupportApi.getAllPendingTickets(page, pageSize, query),
            [page, pageSize, query]
        )
    );
};