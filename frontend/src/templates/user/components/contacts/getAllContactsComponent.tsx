import { useEffect, useState } from "react";
import useContactStore, { Contact } from "../../../../store/userstore/contactStore";
import { convertToNormalTime } from '../../../../utils/utils';
import Pagination from "../../../../components/Pagination";
import EditContact from "./editContactComponent";



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


    const handleSubscriptionToggle = async (uuid: string) => {
        try {
            // Find the contact in the contactData array
            const contactIndex = contactData.findIndex(contact => contact.uuid === uuid);
            if (contactIndex === -1) return;

            // Toggle the is_subscribed status
            const updatedContact = {
                ...contactData[contactIndex],
                is_subscribed: !contactData[contactIndex].is_subscribed
            };

            // Update the contact in your backend
            //   await updateContactSubscription(uuid, updatedContact.is_subscribed);

            // Update the state
            // setContactData(prevData => {
            //     const newData = [...prevData];
            //     newData[contactIndex] = updatedContact;
            //     return newData;
            // });
        } catch (error) {
            console.error('Error toggling subscription:', error);
            // Handle error (e.g., show an error message to the user)
        }
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
                                    checked={selectedIds.length === (contactData?.length ?? 0)}
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
                        {contactData && contactData.length > 0 ? (
                            contactData.map((contact: any) => (
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
                                    <td className="py-4 px-4">
                                        {new Date(contact.created_at).toLocaleString('en-US', {
                                            timeZone: 'UTC',
                                            year: 'numeric',
                                            month: 'long',
                                            day: 'numeric',
                                            hour: 'numeric',
                                            minute: 'numeric',
                                            second: 'numeric'
                                        })}
                                    </td>

                                    <td className="py-4 px-4">
                                        <button
                                            className="text-gray-400 hover:text-gray-600"
                                            onClick={() => openEditModal(contact)}
                                        >
                                            ✏️
                                        </button>
                                    </td>
                                    <td className="py-4 px-4">
                                        <input
                                            type="checkbox"
                                            className="form-checkbox h-4 w-4 text-blue-600"
                                            checked={contact.is_subscribed}
                                            onChange={() => handleSubscriptionToggle(contact.uuid)}
                                        />
                                    </td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={7} className="py-4 px-4  text-center">
                                    No contacts available
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>

            </div>

            <Pagination paginationInfo={paginationInfo} handlePageChange={handlePageChange} item="Contacts" />

            <EditContact
                isOpen={isModalOpen}
                onClose={closeEditModal}
                contact={selectedContact}
            />


        </>
    );
};

export default GetAllContacts;