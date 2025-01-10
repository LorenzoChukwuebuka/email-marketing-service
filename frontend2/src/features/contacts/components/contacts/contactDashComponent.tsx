import { FormEvent, useState } from "react";
import useDebounce from "../../../../hooks/useDebounce";
import CreateContact from './createContact';
import useContactStore from "../../store/contact.store";
import GetAllContacts from "./getAllContacts";
import { Modal } from "antd";
import { useContactQuery } from "../../hooks/useContactQuery";
import ContactUpload from "./contactBatchUploadComponent";
import AddContactsToGroupComponent from './addContactsToGroupComponent';
import { APIResponse } from "../../../../interface/api.interface";
import { PaginatedResponse } from "../../../../interface/pagination.interface";
import { ContactAPIResponse } from "../../interface/contact.interface";

const ContactsDashComponent: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { selectedIds, deleteContact } = useContactStore();
    const [importModalOpen, setImportModalOpen] = useState<boolean>(false);
    const [groupModalOpen, setGroupModalOpen] = useState<boolean>(false);
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query

    // Debounce the search query
    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [currentPage, _setCurrentPage] = useState(1);
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [pageSize, _setPageSize] = useState(20);

    const { data: contacts,isLoading } = useContactQuery(currentPage, pageSize, debouncedSearchQuery)


    const handleDelete = async (e: FormEvent<HTMLButtonElement>) => {
        e.preventDefault();
        Modal.confirm({
            title: "Are you sure?",
            content: "Do you want to delete contact(s)?",
            okText: "Yes",
            cancelText: "No",
            onOk: async () => {
                await deleteContact();
                await new Promise(resolve => setTimeout(resolve, 1000));
                location.reload()
            },
        });
    };

    const addContactToGroupt = () => {
        setGroupModalOpen(true);
    };

    const importContact = () => {
        setImportModalOpen(true);
    };

    // Update search query state
    const handleSearchInput = (query: string) => {
        setSearchQuery(query);
    };

    return (
        <>

            <div className="flex justify-between items-center rounded-md p-2 bg-white mt-10">
                <div className="space-x-1 h-auto w-full p-2 px-2">
                    <button
                        className="bg-gray-300 px-2 py-2 rounded-md transition duration-300"
                        onClick={() => setIsModalOpen(true)}
                    >
                        Create Contact
                    </button>

                    <button
                        className="bg-gray-300 px-2 py-2 rounded-md transition duration-300"
                        onClick={() => importContact()}
                    >
                        Import Contact
                    </button>

                    {selectedIds.length > 0 && (
                        <>
                            <button
                                className="bg-red-200 px-4 py-2 rounded-md transition duration-300"
                                onClick={(e) => handleDelete(e)}
                            >
                                <span className="text-red-500"> Delete Contact </span>
                                <i className="bi bi-trash text-red-500"></i>
                            </button>

                            <button
                                className="bg-blue-200 px-4 py-2 rounded-md transition duration-300"
                                onClick={() => addContactToGroupt()}
                            >
                                <span className="text-blue-700"> Add to Group </span>
                                <i className="bi bi-people text-blue-500"></i>
                            </button>
                        </>
                    )}
                </div>

                <div className="ml-3">
                    <input
                        type="text"
                        placeholder="Search..."
                        className="bg-gray-100 px-4 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-300"
                        onChange={(e) => handleSearchInput(e.target.value)}
                        value={searchQuery}
                    />
                </div>
            </div>
            <CreateContact isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />
            <GetAllContacts contactData={contacts as APIResponse<PaginatedResponse<ContactAPIResponse>> } loading={isLoading} />
            <ContactUpload isOpen={importModalOpen} onClose={() => setImportModalOpen(false)} />
            <AddContactsToGroupComponent isOpen={groupModalOpen} onClose={() => setGroupModalOpen(false)} />
        </>

    )

};

export default ContactsDashComponent;