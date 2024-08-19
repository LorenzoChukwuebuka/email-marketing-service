import React from 'react';

const RecentCampaigns = () => {
    return (
        <div className="font-sans p-6">
            <h2 className="text-2xl font-bold mb-4">Recent campaigns</h2>
            <div className="bg-white rounded-lg shadow-md p-8">
                <div className="text-center">
                    <h3 className="text-xl font-semibold mb-2">Launch your first campaign</h3>
                    <p className="text-gray-600 mb-6">
                        You can boost your business by messaging a wide audience.
                    </p>
                    <button className="bg-gray-900 text-white font-semibold py-2 px-4 rounded-full hover:bg-gray-800 transition duration-300">
                        Create a campaign
                    </button>
                </div>
            </div>
        </div>
    );
};

export default RecentCampaigns;