import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { errResponse } from '../../utils/isError';

interface Sender extends BaseEntity {
    name: string;
    email: string;
    // Add other sender-specific fields here
}

interface SenderFormValues {
    name: string;
    email: string;
}

interface SenderStore {
    senders: Sender[];
    isLoading: boolean;
    error: string | null;
    senderFormValues: SenderFormValues;
    currentPage: number;
    totalPages: number;
    setSenderFormValues: (values: Partial<SenderFormValues>) => void;
    createSender: () => Promise<void>;
    getSenders: (page?: number) => Promise<void>;
    updateSender: (id: string, data: Partial<Sender>) => Promise<void>;
    deleteSender: (id: string) => Promise<void>;
}

const useSenderStore = create<SenderStore>((set, get) => ({
    senders: [],
    isLoading: false,
    error: null,
    senderFormValues: {
        name: '',
        email: '',
    },
    currentPage: 1,
    totalPages: 1,

    setSenderFormValues: (values) => {
        set((state) => ({
            senderFormValues: { ...state.senderFormValues, ...values },
        }));
    },

    createSender: async () => {
        set({ isLoading: true, error: null });
        try {

        } catch (error) {

        } finally {
            set({ isLoading: false });
        }
    },

    getSenders: async (page = 1) => {

        try {


        } catch (error) {

        } finally {
            set({ isLoading: false });
        }
    },

    updateSender: async (id: string, data: Partial<Sender>) => {



    },

    deleteSender: async (id: string) => {

    },
}));

export default useSenderStore;