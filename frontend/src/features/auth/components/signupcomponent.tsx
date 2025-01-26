import { Link } from "react-router-dom";
import { Form, Input, Button, message } from "antd";
import { useState } from "react";
import renderApiName from "../../../utils/render-name";
import useAuthStore from "../store/auth.store";

const SignUpComponent: React.FC = () => {
    const { formValues, setFormValues, registerUser } = useAuthStore();
    const [isLoading, setIsLoading] = useState(false);
    const [form] = Form.useForm();

    const handleSubmit = async () => {
        try {
            setIsLoading(true);
            await form.validateFields();
            await registerUser();
        } catch (error) {
            console.log(error)
        } finally {
            setIsLoading(false)
            message.success("Account created successfully. Please check your email to verify your account.");
            form.resetFields();
        }
    };

    const handleGoogleSignUp = () => {
        window.location.href = `${import.meta.env.VITE_API_URL}/google/signup`;
    }

    return (
        <main className="min-h-screen bg-gradient-to-b from-gray-100 to-gray-200">




            <div className="py-4">
                <h1 className="text-center text-3xl font-bold">{renderApiName()}</h1>
            </div>

            <div className="container mx-auto px-4 py-4">
                <div className="bg-gray-50 shadow-md rounded-lg p-8 max-w-md mx-auto">
                    <h2 className="text-2xl font-semibold text-center text-gray-700 mb-2">
                        Get Started with {import.meta.env.VITE_API_NAME}
                    </h2>

                    <Form
                        form={form}
                        layout="vertical"
                        initialValues={formValues}
                        onFinish={handleSubmit}
                        onValuesChange={(changedValues) =>
                            setFormValues({ ...formValues, ...changedValues })
                        }
                    >
                        <Form.Item
                            label="Full Name"
                            name="fullname"
                            rules={[
                                { required: true, message: "Name is required" },
                                { min: 5, message: "Name must be at least 5 characters" },
                            ]}
                        >
                            <Input placeholder="Full Name" />
                        </Form.Item>

                        <Form.Item
                            label="Email"
                            name="email"
                            rules={[
                                { required: true, message: "Email is required" },
                                { type: "email", message: "Invalid email format" },
                            ]}
                        >
                            <Input type="email" placeholder="Email" />
                        </Form.Item>

                        <Form.Item
                            label="Company"
                            name="company"
                            rules={[{ required: true, message: "Company is required" }]}
                        >
                            <Input placeholder="Company" />
                        </Form.Item>

                        <Form.Item
                            label="Password"
                            name="password"
                            rules={[
                                { required: true, message: "Password is required" },
                                { min: 8, message: "Password must be at least 8 characters" },
                                {
                                    pattern: /[a-zA-Z]/,
                                    message: "Password must contain at least one letter",
                                },
                                {
                                    pattern: /[0-9]/,
                                    message: "Password must contain at least one number",
                                },
                            ]}
                        >
                            <Input.Password placeholder="Password" />
                        </Form.Item>

                        <Form.Item
                            label="Confirm Password"
                            name="confirmPassword"
                            dependencies={['password']}
                            rules={[
                                { required: true, message: "Confirm Password is required" },
                                ({ getFieldValue }) => ({
                                    validator(_, value) {
                                        if (!value || getFieldValue("password") === value) {
                                            return Promise.resolve();
                                        }
                                        return Promise.reject(
                                            new Error("Passwords must match")
                                        );
                                    },
                                }),
                            ]}
                        >
                            <Input.Password placeholder="Confirm Password" />
                        </Form.Item>

                        <Form.Item>
                            <Button
                                type="primary"
                                htmlType="submit"
                                loading={isLoading}
                                className="w-full"
                            >
                                {isLoading ? "Creating account..." : "Create Account"}
                            </Button>
                        </Form.Item>
                    </Form>

                    
            <div className="relative mb-4">
                <div className="absolute inset-0 flex items-center">
                    <div className="w-full border-t border-gray-300"></div>
                </div>
                <div className="relative flex justify-center text-sm">
                    <span className="px-2 bg-white text-gray-500">Or continue with</span>
                </div>
            </div>

            <Button
                onClick={() => handleGoogleSignUp()}
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
                Sign up with Google
            </Button>

                    <div className="mt-2 text-center text-sm text-gray-600">
                        By signing up, you agree to our{" "}
                        <a href="/tos" className="text-indigo-600 hover:underline">
                            Terms of Service
                        </a>{" "}
                        and{" "}
                        <a href="/privacy" className="text-indigo-600 hover:underline">
                            Privacy Policy
                        </a>.
                    </div>
                    <div className="mt-2 text-center">
                        <Link
                            to="/auth/login"
                            className="text-indigo-600 hover:text-indigo-800 transition duration-300 ease-in-out"
                        >
                            Already have an account? Log in
                        </Link>
                    </div>
                </div>
            </div>
        </main>
    );
};

export default SignUpComponent;
