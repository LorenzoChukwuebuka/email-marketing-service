import { useEffect } from "react"
import useAdminSupportStore from "../../../../../store/admin/AdminSupport"
import Pagination from '../../../../../components/Pagination';
import { useNavigate } from "react-router-dom";

const AllSupportTicketComponentTable: React.FC = () => {
    const { supportData, getAllTickets, paginationInfo } = useAdminSupportStore()
    const navigate = useNavigate()
    const handleSearch = (query: string) => { }

    const handlePageChange = () => {
        getAllTickets(paginationInfo.page_size)
    }

    useEffect(() => {
        getAllTickets()
    }, [getAllTickets])
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

        <div className="overflow-x-auto mt-4">
            <h1 className="font-semibold text-lg mt-4 mb-4"> Ticket List </h1>
            <table className="md:min-w-5xl min-w-full w-full rounded-sm bg-white">
                <thead className="bg-gray-50">
                    <tr>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Name
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Email
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Subject
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Description
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Ticket Number
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Status
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Priority
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Last Reply
                        </th>
                        <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            View
                        </th>
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
                                    <td className="py-4 px-4">{ticket?.description || "N/A"}</td>
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
                                    <td className="py-4 cursor-pointer px-4" onClick={() => navigate("/zen/dash/support/details/" + ticket.uuid)}> <i className="bi bi-eye"></i> </td>
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

            {/* Commenting out the pagination component as requested */}
            <Pagination paginationInfo={paginationInfo} handlePageChange={handlePageChange} item={"All Tickets"} />
        </div>


    </>
}

export default AllSupportTicketComponentTable