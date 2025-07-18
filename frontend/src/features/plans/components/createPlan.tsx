import { Modal, Form, Input, Select, InputNumber, Button, Card, Typography, Divider } from "antd";
import { PlusOutlined, DeleteOutlined } from "@ant-design/icons";
import React, { useState } from "react";
import usePlanStore from "../store/plan.store";

const { Title, Text } = Typography;
const { Option } = Select;
const { TextArea } = Input;

interface CreatePlanProps {
    isOpen: boolean;
    onClose: () => void;
}

interface Feature {
    name: string;
    description: string;
    value: string;
}

interface MailingLimits {
    daily_limit: number;
    monthly_limit: number;
    max_recipients_per_mail: number;
}

interface PlanValues {
    plan_name: string;
    description: string;
    price: number;
    billing_cycle: string;
    status: string;
    features: Feature[];
    mailing_limits: MailingLimits;
}

const CreatePlan: React.FC<CreatePlanProps> = ({ isOpen, onClose }) => {
    const { createPlan, setPlanValues } = usePlanStore();
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);

    const initialValues: PlanValues = {
        plan_name: "",
        description: "",
        price: 0,
        billing_cycle: "monthly",
        status: "active",
        features: [],
        mailing_limits: {
            daily_limit: 0,
            monthly_limit: 0,
            max_recipients_per_mail: 0,
        },
    };

    const handleSubmit = async (values: PlanValues) => {
        setLoading(true);
        setPlanValues({ ...values })
        try {
            await createPlan();
            form.resetFields();
            onClose();
        } catch (error) {
            console.error("Error creating plan:", error);
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
            title={
                <div className="flex items-center space-x-2">
                    <Title level={4} className="mb-0">Create New Plan</Title>
                </div>
            }
            open={isOpen}
            onCancel={handleCancel}
            footer={null}
            width={800}
            className="create-plan-modal"
        >
            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
                initialValues={initialValues}
                className="space-y-4"
            >
                {/* Basic Plan Information */}
                <Card className="shadow-sm border-gray-200">
                    <Title level={5} className="text-gray-800 mb-4">Basic Information</Title>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Form.Item
                            label={<Text strong>Plan Name</Text>}
                            name="plan_name"
                            rules={[{ required: true, message: "Please enter plan name" }]}
                        >
                            <Input
                                placeholder="Enter plan name"
                                className="rounded-lg"
                            />
                        </Form.Item>

                        <Form.Item
                            label={<Text strong>Status</Text>}
                            name="status"
                            rules={[{ required: true, message: "Please select status" }]}
                        >
                            <Select className="rounded-lg">
                                <Option value="active">Active</Option>
                                <Option value="inactive">Inactive</Option>
                            </Select>
                        </Form.Item>
                    </div>

                    <Form.Item
                        label={<Text strong>Description</Text>}
                        name="description"
                        rules={[{ required: true, message: "Please enter description" }]}
                    >
                        <TextArea
                            rows={3}
                            placeholder="Enter plan description"
                            className="rounded-lg"
                        />
                    </Form.Item>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Form.Item
                            label={<Text strong>Price</Text>}
                            name="price"
                            rules={[{ required: true, message: "Please enter price" }]}
                        >
                            <InputNumber
                                prefix="$"
                                min={0}
                                step={0.01}
                                placeholder="0.00"
                                className="w-full rounded-lg"
                            />
                        </Form.Item>

                        <Form.Item
                            label={<Text strong>Billing Cycle</Text>}
                            name="billing_cycle"
                            rules={[{ required: true, message: "Please select billing cycle" }]}
                        >
                            <Select className="rounded-lg">
                                <Option value="monthly">Monthly</Option>
                                <Option value="yearly">Yearly</Option>
                                <Option value="quarterly">Quarterly</Option>
                            </Select>
                        </Form.Item>
                    </div>
                </Card>

                {/* Mailing Limits */}
                <Card className="shadow-sm border-gray-200">
                    <Title level={5} className="text-gray-800 mb-4">Mailing Limits</Title>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <Form.Item
                            label={<Text strong>Daily Limit</Text>}
                            name={["mailing_limits", "daily_limit"]}
                            rules={[{ required: true, message: "Please enter daily limit" }]}
                        >
                            <InputNumber
                                min={0}
                                placeholder="100"
                                className="w-full rounded-lg"
                            />
                        </Form.Item>

                        <Form.Item
                            label={<Text strong>Monthly Limit</Text>}
                            name={["mailing_limits", "monthly_limit"]}
                            rules={[{ required: true, message: "Please enter monthly limit" }]}
                        >
                            <InputNumber
                                min={0}
                                placeholder="3000"
                                className="w-full rounded-lg"
                            />
                        </Form.Item>

                        <Form.Item
                            label={<Text strong>Max Recipients per Mail</Text>}
                            name={["mailing_limits", "max_recipients_per_mail"]}
                            rules={[{ required: true, message: "Please enter max recipients" }]}
                        >
                            <InputNumber
                                min={0}
                                placeholder="50"
                                className="w-full rounded-lg"
                            />
                        </Form.Item>
                    </div>
                </Card>

                {/* Features */}
                <Card className="shadow-sm border-gray-200">
                    <Title level={5} className="text-gray-800 mb-4">Features</Title>

                    <Form.List name="features">
                        {(fields, { add, remove }) => (
                            <>
                                {fields.map((field, index) => (
                                    <Card
                                        key={field.key}
                                        className="mb-4 bg-gray-50 border-gray-200"
                                        size="small"
                                    >
                                        <div className="flex justify-between items-start mb-3">
                                            <Text strong className="text-gray-700">
                                                Feature {index + 1}
                                            </Text>
                                            <Button
                                                type="text"
                                                danger
                                                size="small"
                                                icon={<DeleteOutlined />}
                                                onClick={() => remove(field.name)}
                                                className="hover:bg-red-50"
                                            />
                                        </div>

                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                            <Form.Item
                                                label={<Text>Feature Name</Text>}
                                                name={[field.name, "name"]}
                                                rules={[{ required: true, message: "Please enter feature name" }]}
                                            >
                                                <Input
                                                    placeholder="e.g., Emails per day"
                                                    className="rounded-lg"
                                                />
                                            </Form.Item>

                                            <Form.Item
                                                label={<Text>Value</Text>}
                                                name={[field.name, "value"]}
                                                rules={[{ required: true, message: "Please enter feature value" }]}
                                            >
                                                <Input
                                                    placeholder="e.g., 100"
                                                    className="rounded-lg"
                                                />
                                            </Form.Item>
                                        </div>

                                        <Form.Item
                                            label={<Text>Description</Text>}
                                            name={[field.name, "description"]}
                                            rules={[{ required: true, message: "Please enter feature description" }]}
                                        >
                                            <TextArea
                                                rows={2}
                                                placeholder="e.g., Maximum number of emails you can send daily"
                                                className="rounded-lg"
                                            />
                                        </Form.Item>
                                    </Card>
                                ))}

                                <Button
                                    type="dashed"
                                    onClick={() => add()}
                                    icon={<PlusOutlined />}
                                    className="w-full h-10 rounded-lg border-gray-300 hover:border-blue-400 hover:text-blue-400"
                                >
                                    Add Feature
                                </Button>
                            </>
                        )}
                    </Form.List>
                </Card>

                <Divider />

                {/* Action Buttons */}
                <div className="flex justify-end space-x-3">
                    <Button
                        onClick={handleCancel}
                        className="px-6 py-2 h-10 rounded-lg"
                    >
                        Cancel
                    </Button>
                    <Button
                        type="primary"
                        htmlType="submit"
                        loading={loading}
                        className="px-6 py-2 h-10 rounded-lg bg-blue-600 hover:bg-blue-700"
                    >
                        Create Plan
                    </Button>
                </div>
            </Form>
        </Modal>
    );
};

export default CreatePlan;