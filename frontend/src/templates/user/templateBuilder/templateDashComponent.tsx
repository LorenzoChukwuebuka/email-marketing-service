import { useEffect, useState } from "react";
import MarketingTemplateDash from "./marketingDash";
import TransactionalTemplateDash from './transactionalDash';


const TemplateBuilderDashComponent: React.FC = () => {
    const [activeTab, setActiveTab] = useState<"Transactional" | "Marketing">("Transactional");

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as "Transactional" | "Marketing");
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

            <nav className="flex space-x-8  border-b">
                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Transactional"
                        ? "border-blue-500 text-blue-500"
                        : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => setActiveTab("Transactional")}
                >
                    Transactional Templates
                </button>

                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Marketing"
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
                    <TransactionalTemplateDash />
                </>
            )}

            {activeTab === "Marketing" && (
                <>
                    <MarketingTemplateDash />
                </>
            )}
        </div>


    </>
};

export default TemplateBuilderDashComponent;