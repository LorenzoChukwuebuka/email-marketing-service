import { useEffect, useState } from "react";
import TransactionBuilderComponent from "../../../components/templateBuilder";
import MarketingTemplateDash from "../components/templates/marketingDash";


const TemplateBuilderDashComponent: React.FC = () => {
    const [activeTab, setActiveTab] = useState<"Transactional" | "Marketing">("Transactional");

    useEffect(() => {
        // Load the active tab from localStorage on mount
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as "Transactional" | "Marketing");
        }
    }, []);

    useEffect(() => {
        // Save the active tab to localStorage whenever it changes
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);

    useEffect(() => {
        // Clear the activeTab from localStorage when the component is unmounted
        return () => {
            localStorage.removeItem("activeTab");
        };
    }, []);
    return <>

        <div className="p-6 max-w-full">

            <nav className="flex space-x-8  border-b">
                <button
                    className={`py-2 border-b-2 text-xl font-semibold ${activeTab === "Transactional"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Transactional")}
                >
                    Transactional Templates
                </button>

                <button
                    className={`py-2 border-b-2 text-xl font-semibold ${activeTab === "Marketing"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Marketing")}
                >
                    Marketing Templates
                </button>


            </nav>


            {activeTab === "Transactional" && (
                <>
              
                </>
            )}

            {activeTab === "Marketing" && (
                <>
                    <MarketingTemplateDash/>
                </>
            )}
        </div>


    </>
};

export default TemplateBuilderDashComponent;