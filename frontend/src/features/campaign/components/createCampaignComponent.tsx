import { useState } from "react";
import { Modal, Form, Input, Button } from "antd";
import useCampaignStore from "../store/campaign.store";

interface Props {
    isOpen: boolean;
    onClose: () => void;
}

const CreateCampaignComponent: React.FC<Props> = ({ isOpen, onClose }) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);

    const { createCampaign,  setCreateCampaignValues } = useCampaignStore();

    const handleSubmit = async (values: { name: string }) => {
        try {
            setLoading(true);
            setCreateCampaignValues(values);
            await createCampaign();
            await new Promise(resolve => setTimeout(resolve, 500));
            form.resetFields();
            location.reload()
            onClose();
        } catch (error) {
            console.error('Failed to create campaign:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleCancel = () => {
        form.resetFields();
        onClose();
    };

    return (
        <Modal
            title="Create Campaign"
            open={isOpen}
            onCancel={handleCancel}
            footer={null}
            maskClosable={false}
        >
            <p className="mt-2 mb-5">
                Keep subscribers engaged by sharing your latest news, promoting your bestselling products, or announcing an upcoming event.
            </p>

            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
                autoComplete="off"
            >
                <Form.Item
                    name="name"
                    label="Campaign name"
                    rules={[
                        {
                            required: true,
                            message: 'Campaign name is required',
                        },
                        {
                            min: 3,
                            message: 'Campaign name must be at least 3 characters',
                        },
                        {
                            max: 50,
                            message: 'Campaign name cannot exceed 50 characters',
                        }
                    ]}
                >
                    <Input placeholder="Enter campaign name" />
                </Form.Item>

                <Form.Item className="flex justify-end mb-0">
                    <Button className="mr-2" onClick={handleCancel}>
                        Cancel
                    </Button>
                    <Button type="primary" htmlType="submit" loading={loading}>
                        Create
                    </Button>
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default CreateCampaignComponent;