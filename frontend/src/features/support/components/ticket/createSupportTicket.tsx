import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Form, Input, Select, Button, Upload, Card, Typography, Alert } from 'antd';
import { ChevronUp, ChevronDown, ArrowLeft, Upload as UploadIcon } from 'lucide-react';
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

    return (
        <div ref={parentRef} className="max-w-6xl mx-auto p-6 mt-10">
            <div className="flex space-x-6">
                {/* Sidebar */}
                <Card
                    className="w-64"
                    style={{
                        height: isRecentTicketsOpen ? `${parentHeight}px` : 'auto',
                        overflow: 'hidden'
                    }}
                >
                    <Button
                        icon={<ArrowLeft className="h-5 w-5" />}
                        type="text"
                        onClick={() => window.history.back()}
                        className="mb-4"
                    />

                    <Button
                        type="primary"
                        ghost
                        block
                        onClick={toggleRecentTickets}
                        icon={isRecentTicketsOpen ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                        className="mb-4"
                    >
                        Your Recent Tickets
                    </Button>

                    {isRecentTicketsOpen && (
                        <div className="space-y-2">
                            {Array.isArray(stdata) && stdata.length > 0 ? (
                                stdata.map((ticket, index) => (
                                    <Card
                                        key={ticket.uuid}
                                        size="small"
                                        className="cursor-pointer hover:shadow-md transition-shadow"
                                        onClick={() => handleNavigation(ticket.uuid)}
                                    >
                                        <Title level={5} className="m-0">
                                            #{index + 1 + ticket.ticket_number} - {ticket.subject}
                                        </Title>
                                        <Text type="secondary">
                                            {ticket.status.charAt(0).toUpperCase() + ticket.status.slice(1)}
                                            {ticket.last_reply && ` - Last reply: ${new Date(ticket.last_reply).toLocaleDateString()}`}
                                        </Text>
                                    </Card>
                                ))
                            ) : (
                                <Alert message="You have not created any tickets" type="info" />
                            )}
                        </div>
                    )}
                </Card>

                {/* Form */}
                <Card className="flex-1">
                    <Title level={2}>Create new Support Request</Title>
                    <Text className="block mb-6">
                        If you can't find a solution to your problems in our knowledgebase,
                        you can submit a ticket. An agent will respond to you soonest.
                    </Text>

                    <Form
                        form={form}
                        layout="vertical"
                        onFinish={handleSubmit}
                        initialValues={{
                            name: user,
                            email: email,
                            priority: 'Medium'
                        }}
                    >
                        <div className="grid grid-cols-2 gap-6">
                            <Form.Item
                                name="name"
                                label="Name"
                                rules={[{ required: true, message: 'Name is required' }]}
                            >
                                <Input />
                            </Form.Item>

                            <Form.Item
                                name="email"
                                label="Email Address"
                                rules={[
                                    { required: true, message: 'Email is required' },
                                    { type: 'email', message: 'Invalid email format' }
                                ]}
                            >
                                <Input />
                            </Form.Item>
                        </div>

                        <Form.Item
                            name="priority"
                            label="Priority"
                            rules={[{ required: true, message: 'Priority is required' }]}
                        >
                            <Select>
                                <Option value="Low">Low</Option>
                                <Option value="Medium">Medium</Option>
                                <Option value="High">High</Option>
                            </Select>
                        </Form.Item>

                        <Form.Item
                            name="subject"
                            label="Subject"
                            rules={[{ required: true, message: 'Subject is required' }]}
                        >
                            <Input />
                        </Form.Item>

                        <Form.Item
                            name="message"
                            label="Message"
                            rules={[{ required: true, message: 'Message is required' }]}
                        >
                            <TextArea rows={6} />
                        </Form.Item>

                        <Form.Item
                            label={`Attachments (Max ${MAX_FILES} files)`}
                        >
                            <Upload
                                fileList={fileList}
                                onChange={handleFileChange}
                                beforeUpload={() => false}
                                accept={ALLOWED_FILE_TYPES}
                                multiple
                                maxCount={MAX_FILES}
                            >
                                <Button icon={<UploadIcon className="h-4 w-4" />} disabled={fileList.length >= MAX_FILES}>
                                    Select Files
                                </Button>
                            </Upload>
                            <Text type="secondary" className="mt-2 block">
                                Allowed File Extensions: .jpg, .gif, .jpeg, .png, .txt, .pdf (Max file size: 1024MB)
                            </Text>
                        </Form.Item>

                        <Form.Item>
                            <div className="flex space-x-4">
                                <Button type="primary" htmlType="submit">
                                    Submit
                                </Button>
                                <Button onClick={() => form.resetFields()}>
                                    Cancel
                                </Button>
                            </div>
                        </Form.Item>
                    </Form>
                </Card>
            </div>
        </div>
    );
};

export default SupportRequestForm;