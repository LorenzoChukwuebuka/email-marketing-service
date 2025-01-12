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
