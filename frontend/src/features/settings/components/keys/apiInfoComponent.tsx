import { useState } from "react";
import { Button, Card, Typography, Space, Tabs, Row, Col } from 'antd';
import { BookOutlined } from '@ant-design/icons';

type CodeExamples = {
    [key: string]: string;
};

const { Paragraph, Title } = Typography;

const APIInfo: React.FC = () => {
    //@ts-ignore
    const [activeTab, setActiveTab] = useState<string>("curl");

    const codeExamples: CodeExamples = {
        curl: `# Create a campaign
curl -H 'api-key: YOUR_API_KEY' \\
-X POST \\
-H 'Content-Type: application/json' \\
-d '{
  "name": "Campaign sent via the API",
  "subject": "My subject",
  "sender": {
    "name": "From name", 
    "email": "sender@example.com"
  },
  "type": "classic",
  "htmlContent": "Hello! This is a test campaign via API.",
  "recipients": { 
    "listIds": [1, 2] 
  },
  "scheduledAt": "2024-01-01 12:00:00"
}' \\
'https://api.yourdomain.com/v1/campaigns'`,

        ruby: `# Ruby example
require 'net/http'
require 'json'

uri = URI('https://api.yourdomain.com/v1/campaigns')
http = Net::HTTP.new(uri.host, uri.port)
http.use_ssl = true

request = Net::HTTP::Post.new(uri)
request['api-key'] = 'YOUR_API_KEY'
request['Content-Type'] = 'application/json'

request.body = {
  name: "Campaign sent via Ruby API",
  subject: "My subject",
  sender: {
    name: "From name",
    email: "sender@example.com"
  },
  type: "classic",
  htmlContent: "Hello from Ruby!",
  recipients: {
    listIds: [1, 2]
  }
}.to_json

response = http.request(request)
puts response.body`,

        php: `<?php
// PHP example
$curl = curl_init();

curl_setopt_array($curl, array(
  CURLOPT_URL => 'https://api.yourdomain.com/v1/campaigns',
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_ENCODING => '',
  CURLOPT_MAXREDIRS => 10,
  CURLOPT_TIMEOUT => 0,
  CURLOPT_FOLLOWLOCATION => true,
  CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
  CURLOPT_CUSTOMREQUEST => 'POST',
  CURLOPT_POSTFIELDS => json_encode([
    'name' => 'Campaign sent via PHP API',
    'subject' => 'My subject',
    'sender' => [
      'name' => 'From name',
      'email' => 'sender@example.com'
    ],
    'type' => 'classic',
    'htmlContent' => 'Hello from PHP!',
    'recipients' => [
      'listIds' => [1, 2]
    ]
  ]),
  CURLOPT_HTTPHEADER => array(
    'api-key: YOUR_API_KEY',
    'Content-Type: application/json'
  ),
));

$response = curl_exec($curl);
curl_close($curl);
echo $response;
?>`,

        python: `# Python example
import requests
import json

url = "https://api.yourdomain.com/v1/campaigns"

payload = {
    "name": "Campaign sent via Python API",
    "subject": "My subject", 
    "sender": {
        "name": "From name",
        "email": "sender@example.com"
    },
    "type": "classic",
    "htmlContent": "Hello from Python!",
    "recipients": {
        "listIds": [1, 2]
    }
}

headers = {
    'api-key': 'YOUR_API_KEY',
    'Content-Type': 'application/json'
}

response = requests.post(url, json=payload, headers=headers)
print(response.json())`,

        nodejs: `// Node.js example
const axios = require('axios');

const config = {
  method: 'post',
  url: 'https://api.yourdomain.com/v1/campaigns',
  headers: {
    'api-key': 'YOUR_API_KEY',
    'Content-Type': 'application/json'
  },
  data: {
    name: "Campaign sent via Node.js API",
    subject: "My subject",
    sender: {
      name: "From name",
      email: "sender@example.com"
    },
    type: "classic",
    htmlContent: "Hello from Node.js!",
    recipients: {
      listIds: [1, 2]
    }
  }
};

axios(config)
  .then(response => {
    console.log(JSON.stringify(response.data));
  })
  .catch(error => {
    console.log(error);
  });`
    };

    const tabItems = [
        {
            key: 'curl',
            label: 'cURL',
            children: (
                <pre style={{
                    background: '#f6f8fa',
                    padding: '16px',
                    borderRadius: '6px',
                    fontSize: '13px',
                    lineHeight: '1.5',
                    overflow: 'auto'
                }}>
                    <code>{codeExamples.curl}</code>
                </pre>
            )
        },
        {
            key: 'ruby',
            label: 'Ruby',
            children: (
                <pre style={{
                    background: '#f6f8fa',
                    padding: '16px',
                    borderRadius: '6px',
                    fontSize: '13px',
                    lineHeight: '1.5',
                    overflow: 'auto'
                }}>
                    <code>{codeExamples.ruby}</code>
                </pre>
            )
        },
        {
            key: 'php',
            label: 'PHP',
            children: (
                <pre style={{
                    background: '#f6f8fa',
                    padding: '16px',
                    borderRadius: '6px',
                    fontSize: '13px',
                    lineHeight: '1.5',
                    overflow: 'auto'
                }}>
                    <code>{codeExamples.php}</code>
                </pre>
            )
        },
        {
            key: 'python',
            label: 'Python',
            children: (
                <pre style={{
                    background: '#f6f8fa',
                    padding: '16px',
                    borderRadius: '6px',
                    fontSize: '13px',
                    lineHeight: '1.5',
                    overflow: 'auto'
                }}>
                    <code>{codeExamples.python}</code>
                </pre>
            )
        },
        {
            key: 'nodejs',
            label: 'Node.js',
            children: (
                <pre style={{
                    background: '#f6f8fa',
                    padding: '16px',
                    borderRadius: '6px',
                    fontSize: '13px',
                    lineHeight: '1.5',
                    overflow: 'auto'
                }}>
                    <code>{codeExamples.nodejs}</code>
                </pre>
            )
        }
    ];

    return (
        <Row gutter={[24, 24]} style={{ padding: '24px' }}>
            <Col xs={24} lg={8}>
                <Card>
                    <Space direction="vertical" size="large">
                        <div>
                            <Title level={3}>
                                <BookOutlined /> About the API
                            </Title>
                            <Paragraph>
                                The {import.meta.env.VITE_API_NAME || 'Platform'} API makes it easy for developers
                                to integrate many features into other applications.
                                Build powerful integrations with our RESTful API.
                            </Paragraph>
                        </div>

                        <div>
                            <Title level={4}>Features</Title>
                            <ul style={{ paddingLeft: '20px' }}>
                                <li>RESTful API design</li>
                                <li>JSON request/response format</li>
                                <li>Comprehensive error handling</li>
                                <li>Rate limiting protection</li>
                                <li>Detailed documentation</li>
                            </ul>
                        </div>

                        <Button
                            type="primary"
                            icon={<BookOutlined />}
                            size="large"
                            href="#"
                        >
                            Read Full Documentation
                        </Button>
                    </Space>
                </Card>
            </Col>

            <Col xs={24} lg={16}>
                <Card title="API Examples" style={{ height: '100%' }}>
                    <Tabs
                        defaultActiveKey="curl"
                        items={tabItems}
                        onChange={setActiveTab}
                        size="large"
                    />
                </Card>
            </Col>
        </Row>
    );
};

export default APIInfo;
