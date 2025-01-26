import { QueryClient, useQuery } from "@tanstack/react-query";
import adminusersApi from "../api/adminusers.api";
import { APIResponse } from '../../../interface/api.interface';
import { PaginatedResponse } from '../../../interface/pagination.interface';
import { AdminUserDetails, AdminUserStats } from "../interface/adminusers.interface";

/***************** Admin Users Query Keys *******************/
export const ADMIN_USERS_QUERY_KEY = ["admin_users"] as const;

// Reusable query options for all APIs
const createAdminUsersQueryOptions = (
    endpoint: string,
    fetchFn: () => Promise<any>,
    dependencies: any[] = []
) => ({
    queryKey: [ADMIN_USERS_QUERY_KEY, endpoint, ...dependencies],
    queryFn: fetchFn,
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 15 * 60 * 1000, // 15 minutes
    retry: 3,
    refetchOnWindowFocus: true,
});

/***************** Hooks for APIs *******************/

// Get all users
export const useAllUsersQuery = (
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return useQuery<APIResponse<PaginatedResponse<AdminUserDetails>>>(
        createAdminUsersQueryOptions(
            "all_users",
            () => adminusersApi.getAllUsers(page, pageSize, query),
            [page, pageSize, query]
        )
    );
};

// Get verified users
export const useVerifiedUsersQuery = (
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return useQuery<APIResponse<PaginatedResponse<AdminUserDetails>>>(
        createAdminUsersQueryOptions(
            "verified_users",
            () => adminusersApi.getVerifiedUsers(page, pageSize, query),
            [page, pageSize, query]
        )
    );
};

// Get unverified users
export const useUnverifiedUsersQuery = (
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return useQuery<APIResponse<PaginatedResponse<AdminUserDetails>>>(
        createAdminUsersQueryOptions(
            "unverified_users",
            () => adminusersApi.getUnverifiedUsers(page, pageSize, query),
            [page, pageSize, query]
        )
    );
};

// Get single user
export const useSingleUserQuery = (userId: string) => {
    return useQuery<APIResponse<AdminUserDetails>>(
        createAdminUsersQueryOptions(
            "single_user",
            () => adminusersApi.getSingleUser(userId),
            [userId]
        )
    );
};

// Get user stats
export const useUserStatsQuery = (userId: string) => {
    return useQuery<APIResponse<AdminUserStats>>(
        createAdminUsersQueryOptions(
            "user_stats",
            () => adminusersApi.getUserStats(userId),
            [userId]
        )
    );
};

/***************** Query Invalidations *******************/

export const invalidateAdminUsersQuery = (queryClient: QueryClient) => {
    return queryClient.invalidateQueries({ queryKey: ADMIN_USERS_QUERY_KEY });
};

export const invalidateSingleUserQuery = (queryClient: QueryClient, userId: string) => {
    return queryClient.invalidateQueries({ queryKey: [ADMIN_USERS_QUERY_KEY, "single_user", userId] });
};

export const invalidateUserStatsQuery = (queryClient: QueryClient, userId: string) => {
    return queryClient.invalidateQueries({ queryKey: [ADMIN_USERS_QUERY_KEY, "user_stats", userId] });
};



/***************** Prefetch Functions *******************/

export const prefetchAllUsersQuery = (
    queryClient: QueryClient,
    page?: number,
    pageSize?: number,
    query?: string
) => {
    return queryClient.prefetchQuery(
        createAdminUsersQueryOptions(
            "all_users",
            () => adminusersApi.getAllUsers(page, pageSize, query),
            [page, pageSize, query]
        )
    );
};

export const prefetchSingleUserQuery = (
    queryClient: QueryClient,
    userId: string
) => {
    return queryClient.prefetchQuery(
        createAdminUsersQueryOptions(
            "single_user",
            () => adminusersApi.getSingleUser(userId),
            [userId]
        )
    );
};
