import { create } from 'zustand';
import { Sender, VerifySender, SenderFormValues } from '../interface/sender.interface';
import senderApi from '../api/sender.api';
import { handleError } from '../../../utils/isError';
import eventBus from '../../../utils/eventbus';

interface SenderState {
    verifySenderForm: VerifySender
    senderFormValues: SenderFormValues;
}

interface SenderActions {
    setSenderFormValues: (values: SenderFormValues) => void;
    setVerifySender: (values: VerifySender) => void
}

interface SenderAsyncActions {
    createSender: () => Promise<void>;
    updateSender: (id: string, data: Partial<Sender>) => Promise<void>;
    deleteSender: (id: string) => Promise<void>;
    verifySender: () => Promise<void>
}


const InitialState: SenderState = {
    senderFormValues: {
        name: '',
        email: '',
    },
    verifySenderForm: {
        email: "", token: "", user_id: ""
    },
}


type SenderStore = SenderState & SenderActions & SenderAsyncActions

const useSenderStore = create<SenderStore>((set, get) => ({
    ...InitialState,

    setSenderFormValues: (newData) => set({ senderFormValues: newData }),
    setVerifySender: (newData) => set({ verifySenderForm: newData }),

    createSender: async () => {
        try {
            await senderApi.createSender(get().senderFormValues)
        } catch (error) {
            handleError(error)
        }
    },

    updateSender: async (id: string, data: Partial<Sender>) => {
        try {
            const response = await senderApi.updateSender(id, data)
            if (response) {
                eventBus.emit('success', 'sender updated successfully')
            }
        } catch (error) {
            handleError(error)
        }
    },

    deleteSender: async (id: string) => {
        try {
            const response = await senderApi.deleteSender(id)
            if (response) {
                eventBus.emit('success', 'sender deleted successfully')
            }
        } catch (error) {
            handleError(error)
        }
    },
    verifySender: async () => {
        try {
            const response = await senderApi.verifySender(get().verifySenderForm)
            if (response.status === true) {
                eventBus.emit("success", "Sender has been verified successfully")
            }
        } catch (error) {
            handleError(error)
        }
    }
}));

export default useSenderStore;