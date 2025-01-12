import React from 'react';
import { HelmetProvider, Helmet } from 'react-helmet-async';
import useMetadata from '../../hooks/useMetaData';
import NavBar from '../../components/landingpage/navbar';
import Footer from '../../components/landingpage/footer';

const TermsOfService: React.FC = () => {
    const metaData = useMetadata("TermsAndConditions")
    return (
        <HelmetProvider>
            <Helmet {...metaData} />
            <NavBar />
            <div className="bg-gray-50 py-10 px-5">
                <div className="max-w-4xl mx-auto bg-white p-10 rounded-lg">
                    <h1 className="text-3xl font-bold mb-5 text-gray-800">Terms of Service</h1>
                    <p className="mb-6 text-gray-700">
                        <strong>Effective Date:</strong> [Insert Date]
                    </p>
                    <p className="mb-6 text-gray-700">
                        Welcome to [Your Company Name] ("we", "our", "us"). These Terms of Service ("Terms") govern your access to and use of our email marketing platform ("Service"). By using the Service, you agree to comply with and be bound by these Terms.
                    </p>

                    {/* Section 1 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">1. Acceptance of Terms</h2>
                    <p className="mb-6 text-gray-700">
                        By creating an account or using our Service, you agree to be legally bound by these Terms. If you do not agree with any part of these Terms, you may not use the Service.
                    </p>

                    {/* Section 2 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">2. Eligibility</h2>
                    <p className="mb-6 text-gray-700">
                        You must be at least 18 years old to use the Service. By using our Service, you confirm that you have the legal authority to enter into these Terms and comply with all applicable laws and regulations.
                    </p>

                    {/* Section 3 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">3. User Accounts</h2>
                    <p className="mb-6 text-gray-700">
                        To access certain features, you may need to create an account. You are responsible for safeguarding your account credentials, and you agree not to share your credentials with any third party. You are responsible for all activities that occur under your account.
                    </p>

                    {/* Section 4 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">4. Compliance with Email Marketing Laws</h2>
                    <p className="mb-6 text-gray-700">
                        When using our Service, you agree to comply with all applicable email marketing laws, including but not limited to:
                    </p>
                    <ul className="list-disc list-inside mb-6 text-gray-700">
                        <li>CAN-SPAM Act (United States)</li>
                        <li>General Data Protection Regulation (GDPR) (European Union)</li>
                        <li>Canadian Anti-Spam Legislation (CASL)</li>
                    </ul>

                    {/* Section 5 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">5. Prohibited Activities</h2>
                    <p className="mb-6 text-gray-700">
                        You agree not to use the Service for unlawful purposes, including but not limited to:
                    </p>
                    <ul className="list-disc list-inside mb-6 text-gray-700">
                        <li>Sending unsolicited bulk emails or spam</li>
                        <li>Engaging in fraudulent or deceptive practices</li>
                        <li>Violating any privacy or data protection laws</li>
                    </ul>

                    {/* Section 6 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">6. Data Protection and Privacy</h2>
                    <p className="mb-6 text-gray-700">
                        We take the privacy of your data seriously. By using the Service, you agree to our collection and use of your data as described in our Privacy Policy.
                    </p>

                    {/* Section 7 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">7. Subscription and Payment</h2>
                    <p className="mb-6 text-gray-700">
                        Some features of our Service may require a paid subscription. You agree to pay all applicable fees for using paid features as outlined at the time of purchase. All fees are non-refundable unless otherwise stated.
                    </p>

                    {/* Section 8 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">8. Termination</h2>
                    <p className="mb-6 text-gray-700">
                        We reserve the right to suspend or terminate your account if you violate these Terms or engage in activities that harm the Service or other users.
                    </p>

                    {/* Section 9 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">9. Limitation of Liability</h2>
                    <p className="mb-6 text-gray-700">
                        To the fullest extent permitted by law, [Your Company Name] shall not be liable for any indirect, incidental, or consequential damages arising from your use of the Service.
                    </p>

                    {/* Section 10 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">10. Changes to Terms</h2>
                    <p className="mb-6 text-gray-700">
                        We may update these Terms from time to time. We will notify you of any changes by updating the "Effective Date" at the top of these Terms. Your continued use of the Service after any changes constitutes acceptance of those changes.
                    </p>

                    {/* Section 11 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">11. Governing Law</h2>
                    <p className="mb-6 text-gray-700">
                        These Terms will be governed by the laws of [Your Jurisdiction], without regard to its conflict of law provisions.
                    </p>

                    {/* Section 12 */}
                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">12. Contact Us</h2>
                    <p className="mb-6 text-gray-700">
                        If you have any questions or concerns about these Terms, please contact us at:
                    </p>
                    <p className="text-gray-700">
                        [Your Company Name] <br />
                        [Your Address] <br />
                        [Your Email Address] <br />
                        [Your Phone Number]
                    </p>
                </div>
            </div>
            <Footer />
        </HelmetProvider>
    );
};

export default TermsOfService;
