import { APIResponse, ResponseT } from '../../../interface/api.interface';
import axiosInstance from '../../../utils/api';
import { ChangePasswordValues, EditFormValues, ForgetPasswordValues, FormValues, LoginValues, OtpValue, ResetPasswordValues, SignUpAPIData, UserDetails } from '../interface/auth.interface';


class AuthApi {

    async registerUser(formValues: Partial<FormValues>): Promise<APIResponse<SignUpAPIData>> {
        const response = await axiosInstance.post<APIResponse<SignUpAPIData>>("/user-signup", formValues);
        return response.data;
    }

    // async createUserSession(userId: string): Promise<any> {
    //     try {
    //         const userAgent = navigator.userAgent;
    //         const browserInfo = parseUserAgent(userAgent);
    //         const ipResponse = await axios.get<{ ip: string }>("https://api.ipify.org/?format=json");

    //         const userDevice = {
    //             user_id: userId,
    //             device: navigator.platform,
    //             ip_address: ipResponse.data.ip,
    //             browser: browserInfo.name,
    //         };

    //         const response = await axiosInstance.post("/create-session", userDevice);
    //         return response.data;
    //     } catch (error) {
    //       ;
    //     }
    // }

    async verifyUser(otpValue: OtpValue): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.post<APIResponse<ResponseT>>("/verify-user", otpValue);
        return response.data;
    }

    async resendOTP(data: Record<string, string>): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.post("/resend-otp", data);
        return response.data;
    }

    async forgotPass(forgetPasswordValues: ForgetPasswordValues): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.post("/user-forget-password", forgetPasswordValues);
        return response.data;
    }

    async resetPassword(resetPasswordValues: ResetPasswordValues): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.post("/user-reset-password", resetPasswordValues);
        return response.data;
    }

    async changePassword(changePasswordValues: ChangePasswordValues): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.put("/change-user-password", changePasswordValues);
        return response.data;
    }

    async loginUser(loginValues: LoginValues): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.post("/user-login", loginValues);
        return response.data;
    }

    async getUserDetails(): Promise<APIResponse<UserDetails>> {
        const response = await axiosInstance.get<APIResponse<UserDetails>>("/get-user-details");
        return response.data;
    }

    async editUserDetails(editFormValues: EditFormValues): Promise<APIResponse<ResponseT>> {
        const response = await axiosInstance.put("/edit-user-details", editFormValues);
        return response.data;
    }

    async deleteUserAccount() {
        const response = await axiosInstance.delete("/delete-user")
        return response.data
    }

    async cancelDeleteUserAccount() {
        const response = await axiosInstance.put("/cancel-delete")
        return response.data
    }

    async getGoogleAuthLoginDetails(key: string): Promise<APIResponse<UserDetails>> {
        const response = await axiosInstance.get("/google-login-details?key=" + key)
        return response.data
    }
}

export default new AuthApi()