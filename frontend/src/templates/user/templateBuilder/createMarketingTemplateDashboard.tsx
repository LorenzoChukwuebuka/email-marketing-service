import { useEffect, useState } from "react";
import CreateMarketingTemplate from "../components/templates/createMarketingTemplate";

type templateTypes = "Templates Gallery" | "Blank Template" | "Custom HTML";

const CreateMarketingTemplateDashBoard: React.FC = () => {
    const [activeTab, setActiveTab] = useState<templateTypes>("Templates Gallery");
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(false);

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as templateTypes);
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

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setIsLoading(true);
        setTimeout(() => {
            setIsLoading(false);
        }, 30000);
    };

    const handleTabChange = (newTab: templateTypes) => {
        if (isLoading) {
            setIsLoading(false);
        }
        setActiveTab(newTab);
        if (newTab !== "Templates Gallery") {
            setIsModalOpen(true);
        }
    };

    return (
        <div className="p-6 max-w-full">
            <h1 className="text-xl font-semibold mb-5"> Create Marketing Templates </h1>
            <nav className="flex space-x-8 border-b">
                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Templates Gallery"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => handleTabChange("Templates Gallery")}
                >
                    Templates Gallery
                </button>

                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Blank Template"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => handleTabChange("Blank Template")}
                >
                    Blank Template
                </button>

                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Custom HTML"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => handleTabChange("Custom HTML")}
                >
                    Custom HTML
                </button>
            </nav>

            {activeTab === "Templates Gallery" && (
                <>
                    Nothing here
                </>
            )}

            {isLoading && <div className="flex items-center justify-center mt-20"><span className="loading loading-spinner loading-lg"></span></div>}

            <CreateMarketingTemplate isOpen={isModalOpen} onClose={handleCloseModal} editorType={activeTab === "Blank Template" ? "drag-and-drop" : "html-editor"} />
        </div>
    );
};

export default CreateMarketingTemplateDashBoard;