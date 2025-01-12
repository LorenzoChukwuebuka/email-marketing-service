import { Helmet, HelmetProvider } from "react-helmet-async";
import { useState } from "react";
import { Card, Button, Typography, Form } from "antd";
import renderApiName from "../../../../utils/render-name";
import useSenderStore from "../../store/sender.store";

const { Title } = Typography;

const VerifySenderComponent: React.FC = () => {
   
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const { setVerifySender, verifySender } = useSenderStore();

    const handleVerify = async () => {
        setIsLoading(true);
        try {
            // Extract token from URL
            const searchParams = new URLSearchParams(window.location.search);
            const token = searchParams.get("token");
            const userId = searchParams.get("userId");
            const email = searchParams.get("email");

            setVerifySender({
                email: email as string,
                user_id: userId as string,
                token: token as string
            });

            await verifySender();
        } catch (error) {
            console.error("Verification failed:", error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <HelmetProvider>
            <Helmet title="Verify Sender" />

            <div className="container mx-auto mt-40 px-4">
                <div className="max-w-lg mx-auto">
                    <Title level={3} className="text-center mb-4">
                        <a href="/">{renderApiName()}</a>
                    </Title>

                    <Card className="shadow-md">
                        <Title level={3} className="text-center mb-6">
                            Verify Sender Email
                        </Title>

                        <Form onFinish={handleVerify}>
                            <Form.Item className="text-center mb-0">
                                <Button
                                    type="primary"
                                    htmlType="submit"
                                    className="w-full bg-gray-800 hover:bg-gray-700"
                                    loading={isLoading}
                                    size="large"
                                >
                                    {isLoading ? "Please wait" : "Verify Sender"}
                                </Button>
                            </Form.Item>
                        </Form>
                    </Card>
                </div>
            </div>
        </HelmetProvider>
    );
};

export default VerifySenderComponent;