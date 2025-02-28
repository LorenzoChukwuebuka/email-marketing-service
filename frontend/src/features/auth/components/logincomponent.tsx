import { Link } from "react-router-dom";
import { useState, ChangeEvent } from "react";
import { Form, Input, Button, Typography } from "antd";
import useAuthStore from "../store/auth.store";
import renderApiName from "../../../utils/render-name";


const { Title, Text } = Typography;

const LoginComponent: React.FC = () => {
    const {
        loginValues,
        setLoginValues,
        loginUser,

    } = useAuthStore();

    const [form] = Form.useForm();
    const [isLoading, setIsLoading] = useState(false);



    const handleLogin = async (values: { email: string; password: string }) => {
        try {
            setIsLoading(true);
            setLoginValues(values);
            await loginUser();
        } catch (error) {
            console.error("Error logging in:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleGoogleLogin = () => {
        // Redirect to your backend's Google login endpoint
        window.location.href = `${import.meta.env.VITE_API_URL}/google/signup`;

    };

    return (
        <div className="flex justify-center items-center h-screen bg-gray-100">
            <div className="container mx-auto">
                <Title level={3} className="text-center">
                    <a href="/">{renderApiName()}</a>
                </Title>
                <div className="bg-white shadow-lg rounded-lg max-w-md mx-auto p-6">
                    <Title level={3} className="text-center mb-4">
                        Log in
                    </Title>
                    <Form
                        form={form}
                        layout="vertical"
                        onFinish={handleLogin}
                        initialValues={loginValues}
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
                                placeholder="Enter your email"
                                onChange={(event: ChangeEvent<HTMLInputElement>) =>
                                    setLoginValues({
                                        ...loginValues,
                                        email: event.target.value,
                                    })
                                }
                            />
                        </Form.Item>

                        <Form.Item
                            label="Password"
                            name="password"
                            rules={[{ required: true, message: "Password is required" }]}
                        >
                            <Input.Password
                                placeholder="Enter your password"
                                onChange={(event: ChangeEvent<HTMLInputElement>) =>
                                    setLoginValues({
                                        ...loginValues,
                                        password: event.target.value,
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
                                {isLoading ? "Please Wait" : "Login"}
                            </Button>
                        </Form.Item>


                        <div className="relative mb-4">
                            <div className="absolute inset-0 flex items-center">
                                <div className="w-full border-t border-gray-300"></div>
                            </div>
                            <div className="relative flex justify-center text-sm">
                                <span className="px-2 bg-white text-gray-500">Or continue with</span>
                            </div>
                        </div>

                        <Button
                            onClick={() => handleGoogleLogin()}
                            className="w-full flex items-center justify-center gap-2 border border-gray-300 rounded-md px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 mb-4"
                        >
                            <svg className="w-5 h-5" viewBox="0 0 24 24">
                                <path
                                    d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                                    fill="#4285F4"
                                />
                                <path
                                    d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                                    fill="#34A853"
                                />
                                <path
                                    d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                                    fill="#FBBC05"
                                />
                                <path
                                    d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                                    fill="#EA4335"
                                />
                            </svg>
                            Sign in with Google
                        </Button>

                    </Form>

                    <div className="text-center mt-4">
                        <Text>
                            <Link to="/auth/forgot-password">Forgot Password</Link>
                            <Link to="/auth/sign-up" className="ml-4">
                                Create Account
                            </Link>
                        </Text>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default LoginComponent;
