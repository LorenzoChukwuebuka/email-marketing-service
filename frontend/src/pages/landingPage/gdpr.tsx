import React from "react";
import { Helmet, HelmetProvider } from "react-helmet-async";
import NavBar from "./navbar";
import Footer from "./footer";

const GDPRExplanation = () => {
    return (
        <HelmetProvider>
            <Helmet>
                <title>Understanding GDPR - Your Data Rights</title>
                <meta name="description" content="Learn about GDPR and your data protection rights" />
            </Helmet>
            <NavBar/>
            <div className="bg-gray-50 py-10 px-5">
                <div className="max-w-4xl mx-auto bg-white p-10 rounded-lg shadow-md">
                    <h1 className="text-3xl font-bold mb-5 text-gray-800">Understanding GDPR</h1>
                    <p className="mb-6 text-gray-700">
                        The General Data Protection Regulation (GDPR) is a comprehensive data protection law that came into effect on May 25, 2018, in the European Union (EU) and European Economic Area (EEA).
                    </p>

                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        1. What is GDPR?
                    </h2>
                    <p className="mb-6 text-gray-700">
                        GDPR is a legal framework that sets guidelines for the collection and processing of personal information from individuals who live in the EU. It standardizes data protection law across all 27 EU countries and imposes strict rules on controlling and processing personally identifiable information (PII).
                    </p>

                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        2. Your Core Rights Under GDPR
                    </h2>
                    <div className="grid md:grid-cols-2 gap-6 mb-6">
                        {[
                            {
                                title: "Right to Access",
                                description: "You can request access to your personal data and ask how it's being used"
                            },
                            {
                                title: "Right to Rectification",
                                description: "You can request to update or correct your personal information"
                            },
                            {
                                title: "Right to Erasure",
                                description: "Also known as 'Right to be Forgotten' - request deletion of your data"
                            },
                            {
                                title: "Right to Restrict Processing",
                                description: "You can limit how your data is used while still stored"
                            },
                            {
                                title: "Right to Data Portability",
                                description: "You can request your data in a portable format and transfer it"
                            },
                            {
                                title: "Right to Object",
                                description: "You can object to how your data is processed"
                            }
                        ].map((right, index) => (
                            <div key={index} className="bg-gray-50 p-4 rounded-lg">
                                <h3 className="text-lg font-semibold mb-2 text-gray-800">{right.title}</h3>
                                <p className="text-gray-700">{right.description}</p>
                            </div>
                        ))}
                    </div>

                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        3. Key GDPR Principles
                    </h2>
                    <ul className="list-disc list-inside mb-6 space-y-3 text-gray-700">
                        <li><strong>Lawfulness, Fairness, and Transparency:</strong> Personal data must be processed legally, fairly, and transparently.</li>
                        <li><strong>Purpose Limitation:</strong> Data must be collected for specified, explicit, and legitimate purposes.</li>
                        <li><strong>Data Minimization:</strong> Only collect and process data that's necessary for the intended purpose.</li>
                        <li><strong>Accuracy:</strong> Personal data must be accurate and kept up to date.</li>
                        <li><strong>Storage Limitation:</strong> Data should be kept only as long as necessary.</li>
                        <li><strong>Integrity and Confidentiality:</strong> Data must be processed securely.</li>
                    </ul>

                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        4. When GDPR Applies
                    </h2>
                    <div className="bg-blue-50 p-6 rounded-lg mb-6">
                        <h3 className="text-xl font-semibold mb-3 text-gray-800">GDPR applies when:</h3>
                        <ul className="list-disc list-inside space-y-2 text-gray-700">
                            <li>A company processes personal data of EU residents</li>
                            <li>A company offers goods or services to EU residents</li>
                            <li>A company monitors the behavior of EU residents</li>
                            <li>A company has an establishment in the EU</li>
                        </ul>
                    </div>

                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        5. Consent Requirements
                    </h2>
                    <div className="bg-green-50 p-6 rounded-lg mb-6">
                        <p className="mb-4 text-gray-700">Under GDPR, consent must be:</p>
                        <ul className="list-disc list-inside space-y-2 text-gray-700">
                            <li><strong>Freely Given:</strong> No pressure or negative consequences for refusing</li>
                            <li><strong>Specific:</strong> Separate consent for different processing activities</li>
                            <li><strong>Informed:</strong> Clear explanation of what is being consented to</li>
                            <li><strong>Unambiguous:</strong> Clear affirmative action required (no pre-ticked boxes)</li>
                            <li><strong>Withdrawable:</strong> As easy to withdraw consent as to give it</li>
                        </ul>
                    </div>

                    <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                        6. Penalties for Non-Compliance
                    </h2>
                    <div className="bg-red-50 p-6 rounded-lg mb-6">
                        <p className="mb-4 text-gray-700">GDPR violations can result in:</p>
                        <ul className="list-disc list-inside space-y-2 text-gray-700">
                            <li>Fines up to €20 million or 4% of global annual revenue, whichever is higher</li>
                            <li>Regular audits and monitoring</li>
                            <li>Temporary or permanent ban on data processing</li>
                            <li>Compensation claims from affected individuals</li>
                        </ul>
                    </div>

                    <div className="mt-8 p-6 bg-gray-100 rounded-lg">
                        <h2 className="text-2xl font-semibold mb-4 text-gray-800">
                            Need More Information?
                        </h2>
                        <p className="text-gray-700">
                            For more detailed information about GDPR and how it affects your rights, you can:
                        </p>
                        <ul className="list-disc list-inside mt-3 space-y-2 text-gray-700">
                            <li>Visit the official EU GDPR portal</li>
                            <li>Contact your national data protection authority</li>
                            <li>Consult with a data protection officer</li>
                            <li>Read our detailed privacy policy</li>
                        </ul>
                    </div>
                </div>
            </div>
            <Footer/>
        </HelmetProvider>
    );
};

export default GDPRExplanation;