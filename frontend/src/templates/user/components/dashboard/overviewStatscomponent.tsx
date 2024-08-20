import React, { useState } from 'react';
import { CalendarIcon, XIcon } from 'lucide-react';


const OverviewStats = () => {
    const [dateRange, setDateRange] = useState('01 Aug 2024 ~ 09 Aug 2024');

    return (
        <div className="bg-gray-100 p-6 rounded-lg">
            <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-bold text-gray-800">Overview Stats</h2>
                <div className="flex items-center bg-white rounded-md shadow px-3 py-2">
                    <CalendarIcon className="text-gray-400 mr-2" size={16} />
                    <input
                        type="text"
                        value={dateRange}
                        onChange={(e) => setDateRange(e.target.value)}
                        className="text-sm text-gray-600 focus:outline-none"
                    />
                    <XIcon className="text-gray-400 ml-2 cursor-pointer" size={16} onClick={() => setDateRange('')} />
                </div>
            </div>

            <div>
                <h3 className="text-sm font-medium text-gray-500 mb-4">Subscription & Audience</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                    <StatCard title="Total active subscribers" value="0" />
                    <StatCard title="New subscribers" value="0" />
                    <StatCard title="Unsubscribed" value="0" />
                    <StatCard title="Engaged Subscribers" value="0" />
                </div>
            </div>
        </div>
    );
};


type StatProps = {
    title: string
    value: string | number | boolean | null
}

const StatCard = ({ title, value }: StatProps) => (
    <div className="bg-white p-4 rounded-lg shadow transition-transform transform hover:scale-105 hover:shadow-lg hover:bg-gray-50">
        <p className="text-3xl font-bold text-gray-800">{value}</p>
        <h3 className="text-sm text-gray-500 mb-2">{title}</h3>
    </div>
);


export default OverviewStats;