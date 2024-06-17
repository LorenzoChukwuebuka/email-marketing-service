import axios from 'axios'
import Cookies from 'js-cookie'

export const APIURL = import.meta.env.VITE_API_URL

const getToken = () => {
  const cookies = Cookies.get('Cookies')
  return cookies ? cookies.token : null
}

const axiosInstance = axios.create({
  baseURL: APIURL,
  headers: {
    Authorization: `Bearer ${getToken()}`
  }
})

export default axiosInstance
