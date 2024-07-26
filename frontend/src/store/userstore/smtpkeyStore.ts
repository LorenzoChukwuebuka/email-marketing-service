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
}

interface SMTPStore {
    smtpKeyData: SMTPKeyDATA;
    formValues: FormValues;


    setFormValues: (newFormValues: FormValues) => void;
    setSmtpKeyData: (data: SMTPKeyDATA) => void;
    getSMTPKeys: () => void
    generateSMTKey: () => void
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
        smtp_server: ""
    },
    formValues: { key_name: "" },

    //initialize 
    setFormValues: (newFormValues) => set(() => ({ formValues: newFormValues })),
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

    generateSMTKey: async () => { },

    deleteSMTPKey: async (smtpKeyId: string) => {

    }
}));

export default useSMTPKeyStore;