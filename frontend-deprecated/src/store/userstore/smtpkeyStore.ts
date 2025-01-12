import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { errResponse } from '../../utils/isError';

interface FormValues {
    key_name: string;
}

type SMTPKey = BaseEntity & {
    user_id: string;
    key_name: string;
    smtp_login: string;
    password: string;
    status: string;
}

interface SMTPKeyDATA {
    keys: SMTPKey[];
    smtp_login: string;
    smtp_master: string;
    smtp_master_password: string;
    smtp_master_status: string;
    smtp_port: string;
    smtp_server: string;
    smtp_created_at: string
}

interface SMTPStore {
    smtpKeyData: SMTPKeyDATA;
    smtpformValues: FormValues;
    setSmtpFormValues: (newFormValues: FormValues) => void;
    setSmtpKeyData: (data: SMTPKeyDATA) => void;
    getSMTPKeys: () => Promise<void>
    generateSMTPKey: () => Promise<void>
    deleteSMTPKey: (smtpKeyId: string) => Promise<void>
    createSMTPKey: () => Promise<void>
}

const useSMTPKeyStore = create<SMTPStore>((set, get) => ({
    smtpKeyData: {
        keys: [],
        smtp_login: "",
        smtp_master: "",
        smtp_master_password: "",
        smtp_master_status: "",
        smtp_port: "",
        smtp_server: "",
        smtp_created_at: ""
    },
    smtpformValues: { key_name: "" },


    setSmtpFormValues: (newFormValues) => set(() => ({ smtpformValues: newFormValues })),
    setSmtpKeyData: (data) => set(() => ({ smtpKeyData: data })),

    getSMTPKeys: async () => {
        try {
            const { setSmtpKeyData } = get()
            let response = await axiosInstance.get('/smtpkey/get-smtp-keys')
            setSmtpKeyData(response.data.payload)
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

    createSMTPKey: async () => {
        try {
            const { smtpformValues } = get();
            const response = await axiosInstance.post('/smtpkey/create-smtp-key', smtpformValues);
            eventBus.emit('success', 'SMTP key generated successfully');
            return response.data.payload;
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

    deleteSMTPKey: async (smtpKeyId: string) => {
        try {
            await axiosInstance.delete(`/smtpkey/delete-smtp-key/${smtpKeyId}`);
            eventBus.emit('success', 'SMTP key deleted successfully');
            get().getSMTPKeys();
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

    generateSMTPKey: async () => {
        try {
            const response = await axiosInstance.put('/smtpkey/generate-new-smtp-master-password');
            eventBus.emit('success', 'SMTP credentials changed successfully');
            get().getSMTPKeys();
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

export default useSMTPKeyStore;