import { create } from 'zustand';
import { PaymentValue, InitializeData, PaymentMethod } from '../interface/billing.interface';
import { handleError } from '../../../utils/isError';
import eventBus from '../../../utils/eventbus';
import BillingApi from '../api/biiling.api';

interface BillingState {
    paymentValues: PaymentValue
    intializeData: InitializeData
}

interface BillingActions {
    setPaymentValues: (newPaymentValues: PaymentValue) => void;
    setInitalizeData: (newData: InitializeData) => void
}

interface BillingAsyncActions {
    initializePayment: () => Promise<void>
    confirmPayment: (reference: string, paymentMethod: PaymentMethod) => Promise<void>
}

type BillingStore = BillingActions & BillingAsyncActions & BillingState

const InitialState = {
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
} satisfies BillingState

const useBillingStore = create<BillingStore>((set, get) => ({
    ...InitialState,
    setPaymentValues: (newPaymentValues: PaymentValue) => set({ paymentValues: newPaymentValues }),
    setInitalizeData: (newData: InitializeData) => set({ intializeData: newData }),

    initializePayment: async () => {
        try {
            const { paymentValues } = get()
            console.log("Payment Values: ", paymentValues)
            const response = await BillingApi.initializePayment(paymentValues)
            if (response.status == true) {
                window.location.href = response.payload.data.authorization_url
            }
        } catch (error) {
            handleError(error)
        }
    },
    confirmPayment: async (reference: string, paymentMethod: PaymentMethod) => {
        try {
            const response = await BillingApi.confirmPayment(paymentMethod, reference)
            if (response.status == true) {
                window.location.href = "/app/billing"
                eventBus.emit('success', response.payload.data)
            }
        } catch (error) {
            handleError(error)
        }
    },

}));


export default useBillingStore

