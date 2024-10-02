import React, { useEffect, useState } from 'react';
import Cookies from "js-cookie";
import { Link } from "react-router-dom";
import { ArrowRight, Mail, Zap, Shield, Users, ChevronDown, Menu, X } from 'lucide-react';
import mailpicture from "./../../assets/dashboard.jpg"
import { Helmet, HelmetProvider } from 'react-helmet-async';
import useMetadata from '../../hooks/useMetaData';

const IndexLandingPage: React.FC = () => {
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const [activeAccordion, setActiveAccordion] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
        const userCookie = Cookies.get("Cookies"); // Replace "Cookies" with the actual name of your cookie
        if (userCookie) {
            setIsAuthenticated(true); // Set authenticated if cookie exists
        } else {
            setIsAuthenticated(false);
        }
    }, []);

    const features = [
        { icon: <Mail className="w-8 h-8" />, title: 'Powerful Email Campaigns', description: 'Create and send beautiful, responsive emails that drive results.' },
        { icon: <Zap className="w-8 h-8" />, title: 'Smart Automation', description: 'Set up triggers and workflows to nurture leads and customers automatically.' },
        { icon: <Shield className="w-8 h-8" />, title: 'Top-tier Deliverability', description: 'Our advanced infrastructure ensures your emails reach the inbox.' },
        { icon: <Users className="w-8 h-8" />, title: 'Precise Audience Segmentation', description: 'Target the right subscribers with personalized content for maximum engagement.' },
    ];

    const faqs = [
        { question: 'How does the free trial work?', answer: 'Our 14-day free trial gives you full access to all features. No credit card required.' },
        { question: 'Can I integrate with my existing tools?', answer: 'Yes, we offer integrations with popular CRMs, e-commerce platforms, and more.' },
        { question: 'What kind of support do you offer?', answer: '24/7 email support for all plans, with phone and chat support for higher tiers.' },
        { question: 'Is there a limit on subscribers or emails?', answer: 'Plans vary, but we have options for businesses of all sizes, from startups to enterprises.' },
    ];

    const toggleAccordion = (index: any) => {
        setActiveAccordion(activeAccordion === index ? null : index);
    };

    const plans = [
        { name: 'Basic', price: '29', features: ['Feature 1', 'Feature 2', 'Feature 3'] },
        { name: 'Pro', price: '79', features: ['Feature 1', 'Feature 2', 'Feature 3', 'Feature 4'] },
        { name: 'Enterprise', price: '199', features: ['Feature 1', 'Feature 2', 'Feature 3', 'Feature 4', 'Feature 5'] }
    ];

    const apiName = import.meta.env.VITE_API_NAME;
    const firstFourLetters = apiName.slice(0, 4);
    const remainingLetters = apiName.slice(4);

    const metaData = useMetadata()("LandingPage")

    return (

        <HelmetProvider>
            <Helmet {...metaData} />
            <div className="min-h-screen  flex flex-col bg-gray-50">
                <header className="bg-white py-4 shadow-sm sticky top-0 z-50">
                    <div className="container mx-auto px-4 flex justify-between items-center">

                        <h1 className="text-center text-2xl font-bold">
                            <span className="text-indigo-700">{firstFourLetters}</span>
                            <span className="text-gray-700">{remainingLetters}</span>
                            <i className="bi bi-mailbox2-flag text-indigo-700 ml-2"></i>
                        </h1>
                        <nav className="hidden md:flex space-x-6">
                            <a href="#features" className="text-gray-600 hover:text-gray-800">Features</a>
                            {/* <a href="#testimonials" className="text-gray-600 hover:text-gray-800">Testimonials</a> */}
                            <a href="#pricing" className="text-gray-600 hover:text-gray-800">Pricing</a>
                            <a href="#faq" className="text-gray-600 hover:text-gray-800">FAQ</a>
                        </nav>
                        <div className="hidden md:flex items-center">
                            {!isAuthenticated ? (
                                <>
                                    <Link to="/auth/login" className="text-gray-600 hover:text-gray-800 px-3 py-2">Login</Link>
                                    <Link to="/auth/sign-up" className="bg-blue-900 hover:bg-blue-700 text-white px-4 py-2 rounded-md ml-4 transition duration-300">Sign up</Link>
                                </>
                            ) : (
                                <Link to="/user/dash" className="bg-blue-900 hover:bg-blue-700 text-white px-4 py-2 rounded-md ml-4 transition duration-300">My Dashboard</Link>
                            )}
                        </div>
                        <button className="md:hidden" onClick={() => setIsMenuOpen(!isMenuOpen)}>
                            {isMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
                        </button>
                    </div>
                    {isMenuOpen && (
                        <div className="md:hidden bg-white py-4 px-4">
                            <a href="#features" className="block py-2 text-gray-600 hover:text-gray-800">Features</a>
                            <a href="#testimonials" className="block py-2 text-gray-600 hover:text-gray-800">Testimonials</a>
                            <a href="#pricing" className="block py-2 text-gray-600 hover:text-gray-800">Pricing</a>
                            <a href="#faq" className="block py-2 text-gray-600 hover:text-gray-800">FAQ</a>
                            <Link to="/auth/login" className="block py-2 text-gray-600 hover:text-gray-800">Login</Link>
                            <Link to="/auth/sign-up" className="block py-2 text-blue-600 hover:text-blue-700">Sign up</Link>
                        </div>
                    )}
                </header>

                <main className="flex-grow ">
                    {/* Hero Section */}
                    <section className="bg-blue-900 p-4 text-white space-x-5 py-20">
                        <div className="container mx-auto px-4 flex flex-col md:flex-row items-center">
                            <div className="md:w-1/2 mb-10 md:mb-0">
                                <h1 className="text-4xl md:text-5xl font-bold mb-6">Supercharge Your Email Marketing</h1>
                                <p className="text-xl">Reach your audience, drive conversions, and grow</p>
                                <p className="text-xl mb-8"> your business with our powerful email marketing platform.</p>
                                <Link to="/auth/sign-up" className="bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3 rounded-md inline-flex items-center transition duration-300">
                                    Start Free Trial
                                    <ArrowRight className="ml-2 w-5 h-5" />
                                </Link>
                            </div>
                            <div className="md:w-1/2">
                                <img src={mailpicture} alt="Email marketing illustration" className="rounded-lg w-[80%] h-[70%] shadow-xl" />
                            </div>
                        </div>
                    </section>

                    {/* Features Section */}
                    <section id="features" className="py-20  bg-white">
                        <div className="container mx-auto  px-8">
                            <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">Why Choose CrabMailer?</h2>
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
                                {features.map((feature, index) => (
                                    <div key={index} className="bg-gray-50 p-6 rounded-lg shadow-md hover:shadow-lg transition duration-300">
                                        <div className="text-blue-600 mb-4">{feature.icon}</div>
                                        <h3 className="text-xl font-semibold mb-2">{feature.title}</h3>
                                        <p className="text-gray-600">{feature.description}</p>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </section>

                    {/* Testimonials Section */}
                    {/* <section id="testimonials" className="py-20 bg-gray-50">
                    <div className="container mx-auto px-8">
                        <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">What Our Customers Say</h2>
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                            {[1, 2,3].map((_, index) => (
                                <div key={index} className="bg-white p-6 rounded-lg shadow-md">
                                    <p className="text-gray-600 mb-4">"CrabMailer has revolutionized our email marketing strategy. The results speak for themselves!"</p>
                                    <div className="flex items-center">
                                        <img src={`/api/placeholder/50/50`} alt={`Customer ${index + 1}`} className="w-12 h-12 rounded-full mr-4" />
                                        <div>
                                            <p className="font-semibold">John Doe</p>
                                            <p className="text-sm text-gray-500">CEO, Tech Company</p>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </section> */}

                    {/* Pricing Section */}
                    <section id="pricing" className="py-20 bg-white">
                        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
                            <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">Simple, Transparent Pricing</h2>
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                                {plans.map((plan, index) => (
                                    <div key={index} className="bg-gray-50 p-8 rounded-lg shadow-md hover:shadow-lg transition duration-300">
                                        <h3 className="text-2xl font-bold mb-4">{plan.name}</h3>
                                        <p className="text-4xl font-bold mb-6">
                                            &#8358;{plan.price}<span className="text-lg font-normal text-gray-500">/mo</span>
                                        </p>
                                        <ul className="mb-8 space-y-2">
                                            {plan.features.map((feature, featureIndex) => (
                                                <li key={featureIndex} className="flex items-center">
                                                    <ArrowRight className="w-4 h-4 mr-2 text-green-500" />
                                                    {feature}
                                                </li>
                                            ))}
                                        </ul>
                                        <button className="w-full bg-blue-900 hover:bg-blue-700 text-white py-2 px-4 rounded-md transition duration-300">
                                            Choose Plan
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </section>

                    {/* FAQ Section */}
                    <section id="faq" className="py-20 bg-gray-50">
                        <div className="container mx-auto px-8">
                            <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">Frequently Asked Questions</h2>
                            <div className="max-w-3xl mx-auto">
                                {faqs.map((faq, index) => (
                                    <div key={index} className="mb-4">
                                        <button
                                            className="flex justify-between items-center w-full p-4 bg-white hover:bg-gray-100 rounded-lg focus:outline-none"
                                            onClick={() => toggleAccordion(index)}
                                        >
                                            <span className="font-semibold">{faq.question}</span>
                                            <ChevronDown className={`w-5 h-5 transition-transform ${activeAccordion === index ? 'transform rotate-180' : ''}`} />
                                        </button>
                                        {activeAccordion === index && (
                                            <div className="p-4 bg-gray-50">
                                                <p className="text-gray-600">{faq.answer}</p>
                                            </div>
                                        )}
                                    </div>
                                ))}
                            </div>
                        </div>
                    </section>

                    {/* CTA Section */}
                    <section className="bg-blue-900 py-20 text-white">
                        <div className="container mx-auto px-4 text-center">
                            <h2 className="text-3xl font-bold mb-4">Ready to supercharge your email marketing?</h2>
                            <p className="text-xl mb-8">Join thousands of businesses that trust CrabMailer for their email marketing needs.</p>
                            <Link to="/auth/sign-up" className="bg-white text-blue-600 hover:bg-gray-100 text-lg px-8 py-3 rounded-md inline-flex items-center transition duration-300">
                                Start Your Free Trial
                                <ArrowRight className="ml-2 w-5 h-5" />
                            </Link>
                        </div>
                    </section>
                </main>

                <footer className="bg-gray-800 text-white py-12">
                    <div className="container mx-auto px-4">
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
                            <div>
                                <h3 className="text-lg font-semibold mb-4">Product</h3>
                                <ul className="space-y-2">
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">Features</a></li>
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">Pricing</a></li>
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">Integrations</a></li>
                                </ul>
                            </div>
                            <div>
                                <h3 className="text-lg font-semibold mb-4">Resources</h3>
                                <ul className="space-y-2">
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">Blog</a></li>
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">Help Center</a></li>
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">Guides</a></li>
                                </ul>
                            </div>
                            <div>
                                <h3 className="text-lg font-semibold mb-4">Company</h3>
                                <ul className="space-y-2">
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">About Us</a></li>
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">Careers</a></li>
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">Contact</a></li>
                                </ul>
                            </div>
                            <div>
                                <h3 className="text-lg font-semibold mb-4">Legal</h3>
                                <ul className="space-y-2">
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">Privacy Policy</a></li>
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">Terms of Service</a></li>
                                    <li><a href="#" className="hover:text-gray-300 transition duration-300">GDPR</a></li>
                                </ul>
                            </div>
                        </div>
                        <div className="mt-12 pt-8 border-t border-gray-700 text-center">
                            <p>&copy; {new Date().getFullYear()} CrabMailer. All rights reserved.</p>
                        </div>
                    </div>
                </footer>
            </div>
        </HelmetProvider>
    );
};

export default IndexLandingPage;