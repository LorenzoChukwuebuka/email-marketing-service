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
        window.location.href = `${import.meta.env.VITE_API_URL}/google/login`;

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
                    </Form>

                    <button onClick={() => handleGoogleLogin()}>
                        Sign in with Google
                    </button>


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
