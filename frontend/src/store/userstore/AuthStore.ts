import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import parseUserAgent from '../../utils/userAgent'
import axios from 'axios'
import eventBus from '../../utils/eventBus'
import Cookies from 'js-cookie'
import { BaseEntity } from '../../interface/baseentity.interface'
import { APIResponse } from '../../interface/api.interface'
import { errResponse } from '../../utils/isError'

type FormValues = {
    fullname: string;
    company: string;
    email: string;
    password: string;
    confirmPassword: string;
    token: string;
    phonenumber: string;
}

type LoginValues = Pick<FormValues, 'email' | 'password'>

type ForgetPasswordValues = Pick<FormValues, 'email'>

type ResetPasswordValues = Pick<FormValues, 'email' | 'confirmPassword' | 'password' | 'token'>

type EditFormValues = Omit<FormValues, 'password' | 'confirmPassword' | 'token'>;

type UserDetails = {
    fullname: string
    email: string
    company: string
    phonenumber: string
    verified: boolean,
    blocked: boolean,
} & BaseEntity

type ChangePasswordValues = {
    old_password: string;
    new_password: string;
    confirm_password: string;
}
type AuthStore = {
    formValues: Omit<FormValues, 'token' | 'phonenumber'>;
    error: boolean;
    success: boolean;
    errorMessage: string;
    successMessage: string;
    formError: boolean;
    isLoading: boolean;
    redirectToOTP: boolean;
    otpValue: { token: string };
    isVerified: boolean;
    userId: string;
    loginValues: LoginValues;
    forgetPasswordValues: ForgetPasswordValues;
    resetPasswordValues: ResetPasswordValues;
    isLoggedIn: boolean;
    editFormValues: EditFormValues;
    userData: any;
    changePasswordValues: ChangePasswordValues;

    setFormValues: (newFormValues: Omit<FormValues, 'token' | 'phonenumber'>) => void;
    setError: (newError: boolean) => void;
    setSuccess: (newSuccess: boolean) => void;
    setErrorMessage: (newErrorMessage: string) => void;
    setSuccessMessage: (newSuccessMessage: string) => void;
    setUserId: (newUserId: string) => void;
    setIsLoading: (newIsLoading: boolean) => void;
    setRedirectToOTP: (newRedirectOTP: boolean) => void;
    setOTPValue: (newOtpValue: { token: string }) => void;
    setIsVerified: (newVerification: boolean) => void;
    setLoginValues: (newLoginValues: LoginValues) => void;
    setForgetPasswordValues: (newforgetPasswordValues: ForgetPasswordValues) => void;
    setResetPasswordValues: (newResetPasswordValues: ResetPasswordValues) => void;
    setIsLoggedIn: (newIsLoggedIn: boolean) => void;
    setEditFormValues: (newEditFormValues: EditFormValues) => void;
    setUserData: (newUserData: any) => void;
    setChangePasswordValues: (newChangePasswordValues: ChangePasswordValues) => void;

    registerUser: () => Promise<{ userId: string; email: string; fullname: string } | undefined>;
    createUserSession: (userId: string) => Promise<void>;
    verifyUser: () => Promise<void>;
    resendOTP: (data: any) => Promise<void>;
    forgotPass: () => Promise<void>;
    resetPassword: () => Promise<void>;
    changePassword: () => Promise<void>;
    loginUser: () => Promise<void>;
    getUserDetails: () => Promise<void>;
    editUserDetails: () => Promise<void>;
}

const useAuthStore = create<AuthStore>((set, get) => ({
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
    changePasswordValues: {
        old_password: '',
        new_password: '',
        confirm_password: ''
    },

    //initializers
    setFormValues: (newFormValues) => set({ formValues: newFormValues }),
    setError: (newError) => set({ error: newError }),
    setSuccess: (newSuccess) => set({ success: newSuccess }),
    setErrorMessage: (newErrorMessage) => set({ errorMessage: newErrorMessage }),
    setSuccessMessage: (newSuccessMessage) =>
        set({ successMessage: newSuccessMessage }),
    setUserId: (newUserId) => set({ userId: newUserId }),
    setIsLoading: (newIsLoading) => set({ isLoading: newIsLoading }),
    setRedirectToOTP: (newRedirectOTP) => set({ redirectToOTP: newRedirectOTP }),
    setOTPValue: (newOtpValue) => set({ otpValue: newOtpValue }),
    setIsVerified: (newVerification) => set({ isVerified: newVerification }),
    setLoginValues: (newLoginValues) => set({ loginValues: newLoginValues }),
    setForgetPasswordValues: (newforgetPasswordValues) =>
        set({ forgetPasswordValues: newforgetPasswordValues }),
    setResetPasswordValues: (newResetPasswordValues) =>
        set({ resetPasswordValues: newResetPasswordValues }),
    setIsLoggedIn: (newIsLoggedIn) => set({ isLoggedIn: newIsLoggedIn }),
    setEditFormValues: (newEditFormValues) =>
        set({ editFormValues: newEditFormValues }),
    setUserData: (newUserData) => set({ userData: newUserData }),
    setChangePasswordValues: (newChangePasswordValues) =>
        set({ changePasswordValues: newChangePasswordValues }),

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
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            get().setIsLoading(false)
        }
    },
    createUserSession: async (userId: string) => {
        const userAgent = navigator.userAgent
        const browserInfo = parseUserAgent(userAgent)

        let response = await axios.get<Record<string, string>>('https://api.ipify.org/?format=json')

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
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }

        } finally {
            get().setIsLoading(false)
        }
    },

    resendOTP: async (data: any) => {
        const { setIsLoading } = get()
        try {
            setIsLoading(true)
            let response = await axiosInstance.post('/resend-otp', data)
            console.log(response)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
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
        } catch (error: any) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
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
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            get().setIsLoading(false)
        }
    },
    changePassword: async () => {
        try {
            const { setIsLoading, changePasswordValues } = get()
            setIsLoading(true)
            let response = await axiosInstance.put(
                '/change-user-password',
                changePasswordValues
            )

            eventBus.emit('success', response.data.payload)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            get().setIsLoading(false)
        }
    },
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
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            get().setIsLoading(false)
        }
    },

    getUserDetails: async (): Promise<void> => {
        try {
            const { setIsLoading, setUserData } = get()
            setIsLoading(true)
            let response = await axiosInstance.get<APIResponse<UserDetails>>('/get-user-details')
            setUserData(response.data.payload)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            get().setIsLoading(false)
        }
    },

    editUserDetails: async (): Promise<void> => {
        try {
            const { setIsLoading, editFormValues } = get()
            setIsLoading(true)
            let response = await axiosInstance.put(
                '/edit-user-details',
                editFormValues
            )

            eventBus.emit('success', response.data.payload)
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            get().setIsLoading(false)
        }
    }
}))

export default useAuthStore