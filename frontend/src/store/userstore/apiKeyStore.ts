import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { errResponse } from '../../utils/isError';

interface FormValues {
    name: string;
}

interface APIKeyState {
    apiKeyData: any;
    isLoading: boolean;
    formValues: FormValues;

    setAPIKeyData: (newAPIKeyData: any) => void;
    setIsLoading: (newIsLoading: boolean) => void;
    setFormValues: (newFormValue: FormValues) => void;

    getAPIKey: () => Promise<void>;
    generateAPIKey: () => Promise<any>;
    deleteAPIKey: (apiId: string) => Promise<void>;
}

const useAPIKeyStore = create<APIKeyState>((set, get) => ({
    apiKeyData: null,
    isLoading: false,
    formValues: { name: '' },

    setAPIKeyData: newAPIKeyData => set({ apiKeyData: newAPIKeyData }),
    setIsLoading: newIsLoading => set({ isLoading: newIsLoading }),
    setFormValues: newFormValue => set({ formValues: newFormValue }),

    getAPIKey: async () => {
        try {
            const { setAPIKeyData } = get();

            let response = await axiosInstance.get('/get-apikey');
            setAPIKeyData(response.data);
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
    generateAPIKey: async () => {
        const { setIsLoading, formValues } = get();
        try {
            setIsLoading(true);
            let response = await axiosInstance.post('/generate-apikey', formValues);
            return response.data.payload;
        } catch (error) {
            if (errResponse(error)) {
                eventBus.emit('error', error?.response?.data.payload)
            } else if (error instanceof Error) {
                eventBus.emit('error', error.message);
            } else {
                console.error("Unknown error:", error);
            }
        } finally {
            get().setIsLoading(false);
        }
    },
    deleteAPIKey: async (apiId: string) => {
        try {
            let response = await axiosInstance.delete('/delete-apikey/' + apiId);
            eventBus.emit('success', response.data.payload);
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

export default useAPIKeyStore;
