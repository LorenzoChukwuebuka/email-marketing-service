import { useState } from "react";
import ContactsDashComponent from "../components/contacts/contactDashComponent";
import ContactGroupDash from "../components/contactGroup/contactGroupDashComponent";

type TabType = "Contact" | "Contact Group";

const ContactDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "Contact" || storedTab === "Contact Group") ? storedTab : "Contact";
    });

    const handleTabChange = (tab: TabType) => {
        setActiveTab(tab);
        localStorage.setItem("activeTab", tab);
    };

    return (
        <div className="p-6 max-w-full">
            <nav className="flex space-x-8 border-b">
                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${
                        activeTab === "Contact"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                    } transition-colors`}
                    onClick={() => handleTabChange("Contact")}
                >
                    Contact
                </button>
                <button
                    className={`py-2 border-b-2 text-lg font-semibold ${
                        activeTab === "Contact Group"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                    } transition-colors`}
                    onClick={() => handleTabChange("Contact Group")}
                >
                    Contact Group
                </button>
            </nav>

            {activeTab === "Contact" && <ContactsDashComponent />}
            {activeTab === "Contact Group" && <ContactGroupDash />}
        </div>
    );
};

export default ContactDashTemplate;