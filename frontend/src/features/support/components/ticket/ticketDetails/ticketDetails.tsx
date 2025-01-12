import { useMemo } from "react";
import { useParams } from "react-router-dom";
import Cookies from "js-cookie";
import { Helmet, HelmetProvider } from "react-helmet-async";
import { Tooltip } from "antd";
import { QuestionCircleOutlined } from '@ant-design/icons';
import useSupportStore from "../../../store/support.store";
import { TicketFile, Ticket } from '../../../interface/support.interface';
import { useTicketDetailsQuery } from "../../../hooks/useSupporTicketQuery";
import TicketSidebar from "./ticketSidebar";
import TicketReplyForm from "./ticketReplyForm";

const TicketDetails: React.FC = () => {
    const { replyTicket, closeTicket } = useSupportStore();
    const { id } = useParams<{ id: string }>() as { id: string };

    const cookie = Cookies.get("Cookies");
    const user = cookie ? JSON.parse(cookie)?.details?.fullname : "";
    const email = cookie ? JSON.parse(cookie)?.details?.email : "";

    const { data: supportTicketData, isLoading } = useTicketDetailsQuery(id);
    const sTicketData = useMemo(() => supportTicketData?.payload || null, [supportTicketData]);

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: 'numeric',
            minute: 'numeric',
            second: 'numeric',
            timeZone: 'UTC'
        });
    };

    const renderAttachments = (files: TicketFile[]) => {
        if (!files || files.length === 0) return null;
        const baseUrl = import.meta.env.VITE_BASE_API_URL as string;
        return (
            <div className="mt-4">
                <h4 className="text-sm font-semibold mb-2">Attachments:</h4>
                <div className="flex flex-wrap gap-2">
                    {files.map((file, index) => (
                        <a
                            key={index}
                            href={`${baseUrl}/${file.file_path}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-gray-200 text-gray-800 hover:bg-gray-300"
                        >
                            <svg className="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                                <path fillRule="evenodd" d="M8 4a3 3 0 00-3 3v4a5 5 0 0010 0V7a1 1 0 112 0v4a7 7 0 11-14 0V7a5 5 0 0110 0v4a3 3 0 11-6 0V7a1 1 0 012 0v4a1 1 0 102 0V7a3 3 0 00-3-3z" clipRule="evenodd" />
                            </svg>
                            {file.file_name}
                        </a>
                    ))}
                </div>
            </div>
        );
    };

    if (!supportTicketData) {
        return (
            <div className="flex items-center justify-center mt-20">
                <p>No ticket details available.</p>
            </div>
        );
    }

    return (
        <HelmetProvider>
            <Helmet title={`Ticket #${sTicketData?.ticket_number}`} />
            {isLoading ? (
                <div className="flex items-center justify-center mt-20">
                    <span className="loading loading-spinner loading-lg"></span>
                </div>
            ) : (
                <div className="flex flex-col lg:flex-row mb-10 mt-5 p-4 bg-gray-100 min-h-screen">
                    <TicketSidebar ticketData={sTicketData as Ticket} onClose={closeTicket} />
                    <div className="w-full lg:w-3/4 -mt-5 p-6">
                        <div className="bg-white p-6 rounded-lg shadow mb-10">
                            <div className="flex justify-between mb-6">
                                <button
                                    className="text-blue-600 mr-2 tooltip tooltip-right"
                                    data-tip="Go Back"
                                    onClick={() => window.history.back()}
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                                    </svg>
                                </button>
                                <h2 className="text-xl font-semibold">View Ticket #{sTicketData?.ticket_number}</h2>
                                <span className={`${sTicketData?.status === 'closed' ? 'bg-black text-white' :
                                    sTicketData?.status === 'pending' ? 'bg-red-500 text-white' :
                                        sTicketData?.status === 'open' ? 'bg-green-500 text-white' :
                                            sTicketData?.status === 'resolved' ? 'bg-green-500 text-white' : ''
                                    } px-2 py-1 rounded`}>
                                    {sTicketData?.status}
                                </span>
                                <Tooltip title="You can reopen a ticket by replying to the ticket">
                                    <QuestionCircleOutlined />
                                </Tooltip>
                            </div>

                            {sTicketData?.messages &&
                                sTicketData?.messages.map((message) => (
                                    <div key={message.uuid} className="mb-6 border-t border-b py-7 pt-4">
                                        <div className="flex justify-between bg-gray-200 p-2">
                                            <p className="text-sm text-gray-700">
                                                Posted by <span className="font-semibold">
                                                    {message.is_admin ?
                                                        message.admin.firstname + " " + message.admin.lastname :
                                                        sTicketData.name}
                                                </span>
                                                on {formatDate(message.created_at)}
                                            </p>
                                            <span className="text-sm bg-blue-500 text-white rounded-md p-1">
                                                {message.is_admin ? "admin/operator" : "authorized user"}
                                            </span>
                                        </div>
                                        <p className="mt-2">{message.message}</p>
                                        {renderAttachments(message.files)}
                                    </div>
                                ))}
                        </div>

                        <TicketReplyForm
                            user={user}
                            email={email}
                            onSubmit={async (message, files) => {
                                await replyTicket(id, message, files);
                            }}
                        />
                    </div>
                </div>
            )}
        </HelmetProvider>
    );
};

export default TicketDetails;