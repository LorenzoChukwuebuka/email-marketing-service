import { Button, Form, Typography, notification, Input } from "antd";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import useAuthStore from "../store/auth.store";
import renderApiName from "../../../utils/render-name";
import { VerifyLoginFormData } from '../interface/auth.interface';

const { Title, Text } = Typography;

const VerifyLoginComponent: React.FC = () => {
    const navigate = useNavigate();
    const { setVerifyLoginValues, verifyLogin, isLoginVerified } = useAuthStore();
    const [isLoading, setIsLoading] = useState(false);
    const [form] = Form.useForm<VerifyLoginFormData>();

    const handleVerify = async (values: any) => {
        setIsLoading(true);
        try {
            setVerifyLoginValues(values);
            await verifyLogin();
        } catch (error) {
            notification.error({
                message: "Verification Failed",
                description: "Invalid or expired OTP. Please try again.",
            });
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        if (isLoginVerified) {
            notification.success({
                message: "Login Verified",
                description: "You have successfully verified your login. Redirecting to dashboard...",
            });

            const timer = setTimeout(() => {
                navigate("/dashboard");
            }, 1500);

            return () => clearTimeout(timer);
        }
    }, [isLoginVerified, navigate]);

    return (
        <div className="container mx-auto mt-[10em] px-4">
            <div className="max-w-lg mx-auto mt-5">
                <Title level={3} className="text-center">
                    <a href="/">{renderApiName()}</a>
                </Title>
                <div className="bg-white shadow-md rounded-lg p-8">
                    <Title level={4} className="text-center">
                        Verify Login
                    </Title>
                    <Text className="text-center block mb-6 text-gray-600">
                        Enter the 8-digit verification code sent to your email
                    </Text>

                    <Form
                        form={form}
                        layout="vertical"
                        onFinish={handleVerify}
                        className="ant-form"
                    >
                        <Form.Item
                            label="Verification Code"
                            name="token"
                            rules={[
                                {
                                    required: true,
                                    message: "Please enter the verification code",
                                },
                                {
                                    len: 8,
                                    message: "Verification code must be 8 digits",
                                },
                                // {
                                //     pattern: /^\d+$/,
                                //     message: "Verification code must contain only numbers",
                                // },
                            ]}
                        >
                            <Input
                                placeholder="Enter 8-digit code"
                                maxLength={8}
                                size="large"
                                className="text-center text-lg tracking-widest"
                                autoComplete="one-time-code"
                            // inputMode="numeric"
                            // pattern="[0-9]*"
                            />
                        </Form.Item>

                        <Form.Item>
                            <Button
                                type="primary"
                                block
                                htmlType="submit"
                                loading={isLoading}
                                size="large"
                            >
                                Verify Login
                            </Button>
                        </Form.Item>
                    </Form>
                </div>
            </div>
        </div>
    );
};

export default VerifyLoginComponent;