import { useState } from "react";
import { Card, Form, Input, Button, Alert, message } from "antd";
import { PaperClipOutlined, SendOutlined, DeleteOutlined } from '@ant-design/icons';
import { MessageSquare, User, Mail, Upload as UploadIcon } from 'lucide-react';
import * as yup from 'yup';

const { TextArea } = Input;

interface TicketReplyFormProps {
    user: string;
    email: string;
    onSubmit: (message: string, files: File[]) => Promise<void>;
}

const MAX_FILES = 3;

const TicketReplyForm: React.FC<TicketReplyFormProps> = ({ user, email, onSubmit }) => {
    const [form] = Form.useForm();
    const [files, setFiles] = useState<File[]>([]);
    const [loading, setLoading] = useState(false);

    const validationSchema = yup.object().shape({
        message: yup.string().required('Message is required').min(10, 'Message should be at least 10 characters'),
    });

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files) {
            const newFiles = Array.from(e.target.files);
            if (files.length + newFiles.length > MAX_FILES) {
                message.error(`You can only upload a maximum of ${MAX_FILES} files.`);
                return;
            }
            setFiles(prevFiles => [...prevFiles, ...newFiles].slice(0, MAX_FILES));
        }
    };

    const removeFile = (index: number) => {
        setFiles(prevFiles => prevFiles.filter((_, i) => i !== index));
    };

    const handleSubmit = async (values: any) => {
        try {
            setLoading(true);
            await validationSchema.validate(values, { abortEarly: false });
            await onSubmit(values.message, files);
            form.resetFields();
            setFiles([]);
            message.success('Reply sent successfully!');
        } catch (error) {
            if (error instanceof yup.ValidationError) {
                message.error(error.message);
            } else {
                message.error('Failed to send reply. Please try again.');
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <Card
            className="shadow-lg border-0 backdrop-blur-sm bg-white/90"
            title={
                <div className="flex items-center gap-2">
                    <MessageSquare className="h-5 w-5 text-blue-600" />
                    <span>Reply to Ticket</span>
                </div>
            }
            id="replyTicket"
        >
            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
                initialValues={{
                    name: user,
                    email: email,
                    message: ''
                }}
            >
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                    <Form.Item
                        name="name"
                        label="Name"
                        rules={[{ required: true, message: 'Name is required' }]}
                    >
                        <Input
                            prefix={<User className="h-4 w-4 text-gray-400" />}
                            className="h-11"
                            readOnly
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
                            className="h-11"
                            readOnly
                        />
                    </Form.Item>
                </div>

                <Form.Item
                    name="message"
                    label="Message"
                    rules={[
                        { required: true, message: 'Message is required' },
                        { min: 10, message: 'Message should be at least 10 characters' }
                    ]}
                >
                    <TextArea
                        rows={6}
                        placeholder="Type your reply here..."
                        className="resize-none"
                    />
                </Form.Item>

                <Form.Item
                    label={`Attachments (Maximum ${MAX_FILES} files)`}
                >
                    <div className="space-y-4">
                        <div className="flex items-center gap-4">
                            <input
                                type="file"
                                onChange={handleFileChange}
                                multiple
                                className="hidden"
                                id="file-upload"
                                disabled={files.length >= MAX_FILES}
                            />
                            <label
                                htmlFor="file-upload"
                                className={`inline-flex items-center gap-2 px-4 py-2 rounded-lg border-2 border-dashed cursor-pointer transition-all duration-200 ${files.length >= MAX_FILES
                                        ? 'border-gray-300 bg-gray-50 text-gray-400 cursor-not-allowed'
                                        : 'border-blue-300 bg-blue-50 text-blue-600 hover:border-blue-400 hover:bg-blue-100'
                                    }`}
                            >
                                <UploadIcon className="h-4 w-4" />
                                {files.length >= MAX_FILES ? 'Maximum files reached' : 'Choose Files'}
                            </label>
                        </div>

                        {files.length > 0 && (
                            <div className="space-y-2">
                                {files.map((file, index) => (
                                    <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg border">
                                        <div className="flex items-center gap-2">
                                            <PaperClipOutlined className="text-gray-500" />
                                            <span className="text-sm text-gray-700">{file.name}</span>
                                            <span className="text-xs text-gray-500">({(file.size / 1024 / 1024).toFixed(2)} MB)</span>
                                        </div>
                                        <Button
                                            type="text"
                                            icon={<DeleteOutlined />}
                                            onClick={() => removeFile(index)}
                                            className="text-red-500 hover:text-red-700"
                                            size="small"
                                        />
                                    </div>
                                ))}
                            </div>
                        )}

                        <Alert
                            message="File Requirements"
                            description="Supported formats: JPG, GIF, JPEG, PNG, TXT, PDF | Maximum file size: 1024MB each"
                            type="info"
                            showIcon
                            className="border-blue-200 bg-blue-50"
                        />
                    </div>
                </Form.Item>

                <div className="flex justify-end gap-4 pt-4 border-t border-gray-200">
                    <Button
                        size="large"
                        onClick={() => {
                            form.resetFields();
                            setFiles([]);
                        }}
                    >
                        Cancel
                    </Button>
                    <Button
                        type="primary"
                        htmlType="submit"
                        size="large"
                        loading={loading}
                        icon={<SendOutlined />}
                        className="bg-gradient-to-r from-blue-600 to-purple-600 border-none min-w-32"
                    >
                        Send Reply
                    </Button>
                </div>
            </Form>
        </Card>
    );
};

export default TicketReplyForm;