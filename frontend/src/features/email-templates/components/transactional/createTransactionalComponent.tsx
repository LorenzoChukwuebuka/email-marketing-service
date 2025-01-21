import { Modal, Form, Input, Button } from 'antd';
import useTemplateStore from '../../store/template.store';
import { Template } from '../../interface/email-templates.interface';

interface Props {
    isOpen: boolean;
    onClose: () => void;
    editorType: "drag-and-drop" | "html-editor" | "rich-text";
}

const CreateTransactionalTemplateComponent: React.FC<Props> = ({ isOpen, onClose, editorType }) => {
    const [form] = Form.useForm();
    const { setFormValues, createTemplate } = useTemplateStore();

    const handleSubmit = async (values: Omit<Template, 'user_id'>) => {
        try {
            const formData = {
                ...values,
                type: "transactional",
                editor_type: editorType,
                is_editable: true,
                is_gallery_template: false,
                is_published: false,
            };

            setFormValues(formData);
            await createTemplate();
            form.resetFields();
            onClose();
        } catch (error) {
            console.error('Failed to create template:', error);
        }
    };

    return (
        <Modal
            title="Create Template"
            open={isOpen}
            onCancel={onClose}
            footer={null}
        >
            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
            >
                <Form.Item
                    label="Template Name"
                    name="template_name"
                    rules={[
                        { required: true, message: 'Please enter template name' }
                    ]}
                >
                    <Input placeholder="template ..." />
                </Form.Item>

                <Form.Item
                    label="Tags"
                    name="tags"
                    rules={[
                        { required: true, message: 'Please enter tags' },
                        {
                            validator: (_, value) => {
                                if (!value) return Promise.resolve();
                                const isValid = value.split(',').every(tag =>
                                    tag.trim() !== "" && !tag.includes(' ')
                                );
                                return isValid
                                    ? Promise.resolve()
                                    : Promise.reject('Tags must be non-empty, separated by commas, and contain no spaces');
                            }
                        }
                    ]}
                >
                    <Input placeholder="E.g. signup,welcome" />
                </Form.Item>

                <Form.Item className="flex justify-end mb-0">
                    <Button onClick={onClose} className="mr-2">
                        Cancel
                    </Button>
                    <Button type="primary" htmlType="submit">
                        Submit
                    </Button>
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default CreateTransactionalTemplateComponent;