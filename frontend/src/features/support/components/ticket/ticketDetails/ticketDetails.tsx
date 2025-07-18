import { useMemo } from "react";
import { useParams } from "react-router-dom";
import Cookies from "js-cookie";
import { Helmet, HelmetProvider } from "react-helmet-async";
import { Tooltip, Card, Badge, Avatar, Button, Spin, Empty, Tag } from "antd";
import { QuestionCircleOutlined, ArrowLeftOutlined, PaperClipOutlined, UserOutlined, CrownOutlined } from '@ant-design/icons';
import { MessageSquare, Clock, Download } from 'lucide-react';
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
            timeZone: 'UTC'
        });
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

    const renderAttachments = (files: TicketFile[]) => {
        if (!files || files.length === 0) return null;
        const baseUrl = import.meta.env.VITE_BASE_API_URL as string;
        return (
            <div className="mt-4">
                <div className="flex items-center gap-2 mb-3">
                    <PaperClipOutlined className="text-gray-500" />
                    <span className="text-sm font-medium text-gray-700">Attachments</span>
                </div>
                <div className="flex flex-wrap gap-2">
                    {files.map((file, index) => (
                        <a
                            key={index}
                            href={`${baseUrl}/${file.file_path}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center px-3 py-2 rounded-lg text-sm font-medium bg-blue-50 text-blue-700 border border-blue-200 hover:bg-blue-100 hover:border-blue-300 transition-all duration-200"
                        >
                            <Download className="w-4 h-4 mr-2" />
                            {file.file_name}
                        </a>
                    ))}
                </div>
            </div>
        );
    };

    if (!supportTicketData) {
        return (
            <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 flex items-center justify-center">
                <Empty
                    description="No ticket details available"
                    className="text-gray-500"
                />
            </div>
        );
    }

    return (
        <HelmetProvider>
            <Helmet title={`Ticket #${sTicketData?.ticket_number}`} />
            {isLoading ? (
                <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 flex items-center justify-center">
                    <Spin size="large" />
                </div>
            ) : (
                <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
                    <div className="max-w-7xl mx-auto p-6">
                        <div className="flex flex-col lg:flex-row gap-8">
                            <TicketSidebar ticketData={sTicketData as Ticket} onClose={closeTicket} />

                            <div className="flex-1">
                                {/* Header */}
                                <Card className="mb-6 shadow-lg border-0 backdrop-blur-sm bg-white/90">
                                    <div className="flex items-center justify-between mb-6">
                                        <div className="flex items-center gap-4">
                                            <Button
                                                icon={<ArrowLeftOutlined />}
                                                onClick={() => window.history.back()}
                                                className="flex items-center gap-2"
                                            >
                                                Back
                                            </Button>
                                            <div>
                                                <h1 className="text-2xl font-bold text-gray-800 m-0">
                                                    Ticket #{sTicketData?.ticket_number}
                                                </h1>
                                                <p className="text-gray-600 m-0">{sTicketData?.subject}</p>
                                            </div>
                                        </div>
                                        <div className="flex items-center gap-3">
                                            <Badge
                                                status={getStatusColor(sTicketData?.status as string)}
                                                text={sTicketData?.status?.toUpperCase()}
                                                className="text-sm font-medium"
                                            />
                                            <Tooltip title="You can reopen a ticket by replying to the ticket">
                                                <QuestionCircleOutlined className="text-gray-500 cursor-help" />
                                            </Tooltip>
                                        </div>
                                    </div>
                                </Card>

                                {/* Messages */}
                                <Card
                                    className="mb-6 shadow-lg border-0 backdrop-blur-sm bg-white/90"
                                    title={
                                        <div className="flex items-center gap-2">
                                            <MessageSquare className="h-5 w-5 text-blue-600" />
                                            <span>Conversation</span>
                                        </div>
                                    }
                                >
                                    <div className="space-y-6">
                                        {sTicketData?.messages && sTicketData?.messages.length > 0 ? (
                                            sTicketData?.messages.map((message, index) => (
                                                <div key={message.id} className="relative">
                                                    <div className="flex items-start gap-4">
                                                        <Avatar
                                                            icon={message.is_admin ? <CrownOutlined /> : <UserOutlined />}
                                                            className={`${message.is_admin ? 'bg-purple-600' : 'bg-blue-600'} flex-shrink-0`}
                                                        />
                                                        <div className="flex-1">
                                                            <div className="bg-gray-50 rounded-lg p-4 border border-gray-200">
                                                                <div className="flex items-center justify-between mb-3">
                                                                    <div className="flex items-center gap-2">
                                                                        <span className="font-semibold text-gray-800">
                                                                            {message.is_admin ?
                                                                                message.admin.firstname + " " + message.admin.lastname :
                                                                                sTicketData.name}
                                                                        </span>
                                                                        <Tag
                                                                            color={message.is_admin ? 'purple' : 'blue'}
                                                                            className="text-xs"
                                                                        >
                                                                            {message.is_admin ? "Admin" : "User"}
                                                                        </Tag>
                                                                    </div>
                                                                    <div className="flex items-center gap-1 text-gray-500 text-sm">
                                                                        <Clock className="w-4 h-4" />
                                                                        {formatDate(message.created_at)}
                                                                    </div>
                                                                </div>
                                                                <p className="text-gray-700 leading-relaxed whitespace-pre-wrap">
                                                                    {message.message}
                                                                </p>
                                                                {renderAttachments(message.files)}
                                                            </div>
                                                        </div>
                                                    </div>
                                                    {index < sTicketData.messages.length - 1 && (
                                                        <div className="flex justify-center my-6">
                                                            <div className="w-px h-6 bg-gray-300"></div>
                                                        </div>
                                                    )}
                                                </div>
                                            ))
                                        ) : (
                                            <Empty
                                                description="No messages yet"
                                                className="py-8"
                                            />
                                        )}
                                    </div>
                                </Card>

                                {/* Reply Form */}
                                <TicketReplyForm
                                    user={user}
                                    email={email}
                                    onSubmit={async (message, files) => {
                                        await replyTicket(id, message, files);
                                    }}
                                />
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </HelmetProvider>
    );
};

export default TicketDetails;
