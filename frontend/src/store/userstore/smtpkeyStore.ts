import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';

interface FormValues {
    key_name: string;
}

interface SMTPKey {
    uuid: string;
    user_id: string;
    key_name: string;
    smtp_login: string;
    password: string;
    status: string;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
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
    getSMTPKeys: () => void
    generateSMTPKey: () => Promise<void>
    deleteSMTPKey: (smtpKeyId: string) => void
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

    //initialize 
    setSmtpFormValues: (newFormValues) => set(() => ({ smtpformValues: newFormValues })),
    setSmtpKeyData: (data) => set(() => ({ smtpKeyData: data })),

    //functions
    getSMTPKeys: async () => {
        try {
            const { setSmtpKeyData } = get()
            let response = await axiosInstance.get('/get-smtp-keys')
            setSmtpKeyData(response.data.payload)
        } catch (error) {
            eventBus.emit('error', error instanceof Error || 'An unexpected error occured')
        }
    },

    generateSMTPKey: async () => {
        try {
            const { smtpformValues } = get();
            const response = await axiosInstance.post('/create-smtp-key', smtpformValues);
            eventBus.emit('success', 'SMTP key generated successfully');
            return response.data.payload;
        } catch (error) {
            eventBus.emit('error', error instanceof Error ? error.message : 'An unexpected error occurred');
            throw error;
        }
    },

    deleteSMTPKey: async (smtpKeyId: string) => {
        try {
            const response = await axiosInstance.delete(`/delete-smtp-key/${smtpKeyId}`);
            eventBus.emit('success', 'SMTP key deleted successfully');
            get().getSMTPKeys(); 
        } catch (error) {
            eventBus.emit('error', error instanceof Error ? error.message : 'An unexpected error occurred');
        }
    }
}));

export default useSMTPKeyStore;