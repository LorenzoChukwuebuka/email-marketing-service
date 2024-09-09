import React, { useEffect, useState } from "react";
import SendersDashComponent from "../components/senders/sendersDashComponent";
import DomainDashboardComponent from "../components/domain/domainDashComponent";

import { Modal } from "../../../components";
import useDomainStore from "../../../store/userstore/domainStore";
import useSenderStore from "../../../store/userstore/senderStore";

type TabType = "Domain" | "Sender"

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
    const [isKeyModalOpen, setIsKeyModalOpen] = useState<boolean>(false);
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

    const { createDomain, setDomainFormValues, domainformValues,getAllDomain } = useDomainStore();
    const { createSender, setSenderFormValues, senderFormValues } = useSenderStore()

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);

    const openModal = (title: string, content: string): void => {
        setModalContent({ title, content });
        setIsModalOpen(true);
    };

    const handleCreateDomain = (): void => {
        setKeyType("Domain");
        setDomainFormValues({ domain: "" });
        openModal("Add Domain", "Please provide a name for your new domain.");
    };

    const handleCreateSender = (): void => {
        setKeyType("Sender");
        setSenderFormValues({ name: "" });
        openModal("Add Sender", "Please provide a name for your new sender.");
    };

    const handleSubmitForm = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
        e.preventDefault();
        setIsModalOpen(false);

        try {
            if (keyType === "Domain") {
                await createDomain();
                setModalContent({
                    title: "New Domain Added",
                    content: "Your new domain has been added successfully.",
                });
                await getAllDomain()
            } else if (keyType === "Sender") {
                await createSender();
                setModalContent({
                    title: "New Sender Added",
                    content: "Your new sender has been added successfully.",
                });
            }
        } catch (error) {
            setModalContent({
                title: "Error",
                content: `Failed to add ${keyType}. Please try again.`,
            });
        }

        setIsKeyModalOpen(true);
        if (keyType === "Domain") {
            setDomainFormValues({ domain: "" });
        } else {
            setSenderFormValues({ name: "" });
        }
    };

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const { value } = event.target;
        if (keyType === "Domain") {
            setDomainFormValues({ domain: value });
        } else {
            setSenderFormValues({ name: value });
        }
    };

    return (
        <div className="p-6 max-w-7xl">
            <div className="flex justify-between items-center mb-6">
                <div>
                    {activeTab === "Sender" && (
                        <button
                            onClick={handleCreateSender}
                            className="bg-gray-900 text-white px-4 py-2 rounded-full hover:bg-gray-700 transition-colors"
                        >
                            Add Sender
                        </button>
                    )}
                    {activeTab === "Domain" && (
                        <button
                            onClick={handleCreateDomain}
                            className="bg-gray-900 text-white px-4 py-2 rounded-full hover:bg-gray-700 transition-colors"
                        >
                            Add Domain
                        </button>
                    )}
                </div>
            </div>

            <div className="mb-6">
                <nav className="flex space-x-4 border-b">
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Sender"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("Sender")}
                    >
                        Senders
                    </button>
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Domain"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => setActiveTab("Domain")}
                    >
                        Domains
                    </button>
                </nav>
            </div>

            {activeTab === "Sender" && <SendersDashComponent />}
            {activeTab === "Domain" && <DomainDashboardComponent />}

            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title={modalContent.title}
            >
                <form onSubmit={handleSubmitForm}>
                    <p className="mb-4 text-gray-600">{modalContent.content}</p>
                    <div className="mb-4">
                        <label
                            htmlFor="name"
                            className="block text-sm font-medium text-gray-700"
                        >
                            {keyType} Name
                        </label>
                        <input
                            type="text"
                            id="name"
                            value={keyType === "Domain" ? domainformValues.domain : senderFormValues.name}
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
                            Add {keyType}
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

export default DomainTemplateDash;