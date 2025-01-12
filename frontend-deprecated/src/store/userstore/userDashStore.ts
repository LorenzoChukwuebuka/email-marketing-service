import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'
import { APIResponse } from '../../interface/api.interface';
import { errResponse } from '../../utils/isError';


interface MailData {
    remainingMails: number;
    mailsPerDay: number;
    plan?: string;
}


interface DailyUserMailStore {
    mailData: MailData | null;
    setMailData: (newMailData: MailData | null) => void;
    getUserMailData: () => Promise<void>;
    
}

const useDailyUserMailSentCalc = create<DailyUserMailStore>((set, get) => ({
    mailData: null,
    setMailData: newMailData => set({ mailData: newMailData }),

    getUserMailData: async () => {
        const { setMailData } = get()
        try {
            let response = await axiosInstance.get<APIResponse<MailData>>('/get-user-current-sub')
            setMailData(response.data.payload)

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
}))


export default useDailyUserMailSentCalc
