import React from "react";
import { Helmet, HelmetProvider } from "react-helmet-async";
import useMetadata from "../../hooks/useMetaData";
import Footer from "./footer";
import NavBar from "./navbar";

const PrivacyPolicy = () => {
    const metaData = useMetadata()("PrivacyPolicy")
    return (
        <HelmetProvider>
            <Helmet {...metaData} />
            <NavBar />
            <div className="bg-gray-50  py-10 px-5">
                <div className="max-w-4xl mx-auto bg-white p-10 rounded-lg">
                    <h1 className="text-3xl font-bold mb-5 text-gray-800">Privacy Policy</h1>
                    <p className="mb-6 text-gray-700">
                        <strong>Effective Date:</strong> [Insert Date]
                    </p>
                    <p className="mb-6 text-gray-700">
                        [Your Company Name] ("we," "our," or "us") is committed to protecting
                        the privacy and security of our users' personal information. This
                        Privacy Policy outlines how we collect, use, disclose, and safeguard
                        your information when you use our email marketing services
                        ("Services"). By using the Services, you agree to the collection and
                        use of information in accordance with this policy.
                    </p>

                    {/* Section 1 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        1. Information We Collect
                    </h2>
                    <p className="mb-4 text-gray-700">
                        We collect several types of information to provide and improve our
                        Services, including:
                    </p>

                    <h3 className="text-xl font-semibold mb-2 text-gray-800">a) Personal Data</h3>
                    <ul className="list-disc list-inside mb-6 text-gray-700">
                        <li><strong>Email Address:</strong> Required for account creation, sending emails, and delivering content.</li>
                        <li><strong>Name and Contact Information:</strong> To personalize emails and ensure proper communication.</li>
                        <li><strong>Billing Information:</strong> Including credit card details and billing address to process transactions.</li>
                    </ul>

                    <h3 className="text-xl font-semibold mb-2 text-gray-800">b) User Data</h3>
                    <ul className="list-disc list-inside mb-6 text-gray-700">
                        <li><strong>Behavioral Data:</strong> How users interact with emails (open rates, clicks, and other engagement metrics).</li>
                        <li><strong>IP Address and Device Information:</strong> For security purposes, location tracking, and to improve our services.</li>
                        <li><strong>Third-party Integrations Data:</strong> Information from platforms you connect to our service, such as CRM or e-commerce platforms.</li>
                    </ul>

                    {/* Section 2 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        2. How We Use Your Information
                    </h2>
                    <ul className="list-disc list-inside mb-6 text-gray-700">
                        <li><strong>Providing Services:</strong> To facilitate the sending and delivery of marketing emails on behalf of our users.</li>
                        <li><strong>Personalization:</strong> Tailoring email content based on user preferences and behavior.</li>
                        <li><strong>Analytics and Improvements:</strong> Analyzing engagement with emails to improve the effectiveness of campaigns and user experience.</li>
                        <li><strong>Compliance:</strong> Ensuring compliance with legal obligations, including anti-spam regulations (such as CAN-SPAM, GDPR, and others).</li>
                        <li><strong>Security:</strong> Monitoring and protecting our systems against fraud, unauthorized access, and illegal activity.</li>
                    </ul>

                    {/* Section 3 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        3. Data Sharing and Disclosure
                    </h2>
                    <p className="mb-6 text-gray-700">
                        We do not share your personal information with third parties except
                        in the following cases:
                    </p>
                    <ul className="list-disc list-inside mb-6 text-gray-700">
                        <li><strong>Service Providers:</strong> We may share your information with third-party vendors who help us provide our services (e.g., email delivery, data storage, analytics).</li>
                        <li><strong>Legal Requirements:</strong> We may disclose your information if required by law, such as in response to a court order or to comply with legal regulations.</li>
                        <li><strong>Business Transfers:</strong> If we undergo a merger, acquisition, or sale of assets, your information may be transferred as part of that transaction.</li>
                    </ul>

                    {/* Section 4 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        4. Your Rights and Choices
                    </h2>
                    <p className="mb-6 text-gray-700">
                        Depending on your jurisdiction, you may have the following rights
                        regarding your personal data:
                    </p>
                    <ul className="list-disc list-inside mb-6 text-gray-700">
                        <li><strong>Access:</strong> You can request to see the personal information we hold about you.</li>
                        <li><strong>Correction:</strong> You can request that we correct or update inaccurate information.</li>
                        <li><strong>Deletion:</strong> You may request the deletion of your personal data, subject to certain legal restrictions.</li>
                        <li><strong>Opt-Out of Marketing Emails:</strong> You can unsubscribe from marketing communications by following the instructions in those emails or by contacting us directly.</li>
                    </ul>

                    {/* More sections can be added similarly */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        5. Data Security
                    </h2>
                    <p className="mb-6 text-gray-700">
                        We implement appropriate technical and organizational measures to
                        protect your personal data from unauthorized access, disclosure,
                        alteration, or destruction. However, no method of transmission over
                        the Internet or electronic storage is 100% secure, and we cannot
                        guarantee its absolute security.
                    </p>

                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        6. Data Retention
                    </h2>
                    <p className="mb-6 text-gray-700">
                        We retain your personal data only as long as necessary to fulfill the
                        purposes outlined in this Privacy Policy, unless a longer retention
                        period is required by law.
                    </p>

                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        7. Changes to This Privacy Policy
                    </h2>
                    <p className="mb-6 text-gray-700">
                        We may update this Privacy Policy from time to time. We will notify
                        you of any changes by posting the new Privacy Policy on this page and
                        updating the "Effective Date" at the top. We encourage you to review
                        this policy periodically for any changes.
                    </p>

                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        8. Contact Us
                    </h2>
                    <p className="mb-6 text-gray-700">
                        If you have any questions or concerns about this Privacy Policy,
                        please contact us at:
                    </p>
                    <p className="text-gray-700">
                        [Your Company Name]
                        <br />
                        [Your Address]
                        <br />
                        [Your Email Address]
                        <br />
                        [Your Phone Number]
                    </p>
                </div>
            </div>
            <Footer />
        </HelmetProvider>
    );
};

export default PrivacyPolicy;
