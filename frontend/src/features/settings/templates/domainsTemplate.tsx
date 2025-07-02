import React, { useEffect, useState } from "react";
import { Modal, Form, Input, Button, Tabs } from "antd";
// import SendersDashComponent from "../components/senders/sendersDashComponent";
// import DomainDashboardComponent from "../components/domain/domainDashComponent";
import useMetadata from "../../../hooks/useMetaData";
import { Helmet, HelmetProvider } from "react-helmet-async";
import useSenderStore from "../store/sender.store";
import useDomainStore from "../store/domain.store";
import DomainDashboardComponent from "../components/domain/domainDashComponent";
import SendersDashComponent from "../components/senders/sendersDashComponent";

type TabType = "Domain" | "Sender";
interface ModalContent {
    title: string;
    content: string;
}

const DomainTemplateDash: React.FC = () => {
    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "Sender" || storedTab === "Domain") ? storedTab : "Sender";
    });
    const [keyType, setKeyType] = useState<TabType | null>(null);
    const [modalContent, setModalContent] = useState<ModalContent>({ title: "", content: "" });
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [resultModalOpen, setResultModalOpen] = useState<boolean>(false);
    const [form] = Form.useForm();

    const { createDomain, setDomainFormValues } = useDomainStore();
    const { createSender, setSenderFormValues } = useSenderStore();

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);

    const openModal = (title: string, content: string): void => {
        setModalContent({ title, content });
        setIsModalOpen(true);
    };

    const handleCreateDomain = (): void => {
        setKeyType("Domain");
        form.resetFields();
        openModal("Add Domain", "Please provide a name for your new domain.");
    };

    const handleCreateSender = (): void => {
        setKeyType("Sender");
        form.resetFields();
        openModal("Add Sender", "Please provide a name for your new sender.");
    };

    const handleSubmit = async (values: any): Promise<void> => {
        setIsModalOpen(false);

        try {
            if (keyType === "Domain") {
                setDomainFormValues({ ...values })
                await createDomain();
                setModalContent({
                    title: "New Domain Added",
                    content: "Your new domain has been added successfully.",
                });
                location.reload()
            } else if (keyType === "Sender") {
                setSenderFormValues({ ...values })
                await createSender();
                setModalContent({
                    title: "New Sender Added",
                    content: "Your new sender has been added successfully. An email has been sent to you to verify your sender",
                });
            }
            location.reload()
            setResultModalOpen(true);
        } catch (error) {
            console.log(error);
            setModalContent({
                title: "Error",
                content: `Failed to add ${keyType}. Please try again.`,
            });
            setResultModalOpen(true);
        }

        form.resetFields();
    };

    const metaData = useMetadata("Settings");

    const items = [
        {
            key: 'Sender',
            label: 'Senders',
            children: <SendersDashComponent />
        },
        {
            key: 'Domain',
            label: 'Domains',
            children: <DomainDashboardComponent />
        }
    ];

    return (
        <HelmetProvider>
            <Helmet
                {...metaData}
                title={`${activeTab} - CrabMailer`}
            />
            <div className="p-6 max-w-7xl">
                <div className="flex justify-between items-center mb-6">
                    <div>
                        <Button
                            type="primary"
                            onClick={activeTab === "Sender" ? handleCreateSender : handleCreateDomain}
                            className="bg-gray-900 hover:bg-gray-700"
                        >
                            Add {activeTab}
                        </Button>
                    </div>
                </div>

                <Tabs
                    activeKey={activeTab}
                    onChange={(key: string) => setActiveTab(key as TabType)}
                    items={items}
                    className="mb-6"
                />

                <Modal
                    title={modalContent.title}
                    open={isModalOpen}
                    onCancel={() => setIsModalOpen(false)}
                    footer={null}
                >
                    <Form
                        form={form}
                        onFinish={handleSubmit}
                        layout="vertical"
                    >
                        <p className="mb-4 text-gray-600">{modalContent.content}</p>

                        <Form.Item
                            label={`${keyType} Name`}
                            name={keyType === "Domain" ? "domain" : "name"}
                            rules={[{ required: true, message: `Please input your ${keyType} name!` }]}
                        >
                            <Input />
                        </Form.Item>

                        {keyType === "Sender" && (
                            <Form.Item
                                label="Email"
                                name="email"
                                rules={[
                                    { required: true, message: 'Please input your email!' },
                                    { type: 'email', message: 'Please enter a valid email!' }
                                ]}
                            >
                                <Input />
                            </Form.Item>
                        )}

                        <Form.Item className="flex justify-end">
                            <Button onClick={() => setIsModalOpen(false)} className="mr-2">
                                Cancel
                            </Button>
                            <Button type="primary" htmlType="submit" className="bg-blue-500">
                                Add {keyType}
                            </Button>
                        </Form.Item>
                    </Form>
                </Modal>

                <Modal
                    title={modalContent.title}
                    open={resultModalOpen}
                    onCancel={() => setResultModalOpen(false)}
                    footer={[
                        <Button
                            key="close"
                            type="primary"
                            onClick={() => setResultModalOpen(false)}
                            className="bg-blue-500"
                        >
                            Close
                        </Button>
                    ]}
                >
                    <p>{modalContent.content}</p>
                </Modal>
            </div>
        </HelmetProvider>
    );
};

export default DomainTemplateDash;