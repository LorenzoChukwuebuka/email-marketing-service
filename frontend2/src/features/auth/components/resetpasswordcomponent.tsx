import React, { useState } from "react";
import { Form, Input, Button, Typography, notification } from "antd";
import useAuthStore from "../store/auth.store";
import renderApiName from "../../../utils/render-name";

const { Title } = Typography;

interface ResetPasswordValues {
    password: string;
    confirmPassword: string;
}

const ResetPasswordComponent: React.FC = () => {
    const { resetPasswordValues, setResetPasswordValues, resetPassword } = useAuthStore();
    const [isLoading, setIsLoading] = useState(false);


    const handleValidation = async (values: ResetPasswordValues) => {
        try {

            setIsLoading(true)
            const emailFromURL = new URLSearchParams(window.location.search).get("email");
            const tokenFromURL = new URLSearchParams(window.location.search).get("token");

            setResetPasswordValues({
                ...resetPasswordValues,
                ...values,
                token: tokenFromURL || "",
                email: emailFromURL || "",
            });

            console.log(emailFromURL, tokenFromURL);

            await resetPassword();
            notification.success({
                message: "Success",
                description: "Password reset successfully!",
                placement: "bottomRight",
            });
        } catch (error) {
            console.log(error);
        } finally {
            setIsLoading(false);
        }
    };

    const [form] = Form.useForm();

    const handleSubmit = (values: ResetPasswordValues) => {
        handleValidation(values);
    };

    return (
        <div className="container mx-auto px-4">
            <div className="max-w-lg mx-auto mt-10">
                <Title level={3} className="text-center mt-10 mb-4">
                    {renderApiName()}
                </Title>
                <div className="bg-white shadow-md rounded-lg p-8">

                    <Title level={4} className="text-center mb-4">
                        Reset Password
                    </Title>

                    <Form
                        form={form}
                        layout="vertical"
                        onFinish={handleSubmit}
                        initialValues={resetPasswordValues}
                    >
                        <Form.Item
                            label="Password"
                            name="password"
                            rules={[
                                { required: true, message: "Password is required" },
                                { min: 8, message: "Password must be at least 8 characters" },
                            ]}
                        >
                            <Input.Password
                                placeholder="Enter your new password"
                                onChange={(e) =>
                                    setResetPasswordValues({
                                        ...resetPasswordValues,
                                        password: e.target.value,
                                    })
                                }
                            />
                        </Form.Item>

                        <Form.Item
                            label="Confirm Password"
                            name="confirmPassword"
                            dependencies={["password"]}
                            rules={[
                                { required: true, message: "Confirm Password is required" },
                                ({ getFieldValue }) => ({
                                    validator(_, value) {
                                        if (!value || getFieldValue("password") === value) {
                                            return Promise.resolve();
                                        }
                                        return Promise.reject(
                                            new Error("Passwords do not match")
                                        );
                                    },
                                }),
                            ]}
                        >
                            <Input.Password
                                placeholder="Confirm your password"
                                onChange={(e) =>
                                    setResetPasswordValues({
                                        ...resetPasswordValues,
                                        confirmPassword: e.target.value,
                                    })
                                }
                            />
                        </Form.Item>

                        <Form.Item>
                            <Button
                                type="primary"
                                htmlType="submit"
                                block
                                loading={isLoading}
                            >
                                {isLoading ? "Please wait..." : "Submit"}
                            </Button>
                        </Form.Item>
                    </Form>
                </div>
            </div>
        </div>
    );
};

export default ResetPasswordComponent;
