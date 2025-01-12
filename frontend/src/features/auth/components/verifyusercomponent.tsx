import { Button, Form, Typography, notification } from "antd";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import useAuthStore from "../store/auth.store";
import renderApiName from "../../../utils/render-name";

const { Title, Text } = Typography;

const VerifyAccountComponent: React.FC = () => {
    const navigate = useNavigate();
    const { setOTPValue, verifyUser, isVerified, resendOTP } = useAuthStore();
    const [isLoading, setIsLoading] = useState(false);

    const handleVerify = async () => {
        setIsLoading(true);
        const tokenFromURL = new URLSearchParams(window.location.search).get("token");
        if (tokenFromURL) {
            setOTPValue({ token: tokenFromURL });
            await verifyUser();
        }
        setIsLoading(false);
    };

    const handleResendOTP = async () => {
        notification.info({
            message: "OTP Resent",
            description: "You have successfully resent the token. Please check your email.",
        });

        const emailFromURL = new URLSearchParams(window.location.search).get("email");
        const usernameFromURL = new URLSearchParams(window.location.search).get("username");
        const userIdFromURL = new URLSearchParams(window.location.search).get("userId");

        const data = {
            user_id: userIdFromURL || "",
            username: usernameFromURL || "",
            email: emailFromURL || "",
            otp_type: "emailVerify",
        };

        await resendOTP(data);
    };

    useEffect(() => {
        if (isVerified) {
            const timer = setTimeout(() => {
                navigate("/auth/login");
            }, 1500);

            return () => clearTimeout(timer);
        }
    }, [isVerified, navigate]);

    return (
        <div className="container mx-auto mt-[10em] px-4">
            <div className="max-w-lg mx-auto mt-5">
                <Title level={3} className="text-center">
                    <a href="/">{renderApiName()}</a>
                </Title>
                <div className="bg-white shadow-md rounded-lg p-8">
                    <Title level={4} className="text-center">
                        Verify Email
                    </Title>
                    <Form
                        layout="vertical"
                        onFinish={handleVerify}
                        className="ant-form"
                    >
                        <Form.Item>
                            <Button
                                type="primary"
                                block
                                htmlType="submit"
                                loading={isLoading}
                            >
                                Verify Email
                            </Button>
                        </Form.Item>
                    </Form>
                    <div className="text-center mt-4">
                        <Text>
                            Didn’t receive the OTP?{" "}
                            <Button
                                type="link"
                                onClick={handleResendOTP}
                                className="p-0"
                            >
                                Resend OTP
                            </Button>
                        </Text>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default VerifyAccountComponent;
