import { create } from 'zustand';
import { APIKeyFormValues } from '../interface/apikey.interface';
import { handleError } from '../../../utils/isError';
import APIKeyAPI from '../api/apikey.api';
import eventBus from '../../../utils/eventbus';

interface APIKeyState {
    formValues: APIKeyFormValues;
}

interface APIKeyActions {
    setFormValues: (newFormValue: APIKeyFormValues) => void;
}

interface APIKeyAsyncActions {
    generateAPIKey: () => Promise<any>;
    deleteAPIKey: (apiId: string) => Promise<void>;
}

type APIKeyStore = APIKeyState & APIKeyActions & APIKeyAsyncActions

const InitialState: APIKeyState = {
    formValues: {
        name: ""
    }
}

const useAPIKeyStore = create<APIKeyStore>((set, get) => ({
    ...InitialState,
    setFormValues: newFormValue => set({ formValues: newFormValue }),
    generateAPIKey: async () => {
        const { formValues } = get();
        try {
            const response = await APIKeyAPI.generateAPIkey(formValues)
            // return response.payload;
            if (response) {
                eventBus.emit('success', 'Api key generation successful')
            }
        } catch (error) {
            handleError(error)
        }

    },
    deleteAPIKey: async (apiId: string) => {
        try {
            let response = await APIKeyAPI.deleteAPIKey(apiId)
            eventBus.emit('success', response.data.payload);
        } catch (error) {
            handleError(error)
        }
    }
}));

export default useAPIKeyStore;
