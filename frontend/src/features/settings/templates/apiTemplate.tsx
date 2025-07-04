import React, { useEffect, useState } from "react";
import useMetadata from "../../../hooks/useMetaData";
import { Helmet, HelmetProvider } from "react-helmet-async";
import { Modal, Button, Input } from "antd";
import useAPIKeyStore from "../store/apikey.store";
import useSMTPKeyStore from "../store/smtpkey.store";
import APIKeysTableComponent from "../components/keys/apiComponentTable";
import APIInfo from "../components/keys/apiInfoComponent";
import SMTPKeysTableComponent from '../components/keys/smtpComponentTable';
import { generateRandomName } from "../../../utils/generateRandomName";

interface ModalContent {
    title: string;
    content: string;
}

type Tabtype = "API Keys" | "SMTP";

const APISettingsDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<Tabtype>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return storedTab === "API Keys" || storedTab === "SMTP" ? storedTab : "SMTP";
    });
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isKeyModalOpen, setIsKeyModalOpen] = useState(false);
    const [modalContent, setModalContent] = useState<ModalContent>({ title: "", content: "" });
    const [keyType, setKeyType] = useState<"API" | "SMTP" | null>(null);

    const {
        generateAPIKey,
        setFormValues: setAPIFormValues,
        formValues: apiFormValues,
    } = useAPIKeyStore();

    const {
        createSMTPKey,
        setSmtpFormValues,

        smtpformValues,
    } = useSMTPKeyStore();

    const openModal = (title: string, content: string): void => {
        setModalContent({ title, content });
        setIsModalOpen(true);
    };

    const handleGenerateAPIKey = (): void => {
        setKeyType("API");
        setAPIFormValues({ name: generateRandomName() });
        openModal("Generate New API Key", "Please provide a name or description for your new API key.");
    };

    const handleGenerateSMTPKey = (): void => {
        setKeyType("SMTP");
        setSmtpFormValues({ key_name: generateRandomName() });
        openModal("Generate New SMTP Key", "Please provide a name or description for your new SMTP key.");
    };

    const handleSubmitKeyForm = async (): Promise<void> => {
        setIsModalOpen(false);

        try {
            if (keyType === "SMTP") {
                await createSMTPKey();

                setModalContent({
                    title: "New SMTP Key Generated",
                    content: "Your new SMTP key has been generated successfully. You can view it in the SMTP keys table.",
                });
            } else if (keyType === "API") {
                await generateAPIKey();

                setModalContent({
                    title: "New API Key Generated",
                    content: "Your new API key has been generated successfully. You can view it in the API keys table.",
                });
            }
        } catch (error) {
            console.log(error)
            setModalContent({
                title: "Error",
                content: `Failed to generate ${keyType} key. Please try again.`,
            });
        }

        setIsKeyModalOpen(true);
        if (keyType === "API") {
            setAPIFormValues({ name: "" });
        } else {
            setSmtpFormValues({ key_name: "" });
        }
    };

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const { value } = event.target;
        if (keyType === "API") {
            setAPIFormValues({ name: value });
        } else {
            setSmtpFormValues({ key_name: value });
        }
    };

    const metaData = useMetadata("Settings");

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);

    return (
        <HelmetProvider>
            <Helmet
                {...metaData}
                title={activeTab === "SMTP" ? "SMTP - CrabMailer" : "API Key - CrabMailer"}
            />
            <div className="p-6 max-w-5xl">
                <div className="flex justify-between items-center mb-6">
                    <div>
                        {activeTab === "API Keys" && (
                            <Button
                                type="primary"
                                onClick={handleGenerateAPIKey}
                            >
                                Generate a new API Key
                            </Button>
                        )}
                        {activeTab === "SMTP" && (
                            <Button type="primary" onClick={handleGenerateSMTPKey}>
                                Generate a new SMTP key
                            </Button>
                        )}
                    </div>
                </div>

                <div className="mb-6">
                    <nav className="flex space-x-4 border-b">
                        <button
                            className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "SMTP" ? "border-blue-500 text-blue-500" : "border-transparent"
                                }`}
                            onClick={() => setActiveTab("SMTP")}
                        >
                            SMTP
                        </button>
                        <button
                            className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "API Keys" ? "border-blue-500 text-blue-500" : "border-transparent"
                                }`}
                            onClick={() => setActiveTab("API Keys")}
                        >
                            API Keys
                        </button>
                    </nav>
                </div>

                {activeTab === "API Keys" && (
                    <>
                        <APIKeysTableComponent />
                        <APIInfo />
                    </>
                )}
                {activeTab === "SMTP" && <SMTPKeysTableComponent />}

                <Modal
                    title={modalContent.title}
                    open={isModalOpen}
                    onCancel={() => setIsModalOpen(false)}
                    onOk={handleSubmitKeyForm}
                    okText={`Generate ${keyType} Key`}
                >
                    <p>{modalContent.content}</p>
                    <Input
                        placeholder={`${keyType} Key Name or Description`}
                        value={keyType === "API" ? apiFormValues.name : smtpformValues.key_name}
                        onChange={handleInputChange}
                        required
                    />
                </Modal>

                <Modal
                    title={modalContent.title}
                    open={isKeyModalOpen}
                    onCancel={() => setIsKeyModalOpen(false)}
                    footer={[
                        <Button key="close" type="primary" onClick={() => setIsKeyModalOpen(false)}>
                            Close
                        </Button>,
                    ]}
                >
                    <p>{modalContent.content}</p>
                </Modal>
            </div>
        </HelmetProvider>
    );
};

export default APISettingsDashTemplate;
