import React from 'react';
import { Mail, Book, HelpCircle, PlayCircle } from 'lucide-react';

const HelpAndSupport:React.FC = () => {
    return (
        <div className="max-w-4xl mx-auto p-6 mt-10">
            <h1 className="text-3xl font-bold text-center text-indigo-800 mb-8">Help & Support</h1>

            <p className="text-center text-gray-700 mb-8">
                You can reach us at <a href="" className="text-indigo-600 hover:underline"> hello@hello.com </a>.
                We'll get back to you as soon as we can, typically within a few hours.
            </p>

            <div className="grid md:grid-cols-2 gap-6">
                <SupportCard
                    icon={<HelpCircle className="w-8 h-8 text-indigo-600" />}
                    title="FAQs"
                    description="Read our frequently asked questions here, this is a quick starting point to answers common questions"
                />

                <SupportCard
                    icon={<Book className="w-8 h-8 text-indigo-600" />}
                    title="Help Articles"
                    description="Easy short advice, answers & best practices from the crabmailer team & contributors"
                />

                <SupportCard
                    icon={<PlayCircle className="w-8 h-8 text-indigo-600" />}
                    title="Video Tutorials"
                    description="Watch step-by-step guides on how to use crabmailer effectively"
                />

                <SupportCard
                    icon={<Mail className="w-8 h-8 text-indigo-600" />}
                    title="Contact Support"
                    description="Get in touch with our support team for personalized assistance"
                />
            </div>
        </div>
    );
};

type Props = { icon: any; title: string; description: string }

const SupportCard = ({ icon, title, description }: Props) => {
    return (
        <div className="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300">
            <div className="flex items-center mb-4">
                {icon}
                <h2 className="text-xl font-semibold text-gray-800 ml-3">{title}</h2>
            </div>
            <p className="text-gray-600">{description}</p>
        </div>
    );
};

export default HelpAndSupport;