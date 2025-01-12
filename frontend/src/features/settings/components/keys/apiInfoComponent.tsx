import { useState } from "react";

type CodeExamples = {
    [key: string]: string;
};

const APIInfo: React.FC = () => {
    const [activeTab, setActiveTab] = useState<string>("Curl");

    const tabs: string[] = ["Curl", "Ruby", "Php", "Python", "Node Js"];

    const codeExamples: CodeExamples = {
        Curl: `
# ------------------
# Create a campaign
# ------------------
curl -H 'api-key: YOUR_API_V3_KEY'
-X POST -d '{
# Define the campaign settings
"name":"Campaign sent via the API",
"subject":"My subject",
"sender": {"name":"From name", "email":"myfrommail@mycompany.com"},
"type":"classic",
# Content that will be sent
"htmlContent":"Congratulations! You successfully sent this example campaign via the Brevo API.",
# Select the recipients
"recipients": { "listIds": [2,7] },
# Schedule the sending in one hour
"scheduledAt": "2018-01-01 00:00:01",
}'
'https://api.brevo.com/v3/emailCampaigns'`,


    };

    return (
        <div className="flex p-4 bg-gray-100">
            <div className="w-1/3 pr-4">
                <h2 className="text-xl font-bold mb-2">About the API</h2>
                <p className="mb-4">
                    The {import.meta.env.VITE_API_NAME} API makes it easy for programmers
                    to integrate many of
                    {import.meta.env.VITE_API_NAME}`s features into other applications.
                    Interested in learning more?
                </p>
                <a href="#" className="text-blue-600 hover:underline">
                    Read our API documentation
                </a>
            </div>
            <div className="w-2/3 bg-white rounded-lg shadow-sm p-4">
                <div className="flex mb-4 border-b">
                    {tabs.map((tab) => (
                        <button
                            key={tab}
                            className={`px-4 py-2 ${activeTab === tab
                                ? "text-blue-600 border-b-2 border-blue-600"
                                : "text-gray-600"
                                }`}
                            onClick={() => setActiveTab(tab)}
                        >
                            {tab}
                        </button>
                    ))}
                </div>
                <pre className="bg-gray-100 p-4 rounded text-sm overflow-x-auto">
                    <code>{codeExamples[activeTab] || "Code example not available"}</code>
                </pre>
            </div>
        </div>
    );
};

export default APIInfo;
