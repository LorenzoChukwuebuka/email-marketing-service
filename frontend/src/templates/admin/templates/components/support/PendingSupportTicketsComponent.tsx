const PendingSupportTicketComponentTable: React.FC = () => {
    const handleSearch = (query: string) => { }
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
            <h1 className="font-semibold text-lg mt-4 mb-4">Pending Ticket List </h1>
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
                    </tr>
                </thead>

                {/* Commenting out the table body as requested */}
                {/* <tbody className="divide-y divide-gray-200">
            {Array.isArray(ticketData) && ticketData && ticketData.length > 0 ? (
                ticketData.map((ticket, index) => {
                    return (
                        <tr key={index} className="hover:bg-gray-100">
                            <td className="py-4 px-4">{ticket?.name}</td>
                            <td className="py-4 px-4">{ticket?.email}</td>
                            <td className="py-4 px-4">{ticket?.subject}</td>
                            <td className="py-4 px-4">{ticket?.description}</td>
                            <td className="py-4 px-4">{ticket?.ticket_number}</td>
                            <td className="py-4 px-4">{ticket?.status}</td>
                            <td className="py-4 px-4">{ticket?.priority}</td>
                            <td className="py-4 px-4">{ticket?.last_reply}</td>
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
        </tbody> */}

            </table>

            {/* Commenting out the pagination component as requested */}
            {/* <Pagination paginationInfo={paginationInfo} handlePageChange={handlePageChange} item={"Pending Tickets"} /> */}
        </div>


    </>
}

export default PendingSupportTicketComponentTable