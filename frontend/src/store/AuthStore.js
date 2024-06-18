import { create } from 'zustand'
import axiosInstance from '../utils/api'
import parseUserAgent from '../utils/userAgent'
import axios from 'axios'
import eventBus from '../utils/eventBus'
const useAuthStore = create((set, get) => ({
  formValues: {
    fullname: '',
    company: '',
    email: '',
    password: '',
    confirmPassword: ''
  },
  error: false,
  success: false,
  errorMessage: '',
  successMessage: '',
  formError: false,
  isLoading: false,
  redirectToOTP: false,
  otpValue: { token: '' },
  isVerified: false,
  userId: '',
  loginValues: {
    email: '',
    password: ''
  },

  setFormValues: newFormValues => set({ formValues: newFormValues }),
  setError: newError => set({ error: newError }),
  setSuccess: newSuccess => set({ success: newSuccess }),
  setErrorMessage: newErrorMessage => set({ errorMessage: newErrorMessage }),
  setSuccessMessage: newSuccessMessage =>
    set({ successMessage: newSuccessMessage }),
  setUserId: newUserId => set({ userId: newUserId }),
  setIsLoading: newIsLoading => set({ isLoading: newIsLoading }),
  setRedirectToOTP: newRedirectOTP => set({ redirectToOTP: newRedirectOTP }),
  setOTPValue: newOtpValue => set({ otpValue: newOtpValue }),
  setIsVerified: newVerification => set({ isVerified: newVerification }),
  setLoginValues: newLoginValues => set({ loginValues: newLoginValues }),

  registerUser: async () => {
    const {
      setIsLoading,
      setFormValues,
      setRedirectToOTP,
      createUserSession,
      setUserId
    } = get()
    try {
      setIsLoading(true)
      const { formValues } = get()
      let response = await axiosInstance.post('/user-signup', formValues)
      if (response.data.status === true) {
        setRedirectToOTP(true)
        await createUserSession(response.data.payload.userId)
        setUserId(response.data.payload.userId)
        setFormValues({
          fullname: '',
          company: '',
          email: '',
          password: '',
          confirmPassword: ''
        })
      }
    } catch (error) {
      console.log(error)
    } finally {
      get().setIsLoading(false)
    }
  },

  createUserSession: async userId => {
    const userAgent = navigator.userAgent
    const browserInfo = parseUserAgent(userAgent)

    let response = await axios.get('https://api.ipify.org/?format=json')

    let userDevice = {
      user_id: userId,
      device: navigator.platform,
      ip_address: response.data.ip,
      browser: browserInfo.name
    }

    await axiosInstance.post('/create-session', userDevice)
  },
  verifyUser: async () => {
    const { otpValue, setIsLoading, setIsVerified } = get()

    try {
      setIsLoading(true)
      let response = await axiosInstance.post('/verify-user', otpValue)
      if (response.data.status === true) {
        setIsVerified(true)
        eventBus.emit(
          'success',
          'Your account has been successfully verified. Redirecting to login'
        )
      }
    } catch (error) {
      eventBus.emit(
        'error',
        error.response.data.payload || 'An unexpected error occured'
      )
    } finally {
      get().setIsLoading(false)
    }
  },

  resendOTP: async data => {
    const { setIsLoading } = get()
    try {
      setIsLoading(true)
      let response = await axiosInstance.post('/resend-otp', data)
      console.log(response)
    } catch (error) {
      console.log(error)
      eventBus.emit(
        'error',
        error.response.data.payload || 'An unexpected error occured'
      )
    } finally {
      get().setIsLoading(false)
    }
  },
  forgotPass: async () => {},
  resetPassword: async () => {},
  changePassword: async () => {},
  loginUser: async () => {
    const { loginValues, setIsLoading } = get()
    try {
      setIsLoading(true)

      let response = await axiosInstance.post('/user-login', loginValues)
      console.log(response)
    } catch (error) {
      eventBus.emit(
        'error',
        error.response.data.payload || 'An unexpected error occured'
      )
    } finally {
      get().setIsLoading(false)
    }
  }
}))

export default useAuthStore
