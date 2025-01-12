import React, { useState } from "react";
import { Form, Input, Button, Space } from "antd";
import useAuthStore from "../../../auth/store/auth.store";

const ChangePasswordComponent: React.FC = () => {
    const [form] = Form.useForm();
    const {
        changePassword,
        setChangePasswordValues,
    } = useAuthStore();

    const [isLoading, setIsLoading] = useState<boolean>(false)

    const handleSubmit = async (values: {
        old_password: string;
        new_password: string;
        confirm_password: string;
    }) => {
        try {
            setIsLoading(true)
            setChangePasswordValues(values);
            await changePassword();
            form.resetFields();
        } catch (error) {
            console.error("Failed to change password:", error);
        } finally {
            setIsLoading(false)
        }
    };

    return (
        <div className="mt-8 p-4 mb-5">
            <Form
                form={form}
                onFinish={handleSubmit}
                layout="vertical"
                className="w-full max-w-xs"
            >
                <Space direction="vertical" className="w-full">
                    <Form.Item
                        name="old_password"
                        label="Current Password"
                        rules={[
                            {
                                required: true,
                                message: "Current password is required",
                            }
                        ]}
                    >
                        <Input.Password
                            className="h-10"
                        />
                    </Form.Item>

                    <Form.Item
                        name="new_password"
                        label="New Password"
                        rules={[
                            {
                                required: true,
                                message: "New password is required",
                            }
                        ]}
                    >
                        <Input.Password
                            className="h-10"
                        />
                    </Form.Item>

                    <Form.Item
                        name="confirm_password"
                        label="Confirm New Password"
                        dependencies={["new_password"]}
                        rules={[
                            {
                                required: true,
                                message: "Please confirm your password",
                            },
                            ({ getFieldValue }) => ({
                                validator(_, value) {
                                    if (!value || getFieldValue("new_password") === value) {
                                        return Promise.resolve();
                                    }
                                    return Promise.reject(new Error("Passwords must match"));
                                },
                            }),
                        ]}
                    >
                        <Input.Password
                            className="h-10"
                        />
                    </Form.Item>

                    <Form.Item>
                        <Button
                            type="default"
                            htmlType="submit"
                            loading={isLoading}
                            className="mt-6 bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold"
                        >
                            {isLoading ? "Please wait" : "Change Password"}
                        </Button>
                    </Form.Item>
                </Space>
            </Form>
        </div>
    );
};

export default ChangePasswordComponent;