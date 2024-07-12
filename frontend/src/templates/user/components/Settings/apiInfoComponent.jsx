import { useState } from "react";

const APIInfo = () => {
  const [activeTab, setActiveTab] = useState("Curl");

  const tabs = ["Curl", "Ruby", "Php", "Python", "Node Js"];

  const codeExamples = {
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

    Ruby: `
# ------------------
# Create a campaign
# ------------------
require 'uri'
require 'net/http'
require 'json'

url = URI("https://api.brevo.com/v3/emailCampaigns")

http = Net::HTTP.new(url.host, url.port)
http.use_ssl = true

request = Net::HTTP::Post.new(url)
request["api-key"] = 'YOUR_API_V3_KEY'
request["content-type"] = 'application/json'

request.body = JSON.dump({
  "name": "Campaign sent via the API",
  "subject": "My subject",
  "sender": { "name": "From name", "email": "myfrommail@mycompany.com" },
  "type": "classic",
  "htmlContent": "Congratulations! You successfully sent this example campaign via the Brevo API.",
  "recipients": { "listIds": [2, 7] },
  "scheduledAt": "2018-01-01 00:00:01"
})

response = http.request(request)
puts response.read_body`,

    "Node Js": `// ------------------
// Create a campaign
// ------------------
const https = require('https');

let data = JSON.stringify({
  "name": "Campaign sent via the API",
  "subject": "My subject",
  "sender": { "name": "From name", "email": "myfrommail@mycompany.com" },
  "type": "classic",
  "htmlContent": "Congratulations! You successfully sent this example campaign via the Brevo API.",
  "recipients": { "listIds": [2, 7] },
  "scheduledAt": "2018-01-01 00:00:01"
});

let options = {
  hostname: 'api.brevo.com',
  port: 443,
  path: '/v3/emailCampaigns',
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'api-key': 'YOUR_API_V3_KEY'
  }
};

let req = https.request(options, (res) => {
  let body = '';
  res.on('data', (chunk) => {
    body += chunk;
  });
  res.on('end', () => {
    console.log(body);
  });
});

req.write(data);
req.end();`,
  };

  return (
    <div className="flex  p-4 bg-gray-100">
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
              className={`px-4 py-2 ${
                activeTab === tab
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
