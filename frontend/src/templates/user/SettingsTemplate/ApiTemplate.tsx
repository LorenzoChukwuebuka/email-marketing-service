import React, { useEffect, useState } from "react";
import {
    APIInfo,
    APIKeysComponentTable,
    SMTPKeysTableComponent,
} from "../components";
import useAPIKeyStore from "../../../store/userstore/apiKeyStore";
import { Modal } from "../../../components";
import useSMTPKeyStore from "../../../store/userstore/smtpkeyStore";

interface ModalContent {
    title: string;
    content: string;
}

const APISettingsDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<"API Keys" | "SMTP">("SMTP");
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [isKeyModalOpen, setIsKeyModalOpen] = useState<boolean>(false);
    const [modalContent, setModalContent] = useState<ModalContent>({ title: "", content: "" });
    const [keyType, setKeyType] = useState<"API" | "SMTP" | null>(null);

    const {
        isLoading,
        generateAPIKey,
        setFormValues: setAPIFormValues,
        formValues: apiFormValues,
        getAPIKey,
    } = useAPIKeyStore();

    const {
        createSMTPKey,
        setSmtpFormValues,
        getSMTPKeys,
        smtpformValues
    } = useSMTPKeyStore();

    const openModal = (title: string, content: string): void => {
        setModalContent({ title, content });
        setIsModalOpen(true);
    };

    const handleGenerateAPIKey = (): void => {
        setKeyType("API");
        setAPIFormValues({ name: "" });
        openModal("Generate New API Key", "Please provide a name or description for your new API key.");
    };

    const handleGenerateSMTPKey = (): void => {
        setKeyType("SMTP");
        setSmtpFormValues({ key_name: "" });
        openModal("Generate New SMTP Key", "Please provide a name or description for your new SMTP key.");
    };

    const handleSubmitKeyForm = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
        e.preventDefault();
        setIsModalOpen(false);

        try {
            if (keyType === "SMTP") {
                await createSMTPKey();
                await getSMTPKeys();
                setModalContent({
                    title: "New SMTP Key Generated",
                    content: "Your new SMTP key has been generated successfully. You can view it in the SMTP keys table.",
                });
            } else if (keyType === "API") {
                await generateAPIKey();
                getAPIKey();
                setModalContent({
                    title: "New API Key Generated",
                    content: "Your new API key has been generated successfully. You can view it in the API keys table.",
                });
            }
        } catch (error) {
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


    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as "API Keys" | "SMTP");
        }
    }, []);

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);

    useEffect(() => {
        return () => {
            localStorage.removeItem("activeTab");
        };
    }, []);

    return (
        <div className="p-6 max-w-5xl">
            <div className="flex justify-between items-center mb-6">

                <div>
                    {activeTab === "API Keys" && (
                        <button
                            onClick={handleGenerateAPIKey}
                            className="bg-gray-900 text-white px-4 py-2 rounded-full hover:bg-gray-700 transition-colors"
                        >
                            {!isLoading ? (
                                <>Generate a new API Key</>
                            ) : (
                                <>
                                    Please wait
                                    <span className="loading loading-dots loading-sm"></span>
                                </>
                            )}
                        </button>
                    )}
                    {activeTab === "SMTP" && (
                        <button
                            onClick={handleGenerateSMTPKey}
                            className="bg-gray-900 text-white px-4 py-2 rounded-full hover:bg-gray-700 transition-colors"
                        >
                            Generate a new SMTP key
                        </button>
                    )}
                </div>
            </div>

            <div className="mb-6">
                <nav className="flex space-x-4 border-b">
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "SMTP"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("SMTP")}
                    >
                        SMTP
                    </button>
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "API Keys"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("API Keys")}
                    >
                        API Keys
                    </button>
                </nav>
            </div>

            {activeTab === "API Keys" && (
                <>
                    <APIKeysComponentTable />
                    <APIInfo />
                </>
            )}
            {activeTab === "SMTP" && (
                <>
                    <SMTPKeysTableComponent />
                </>
            )}

            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title={modalContent.title}
            >
                <form onSubmit={handleSubmitKeyForm}>
                    <p className="mb-4 text-gray-600">{modalContent.content}</p>
                    <div className="mb-4">
                        <label
                            htmlFor="keyName"
                            className="block text-sm font-medium text-gray-700"
                        >
                            {keyType} Key Name or Description
                        </label>
                        <input
                            type="text"
                            id="keyName"
                            value={keyType === "API" ? apiFormValues.name : smtpformValues.key_name}
                            onChange={handleInputChange}
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            required
                        />
                    </div>
                    <div className="flex justify-end space-x-2">
                        <button
                            type="button"
                            onClick={() => setIsModalOpen(false)}
                            className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                        >
                            Generate {keyType} Key
                        </button>
                    </div>
                </form>
            </Modal>

            <Modal
                isOpen={isKeyModalOpen}
                onClose={() => setIsKeyModalOpen(false)}
                title={modalContent.title}
            >
                <p className="mb-4">{modalContent.content}</p>
                <div className="flex justify-end space-x-2">
                    <button
                        onClick={() => setIsKeyModalOpen(false)}
                        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                    >
                        Close
                    </button>
                </div>
            </Modal>
        </div>
    );
};

export default APISettingsDashTemplate;