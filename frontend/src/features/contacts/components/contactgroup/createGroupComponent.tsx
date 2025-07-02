import { Modal, Form, Input, Button } from "antd";
import { useState } from "react";
import useContactGroupStore from "../../store/contactgroup.store";


interface CreateGroupProps {
    isOpen: boolean;
    onClose: () => void;
    refetch: () => void; 
}

const CreateGroup: React.FC<CreateGroupProps> = ({ isOpen, onClose,refetch }) => {
    const [form] = Form.useForm();
    const { createGroup, setFormValues } = useContactGroupStore();

    const [isLoading, setIsLoading] = useState<boolean>(false)

    const handleSubmit = async (values: { group_name: string; description: string }) => {
        try {
            setIsLoading(true)
            setFormValues(values);
            await createGroup();
            await new Promise((resolve) => setTimeout(resolve, 3000));
            onClose();
            form.resetFields(); // Reset the form after successful submission
            refetch();
        } catch (error) {
            console.error("Error creating group:", error);
        } finally {
            setIsLoading(false)
        }
    };

    return (
        <Modal
            title="Create Group"
            open={isOpen}
            onCancel={onClose}
            footer={null}
            destroyOnClose
        >
            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
            >
                <Form.Item
                    label="Group Name"
                    name="group_name"
                    rules={[
                        { required: true, message: "Group name is required" },
                    ]}
                >
                    <Input placeholder="Enter group name" />
                </Form.Item>

                <Form.Item
                    label="Description"
                    name="description"
                    rules={[
                        { required: true, message: "Description is required" },
                    ]}
                >
                    <Input.TextArea placeholder="Enter group description" rows={4} />
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

export default CreateGroup;
