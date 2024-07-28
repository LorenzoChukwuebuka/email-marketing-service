import { useEffect, useState } from "react";
import useContactStore from "../../../../store/userstore/contactStore";
import { convertToNormalTime } from '../../../../utils/utils';



interface Contact {
    uuid: string;
    first_name: string;
    last_name: string;
    email: string;
    from: string;
    user_id: string;
    created_at: string;
    updated_at: string;
    deleted_at: string;
}

const GetAllContacts: React.FC = () => {
    const { getAllContacts, contactData, selectedIds, setSelectedId, paginationInfo } = useContactStore();
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [selectedContact, setSelectedContact] = useState<Contact | null>(null);

    useEffect(() => {
        getAllContacts();
    }, [getAllContacts]);


    const handlePageChange = (newPage: number) => {
        getAllContacts(newPage, paginationInfo.page_size);
    };


    const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.checked) {
            const allIds = contactData.map((contact) => contact.uuid);
            setSelectedId(allIds);
        } else {
            setSelectedId([]);
        }
    };

    const handleSelect = (uuid: string) => {
        if (selectedIds.includes(uuid)) {
            setSelectedId(selectedIds.filter((id) => id !== uuid));
        } else {
            setSelectedId([...selectedIds, uuid]);
        }
    };

    const openEditModal = (contact: Contact) => {
        setSelectedContact(contact);
        setIsModalOpen(true);
    };

    const closeEditModal = () => {
        setIsModalOpen(false);
        setSelectedContact(null);
    };

    return (
        <>
            <div className="overflow-x-auto mt-8">
                <table className="min-w-full w-full rounded-sm bg-white">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="py-3 px-4 text-left">
                                <input
                                    type="checkbox"
                                    className="form-checkbox h-4 w-4 text-blue-600"
                                    onChange={handleSelectAll}
                                    checked={selectedIds.length === contactData.length}
                                />
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                First Name
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Last Name
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Email
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                From
                            </th>

                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Created At
                            </th>

                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Edit
                            </th>

                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Toggle subscription
                            </th>

                            <th className="py-3 px-4"></th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200 ">
                        {contactData.map((contact: any) => (
                            <tr key={contact.uuid}>
                                <td className="py-4 px-4">
                                    <input
                                        type="checkbox"
                                        className="form-checkbox h-4 w-4 text-blue-600"
                                        checked={selectedIds.includes(contact.uuid)}
                                        onChange={() => handleSelect(contact.uuid)}
                                    />
                                </td>
                                <td className="py-4 px-4">{contact.first_name}</td>
                                <td className="py-4 px-4">{contact.last_name}</td>
                                <td className="py-4 px-4">{contact.email}</td>
                                <td className="py-4 px-4">{contact.from}</td>

                                <td className="py-4 px-4">{new Date(contact.created_at).toLocaleString('en-US', { timeZone: 'UTC', year: 'numeric', month: 'long', day: 'numeric', hour: 'numeric', minute: 'numeric', second: 'numeric' })}
                                </td>

                                <td className="py-4 px-4">
                                    <button
                                        className="text-gray-400 hover:text-gray-600"
                                        onClick={() => openEditModal(contact)}
                                    >
                                        ✏️
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>

            </div>

            <div className="mt-4 flex justify-between items-center">
                <div>Total Contacts: {paginationInfo.total_count}</div>
                <div className="flex space-x-2">
                    <button
                        onClick={() => handlePageChange(paginationInfo.current_page - 1)}
                        disabled={paginationInfo.current_page === 1}
                        className="px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300"
                    >
                        Previous
                    </button>
                    <span className="py-2">
                        Page {paginationInfo.current_page} of {paginationInfo.total_pages}
                    </span>
                    <button
                        onClick={() => handlePageChange(paginationInfo.current_page + 1)}
                        disabled={paginationInfo.current_page === paginationInfo.total_pages}
                        className="px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300"
                    >
                        Next
                    </button>
                </div>
            </div>

            {/* <EditContact
                isOpen={isModalOpen}
                onClose={closeEditModal}
                contact={selectedContact}
            /> */}
        </>
    );
};

export default GetAllContacts;