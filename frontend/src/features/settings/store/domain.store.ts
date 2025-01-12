import { create } from 'zustand';
import { DomainFormValues } from '../interface/domain.interface';
import eventBus from '../../../utils/eventbus';
import { handleError } from '../../../utils/isError';
import { downloadFile } from '../../../utils/utils';
import DomainAPI from '../api/domain.api';

type DomainStore = {
    domainformValues: DomainFormValues;
    setDomainFormValues: (newData: DomainFormValues) => void;
    createDomain: () => Promise<void>;
    deleteDomain: (uuid: string) => Promise<void>;
    authenticateDomain: (uuid: string) => Promise<void>;
};

const useDomainStore = create<DomainStore>((set, get) => ({
    domainformValues: {
        domain: ""
    },

    setDomainFormValues: (newData) => set({ domainformValues: newData }),

    createDomain: async () => {
        try {
            const response = await DomainAPI.createDomain(get().domainformValues)
            if (response.status === true) {
                downloadFile(response.payload.downloadable_records)
            }
        } catch (error) {
            handleError(error)
        }
    },

    deleteDomain: async (uuid: string) => {
        try {
            const response = await DomainAPI.deleteDomain(uuid)
            if (response.status === true) {
                eventBus.emit('success', "Domain deleted successfully")
            }
        } catch (error) {
            handleError(error)
        }
    },

    authenticateDomain: async (uuid: string) => {
        try {
            const response = await DomainAPI.authenticateDomain(uuid)
            if (response.status === true) {
                eventBus.emit('success', "Domain authenticated successfully")
            }
        } catch (error) {
            handleError(error)
        }
    }

}));

export default useDomainStore;
