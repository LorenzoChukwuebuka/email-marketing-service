
import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';

interface ContactFormValues {
    first_name: string
    last_name: string
    email: string
    from: string
}

interface Contact {
    uuid: string;
    first_name: string;
    last_name: string;
    email: string;
    from: string;
    user_id: string;
    created_at: string;
    updated_at: string | null;
    deleted_at: string | null;
    groups: any[] | null;
}

interface PaginationInfo {
    total_count: number;
    total_pages: number;
    current_page: number;
    page_size: number;
}

interface ContactStore {
    contactFormValues: ContactFormValues;
    contactData: Contact[];
    selectedIds: string[];
    paginationInfo: PaginationInfo;
    setContactData: (newContactData: Contact[]) => void;
    setContactFormValues: (newFormValues: ContactFormValues) => void;
    setSelectedId: (newSelectedId: string[]) => void;
    setPaginationInfo: (newPaginationInfo: PaginationInfo) => void;
    createContact: () => Promise<void>;
    deleteContact: () => Promise<void>;
    editContact: () => Promise<void>;
    getAllContacts: (page?: number, pageSize?: number) => Promise<void>;
}

const useContactStore = create<ContactStore>((set, get) => ({
    contactFormValues: {
        first_name: '',
        last_name: '',
        email: '',
        from: '',
    },
    contactData: [],
    selectedIds: [],
    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },
    setContactData: (newContactData) => set({ contactData: newContactData }),
    setContactFormValues: (newFormValues) => set({ contactFormValues: newFormValues }),
    setSelectedId: (newSelectedId) => set({ selectedIds: newSelectedId }),
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),

    createContact: async () => { /* implementation */ },
    deleteContact: async () => { /* implementation */ },
    editContact: async () => { /* implementation */ },
    getAllContacts: async (page = 1, pageSize = 10) => {
        try {
            const response = await axiosInstance.get(`/get-all-contacts?page=${page}&page_size=${pageSize}`);
            const { data, ...paginationInfo } = response.data.payload;
            get().setContactData(data);
            get().setPaginationInfo(paginationInfo);
        } catch (error) {
            eventBus.emit('error', error instanceof Error ? error.message : 'An unexpected error occurred');
        }
    }
}));

export default useContactStore;

