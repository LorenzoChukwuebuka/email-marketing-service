import React, { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import useDebounce from '../../../hooks/useDebounce';
import { useAllUsersQuery } from '../hooks/useAdminUsersQueryHook';
import useAdminUserStore from '../store/adminuser.store';
import { Modal, Pagination } from 'antd';
import LoadingSpinnerComponent from '../../../components/loadingSpinnerComponent';

const AllUsersTable: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query
    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const navigate = useNavigate()

    const { blockUser, verifyUser, unBlockUser } = useAdminUserStore()

    const { data: userData, isLoading } = useAllUsersQuery(currentPage, pageSize, debouncedSearchQuery)

    const handleSearch = (query: string) => {
        setSearchQuery(query)
    }

    const userdetailsData = useMemo(() => userData?.payload.data, [userData])

    const onPageChange = (page: number, size: number) => {
        setCurrentPage(page);
        setPageSize(size);
    };

    const handleToggle = async (userId: string, field: "verified" | "blocked", isChecked: boolean) => {

        const actionText = isChecked ?
            (field === "blocked" ? "unblock" : "unverify") :
            (field === "blocked" ? "block" : "verify");

        try {

            Modal.confirm({
                title: "Are you sure?",
                content: `Do you want to ${actionText} ${userId}?`,
                okText: "Yes",
                cancelText: "No",
                onOk: async () => {
                    if (field === "blocked") {
                        if (isChecked) {
                            await unBlockUser(userId);
                            new Promise(resolve => setTimeout(resolve, 1000))
                            location.reload()
                        } else {
                            await blockUser(userId);
                            new Promise(resolve => setTimeout(resolve, 1000))
                            location.reload()
                        }
                    } else if (field === "verified") {
                        await verifyUser(userId);
                    }
                },
            });


        } catch (error) {
            console.error(error)
        }
    };

    return (
        <>
            <div className="flex justify-between items-center rounded-md p-2 bg-white mt-10">
                <div className="ml-3">
                    <input
                        type="text"
                        placeholder="Search..."
                        className="bg-gray-100 px-4 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-300"
                        onChange={(e) => handleSearch(e.target.value)}
                    />
                </div>
            </div>
            {isLoading ? <LoadingSpinnerComponent /> :
                <div className="overflow-x-auto mt-4">
                    <h1 className="font-semibold text-lg mt-4 mb-4"> User List </h1>
                    <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                        <thead className="bg-gray-50">
                            <tr>
                                {['Full Name', 'Email', 'Company', 'Phone', 'Verified', 'Blocked', 'Created At', 'Verified At', 'Details'].map((header, index) =>
                                    <th key={index} className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        {header}
                                    </th>
                                )}
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-200">
                            {userdetailsData && userdetailsData.length > 0 ? (
                                userdetailsData.map((user, index) => {
                                    return (
                                        <tr key={index} className="hover:bg-gray-100">
                                            <td className="py-4 px-4">{user?.fullname}</td>
                                            <td className="py-4 px-4">{user?.email}</td>
                                            <td className="py-4 px-4">{user?.company || 'N/A'}</td>
                                            <td className="py-4 px-4">{user?.phonenumber || 'N/A'}</td>
                                            <td className="py-4 px-4">
                                                <label className="inline-flex items-center">
                                                    <input
                                                        type="checkbox"
                                                        className="form-checkbox h-5 w-5 text-blue-600"
                                                        checked={user?.verified}
                                                        onChange={() => handleToggle(user.uuid, 'verified', user.verified)}
                                                    />
                                                    <span className="ml-2">{user?.verified ? 'Verified' : 'Not Verified'}</span>
                                                </label>
                                            </td>
                                            <td className="py-4 px-4">
                                                <label className="inline-flex items-center">
                                                    <input
                                                        type="checkbox"
                                                        className="form-checkbox h-5 w-5 text-blue-600"
                                                        checked={user?.blocked}
                                                        onChange={() => handleToggle(user.uuid, 'blocked', user.blocked)}
                                                    />
                                                    <span className="ml-2">{user?.blocked ? 'Blocked' : 'Not Blocked'}</span>
                                                </label>
                                            </td>
                                            <td className="py-4 px-4">{user.created_at && new Date(user.created_at).toLocaleString('en-US', {
                                                timeZone: 'UTC',
                                                year: 'numeric',
                                                month: 'long',
                                                day: 'numeric',
                                                hour: 'numeric',
                                                minute: 'numeric',
                                                second: 'numeric'
                                            }) || "Not available"}</td>
                                            <td className="py-4 px-4">{user.verified_at && new Date(user.verified_at).toLocaleString('en-US', {
                                                timeZone: 'UTC',
                                                year: 'numeric',
                                                month: 'long',
                                                day: 'numeric',
                                                hour: 'numeric',
                                                minute: 'numeric',
                                                second: 'numeric'
                                            }) || "Not verified"}</td>
                                            <td className="py-4 px-4" onClick={() => navigate("/zen/users/detail/" + user.uuid)}>   <i className="bi bi-eye"></i></td>
                                        </tr>
                                    );
                                })
                            ) : (
                                <tr>
                                    <td colSpan={8} className="py-4 px-4 text-center">
                                        No user data available
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>

                    <div className="mt-4 flex justify-center items-center mb-5">
                        {/* Pagination */}
                        <Pagination
                            current={currentPage}
                            pageSize={pageSize}
                            total={userData?.payload?.total_count || 0} // Ensure your API returns a total count
                            onChange={onPageChange}
                            showSizeChanger
                            pageSizeOptions={["10", "20", "50", "100"]}
                            showTotal={(total) => `Total ${total} users`}
                        />
                    </div>
                </div>
            }
        </>
    )
}

export default AllUsersTable;
