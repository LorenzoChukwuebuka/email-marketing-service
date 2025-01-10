import { APIResponse } from "../interface/api.interface";
import { MailData } from "../interface/mailcalc.interface";
import axiosInstance from "../utils/api";

class DailyMailCalculationApi {
    async getDailyMailCalculation(): Promise<APIResponse<MailData>> {
        const response = await axiosInstance.get<APIResponse<MailData>>("/get-user-current-sub")
        return response.data
    }
}

export default new DailyMailCalculationApi()