import { useState } from "react";
import {
  APIKeysComponentTable,
  SMTPKeysTableComponent,
  Modal,
} from "../components";
import useAPIKeyStore from "../../../store/userstore/apiKeyStore";


const APISettingsDashTemplate = () => {
  const [activeTab, setActiveTab] = useState("API Keys");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalContent, setModalContent] = useState({ title: "", content: "" });

  const { isLoading, generateAPIKey } = useAPIKeyStore();



  const openModal = (title, content) => {
    setModalContent({ title, content });
    setIsModalOpen(true);
  };

  const handleGenerateAPIKey = async (e) => {
    e.preventDefault();
    let key = await generateAPIKey();

    if (key) {
      openModal(
        "New API Key Generated",
        `Your API key is ${key?.apiKey}. Note that this will be displayed only once.`
      );
    }
  };

  const handleGenerateSMTPKey = async () => {
    const newSmtpKey = await generateNewSMTPKey();
    openModal(
      "New SMTP Key Generated",
      `Your SMTP key is ${newSmtpKey}. Note that this will be displayed only once.`
    );
  };

  const generateNewSMTPKey = async () => {
    // Simulate API call or key generation logic
    return "SMTP_" + Math.random().toString(36).substr(2, 9);
  };

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">SMTP & API</h1>
        <div>
          {activeTab === "API Keys" && (
            <button
              onClick={handleGenerateAPIKey}
              className="bg-gray-900 text-white px-4 py-2 rounded-full hover:bg-gray-700 transition-colors"
            >
              {!isLoading ? (
                <>Generate a new API Key </>
              ) : (
                <>
                  Please wait
                  <span className="loading loading-dots loading-sm"></span>{" "}
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
            className={`py-2 border-b-2 ${
              activeTab === "SMTP"
                ? "border-blue-500 text-blue-500"
                : "border-transparent hover:border-gray-300"
            } transition-colors`}
            onClick={() => setActiveTab("SMTP")}
          >
            SMTP
          </button>
          <button
            className={`py-2 border-b-2 ${
              activeTab === "API Keys"
                ? "border-blue-500 text-blue-500"
                : "border-transparent hover:border-gray-300"
            } transition-colors`}
            onClick={() => setActiveTab("API Keys")}
          >
            API Keys
          </button>
        </nav>
      </div>

      {activeTab === "API Keys" && <APIKeysComponentTable />}
      {activeTab === "SMTP" && <SMTPKeysTableComponent />}

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={modalContent.title}
      >
        <p className="mb-4">{modalContent.content}</p>
        <div className="flex justify-end space-x-2">
          <button
            onClick={() => setIsModalOpen(false)}
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
