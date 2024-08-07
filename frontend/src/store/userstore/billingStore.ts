import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { errResponse } from '../../utils/isError';
import { APIResponse } from '../../interface/api.interface';

type PaymentMethod = "Paystack" | "FlutterWave"

type PaymentValue = {
    plan_id: string
    amount_to_pay: number
    duration: string
    payment_method: PaymentMethod
}

export type InitializeData = {
    data: {
        access_code: string
        authorization_url: string
        reference: string
    }
}

type BillingStore = {
    paymentValues: PaymentValue
    intializeData: InitializeData
    isLoading: boolean
    setPaymentValues: (newPaymentValues: PaymentValue) => void;
    setInitalizeData: (newData: InitializeData) => void
    setIsLoading: (newIsLoading: boolean) => void;
    initializePayment: () => Promise<void>
    confirmPayment: () => Promise<void>
}

const useBillingStore = create<BillingStore>((set, get) => ({
    isLoading: false,
    paymentValues: {
        plan_id: "",
        amount_to_pay: 0,
        duration: "",
        payment_method: "Paystack" as PaymentMethod
    },
    intializeData: {
        data: {
            access_code: "",
            authorization_url: "",
            reference: ""
        }
    },
    setPaymentValues: (newPaymentValues: PaymentValue) => set({ paymentValues: newPaymentValues }),
    setInitalizeData: (newData: InitializeData) => set({ intializeData: newData }),
    setIsLoading: (newIsLoading) => set({ isLoading: newIsLoading }),
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
    confirmPayment: async () => {
        try {

        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        }
    }
}));


export default useBillingStore

