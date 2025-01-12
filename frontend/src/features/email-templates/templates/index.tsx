import { useEffect, useState } from "react";
import useMetadata from "../../../hooks/useMetaData";
import { Helmet, HelmetProvider } from "react-helmet-async";
import MarketingTemplateDash from './marketing/marketingDashTemplate';
import TransactionalTemplateDash from './transactional/transactionalDashTemplate';

type Tabtype = 'Transactional' | 'Marketing'

const TemplateBuilderDashComponent: React.FC = () => {
    const [activeTab, setActiveTab] = useState<Tabtype>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "Transactional" || storedTab === "Marketing") ? storedTab : "Transactional";
    });

    const metaData = useMetadata("TransactionalTemplates")

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as "Transactional" | "Marketing");
        }
    }, []);

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);


    return <>
        <HelmetProvider>
            <Helmet {...metaData} title={activeTab === "Marketing" ? "Marketing Templates - CrabMailer" : "Transactional Templates - CrabMailer"} />

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

                {activeTab === "Transactional" &&
                    <TransactionalTemplateDash />
                }

                {activeTab === "Marketing" &&
                    <MarketingTemplateDash />
                }
            </div>

        </HelmetProvider>
    </>
};

export default TemplateBuilderDashComponent;