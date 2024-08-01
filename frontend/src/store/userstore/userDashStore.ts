import { create } from 'zustand'
import axiosInstance from '../../utils/api'
import eventBus from '../../utils/eventBus'
import { APIResponse } from '../../interface/api.interface';


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

            console.log(response.data)
        } catch (error) {
            eventBus.emit('error', error)
        }
    }
}))


export default useDailyUserMailSentCalc
