import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import useContactGroupStore, { ContactGroupData } from "../../../../store/userstore/contactGroupStore";

const GroupContactList: React.FC = () => {
    const { isLoading, contactgroupData, getSingleGroup } = useContactGroupStore();
    const [group, setGroup] = useState<ContactGroupData | null>(null);

    const location = useLocation();

    useEffect(() => {
        const fetchGroup = async () => {
            const stateData = location.state as { groupId: string };
            if (stateData && stateData.groupId) {
                await getSingleGroup(stateData.groupId);
            }
        };

        fetchGroup();
    }, [location.state, getSingleGroup]);

    useEffect(() => {
        if (contactgroupData) {
            if (Array.isArray(contactgroupData) && contactgroupData.length > 0) {
                setGroup(contactgroupData[0]);
            } else if (!Array.isArray(contactgroupData)) {
                setGroup(contactgroupData);
            }
        }
    }, [contactgroupData]);

    return (
        <main className="mt-10 min-h-screen">
            <div className="p-6">
                {isLoading ? (
                    <p>Loading...</p>
                ) : group ? (
                    <>
                        <h2 className="text-xl font-bold text-gray-800">Group Name:  {group.group_name}</h2>
                        <p className="text-md text-gray-500 p-2 px-2  mt-2">Description: {group.description}</p>

                        <h1 className="text-xl font-semibold mt-5"> Contacts</h1>

                        <table className="md:min-w-5xl min-w-full w-full mt-5 rounded-sm bg-white">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="py-3 px-4 text-left">
                                        {/* <input
                                    type="checkbox"
                                    className="form-checkbox h-4 w-4 text-blue-600"
                                    onChange={handleSelectAll}
                                    checked={selectedIds.length === (contactData?.length ?? 0)}
                                /> */}
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

                                   

                                    <th className="py-3 px-4"></th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200 ">
                                {group.contacts && group.contacts.length > 0 ? (
                                    group.contacts.map((contact: any) => (
                                        <tr key={contact.uuid}>
                                            <td className="py-4 px-4">
                                                {/* <input
                                            type="checkbox"
                                            className="form-checkbox h-4 w-4 text-blue-600"
                                            checked={selectedIds.includes(contact.uuid)}
                                            onChange={() => handleSelect(contact.uuid)}
                                        /> */}
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

                                            {/* <td className="py-4 px-4">
                                        <button
                                            className="text-gray-400 hover:text-gray-600"
                                            onClick={() => openEditModal(contact)}
                                        >
                                            ✏️
                                        </button>
                                    </td> */}
                                            {/* <td className="py-4 px-4">
                                        <input
                                            type="checkbox"
                                            className="form-checkbox h-4 w-4 text-blue-600"
                                            checked={contact.is_subscribed}
                                            onChange={() => handleSubscriptionToggle(contact.uuid)}
                                        />
                                    </td> */}
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
                    </>
                ) : (
                    <p>No group data available.</p>
                )}
            </div>
        </main>
    );
};

export default GroupContactList;