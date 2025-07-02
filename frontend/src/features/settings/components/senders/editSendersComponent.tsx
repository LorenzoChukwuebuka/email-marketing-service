import React, { useEffect } from 'react';
import { Modal, Form, Input, Button, message } from 'antd';
import useSenderStore from '../../store/sender.store';
import { Sender } from '../../interface/sender.interface';

type Props = {
    isOpen: boolean;
    onClose: () => void;
    Sender: Sender;
}

const EditSenderComponent: React.FC<Props> = ({ isOpen, onClose, Sender }) => {
    const [form] = Form.useForm();
    const { updateSender } = useSenderStore();

    useEffect(() => {
        // Reset form fields when Sender changes
        form.setFieldsValue({
            name: Sender.name,
            email: Sender.email
        });
    }, [Sender, form]);

    const handleSubmit = async (values: { name: string; email: string }) => {
        try {
            await updateSender(Sender.id, values);
            message.success('Sender updated successfully');
            onClose();
            await new Promise(resolve => setTimeout(resolve, 500));

        } catch (error) {
            console.log(error)
            message.error('Failed to update sender');
        }
    };

    return (
        <Modal
            title="Edit Sender"
            open={isOpen}
            onCancel={onClose}
            footer={null}
        >
            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
                initialValues={{
                    name: Sender.name,
                    email: Sender.email
                }}
            >
                <p className="mb-4 text-gray-600">Edit sender details</p>

                <Form.Item
                    name="name"
                    label="Sender Name"
                    rules={[
                        { required: true, message: 'Please input the sender name!' }
                    ]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    name="email"
                    label="Sender Email"
                    rules={[
                        { required: true, message: 'Please input the sender email!' },
                        { type: 'email', message: 'Please enter a valid email!' }
                    ]}
                >
                    <Input />
                </Form.Item>

                <Form.Item className="flex justify-end mb-0">
                    <Button onClick={onClose} className="mr-2">
                        Cancel
                    </Button>
                    <Button type="primary" htmlType="submit" className="bg-blue-500">
                        Save Changes
                    </Button>
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default EditSenderComponent;