import { useState } from "react"
import { useNavigate } from "react-router-dom";
import useDebounce from "../../../../hooks/useDebounce";
import { useAdminClosedTicketsQuery } from "../../hooks/useAdminSupportTicketQuery";
import { Pagination } from "antd";
import LoadingSpinnerComponent from "../../../../components/loadingSpinnerComponent";

const tableHeaders = [
    "Name",
    "Email",
    "Subject",
    "Description",
    "Ticket Number",
    "Status",
    "Priority",
    "Last Reply",
    "View"
];

const ClosedSupportTicketComponentTable: React.FC = () => {
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [searchQuery, setSearchQuery] = useState<string>(""); // New state for search query
    // Debounce the search query
    const debouncedSearchQuery = useDebounce(searchQuery, 300); // 300ms delay
    const navigate = useNavigate()

    const handleSearch = (query: string) => {
        setSearchQuery(query);
    }

    const { data: supportData, isLoading } = useAdminClosedTicketsQuery(currentPage, pageSize, debouncedSearchQuery)

    const onPageChange = (page: number, size: number) => {
        setCurrentPage(page);
        setPageSize(size);
    };

    return <>

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
        {isLoading ? <LoadingSpinnerComponent /> : (
            <div className="overflow-x-auto mt-4">
                <h1 className="font-semibold text-lg mt-4 mb-4"> Closed Ticket List </h1>
                <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                    <thead className="bg-gray-50">
                        <tr>
                            {tableHeaders.map((header, index) => (
                                <th key={index} className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    {header}
                                </th>
                            ))}

                        </tr>
                    </thead>

                    {/* Commenting out the table body as requested */}
                    <tbody className="divide-y divide-gray-200">
                        {Array.isArray(supportData) && supportData && supportData.length > 0 ? (
                            supportData.map((ticket, index) => {
                                return (
                                    <tr key={index} className="hover:bg-gray-100">
                                        <td className="py-4 px-4">{ticket?.name}</td>
                                        <td className="py-4 px-4">{ticket?.email}</td>
                                        <td className="py-4 px-4">{ticket?.subject}</td>
                                        <td className="py-4 px-4">{ticket?.description}</td>
                                        <td className="py-4 px-4">{"#" + ticket?.ticket_number}</td>
                                        <td
                                            className={`py-4 px-4 ${ticket?.status === 'closed' ? 'bg-black text-white' :
                                                ticket?.status === 'pending' ? 'bg-red-500 text-white' :
                                                    ticket?.status === 'open' ? 'bg-green-500 text-white' :
                                                        ticket?.status === 'resolved' ? 'bg-green-500 text-white' : ''
                                                }`}
                                        >
                                            {ticket?.status}
                                        </td>

                                        <td
                                            className={`py-4 px-4 ${ticket?.priority === 'high' ? 'bg-red-500 text-white' :
                                                ticket?.priority === 'medium' ? 'bg-yellow-500 text-black' :
                                                    ticket?.priority === 'low' ? 'bg-blue-500 text-white' : '' // Use any suitable color for 'low'
                                                }`}
                                        >
                                            {ticket?.priority}
                                        </td>

                                        <td className="py-4 px-4">{ticket.last_reply && new Date(ticket.last_reply).toLocaleString('en-US', {
                                            timeZone: 'UTC',
                                            year: 'numeric',
                                            month: 'long',
                                            day: 'numeric',
                                            hour: 'numeric',
                                            minute: 'numeric',
                                            second: 'numeric'
                                        }) || "Not available"}</td>
                                        <td className="py-4 px-4" onClick={() => navigate("/zen/support/details/" + ticket.uuid)}> <i className="bi bi-eye"></i> </td>
                                    </tr>
                                );
                            })
                        ) : (
                            <tr>
                                <td colSpan={8} className="py-4 px-4 text-center">
                                    No ticket data available
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
                        total={supportData?.payload?.total_count || 0} // Ensure your API returns a total count
                        onChange={onPageChange}
                        showSizeChanger
                        pageSizeOptions={["10", "20", "50", "100"]}
                        // showTotal={(total) => `Total ${total} Contacts`}
                    />
                </div>

            </div>
        )}



    </>
}

export default ClosedSupportTicketComponentTable