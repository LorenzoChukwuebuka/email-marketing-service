import { Form, Input, Modal, Checkbox, Button } from 'antd';
import { useEffect, useState } from 'react';
import useContactStore from "./../../store/contact.store";


interface CreateContactProps {
    isOpen: boolean;
    onClose: () => void;
    refetch: () => void; // Optional refetch function to refresh contacts after creation   
}

const CreateContact: React.FC<CreateContactProps> = ({ isOpen, onClose, refetch }) => {
    const [form] = Form.useForm();
    const { createContact, setContactFormValues } = useContactStore();
    const [isLoading, setIsLoading] = useState<boolean>(false)

    // Reset form when modal closes
    useEffect(() => {
        if (!isOpen) {
            form.resetFields();
        }
    }, [isOpen, form]);

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const handleSubmit = async (values: any) => {
        console.log(values)
        try {
            setIsLoading(true)
            setContactFormValues({ ...values })
            await createContact();
            await new Promise((resolve) => setTimeout(resolve, 3000));
            refetch()
            onClose();
            form.resetFields();
        } catch (error) {
            console.error('Failed to create contact:', error);
        } finally {
            setIsLoading(false)
        }
    };

    return (
        <Modal
            title="Create Contact"
            open={isOpen}
            onCancel={onClose}
            footer={null}
        >
            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
                initialValues={{
                    is_subscribed: false,
                }}
            >
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

                <Form.Item
                    name="is_subscribed"
                    valuePropName="checked"
                    rules={[
                        {
                            validator: (_, value) =>
                                value
                                    ? Promise.resolve()
                                    : Promise.reject(new Error('You must agree to subscribe')),
                        },
                    ]}
                >
                    <Checkbox>
                        I confirm that this person gave me permission to email them
                    </Checkbox>
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

export default CreateContact;