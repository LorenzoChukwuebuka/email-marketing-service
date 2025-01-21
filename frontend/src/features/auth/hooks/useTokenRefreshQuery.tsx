import { useQuery, QueryClient } from '@tanstack/react-query';
import Cookies from 'js-cookie';
import axiosInstance from '../../../utils/api';

export const TOKEN_REFRESH_KEY = ['token_refresh'] as const;

interface CookieData {
    token: string;
    refresh_token: string;
    credentials?: any;
}

// Function to update only the access token in cookies
const updateAccessToken = (newAccessToken: string) => {
    const cookies = Cookies.get('Cookies');
    if (cookies) {
        const cookieData: CookieData = JSON.parse(cookies);
        cookieData.token = newAccessToken;

        Cookies.set('Cookies', JSON.stringify(cookieData), {
            expires: 1,
            path: '/',
            sameSite: 'Strict',
            secure: true
        });
    }
};

// Token refresh function
const refreshToken = async () => {
    const cookies = Cookies.get('Cookies');
    if (!cookies) throw new Error('No refresh token available');
    

    const { refresh_token } = JSON.parse(cookies) as CookieData;

    const data = {
        refresh_token: refresh_token
    }

    try {
        const response = await axiosInstance.post('/refresh-access-token',
            data
        );

      //  console.log(response.data.payload.access_token)

        if (response.data?.payload?.access_token) {
            updateAccessToken(response.data.payload.access_token);
            return response.data.payload;
        }
    } catch (error) {
        // If refresh fails, throw error for query to handle
        console.log(error)
        throw new Error('Token refresh failed');
    }
};

// Token refresh query options
export const tokenRefreshOptions = {
    queryKey: TOKEN_REFRESH_KEY,
    queryFn: refreshToken,
    refetchInterval: 5 * 60 * 1000, // Refetch every 5 minutes
    refetchIntervalInBackground: true,
    retry: 10,
    enabled: !!Cookies.get('Cookies'), // Only run if we have cookies
};

// Custom hook for token refresh
export const useTokenRefresh = () => {
    return useQuery(tokenRefreshOptions);
};

// Helper function to manually trigger a token refresh
export const triggerTokenRefresh = async (queryClient: QueryClient) => {
    await queryClient.invalidateQueries({ queryKey: TOKEN_REFRESH_KEY });
};