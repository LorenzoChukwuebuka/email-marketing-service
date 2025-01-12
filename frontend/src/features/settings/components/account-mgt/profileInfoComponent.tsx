import React, { useEffect } from "react";
import { Form, Input, Button, Typography, Space } from "antd";
import useAuthStore from "../../../auth/store/auth.store";
import { useUserDetailsQuery } from "../../../auth/hooks/useUserDetailsQuery";
const { Title, Text } = Typography;

const ProfileInformationComponent: React.FC = () => {
    const [form] = Form.useForm();
    const {

        setEditFormValues,
        editFormValues,
        editUserDetails,
    } = useAuthStore();

    const { data: userData } = useUserDetailsQuery()

    const initEdit = () => {
        setEditFormValues({
            fullname: userData?.payload.fullname || "",
            email: userData?.payload.email || "",
            company: userData?.payload.company || "",
            phonenumber: userData?.payload.phonenumber || "",
        });
        form.setFieldsValue({
            fullname: userData?.payload.fullname || "",
            email: userData?.payload.email || "",
            company: userData?.payload.company || "",
            phonenumber: userData?.payload.phonenumber || "",
        });
    };

    const handleEditInformation = async (values: any) => {
        try {
            setEditFormValues(values);
            await editUserDetails();
            initEdit();
        } catch (error) {
            console.error("Failed to update user details:", error);
        }
    };


    useEffect(() => {
        if (userData) {
            initEdit();
        }
    }, [userData, form]);

    return (
        <div className="max-w-3xl ml-5 p-6">
            <Title level={2} className="mb-4">
                Profile Information
            </Title>
            <Text className="block mb-6" type="secondary">
                This is the information we have associated with your Crabmailer profile,
                which you can use to access multiple Crabmailer accounts.
                <Text strong className="block">
                    All contact information is kept strictly confidential.
                </Text>
            </Text>

            <Form
                form={form}
                layout="vertical"
                onFinish={handleEditInformation}
                initialValues={editFormValues}
            >
                <Space direction="vertical" className="w-full">
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <Form.Item
                            label="Username"
                            name="email"
                        >
                            <Input
                                readOnly
                                className="bg-gray-100"
                            />
                        </Form.Item>

                        <div className="grid grid-cols-2 gap-4">
                            <Form.Item
                                label="FullName"
                                name="fullname"
                            >
                                <Input />
                            </Form.Item>

                            <Form.Item
                                label="Company"
                                name="company"
                            >
                                <Input />
                            </Form.Item>
                        </div>
                    </div>

                    <Form.Item
                        label="Email"
                        name="email"
                    >
                        <Input
                            readOnly
                            className="bg-gray-100 w-[20em]"
                        />
                    </Form.Item>

                    <Form.Item
                        label="Phone Number"
                        name="phonenumber"
                        rules={[
                            { required: true, message: "Phone number is required" },
                            { len: 11, message: "Phone number must be exactly 11 digits" },
                            { pattern: /^\d+$/, message: "Phone number must contain only digits" }
                        ]}
                    >
                        <Input
                            maxLength={11}
                            className="w-[20em]"
                        />
                    </Form.Item>

                    <Form.Item>
                        <Button
                            type="default"
                            htmlType="submit"
                            className="bg-gray-200 hover:bg-gray-300"
                        >
                            Edit Information
                        </Button>
                    </Form.Item>
                </Space>
            </Form>
        </div>
    );
};

export default ProfileInformationComponent;