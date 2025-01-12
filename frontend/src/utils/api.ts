import axios, { AxiosInstance } from 'axios';
import Cookies from 'js-cookie';

export const APIURL: string = import.meta.env.VITE_API_URL as string;

interface CookieData {
    token: string;
    // Add other properties if needed
}

const getToken = (): string | null => {
    const cookies = Cookies.get('Cookies');

    if (!cookies) {
        return null;
    }

    try {
        const cookieData: CookieData = JSON.parse(cookies);
        return cookieData.token;
    } catch (error) {
        console.error('Failed to parse cookies:', error);
        return null;
    }
};

const axiosInstance: AxiosInstance = axios.create({
    baseURL: APIURL,
    headers: {
        Authorization: `Bearer ${getToken()}`
    }
});


axiosInstance.interceptors.request.use(
    (config) => {
        const token = getToken();
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);


export default axiosInstance;