import { create } from 'zustand';
import { SMTPKeyFormValues } from '../interface/smtpkey.interface';
import { handleError } from '../../../utils/isError';
import SMTPKeyAPI from '../api/smtpkey.api';
import eventBus from '../../../utils/eventbus';

interface SMTPState {
    smtpformValues: SMTPKeyFormValues;
}

interface SMTPAsyncActions {
    generateSMTPKey: () => Promise<void>
    deleteSMTPKey: (smtpKeyId: string) => Promise<void>
    createSMTPKey: () => Promise<void>
}

interface SMTPActions {
    setSmtpFormValues: (newFormValues: SMTPKeyFormValues) => void;
}

type SMTPStore = SMTPActions & SMTPAsyncActions & SMTPState


const InitialState: SMTPState = {
    smtpformValues: { key_name: "" },
}

const useSMTPKeyStore = create<SMTPStore>((set, get) => ({
    ...InitialState,


    setSmtpFormValues: (newFormValues) => set(() => ({ smtpformValues: newFormValues })),

    createSMTPKey: async () => {
        try {
            const { smtpformValues } = get();
            const response = await SMTPKeyAPI.createSMTPKey(smtpformValues)
            if (response) {
                eventBus.emit('success', 'SMTP key generated successfully');
            }


        } catch (error) {
            handleError(error)
        }
    },

    deleteSMTPKey: async (smtpKeyId: string) => {
        try {
            await SMTPKeyAPI.deleteSMTPKey(smtpKeyId)
            eventBus.emit('success', 'SMTP key deleted successfully');
        } catch (error) {
            handleError(error)
        }
    },

    generateSMTPKey: async () => {
        try {
            const response = await SMTPKeyAPI.generateSMTPKey()
            if (response) {
                eventBus.emit('success', 'SMTP credentials changed successfully');
            }


        } catch (error) {
            handleError(error)
        }
    }
}));

export default useSMTPKeyStore;