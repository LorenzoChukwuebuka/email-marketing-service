import { create } from 'zustand'
import axiosInstance from '../utils/api'
import parseUserAgent from '../utils/userAgent'
import axios from 'axios'
import eventBus from '../utils/eventBus'
import Cookies from 'js-cookie'

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
  forgetPasswordValues: { email: '' },
  resetPasswordValues: {
    password: '',
    confirmPassword: '',
    token: '',
    email: ''
  },
  isLoggedIn: false,
  editFormValues: { fullname: '', company: '', email: '', phonenumber: '' },
  userData: null,

  //initializers
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
  setForgetPasswordValues: newforgetPasswordValues =>
    set({ forgetPasswordValues: newforgetPasswordValues }),
  setResetPasswordValues: newResetPasswordValues =>
    set({ resetPasswordValues: newResetPasswordValues }),
  setIsLoggedIn: newIsLoggedIn => set({ isLoggedIn: newIsLoggedIn }),
  setEditFormValues: newEditFormValues =>
    set({ editFormValues: newEditFormValues }),
  setUserData: newUserData => set({ userData: newUserData }),

  //functions
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
        setUserId(response?.data?.payload?.userId)
        const registeredData = {
          userId: response?.data?.payload?.userId,
          email: formValues.email,
          fullname: formValues.fullname
        }
        setFormValues({
          fullname: '',
          company: '',
          email: '',
          password: '',
          confirmPassword: ''
        })
        return registeredData // Return the necessary data
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
  forgotPass: async () => {
    try {
      const { setIsLoading, forgetPasswordValues, setForgetPasswordValues } =
        get()
      setIsLoading(true)
      let response = await axiosInstance.post(
        '/user-forget-password',
        forgetPasswordValues
      )

      eventBus.emit('message', response.data.payload)
      setForgetPasswordValues({ email: '' })
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
  resetPassword: async () => {
    const { setIsLoading, resetPasswordValues } = get()
    try {
      setIsLoading(true)
      let response = await axiosInstance.post(
        '/user-reset-password',
        resetPasswordValues
      )
      eventBus.emit('success', response.data.payload)
    } catch (error) {
      console.log(error)
    } finally {
      get().setIsLoading(false)
    }
  },
  changePassword: async () => {},
  loginUser: async () => {
    const { loginValues, setIsLoading, setLoginValues, setIsLoggedIn } = get()
    try {
      setIsLoading(true)

      let response = await axiosInstance.post('/user-login', loginValues)
      if (response.data.message === 'success') {
        //save the user Credentials to a cookie

        Cookies.set('Cookies', JSON.stringify(response.data.payload), {
          expires: 7,
          sameSite: 'Strict',
          secure: true
        })

        setIsLoggedIn(true)
      }

      setLoginValues({
        email: '',
        password: ''
      })
    } catch (error) {
      eventBus.emit(
        'error',
        error.response.data.payload || 'An unexpected error occured'
      )

      console.log(error)
    } finally {
      get().setIsLoading(false)
    }
  },

  getUserDetails: async () => {
    try {
      const { setIsLoading, setUserData } = get()
      setIsLoading(true)
      let response = await axiosInstance.get('/get-user-details')
      setUserData(response.data.payload)
    } catch (error) {
      eventBus.emit(
        'error',
        error.response.data.payload || 'An unexpected error occured'
      )
    } finally {
      get().setIsLoading(false)
    }
  },

  editUserDetails: async () => {
    try {
      const { setIsLoading, editFormValues } = get()
      setIsLoading(true)
      let response = await axiosInstance.put(
        '/edit-user-details',
        editFormValues
      )

      eventBus.emit('success', response.data.payload)
    } catch (error) {
      eventBus.emit('error', error)
    } finally {
      get().setIsLoading(false)
    }
  }
}))

export default useAuthStore
