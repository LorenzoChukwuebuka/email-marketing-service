import { useEffect, useState } from "react";
import EditorSelection from "../../../components/editorSelectorComponent";
import CreateMarketingTemplate from "../components/templates/createMarketingTemplate";



type templateTypes = "Templates Gallery" | "My Templates" | "Blank Template" | "Code your own"

const CreateMarketingTemplateDashBoard: React.FC = () => {

    const [activeTab, setActiveTab] = useState<templateTypes>("Templates Gallery");
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

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

    return <>

        <div className="p-6 max-w-full">
            <h1 className="text-xl font-semibold mb-5"> Create Marketing Templates </h1>
            <nav className="flex space-x-8  border-b">
                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Templates Gallery"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Templates Gallery")}
                >
                    Templates Gallery
                </button>

                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Blank Template"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => {
                        setActiveTab("Blank Template");
                        setIsModalOpen(true);
                    }}
                >
                    Blank Template
                </button>

                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "My Templates"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("My Templates")}
                >
                    My Templates
                </button>

                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Code your own"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Code your own")}
                >
                    Code your own
                </button>


            </nav>


            {activeTab === "Templates Gallery" && (
                <>
                    kai
                </>
            )}

            {activeTab === "Code your own" && (
                <>
                    <EditorSelection />
                </>
            )}


            <CreateMarketingTemplate isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />


        </div>

    </>
}


export default CreateMarketingTemplateDashBoard