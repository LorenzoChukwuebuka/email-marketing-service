import { useState } from "react";
import EmptyState from "../../../../components/emptyStateComponent";
import useContactStore from "../../store/contact.store";
import { Pagination } from 'antd';
import { ContactAPIResponse } from "../../interface/contact.interface";
import EditContact from './editContactComponent';
import useContactGroupStore from "../../store/contactgroup.store";
import { APIResponse } from "../../../../interface/api.interface";
import { PaginatedResponse } from "../../../../interface/pagination.interface";
import LoadingSpinnerComponent from "../../../../components/loadingSpinnerComponent";

interface Props {
    contactData: APIResponse<PaginatedResponse<ContactAPIResponse>>
    loading: boolean
    currentPage:number 
    pageSize:number
    onPageChange: (page: number, size: number) => void 
}

const GetAllContacts: React.FC<Props> = ({ contactData, loading,currentPage,pageSize,onPageChange }) => {
    const { selectedIds, setSelectedId } = useContactStore();
    const { setSelectedContactIds } = useContactGroupStore()
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [selectedContact, setSelectedContact] = useState<ContactAPIResponse | null>(null);

    const cdata = contactData?.payload?.data || []

    const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.checked) {
            const allIds = cdata?.map((contact) => contact.uuid);
            setSelectedId(allIds as string[]);
            setSelectedContactIds(allIds as string[])
        } else {
            setSelectedId([]);
            setSelectedContactIds([])
        }
    };

    const handleSelect = (uuid: string) => {
        if (selectedIds.includes(uuid)) {
            setSelectedId(selectedIds.filter((id) => id !== uuid));
            setSelectedContactIds(selectedIds.filter((id) => id !== uuid));
        } else {
            setSelectedId([...selectedIds, uuid]);
            setSelectedContactIds([...selectedIds, uuid])
        }
    };

    const openEditModal = (contact: ContactAPIResponse) => {
        setSelectedContact(contact);
        setIsModalOpen(true);
    };

    const closeEditModal = () => {
        setIsModalOpen(false);
        setSelectedContact(null);
    };

    return (
        <>

            {loading ? <LoadingSpinnerComponent /> : (
                <>
                    <div className="overflow-x-auto mt-8">
                        {contactData && contactData?.payload?.data?.length > 0 ? (
                            <>
                                <table className="md:min-w-5xl min-w-full w-full rounded-sm  bg-white">
                                    <thead className="bg-gray-50">
                                        <tr>
                                            <th className="py-3 px-4 text-left">
                                                <input
                                                    type="checkbox"
                                                    className="form-checkbox h-4 w-4 text-blue-600"
                                                    onChange={handleSelectAll}
                                                    checked={selectedIds.length === (contactData?.payload?.data?.length ?? 0)}
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
                                                Created On
                                            </th>
                                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Edit
                                            </th>

                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-gray-200">
                                        {cdata?.map((contact) => (
                                            <tr key={contact.uuid} className="hover:bg-slate-100">
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
                                                        second: 'numeric',
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
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>

                                <div className="mt-4 flex justify-center items-center mb-5">
                                    {/* Pagination */}
                                    <Pagination
                                        current={currentPage}
                                        pageSize={pageSize}
                                        total={contactData?.payload?.total_count || 0} // Ensure your API returns a total count
                                        onChange={onPageChange}
                                        showSizeChanger
                                        pageSizeOptions={["10", "20", "50", "100"]}
                                        showTotal={(total) => `Total ${total} Contacts`}
                                    />
                                </div>


                            </>
                        ) : (
                            <div className="py-4 px-4 text-center">
                                <EmptyState
                                    title="You have not created any Contacts"
                                    description="Create contacts"
                                    icon={<i className="bi bi-emoji-frown text-xl"></i>}
                                />
                            </div>
                        )}
                    </div>

                    <EditContact
                        isOpen={isModalOpen}
                        onClose={closeEditModal}
                        contact={selectedContact}
                    />
                </>
            )}
        </>
    );
};

export default GetAllContacts;