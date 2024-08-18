import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { errResponse } from '../../utils/isError';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { PlanData } from '../admin/planStore';
import { BaseEntity } from '../../interface/baseentity.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';

type PaymentMethod = "Paystack" | "FlutterWave"

type PaymentValue = {
    plan_id: string
    amount_to_pay: number
    duration: string
    payment_method: PaymentMethod
}

export interface User extends BaseEntity {
    fullname: string;
    email: string;
    company: string;
    phonenumber: string;
    verified: boolean;
    blocked: boolean;
    verified_at: string | null;
}

export interface BillingData extends BaseEntity {
    user_id: number;
    amount_paid: number;
    plan_id: number;
    duration: string;
    expiry_date: string;
    reference: string;
    transaction_id: string;
    payment_method: string;
    status: string;
    user: User;
    plan: PlanData;
}

export type InitializeData = {
    data: {
        access_code: string
        authorization_url: string
        reference: string
    }
}


type BillingAPIResponse = APIResponse<PaginatedResponse<BillingData>>;

type BillingStore = {
    paymentValues: PaymentValue
    intializeData: InitializeData
    isLoading: boolean
    paginationInfo: Omit<PaginatedResponse<BillingData>, 'data'>;
    billingData: BillingData[]
    setPaymentValues: (newPaymentValues: PaymentValue) => void;
    setInitalizeData: (newData: InitializeData) => void
    setIsLoading: (newIsLoading: boolean) => void;
    setBillingData: (newBillingData: BillingData[]) => void
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<BillingData>, 'data'>) => void;
    initializePayment: () => Promise<void>
    confirmPayment: (reference: string, paymentMethod: string) => Promise<void>
    fetchBillingData: (page?: number, pageSize?: number) => Promise<void>
}

const useBillingStore = create<BillingStore>((set, get) => ({
    isLoading: false,
    paymentValues: {
        plan_id: "",
        amount_to_pay: 0,
        duration: "",
        payment_method: "Paystack" as PaymentMethod
    },
    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },
    intializeData: {
        data: {
            access_code: "",
            authorization_url: "",
            reference: ""
        }
    },
    billingData: [],
    setPaymentValues: (newPaymentValues: PaymentValue) => set({ paymentValues: newPaymentValues }),
    setBillingData: (newData: BillingData[]) => set({ billingData: newData }),
    setInitalizeData: (newData: InitializeData) => set({ intializeData: newData }),
    setIsLoading: (newIsLoading) => set({ isLoading: newIsLoading }),
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),
    initializePayment: async () => {
        try {
            const { paymentValues } = get()

            let response = await axiosInstance.post<APIResponse<InitializeData>>("/initialize-transaction", paymentValues)
            if (response.data.status == true) {
                window.location.href = response.data.payload.data.authorization_url
            }
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },
    confirmPayment: async (reference: string, paymentMethod: string) => {
        try {
            let response = await axiosInstance.get<ResponseT>("/verify-transaction/" + paymentMethod + "/" + reference)
            if (response.data.status == true) {
                window.location.href = "/user/dash/billing"
                eventBus.emit('success', response.data.payload.data)
            }
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    },
    fetchBillingData: async (page = 1, pageSize = 10) => {
        try {
            const { setBillingData, setPaginationInfo,setIsLoading } = get()
            setIsLoading(true)
            let response = await axiosInstance.get<BillingAPIResponse>(`/get-all-billing?page=${page}&page_size=${pageSize}`)
            const { data, ...paginationInfo } = response.data.payload;
            setPaginationInfo(paginationInfo);
            setBillingData(data)
            
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }finally{
            get().setIsLoading(false)
        }
    }
}));


export default useBillingStore

