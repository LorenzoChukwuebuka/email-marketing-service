import { create } from 'zustand';
import { FormValues, LoginValues, ForgetPasswordValues, ResetPasswordValues, EditFormValues, ChangePasswordValues, OtpValue, VerifyLoginFormData } from '../interface/auth.interface';
import authApi from '../api/auth.api';
import { persist } from 'zustand/middleware';
import Cookies from 'js-cookie';
import { handleError } from '../../../utils/isError';
import eventBus from '../../../utils/eventbus';

interface AuthState {
    formValues: Omit<FormValues, 'token' | 'phonenumber'>;
    redirectToOTP: boolean;
    otpValue: OtpValue;
    isVerified: boolean;
    userId: string;
    isLoginVerified: boolean
    loginValues: LoginValues;
    forgetPasswordValues: ForgetPasswordValues;
    resetPasswordValues: ResetPasswordValues;
    isLoggedIn: boolean;
    editFormValues: EditFormValues;
    //  userData: any;
    token: string | undefined;
    changePasswordValues: ChangePasswordValues;
    verifyloginValues: VerifyLoginFormData
}

interface AuthAsyncActions {
    registerUser: () => Promise<void>;
    createUserSession: (userId: string) => Promise<void>;
    verifyUser: () => Promise<void>;
    resendOTP: (data: Record<string, string>) => Promise<void>;
    forgotPass: () => Promise<void>;
    resetPassword: () => Promise<void>;
    changePassword: () => Promise<void>;
    loginUser: () => Promise<void>;
    editUserDetails: () => Promise<void>;
    deleteUser: () => Promise<void>
    cancelDelete: () => Promise<void>
    verifyLogin: () => Promise<void>
}

interface AuthActions {
    setFormValues: (newFormValues: Omit<FormValues, 'token' | 'phonenumber'>) => void;
    setUserId: (newUserId: string) => void;
    setRedirectToOTP: (newRedirectOTP: boolean) => void;
    setOTPValue: (newOtpValue: { token: string }) => void;
    setIsVerified: (newVerification: boolean) => void;
    setLoginValues: (newLoginValues: LoginValues) => void;
    setForgetPasswordValues: (newforgetPasswordValues: ForgetPasswordValues) => void;
    setResetPasswordValues: (newResetPasswordValues: ResetPasswordValues) => void;
    setIsLoggedIn: (newIsLoggedIn: boolean) => void;
    setEditFormValues: (newEditFormValues: EditFormValues) => void;
    //  setUserData: (newUserData: unknown) => void;
    setChangePasswordValues: (newChangePasswordValues: ChangePasswordValues) => void;
    setVerifyLoginValues: (newVerifyLoginValues: VerifyLoginFormData) => void;
}

type AuthStore = AuthState & AuthActions & AuthAsyncActions

const InitialState: AuthState = {
    formValues: {
        fullname: '',
        company: '',
        email: '',
        password: '',
        confirmPassword: ''
    },
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
    // userData: null,
    changePasswordValues: {
        old_password: '',
        new_password: '',
        confirm_password: ''
    },
    token: Cookies.get('Cookies'),
    verifyloginValues: {
        user_id: '',
        token: ''
    },
    isLoginVerified: false
}

const useAuthStore = create<AuthStore>()(
    persist(
        (set, get) => ({
            ...InitialState,

            //initializers
            setFormValues: (newFormValues) => set({ formValues: newFormValues }),
            setUserId: (newUserId) => set({ userId: newUserId }),
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
            // setUserData: (newUserData) => set({ userData: newUserData }),
            setChangePasswordValues: (newChangePasswordValues) =>
                set({ changePasswordValues: newChangePasswordValues }),
            setVerifyLoginValues: (newVerifyLoginValues) => set({ verifyloginValues: newVerifyLoginValues }),


            registerUser: async () => {
                const {
                    formValues,
                    setUserId
                } = get()
                try {
                    const response = await authApi.registerUser(formValues)
                    setUserId(response?.payload?.userId || '')
                } catch (error) {
                    handleError(error)
                }
            },

            createUserSession: async (userId: string) => {
                try {
                    // Call your API here
                    console.log(`Creating session for user ID: ${userId}`);
                } catch (error) {
                    console.error("Error creating user session:", error);
                }
            },

            verifyUser: async () => {
                const { otpValue } = get()
                try {
                    const response = await authApi.verifyUser(otpValue)
                    if (response) {
                        eventBus.emit('success', 'Account verified successfully. Please login to continue.')
                        await new Promise((resolve) => setTimeout(resolve, 500))
                        set({ isVerified: true })
                    }
                } catch (error) {
                    handleError(error)
                }
            },

            resendOTP: async (data: Record<string, string>) => {
                try {
                    await authApi.resendOTP(data)
                } catch (error) {
                    console.error("Error resending OTP:", error);
                }
            },

            forgotPass: async () => {
                const { forgetPasswordValues } = get()
                try {
                    await authApi.forgotPass(forgetPasswordValues)
                } catch (error) {
                    handleError(error)
                }
            },

            resetPassword: async () => {
                try {
                    // Call your API here
                    console.log("Resetting password...");
                } catch (error) {
                    console.error("Error resetting password:", error);
                }
            },

            changePassword: async () => {
                const { changePasswordValues } = get()
                try {
                    const response = await authApi.changePassword(changePasswordValues)
                    if (response) {
                        eventBus.emit('success', 'Password changed successfully... Redirecting to login')
                    }
                    new Promise(resolve => setTimeout(resolve, 2000))

                    window.location.href = '/auth/login'
                } catch (error) {
                    console.error("Error changing password:", error);
                }
            },

            loginUser: async () => {
                const { loginValues } = get()
                try {
                    const response = await authApi.loginUser(loginValues)
                    console.log(response)
                    //@ts-ignore
                    set({ userId: response?.payload?.details?.id })

                    window.location.href = '/auth/verify-login'

                } catch (error) {
                    handleError(error)
                }
            },

            verifyLogin: async () => {
                const { verifyloginValues, userId } = get()
                let formatData = {
                    ...verifyloginValues,
                    userId: userId
                }
                let response = await authApi.verifyLogin(formatData as any)
                if (response) {
                    Cookies.set('Cookies', JSON.stringify(response.payload), {
                        expires: 3,
                        sameSite: 'Strict',
                        secure: true
                    })
                    window.location.href = '/app'
                }
            },

            editUserDetails: async () => {
                const { editFormValues } = get()
                try {
                    const response = await authApi.editUserDetails(editFormValues)
                    if (response) {
                        eventBus.emit('success', 'User details edited successfully')
                    }
                } catch (error) {
                    console.error("Error editing user details:", error);
                }
            },

            deleteUser: async () => {
                try {
                    const response = await authApi.deleteUserAccount()
                    if (response) {
                        eventBus.emit('success', 'Your account has been scheduled for deletion')
                    }
                } catch (error) {
                    handleError(error)
                }
            },
            cancelDelete: async () => {
                try {
                    const response = await authApi.cancelDeleteUserAccount()
                    if (response) {
                        eventBus.emit('success', 'The delete proces has been cancelled?')
                    }
                } catch (error) {
                    handleError(error)
                }
            }
        }),

        {
            name: 'auth-store',
            partialize: (state) => ({
                userId: state.userId,
            })
        }
    )


)


export default useAuthStore;