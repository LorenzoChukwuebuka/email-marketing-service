import { APIResponse, ResponseT } from '../../../interface/api.interface';
import axiosInstance from '../../../utils/api';
import { ChangePasswordValues, EditFormValues, ForgetPasswordValues, FormValues, LoginValues, OtpValue, ResetPasswordValues, SignUpAPIData, UserDetails } from '../interface/auth.interface';
class AuthApi {

    private static instance: AuthApi;
    private baseurl = "/auth"
    private baseurl2 = "/users"

    static getInstance(): AuthApi {
        if (!AuthApi.instance) {
            AuthApi.instance = new AuthApi();
        }
        return AuthApi.instance;
    }

    async registerUser(formValues: Partial<FormValues>): Promise<APIResponse<SignUpAPIData>> {
        const response = await axiosInstance.post<APIResponse<SignUpAPIData>>(`${this.baseurl}/signup`, formValues);
        return response.data;
    }

    async verifyUser(otpValue: OtpValue): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.post<APIResponse<ResponseT>>(`${this.baseurl}/verify`, otpValue);
        return response.data;
    }

    async resendOTP(data: Record<string, string>): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.post("/resend-otp", data);
        return response.data;
    }

    async forgotPass(forgetPasswordValues: ForgetPasswordValues): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.post(`${this.baseurl}/forget-password`, forgetPasswordValues);
        return response.data;
    }

    async resetPassword(resetPasswordValues: ResetPasswordValues): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.post(`${this.baseurl}/reset-password`, resetPasswordValues);
        return response.data;
    }

    async changePassword(changePasswordValues: ChangePasswordValues): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.put(`${this.baseurl}/change-password`, changePasswordValues);
        return response.data;
    }

    async loginUser(loginValues: LoginValues): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.post(`${this.baseurl}/login`, loginValues);
        return response.data;
    }

    async getUserDetails(): Promise<APIResponse<UserDetails>> {
        const response = await axiosInstance.get<APIResponse<UserDetails>>(`${this.baseurl2}/details`);
        return response.data;
    }

    async editUserDetails(editFormValues: EditFormValues): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.put("/edit-user-details", editFormValues);
        return response.data;
    }

    async deleteUserAccount() {
        const response = await axiosInstance.delete(`${this.baseurl2}/delete-user`)
        return response.data
    }

    async cancelDeleteUserAccount() {
        const response = await axiosInstance.put(`${this.baseurl2}/cancel-delete`)
        return response.data
    }

    async getGoogleAuthLoginDetails(key: string): Promise<APIResponse<UserDetails>> {
        const response = await axiosInstance.get("/google-login-details?key=" + key)
        return response.data
    }
}

const authApi = AuthApi.getInstance();

export default authApi