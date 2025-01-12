import React from "react";
import { Form, Input, Button, Typography, notification } from "antd";
import { Link } from "react-router-dom";
import useAuthStore from "../store/auth.store";
import renderApiName from "../../../utils/render-name";

const { Title, Text } = Typography;

const ForgotPasswordComponent: React.FC = () => {
    const {
        forgetPasswordValues,
        setForgetPasswordValues,
        forgotPass,
    } = useAuthStore();

    const [form] = Form.useForm();

    const [isLoading, setIsLoading] = React.useState(false);

    const handleSubmit = async (values: { email: string }) => {
        try {
            setIsLoading(true)
            setForgetPasswordValues(values);
            await forgotPass();
            notification.success({
                message: "Request Sent",
                description: "If your email is registered, you will receive a reset link shortly.",
            });
            form.resetFields();
        } catch (error) {
            console.error("Error sending request:", error);
        } finally {
            setIsLoading(false)
        }
    };

    return (
        <div className="container mx-auto px-4">
            <div className="max-w-lg mx-auto mt-[10em]">
                <Title level={3} className="text-center">
                    {renderApiName()}
                </Title>
                <div className="bg-white shadow-md rounded-lg p-8">
                    <Title level={4} className="text-center mb-4">
                        Forgot Password
                    </Title>
                    <Text type="secondary" className="text-center block mb-4">
                        You will receive an email if your mail is registered with us
                    </Text>

                    <Form
                        form={form}
                        layout="vertical"
                        onFinish={handleSubmit}
                        initialValues={forgetPasswordValues}
                    >
                        <Form.Item
                            label="Email"
                            name="email"
                            rules={[
                                { required: true, message: "Email is required" },
                                { type: "email", message: "Invalid email format" },
                            ]}
                        >
                            <Input
                                placeholder="Enter your registered email"
                                onChange={(event) =>
                                    setForgetPasswordValues({
                                        ...forgetPasswordValues,
                                        email: event.target.value,
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
                                Submit
                            </Button>
                        </Form.Item>
                    </Form>
                    <div className="text-center mt-4">
                        <Text>
                            Remember your password?{" "}
                            <Link to="/auth/login">Login</Link>
                        </Text>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ForgotPasswordComponent;
