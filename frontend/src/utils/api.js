import axios from 'axios';
import Cookies from 'js-cookie';

export const APIURL = import.meta.env.VITE_API_URL;

const getToken = () => {
  let cookies = Cookies.get('Cookies');

  if (!cookies) {
    return null;
  }

  try {
    let cookieData = JSON.parse(cookies);
    return cookieData.token;
  } catch (error) {
    console.error('Failed to parse cookies:', error);
    return null;
  }
};

const axiosInstance = axios.create({
  baseURL: APIURL,
  headers: {
    Authorization: `Bearer ${getToken()}`
  }
});

export default axiosInstance;
