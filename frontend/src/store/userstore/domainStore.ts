import { create } from 'zustand';
import axiosInstance from '../../utils/api';
import eventBus from '../../utils/eventBus';
import { BaseEntity } from '../../interface/baseentity.interface';
import { PaginatedResponse } from '../../interface/pagination.interface';
import { APIResponse, ResponseT } from '../../interface/api.interface';
import { errResponse } from '../../utils/isError';
import { downloadFile } from '../../utils/utils';

export type DomainRecord = {
    user_id: string;
    domain: string;
    txt_record: string;
    dmarc_record: string;
    dkim_selector: string;
    dkim_public_key: string;
    verified: boolean;
} & BaseEntity;

type FormValues = Pick<DomainRecord, "domain">;

type DomainStore = {
    domainformValues: FormValues;
    domainData: DomainRecord[] | DomainRecord;
    setDomainFormValues: (newData: FormValues) => void;
    setDomainData: (newData: DomainRecord[] | DomainRecord) => void;
    createDomain: () => Promise<void>;
    paginationInfo: Omit<PaginatedResponse<DomainRecord>, 'data'>;
    setPaginationInfo: (newPaginationInfo: Omit<PaginatedResponse<DomainRecord>, 'data'>) => void;
    deleteDomain: (uuid: string) => Promise<void>;
    authenticateDomain: (uuid: string) => Promise<void>;
    getDomain: (uuid: string) => Promise<void>;
    getAllDomain: (page?: number, pageSize?: number, search?: string) => Promise<void>;
    searchDomain: (query?: string) => void
};

const useDomainStore = create<DomainStore>((set, get) => ({
    domainformValues: {
        domain: ""
    },
    domainData: [],
    paginationInfo: {
        total_count: 0,
        total_pages: 0,
        current_page: 1,
        page_size: 10,
    },
    setDomainFormValues: (newData) => set({ domainformValues: newData }),
    setDomainData: (newData) => set({ domainData: newData }),
    setPaginationInfo: (newPaginationInfo) => set({ paginationInfo: newPaginationInfo }),
    createDomain: async () => {
        try {
            let response = await axiosInstance.post<ResponseT>("/domain/create-domain", get().domainformValues)
            if (response.data.status === true) {
                downloadFile(response.data.payload.downloadable_records)
            }
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

    deleteDomain: async (uuid: string) => {
        try {
            let response = await axiosInstance.delete<ResponseT>("/domain/delete-domain/" + uuid)

            if (response.data.status === true) {
                eventBus.emit('success', "Domain deleted successfully")
            }
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

    authenticateDomain: async (uuid: string) => {
        try {
            let response = await axiosInstance.put<ResponseT>("/domain/authenticate-domain/" + uuid)

            if (response.data.status === true) {
                eventBus.emit('success', "Domain authenticated successfully")
            }
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

    getDomain: async (uuid: string) => {
        try {
            let response = await axiosInstance.get<APIResponse<DomainRecord>>("/domain/get-domain/" + uuid)
            if (response.data.status === true) {
                get().setDomainData(response.data.payload)
            }
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

    getAllDomain: async (page = 1, pageSize = 10, query = "") => {
        try {
            let response = await axiosInstance.get<APIResponse<PaginatedResponse<DomainRecord>>>("/domain/get-all-domains", {
                params: {
                    page: page || undefined,
                    page_size: pageSize || undefined,
                    search: query || undefined
                }
            })
            const { data, ...paginationInfo } = response.data.payload;
            get().setDomainData(data)
            get().setPaginationInfo(paginationInfo)

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

    searchDomain: async (query?: string) => {

        const { getAllDomain } = get();

        if (!query) {
            await getAllDomain();
            return;
        }

        await getAllDomain(1, 10, query)
    },

}));

export default useDomainStore;
