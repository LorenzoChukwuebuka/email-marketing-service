import axios from 'axios'
import Cookies from 'js-cookie'

export const APIURL = import.meta.env.VITE_API_URL

const getToken = () => {
  let cookies = Cookies.get('Cookies')

  let cookieData = JSON.parse(cookies)

  return cookies ? cookieData.token : null
}

const axiosInstance = axios.create({
  baseURL: APIURL,
  headers: {
    Authorization: `Bearer ${getToken()}`
  }
})

export default axiosInstance
