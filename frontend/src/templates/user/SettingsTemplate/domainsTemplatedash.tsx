import { useEffect, useState } from "react";
import SendersDashComponent from "../components/senders/sendersDashComponent";

type TabType = "Domain" | "Sender"
const DomainTemplateDash: React.FC = () => {

    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "Sender" || storedTab === "Domain") ? storedTab : "Sender";
    });
    const [keyType, setKeyType] = useState<TabType | null>(null);

    useEffect(() => {
        const storedActiveTab = localStorage.getItem("activeTab");
        if (storedActiveTab) {
            setActiveTab(storedActiveTab as "Domain" | "Sender");
        }
    }, []);

    useEffect(() => {
        localStorage.setItem("activeTab", activeTab);
    }, [activeTab]);



    return <>
        <div className="mb-6 mt-10 p-4">
            <nav className="flex space-x-4 mt-5 border-b">
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


        {activeTab === "Sender" && (
            <>
                <SendersDashComponent />
            </>
        )}

        {activeTab === "Domain" && (
            <>
                afasdfdf
            </>
        )}

    </>
}

export default DomainTemplateDash