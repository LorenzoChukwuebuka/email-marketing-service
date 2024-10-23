import React, { useEffect, useState } from 'react';
import { Link } from "react-router-dom";
import { ArrowRight, Mail, Zap, Shield, Users, ChevronDown, Menu, X } from 'lucide-react';
import mailpicture from "./../../assets/dashboard.jpg"
import { Helmet, HelmetProvider } from 'react-helmet-async';
import useMetadata from '../../hooks/useMetaData';
import Footer from './footer';
import NavBar from './navbar';
import { Slider } from '@mui/material';

const IndexLandingPage: React.FC = () => {

    const [activeAccordion, setActiveAccordion] = useState(null);

    const features = [
        { icon: <Mail className="w-8 h-8" />, title: 'Powerful Email Campaigns', description: 'Create and send beautiful, responsive emails that drive results.' },
        { icon: <Zap className="w-8 h-8" />, title: 'Smart Automation', description: 'Set up triggers and workflows to nurture leads and customers automatically.' },
        { icon: <Shield className="w-8 h-8" />, title: 'Top-tier Deliverability', description: 'Our advanced infrastructure ensures your emails reach the inbox.' },
        { icon: <Users className="w-8 h-8" />, title: 'Precise Audience Segmentation', description: 'Target the right subscribers with personalized content for maximum engagement.' },
    ];

    const faqs = [
        { question: 'How does the free trial work?', answer: 'Our free plan lets you send up to 200 emails per day (3,000 monthly). You`ll get full access to all features, and the best part? No credit card required to get started!' },
        { question: 'Can I integrate with my existing tools?', answer: 'Yes, we offer integrations with popular CRMs, e-commerce platforms, and more.' },
        { question: 'What kind of support do you offer?', answer: '24/7 email support for all plans, with phone and chat support for higher tiers.' },
        { question: 'Is there a limit on subscribers or emails?', answer: 'Plans vary, but we have options for businesses of all sizes, from startups to enterprises.' },
    ];

    const toggleAccordion = (index: any) => {
        setActiveAccordion(activeAccordion === index ? null : index);
    };

    const [plan, setPlan] = useState<number>(1); // 1 = Basic, 2 = Pro, 3 = Enterprise
    const [subscriberCount, setSubscriberCount] = useState<number>(1000); // Default number of subscribers

    const pricingData = {
        basic: { basePrice: 7000, perThousandSubscribers: 200 },
        pro: { basePrice: 15000, perThousandSubscribers: 300 },
        enterprise: { basePrice: 30000, perThousandSubscribers: 500 },
    };

    const plans = [
        { name: 'Basic', value: 1, basePrice: pricingData.basic.basePrice, description: 'Basic plan description.' },
        { name: 'Pro', value: 2, basePrice: pricingData.pro.basePrice, description: 'Pro plan description.' },
        { name: 'Enterprise', value: 3, basePrice: pricingData.enterprise.basePrice, description: 'Enterprise plan description.' }
    ];

    const calculatePrice = (basePrice: number, subscriberCount: number, perThousandRate: number) => {
        const extraCost = Math.ceil(subscriberCount / 1000) * perThousandRate;
        return basePrice + extraCost;
    };

    const handleSliderChange = (event: any, newValue: number | number[]) => {
        setSubscriberCount(newValue as number);
    };


    const metaData = useMetadata()("LandingPage")

    return (

        <HelmetProvider>
            <Helmet {...metaData} />
            <div className="min-h-screen landing-page  flex flex-col bg-gray-50">
                <NavBar />
                <main className="flex-grow ">
                    {/* Hero Section */}
                    <section className="bg-blue-900 p-4 text-white space-x-5 py-20">
                        <div className="container mx-auto px-4 flex flex-col md:flex-row items-center">
                            <div className="md:w-1/2 mb-10 md:mb-0">
                                <h1 className="text-4xl md:text-7xl font-semibold mb-6">Supercharge Your Email Strategy</h1>
                                <p className="text-xl mb-4">Elevate your customer communications, drive conversions, and grow your business with our comprehensive email platform.</p>
                                <ul className="text-lg mb-8 space-y-2">
                                    <li className="flex items-start">
                                        <ArrowRight className="mr-2 w-5 h-5 text-blue-600 flex-shrink-0 mt-1" />
                                        <span><strong>Marketing Campaigns:</strong> Reach and engage your audience effectively</span>
                                    </li>
                                    <li className="flex items-start">
                                        <ArrowRight className="mr-2 w-5 h-5 text-blue-600 flex-shrink-0 mt-1" />
                                        <span><strong>Transactional Emails:</strong> Deliver timely, personalized notifications</span>
                                    </li>
                                    <li className="flex items-start">
                                        <ArrowRight className="mr-2 w-5 h-5 text-blue-600 flex-shrink-0 mt-1" />
                                        <span><strong>Automated Workflows:</strong> Streamline your email campaigns for maximum impact</span>
                                    </li>
                                </ul>
                                <p className="text-xl font-semibold mb-8">Power your entire email ecosystem with one robust solution.</p>
                                <Link to="/auth/sign-up" className="bg-blue-600 text-white hover:bg-blue-700 text-lg px-8 py-3 rounded-md inline-flex items-center transition duration-300">
                                    Get Started
                                    <ArrowRight className="ml-2 w-5 h-5" />
                                </Link>
                            </div>

                            <div className="md:w-1/2">
                                <img src={mailpicture} alt="Email marketing illustration" className="rounded-lg w-[100%] h-[100%] shadow-xl" />
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

                    <main className="flex-grow">
                        {/* Pricing Section */}
                        <section id="pricing" className="py-20 bg-white">
                            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
                                <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">Simple, Transparent Pricing</h2>

                                {/* Plan Selector */}
                                <div className="text-center mb-8">
                                    <div className="inline-flex rounded-md shadow-sm" role="group">
                                        {plans.map((p) => (
                                            <button
                                                key={p.value}
                                                className={`px-4 py-2 border text-lg font-medium ${plan === p.value ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-900'
                                                    } focus:outline-none`}
                                                onClick={() => setPlan(p.value)}
                                            >
                                                {p.name}
                                            </button>
                                        ))}
                                    </div>
                                </div>

                                {/* Subscriber Count Slider */}
                                <div className="text-center mb-12">
                                    <h3 className="text-xl font-semibold mb-4">Select Subscriber Count</h3>
                                    <Slider
                                        value={subscriberCount}
                                        onChange={handleSliderChange}
                                        aria-labelledby="subscriber-slider"
                                        min={1000}
                                        max={100000}
                                        step={1000}
                                        valueLabelDisplay="auto"
                                        marks={[
                                            { value: 1000, label: '1k' },
                                            { value: 50000, label: '50k' },
                                            { value: 100000, label: '100k' },
                                        ]}
                                    />
                                    <p className="text-lg font-semibold mt-4">Subscribers: {subscriberCount.toLocaleString()}</p>
                                </div>

                                {/* Display Price */}
                                <div className="text-center">
                                    {plan === 1 && (
                                        <p className="text-2xl font-bold">
                                            Price: &#8358;{calculatePrice(pricingData.basic.basePrice, subscriberCount, pricingData.basic.perThousandSubscribers).toLocaleString()}
                                        </p>
                                    )}
                                    {plan === 2 && (
                                        <p className="text-2xl font-bold">
                                            Price: &#8358;{calculatePrice(pricingData.pro.basePrice, subscriberCount, pricingData.pro.perThousandSubscribers).toLocaleString()}
                                        </p>
                                    )}
                                    {plan === 3 && (
                                        <p className="text-2xl font-bold">
                                            Price: &#8358;{calculatePrice(pricingData.enterprise.basePrice, subscriberCount, pricingData.enterprise.perThousandSubscribers).toLocaleString()}
                                        </p>
                                    )}
                                </div>

                                <div className="text-center mt-8">
                                    <Link to="/auth/sign-up" className="bg-blue-600 text-white hover:bg-blue-700 text-lg px-8 py-3 rounded-md inline-flex items-center transition duration-300">
                                        Choose Plan
                                        <ArrowRight className="ml-2 w-5 h-5" />
                                    </Link>
                                </div>
                            </div>
                        </section>
                    </main>

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
                                Sign Up
                                <ArrowRight className="ml-2 w-5 h-5" />
                            </Link>
                        </div>
                    </section>
                </main>

                <Footer />
            </div>
        </HelmetProvider>
    );
};

export default IndexLandingPage;