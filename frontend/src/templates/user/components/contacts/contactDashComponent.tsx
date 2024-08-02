import { FormEvent, useState } from "react";
import useContactStore from "../../../../store/userstore/contactStore";
import CreateContact from "./CreateContactComponent";
import GetAllContacts from "./getAllContactsComponent";
import ContactUpload from "./contactBatchUploadComponent";

// Define the type for the component props if needed, but for now it's empty
type ContactsDashTemplateProps = {};

const ContactsDashComponent: React.FC<ContactsDashTemplateProps> = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { selectedIds, deleteContact, getAllContacts } = useContactStore();
    const [importModalOpen, setImportModalOpen] = useState<boolean>(false)

    const handleDelete = (e: FormEvent<HTMLButtonElement>) => {
        e.preventDefault();
        deleteContact();
        getAllContacts();
    };

    const addContactToGroupt = (e: FormEvent<HTMLButtonElement>) => {
        e.preventDefault()
    }

    const importContact = () => {
        setImportModalOpen(true)
    }

    return (
        <>

            <div className="flex justify-between items-center mt-10">
                <div className="space-x-2">
                    <button
                        className="bg-gray-300 px-4 py-2 rounded-md transition duration-300"
                        onClick={() => setIsModalOpen(true)}
                    >
                        Create Contact
                    </button>

                    <button
                        className="bg-gray-300 px-4 py-2 rounded-md transition duration-300"
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
                                onClick={(e) => addContactToGroupt(e)}
                            >
                                <span className="text-blue-700"> Add to Group </span>
                                <i className="bi bi-people text-blue-500"></i>
                            </button>
                        </>

                    )}
                </div>
            </div>

            <CreateContact isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />

            <GetAllContacts />

            <ContactUpload isOpen={importModalOpen} onClose={() => setImportModalOpen(false)} />
        </>
    );
};

export default ContactsDashComponent;