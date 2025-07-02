import { useEffect, useState } from "react";
import { Modal, Form, Input, Button } from "antd";
import useContactStore from "./../../store/contact.store";
import { ContactAPIResponse } from "../../interface/contact.interface";

interface EditContactProps {
    isOpen: boolean;
    onClose: () => void;
    contact: ContactAPIResponse | null;
}

const EditContact: React.FC<EditContactProps> = ({ isOpen, onClose, contact }) => {
    const [form] = Form.useForm();
    const { editContact, setEditContactValues } = useContactStore();
    const [isLoading, setIsLoading] = useState(false)

    useEffect(() => {
        if (contact) {
            form.setFieldsValue({
                id: contact.contact_id,
                first_name: contact.first_name,
                last_name: contact.last_name,
                email: contact.email,
                from: contact.from_origin,
                is_subscribed: contact.is_subscribed
            });
        }
    }, [contact, form]);

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const handleSubmit = async (values: any) => {
        setIsLoading(true)
        console.log(values)
        try {
            setEditContactValues({ ...values })
            await editContact();
            onClose();
            form.resetFields();
        } catch (error) {
            console.error('Failed to edit contact:', error);
        } finally {
            setIsLoading(false)
        }
    };

    return (
        <Modal
            title="Edit Contact"
            open={isOpen}
            onCancel={onClose}
            footer={null}
        >
            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
                initialValues={{
                    id: contact?.id,
                    first_name: contact?.first_name,
                    last_name: contact?.last_name,
                    email: contact?.email,
                    from: contact?.from_origin,
                    is_subscribed: contact?.is_subscribed
                }}
            >
                <Form.Item name="id" hidden>
                    <Input />
                </Form.Item>

                <Form.Item
                    name="first_name"
                    label="First Name"
                    rules={[
                        {
                            required: true,
                            message: 'First name is required',
                        },
                    ]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    name="last_name"
                    label="Last Name"
                    rules={[
                        {
                            required: true,
                            message: 'Last name is required',
                        },
                    ]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    name="email"
                    label="Email"
                    rules={[
                        {
                            required: true,
                            message: 'Email is required',
                        },
                        {
                            type: 'email',
                            message: 'Invalid email format',
                        },
                    ]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    name="from"
                    label="From"
                    rules={[
                        {
                            required: true,
                            message: 'From field is required',
                        },
                    ]}
                >
                    <Input />
                </Form.Item>

                <div className="flex justify-end space-x-2">

                    <Form.Item >
                        <Button onClick={onClose}>
                            Cancel
                        </Button>
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" loading={isLoading} htmlType="submit">
                            {isLoading ? "Please wait..." : "Submit"}
                        </Button>
                    </Form.Item>
                </div>
            </Form>
        </Modal>
    );
};

export default EditContact;