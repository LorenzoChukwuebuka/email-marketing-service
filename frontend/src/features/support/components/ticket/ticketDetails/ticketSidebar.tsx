// TicketSidebar.tsx
import { Modal } from "antd";
import { Ticket } from '../../../interface/support.interface';

interface TicketSidebarProps {
    ticketData: Ticket;
    onClose: (id: string) => Promise<void>;
}

const TicketSidebar: React.FC<TicketSidebarProps> = ({ ticketData, onClose }) => {
    const closeTkt = async () => {
        Modal.confirm({
            title: "Are you sure?",
            content: "Do you want to close this ticket?",
            okText: "Yes",
            cancelText: "No",
            onOk: async () => {
                await onClose(ticketData.uuid);
                await new Promise(resolve => setTimeout(resolve, 1000));
                location.reload();
            },
        });
    };

    return (
        <div className="w-full lg:w-1/4 bg-white rounded-lg p-4 border-r lg:sticky lg:top-5 h-auto lg:h-full">
            <div className="mb-4">
                <h2 className="text-lg font-bold">Ticket Information</h2>
                <p className="text-sm">Requestor</p>
                <p className="font-semibold">{ticketData.name} <span className="text-xs text-gray-500">Authorized User</span></p>
            </div>
            <div className="mb-4">
                <p className="text-sm">Submitted</p>
                <p className="font-semibold">{new Date(ticketData.created_at).toLocaleString('en-US', {
                    timeZone: 'UTC',
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                    hour: 'numeric',
                    minute: 'numeric',
                    second: 'numeric',
                })}</p>
            </div>
            <div className="mb-4">
                <p className="text-sm">Last Updated</p>
                <p className="font-semibold">
                    {ticketData.last_reply != null
                        ? new Date(ticketData.last_reply).toLocaleString('en-US', {
                            timeZone: 'UTC',
                            year: 'numeric',
                            month: 'long',
                            day: 'numeric',
                            hour: 'numeric',
                            minute: 'numeric',
                            second: 'numeric',
                        })
                        : "Ticket has not been replied to"}
                </p>
            </div>
            <div className="mb-4">
                <p className="text-sm text-gray-600">Status</p>
                <p className={`font-semibold text-yellow-500 flex items-center space-x-2`}>
                    {ticketData.status}
                </p>
                <p className="text-sm text-gray-600 mb-1">Priority</p>
                <p className="font-semibold flex items-center space-x-2">
                    <span className={`text-sm text-white bg-black  ${ticketData.priority === "high" ? "bg-red-500" :
                        ticketData.priority === "medium" ? "bg-yellow-500" :
                            ticketData.priority === "low" ? "bg-blue-500" : ""
                        } py-1 px-2 rounded-lg`}>
                        {ticketData.priority}
                    </span>
                </p>
            </div>
            <div className="flex space-x-4">
                <button className="bg-green-500 text-white py-2 px-4 rounded-lg">
                    <a href="#replyTicket">Reply</a>
                </button>
                <button className="bg-red-500 text-white py-2 px-4 rounded-lg" onClick={closeTkt}>Close</button>
            </div>
        </div>
    );
};

export default TicketSidebar