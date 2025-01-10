import { useState } from "react";
// import ContactsDashComponent from "../components/contacts/contactDashComponent";
import ContactGroupDash from '../components/contactgroup/contactgroupdash';
// import ContactGroupDash from "../components/contactGroup/contactGroupDashComponent";
import useMetadata from "../../../hooks/useMetaData";
import { Helmet, HelmetProvider } from "react-helmet-async";
import ContactsDashComponent from "../components/contacts/contactDashComponent";

type TabType = "Contact" | "Contact Group" | "Segments";

const ContactDashTemplate: React.FC = () => {
    const [activeTab, setActiveTab] = useState<TabType>(() => {
        const storedTab = localStorage.getItem("activeTab");
        return (storedTab === "Contact" || storedTab === "Contact Group" || storedTab === "Segments") ? storedTab : "Contact";
    });

    const handleTabChange = (tab: TabType) => {
        setActiveTab(tab);
        localStorage.setItem("activeTab", tab);
    };


    const metaData = useMetadata("Contact")

    return (
        <HelmetProvider>
            <Helmet {...metaData} title={activeTab === "Contact" ? "Contacts - CrabMailer" : "Contact Groups - CrabMailer"} />
            <div className="p-6 max-w-full">
                <nav className="flex space-x-8 border-b">
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Contact"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("Contact")}
                    >
                        Contact
                    </button>
                    <button
                        className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Contact Group"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                            } transition-colors`}
                        onClick={() => handleTabChange("Contact Group")}
                    >
                        Contact Group
                    </button>

                    {/* <button
                    className={`py-2 border-b-2 text-lg font-semibold ${activeTab === "Segments"
                            ? "border-blue-500 text-blue-500"
                            : "border-transparent hover:border-gray-300"
                        } transition-colors`}
                    onClick={() => handleTabChange("Segments")}
                >
                    Segments
                </button> */}
                </nav>

                 {activeTab === "Contact" && <ContactsDashComponent/>}

                {activeTab === "Contact Group" && <ContactGroupDash />}  

                {activeTab === "Segments" && <> hello world </>}
            </div>
        </HelmetProvider>
    );
};

export default ContactDashTemplate;