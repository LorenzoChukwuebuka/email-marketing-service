import { Form, Input, Button } from 'antd';
import { useEffect, useState } from "react";
import { Helmet, HelmetProvider } from "react-helmet-async";
import useAdminAuthStore from '../../store/admin.auth.store';
import renderApiName from '../../../../utils/render-name';
import { AdminLoginValues } from '../../interface/admin.auth.interface';

 

const AdminLoginComponent: React.FC = () => {
    const [form] = Form.useForm();
    const {
        loginAdmin,
        isLoggedIn,
        setLoginValues
    } = useAdminAuthStore();

    const [isLoading, setIsLoading] = useState(false)

    const handleLogin = async (values: AdminLoginValues) => {
        try {
            setIsLoading(true)
            setLoginValues(values);
            await loginAdmin();
        } catch (error) {
            console.log(error);
        } finally {
            setIsLoading(false)
        }
    };

    useEffect(() => {
        if (isLoggedIn) {
            location.href = "/zen/dash";
        }
    }, [isLoggedIn]);

    return (
        <HelmetProvider>
            <Helmet title="Admin Login - CrabMailer" />
            <div className="flex justify-center items-center h-screen bg-gray-100">
                <div className="container mx-auto">
                    <h3 className="text-2xl font-bold text-center mb-4">{renderApiName()}</h3>
                    <div className="bg-white shadow-lg rounded-lg max-w-lg mx-auto mt-2 p-6">
                        <h3 className="text-2xl font-semibold text-center mb-4">Log in</h3>
                        <Form
                            form={form}
                            onFinish={handleLogin}
                            layout="vertical"
                        >
                            <Form.Item
                                label="Email"
                                name="email"
                                required
                                rules={[
                                    { required: true, message: 'Email is required' },
                                    { type: 'email', message: 'Invalid email format' }
                                ]}
                            >
                                <Input className="w-full p-2 border border-gray-300 rounded-md" />
                            </Form.Item>

                            <Form.Item
                                label="Password"
                                name="password"
                                required
                                rules={[
                                    { required: true, message: 'Password is required' }
                                ]}
                            >
                                <Input.Password className="w-full p-2 border border-gray-300 rounded-md" />
                            </Form.Item>

                            <Form.Item className="text-center">
                                <Button
                                    type="primary"
                                    htmlType="submit"
                                    loading={isLoading}
                                    className="bg-black text-white py-2 px-4 rounded-md mt-3 hover:bg-gray-800"
                                >
                                    {isLoading ? 'Please wait...' : 'Login'}
                                </Button>
                            </Form.Item>
                        </Form>
                    </div>
                </div>
            </div>
        </HelmetProvider>
    );
};

export default AdminLoginComponent;