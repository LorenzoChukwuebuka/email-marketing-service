import { Modal, Card, Button, Tag, Divider, Avatar} from "antd";
import { CloseOutlined, MessageOutlined, UserOutlined } from '@ant-design/icons';
import { MessageSquare, User, Calendar, Clock, AlertCircle } from 'lucide-react';
import { Ticket } from '../../../interface/support.interface';

interface TicketSidebarProps {
    ticketData: Ticket;
    onClose: (id: string) => Promise<void>;
}

const TicketSidebar: React.FC<TicketSidebarProps> = ({ ticketData, onClose }) => {
    const closeTkt = async () => {
        Modal.confirm({
            title: "Close Ticket",
            content: "Are you sure you want to close this ticket? This action cannot be undone.",
            okText: "Yes, Close",
            cancelText: "Cancel",
            okButtonProps: { danger: true },
            onOk: async () => {
                await onClose(ticketData.id);
                await new Promise(resolve => setTimeout(resolve, 1000));
                location.reload();
            },
        });
    };

    const getPriorityColor = (priority: string) => {
        switch (priority?.toLowerCase()) {
            case 'high': return 'error';
            case 'medium': return 'warning';
            case 'low': return 'success';
            default: return 'default';
        }
    };

    const getStatusColor = (status: string) => {
        switch (status?.toLowerCase()) {
            case 'open': return 'processing';
            case 'closed': return 'default';
            case 'pending': return 'warning';
            case 'resolved': return 'success';
            default: return 'default';
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString('en-US', {
            timeZone: 'UTC',
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    return (
        <div className="w-full lg:w-80">
            <Card 
                className="shadow-lg border-0 backdrop-blur-sm bg-white/90 sticky top-6"
                title={
                    <div className="flex items-center gap-2">
                        <MessageSquare className="h-5 w-5 text-blue-600" />
                        <span>Ticket Information</span>
                    </div>
                }
            >
                <div className="space-y-6">
                    {/* Requestor */}
                    <div>
                        <div className="flex items-center gap-2 mb-3">
                            <User className="h-4 w-4 text-gray-500" />
                            <span className="text-sm font-medium text-gray-700">Requestor</span>
                        </div>
                        <div className="flex items-center gap-3">
                            <Avatar icon={<UserOutlined />} className="bg-blue-600" />
                            <div>
                                <p className="font-semibold text-gray-800 m-0">{ticketData.name}</p>
                                <Tag color="blue" className="text-xs mt-1">Authorized User</Tag>
                            </div>
                        </div>
                    </div>

                    <Divider className="my-4" />

                    {/* Status & Priority */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <div className="flex items-center gap-2 mb-2">
                                <AlertCircle className="h-4 w-4 text-gray-500" />
                                <span className="text-sm font-medium text-gray-700">Status</span>
                            </div>
                            <Tag 
                                color={getStatusColor(ticketData.status)} 
                                className="text-sm font-medium"
                            >
                                {ticketData.status?.toUpperCase()}
                            </Tag>
                        </div>
                        <div>
                            <div className="flex items-center gap-2 mb-2">
                                <AlertCircle className="h-4 w-4 text-gray-500" />
                                <span className="text-sm font-medium text-gray-700">Priority</span>
                            </div>
                            <Tag 
                                color={getPriorityColor(ticketData.priority)} 
                                className="text-sm font-medium"
                            >
                                {ticketData.priority?.toUpperCase()}
                            </Tag>
                        </div>
                    </div>

                    <Divider className="my-4" />

                    {/* Dates */}
                    <div className="space-y-4">
                        <div>
                            <div className="flex items-center gap-2 mb-2">
                                <Calendar className="h-4 w-4 text-gray-500" />
                                <span className="text-sm font-medium text-gray-700">Submitted</span>
                            </div>
                            <p className="text-sm text-gray-600 bg-gray-50 p-2 rounded">
                                {formatDate(ticketData.created_at)}
                            </p>
                        </div>

                        <div>
                            <div className="flex items-center gap-2 mb-2">
                                <Clock className="h-4 w-4 text-gray-500" />
                                <span className="text-sm font-medium text-gray-700">Last Updated</span>
                            </div>
                            <p className="text-sm text-gray-600 bg-gray-50 p-2 rounded">
                                {ticketData.last_reply != null
                                    ? formatDate(ticketData.last_reply)
                                    : "No replies yet"}
                            </p>
                        </div>
                    </div>

                    <Divider className="my-4" />

                    {/* Actions */}
                    <div className="space-y-3">
                        <Button 
                            type="primary"
                            block
                            size="large"
                            icon={<MessageOutlined />}
                            className="bg-gradient-to-r from-blue-600 to-purple-600 border-none"
                            onClick={() => {
                                const element = document.getElementById('replyTicket');
                                if (element) {
                                    element.scrollIntoView({ behavior: 'smooth' });
                                }
                            }}
                        >
                            Reply to Ticket
                        </Button>
                        
                        {ticketData.status !== 'closed' && (
                            <Button 
                                danger
                                block
                                size="large"
                                icon={<CloseOutlined />}
                                onClick={closeTkt}
                            >
                                Close Ticket
                            </Button>
                        )}
                    </div>
                </div>
            </Card>
        </div>
    );
};

export default TicketSidebar;