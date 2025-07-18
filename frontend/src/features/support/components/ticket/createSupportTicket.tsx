import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Form, Input, Select, Button, Upload, Card, Typography, Alert, Badge, Divider } from 'antd';
import { ChevronUp, ChevronDown, ArrowLeft, Upload as UploadIcon, MessageSquare, Clock, User, Mail, AlertCircle, FileText, Paperclip } from 'lucide-react';
import type { UploadFile } from 'antd/es/upload/interface';
import Cookies from "js-cookie";
import { useNavigate } from 'react-router-dom';
import useSupportStore from '../../store/support.store';
import { useSupportTicketQuery } from '../../hooks/useSupporTicketQuery';

const { Title, Text } = Typography;
const { TextArea } = Input;
const { Option } = Select;

const MAX_FILES = 3;
const ALLOWED_FILE_TYPES = '.jpg,.gif,.jpeg,.png,.txt,.pdf';

const SupportRequestForm: React.FC = () => {
    const [form] = Form.useForm();
    const [isRecentTicketsOpen, setIsRecentTicketsOpen] = useState(false);
    const { createTicket } = useSupportStore();
    const [fileList, setFileList] = useState<UploadFile[]>([]);
    const parentRef = useRef<HTMLDivElement>(null);
    const [parentHeight, setParentHeight] = useState<number | null>(null);
    const navigate = useNavigate();

    // Get user details from cookies
    const cookie = Cookies.get("Cookies");
    const user = cookie ? JSON.parse(cookie)?.details?.fullname : "";
    const email = cookie ? JSON.parse(cookie)?.details?.email : "";

    const { data: supportTicketData } = useSupportTicketQuery()

    const stdata = useMemo(() => supportTicketData?.payload || [], [supportTicketData])

    useEffect(() => {
        if (parentRef.current) {
            setParentHeight(parentRef.current.clientHeight);
        }
    }, [isRecentTicketsOpen]);

    const handleSubmit = async (values: any) => {
        try {
            const files = fileList.map(file => file.originFileObj).filter(Boolean) as File[];
            await createTicket({
                ...values,
                priority: values.priority.toLowerCase()
            }, files);
            await new Promise(resolve => setTimeout(resolve, 700));
            form.resetFields();
            setFileList([]);
        } catch (error) {
            console.error('Submit error:', error);
        }
    };

    const handleFileChange = ({ fileList: newFileList }: any) => {
        setFileList(newFileList.slice(0, MAX_FILES));
    };

    const toggleRecentTickets = () => {
        setIsRecentTicketsOpen(!isRecentTicketsOpen);
    };

    const handleNavigation = (uuid: string) => {
        navigate("/app/support/ticket/details/" + uuid);
    };

    const getPriorityColor = (priority: string) => {
        switch (priority.toLowerCase()) {
            case 'high': return 'error';
            case 'medium': return 'warning';
            case 'low': return 'success';
            default: return 'default';
        }
    };

    const getStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case 'open': return 'processing';
            case 'closed': return 'success';
            case 'pending': return 'warning';
            default: return 'default';
        }
    };

    return (
        <div ref={parentRef} className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 p-6">
            <div className="max-w-7xl mx-auto">
                <div className="flex gap-8">
                    {/* Sidebar */}
                    <div className="w-80">
                        <Card
                            className="shadow-lg border-0 backdrop-blur-sm bg-white/80"
                            style={{
                                height: isRecentTicketsOpen ? `${parentHeight}px` : 'auto',
                                overflow: 'hidden'
                            }}
                        >
                            <div className="space-y-4">
                                <Button
                                    icon={<ArrowLeft className="h-4 w-4" />}
                                    type="text"
                                    onClick={() => window.history.back()}
                                    className="flex items-center gap-2 text-gray-600 hover:text-blue-600 transition-colors"
                                >
                                    Back
                                </Button>

                                <Divider className="my-4" />

                                <Button
                                    type="primary"
                                    ghost
                                    block
                                    onClick={toggleRecentTickets}
                                    icon={isRecentTicketsOpen ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                                    className="h-12 text-left flex items-center justify-between border-blue-200 hover:border-blue-400 transition-all duration-200"
                                >
                                    <span className="flex items-center gap-2">
                                        <MessageSquare className="h-4 w-4" />
                                        Recent Tickets
                                    </span>
                                </Button>

                                {isRecentTicketsOpen && (
                                    <div className="space-y-3 mt-4 max-h-96 overflow-y-auto">
                                        {Array.isArray(stdata) && stdata.length > 0 ? (
                                            stdata.map((ticket, index) => (
                                                <Card
                                                    key={ticket.id}
                                                    size="small"
                                                    className="cursor-pointer hover:shadow-md transition-all duration-200 border-l-4 border-l-blue-500 hover:border-l-blue-600"
                                                    onClick={() => handleNavigation(ticket.id)}
                                                >
                                                    <div className="space-y-2">
                                                        <div className="flex items-start justify-between">
                                                            <Title level={5} className="m-0 text-gray-800 text-sm">
                                                                #{index + 1 + ticket.ticket_number}
                                                            </Title>
                                                            <Badge 
                                                                status={getStatusColor(ticket.status)} 
                                                                text={ticket.status.charAt(0).toUpperCase() + ticket.status.slice(1)}
                                                                className="text-xs"
                                                            />
                                                        </div>
                                                        <Text className="text-gray-600 text-sm line-clamp-2">
                                                            {ticket.subject}
                                                        </Text>
                                                        {ticket.last_reply && (
                                                            <div className="flex items-center gap-1 text-gray-500">
                                                                <Clock className="h-3 w-3" />
                                                                <Text className="text-xs">
                                                                    {new Date(ticket.last_reply).toLocaleDateString()}
                                                                </Text>
                                                            </div>
                                                        )}
                                                    </div>
                                                </Card>
                                            ))
                                        ) : (
                                            <Alert 
                                                message="No tickets found" 
                                                description="You haven't created any support tickets yet."
                                                type="info" 
                                                showIcon
                                                className="border-blue-200 bg-blue-50"
                                            />
                                        )}
                                    </div>
                                )}
                            </div>
                        </Card>
                    </div>

                    {/* Main Form */}
                    <div className="flex-1">
                        <Card className="shadow-xl border-0 backdrop-blur-sm bg-white/90">
                            <div className="space-y-6">
                                {/* Header */}
                                <div className="bg-gradient-to-r from-blue-600 to-purple-600 -m-6 mb-6 p-8 rounded-t-lg">
                                    <Title level={1} className="text-white m-0 flex items-center gap-3">
                                        <MessageSquare className="h-8 w-8" />
                                        Create Support Request
                                    </Title>
                                    <Text className="text-blue-100 text-lg mt-2">
                                        Can't find a solution? Submit a ticket and our team will help you out.
                                    </Text>
                                </div>

                                <Form
                                    form={form}
                                    layout="vertical"
                                    onFinish={handleSubmit}
                                    initialValues={{
                                        name: user,
                                        email: email,
                                        priority: 'Medium'
                                    }}
                                    className="space-y-6"
                                >
                                    {/* Personal Information */}
                                    <Card 
                                        title={
                                            <span className="flex items-center gap-2 text-gray-700">
                                                <User className="h-5 w-5" />
                                                Personal Information
                                            </span>
                                        }
                                        className="border-gray-200 shadow-sm"
                                    >
                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                            <Form.Item
                                                name="name"
                                                label="Full Name"
                                                rules={[{ required: true, message: 'Name is required' }]}
                                            >
                                                <Input 
                                                    prefix={<User className="h-4 w-4 text-gray-400" />}
                                                    placeholder="Enter your full name"
                                                    className="h-11"
                                                />
                                            </Form.Item>

                                            <Form.Item
                                                name="email"
                                                label="Email Address"
                                                rules={[
                                                    { required: true, message: 'Email is required' },
                                                    { type: 'email', message: 'Invalid email format' }
                                                ]}
                                            >
                                                <Input 
                                                    prefix={<Mail className="h-4 w-4 text-gray-400" />}
                                                    placeholder="Enter your email address"
                                                    className="h-11"
                                                />
                                            </Form.Item>
                                        </div>
                                    </Card>

                                    {/* Ticket Details */}
                                    <Card 
                                        title={
                                            <span className="flex items-center gap-2 text-gray-700">
                                                <FileText className="h-5 w-5" />
                                                Ticket Details
                                            </span>
                                        }
                                        className="border-gray-200 shadow-sm"
                                    >
                                        <div className="space-y-6">
                                            <Form.Item
                                                name="priority"
                                                label="Priority Level"
                                                rules={[{ required: true, message: 'Priority is required' }]}
                                            >
                                                <Select 
                                                    placeholder="Select priority level"
                                                    className="h-11"
                                                    suffixIcon={<AlertCircle className="h-4 w-4 text-gray-400" />}
                                                >
                                                    <Option value="Low">
                                                        <Badge status="success" text="Low Priority" />
                                                    </Option>
                                                    <Option value="Medium">
                                                        <Badge status="warning" text="Medium Priority" />
                                                    </Option>
                                                    <Option value="High">
                                                        <Badge status="error" text="High Priority" />
                                                    </Option>
                                                </Select>
                                            </Form.Item>

                                            <Form.Item
                                                name="subject"
                                                label="Subject"
                                                rules={[{ required: true, message: 'Subject is required' }]}
                                            >
                                                <Input 
                                                    placeholder="Brief description of your issue"
                                                    className="h-11"
                                                />
                                            </Form.Item>

                                            <Form.Item
                                                name="message"
                                                label="Detailed Message"
                                                rules={[{ required: true, message: 'Message is required' }]}
                                            >
                                                <TextArea 
                                                    rows={6} 
                                                    placeholder="Please provide detailed information about your issue..."
                                                    className="resize-none"
                                                />
                                            </Form.Item>
                                        </div>
                                    </Card>

                                    {/* Attachments */}
                                    <Card 
                                        title={
                                            <span className="flex items-center gap-2 text-gray-700">
                                                <Paperclip className="h-5 w-5" />
                                                Attachments
                                            </span>
                                        }
                                        className="border-gray-200 shadow-sm"
                                    >
                                        <Form.Item
                                            label={`Upload Files (Maximum ${MAX_FILES} files)`}
                                        >
                                            <Upload
                                                fileList={fileList}
                                                onChange={handleFileChange}
                                                beforeUpload={() => false}
                                                accept={ALLOWED_FILE_TYPES}
                                                multiple
                                                maxCount={MAX_FILES}
                                                className="w-full"
                                            >
                                                <Button 
                                                    icon={<UploadIcon className="h-4 w-4" />} 
                                                    disabled={fileList.length >= MAX_FILES}
                                                    className="h-12 border-dashed border-2 border-blue-300 text-blue-600 hover:border-blue-400 hover:text-blue-700"
                                                    block
                                                >
                                                    {fileList.length >= MAX_FILES ? 'Maximum files reached' : 'Click to upload or drag files here'}
                                                </Button>
                                            </Upload>
                                            <div className="mt-3 p-3 bg-blue-50 rounded-lg">
                                                <Text className="text-sm text-blue-700">
                                                    <strong>Supported formats:</strong> JPG, GIF, JPEG, PNG, TXT, PDF
                                                </Text>
                                                <br />
                                                <Text className="text-sm text-blue-600">
                                                    <strong>Maximum file size:</strong> 1024MB per file
                                                </Text>
                                            </div>
                                        </Form.Item>
                                    </Card>

                                    {/* Submit Actions */}
                                    <div className="flex justify-end space-x-4 pt-6 border-t border-gray-200">
                                        <Button 
                                            size="large"
                                            onClick={() => form.resetFields()}
                                            className="min-w-24"
                                        >
                                            Cancel
                                        </Button>
                                        <Button 
                                            type="primary" 
                                            htmlType="submit"
                                            size="large"
                                            className="min-w-32 bg-gradient-to-r from-blue-600 to-purple-600 border-none"
                                        >
                                            Submit Request
                                        </Button>
                                    </div>
                                </Form>
                            </div>
                        </Card>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default SupportRequestForm;