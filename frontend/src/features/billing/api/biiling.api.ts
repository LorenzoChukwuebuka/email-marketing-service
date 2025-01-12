import axiosInstance from "../../../utils/api";
import { PaymentValue, InitializeData, PaymentMethod, BillingAPIResponse } from '../interface/billing.interface';
import { APIResponse, ResponseT } from '../../../interface/api.interface';

class BillingApi {
    static async initializePayment(paymentValues: PaymentValue) {
        const response = await axiosInstance.post<APIResponse<InitializeData>>("/transaction/initialize-transaction", paymentValues)
        return response.data
    }

    static async confirmPayment(paymentMethod: PaymentMethod, reference: string): Promise<ResponseT> {
        const response = await axiosInstance.get<ResponseT>("/transaction/verify-transaction/" + paymentMethod + "/" + reference)
        return response.data
    }

    static async fetchBilling(page?: number, pageSize?: number):Promise<BillingAPIResponse> {
        const response = await axiosInstance.get<BillingAPIResponse>(`/transaction/get-all-billing?page=${page}&page_size=${pageSize}`)
        return response.data
    }
}


export default BillingApi 