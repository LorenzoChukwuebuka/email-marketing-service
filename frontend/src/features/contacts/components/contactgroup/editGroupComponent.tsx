import { useEffect, useState } from "react";
import { Modal, Form, Input, Button } from "antd";
import useContactGroupStore from '../../store/contactgroup.store';
import { EditGroupValues } from "../../interface/contactgroup.interface";
interface EditContactProps {
    isOpen: boolean;
    onClose: () => void;
    group: EditGroupValues
    refetch: () => void; // Optional refetch function to refresh contacts after creation

}

const EditGroupComponent: React.FC<EditContactProps> = ({ isOpen, onClose, group, refetch }) => {
    const [form] = Form.useForm();
    const { setEditValues, editValues, updateGroup } = useContactGroupStore();
    const [isLoading, setIsLoading] = useState<boolean>(false)

    useEffect(() => {
        console.log(group)
        if (group) {
            // Initialize form values when the group data is provided
            form.setFieldsValue({
                group_name: group.group_name,
                description: group.description,
            });
            setEditValues({
                group_id: group.group_id,
                group_name: group.group_name,
                description: group.description,
            });
        }
    }, [group, form, setEditValues]);

    const handleSubmit = async (values: { group_name: string; description: string }) => {
        try {
            setIsLoading(true)
            console.log(values)
            setEditValues({
                ...editValues,
                group_name: values.group_name,
                group_id: group.group_id,
                description: values.description,
            });
            await updateGroup();
            new Promise(resolve => setTimeout(resolve, 2000))
            onClose();
            form.resetFields(); // Reset form fields after successful submission
            refetch()
        } catch (error) {
            console.error("Error updating group:", error);
        } finally {
            setIsLoading(false)
        }
    };

    return (
        <Modal
            title="Edit Group"
            open={isOpen}
            onCancel={onClose}
            footer={null}
            destroyOnClose
        >
            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
                initialValues={{
                    group_name: group?.group_name || "",
                    description: group?.description || "",
                }}
            >
                <Form.Item
                    label="Group Name"
                    name="group_name"
                    rules={[{ required: true, message: "Group name is required" }]}
                >
                    <Input placeholder="Enter group name" />
                </Form.Item>

                <Form.Item
                    label="Description"
                    name="description"
                    rules={[{ required: true, message: "Description is required" }]}
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

export default EditGroupComponent;
