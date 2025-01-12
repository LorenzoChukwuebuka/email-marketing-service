import { FormEvent, useEffect, useState } from "react";
import useContactStore from "../../../../store/userstore/contactStore";
import CreateContact from "./CreateContactComponent";
import GetAllContacts from "./getAllContactsComponent";
import ContactUpload from "./contactBatchUploadComponent";
import AddContactsToGroupComponent from "./addContactsTogroupComponent";
import useDebounce from "../../../../hooks/useDebounce";


type ContactsDashTemplateProps = {};

const ContactsDashComponent: React.FC<ContactsDashTemplateProps> = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { selectedIds, deleteContact, getAllContacts, searchContacts } = useContactStore();
    const [importModalOpen, setImportModalOpen] = useState<boolean>(false);
    const [groupModalOpen, setGroupModalOpen] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query

    // Debounce the search query
    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay

    const handleDelete = async (e: FormEvent<HTMLButtonElement>) => {
        e.preventDefault();
        let confirmResult = confirm("Do you want to delete contact(s)?");

        if (confirmResult) {
            await deleteContact();
        }
        await new Promise(resolve => setTimeout(resolve, 1000));
        await getAllContacts();
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

    useEffect(() => {
        const fetchC = async () => {
            setIsLoading(true);
            await getAllContacts();
            await new Promise(resolve => setTimeout(resolve, 1000));
            setIsLoading(false);
        };

        fetchC();
    }, [getAllContacts]);

    // Effect for debounced search
    useEffect(() => {
        if (debouncedSearchQuery !== "") {
            searchContacts(debouncedSearchQuery);
        } else {
            getAllContacts(); // Reset to all contacts when search query is empty
        }
    }, [debouncedSearchQuery, searchContacts, getAllContacts]);

    return (
        <>
            {isLoading ? (
                <div className="flex items-center justify-center mt-20">
                    <span className="loading loading-spinner loading-lg"></span>
                </div>
            ) : (
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

                    <GetAllContacts />

                    <ContactUpload isOpen={importModalOpen} onClose={() => setImportModalOpen(false)} />

                    <AddContactsToGroupComponent isOpen={groupModalOpen} onClose={() => setGroupModalOpen(false)} />
                </>
            )}
        </>
    );
};

export default ContactsDashComponent;