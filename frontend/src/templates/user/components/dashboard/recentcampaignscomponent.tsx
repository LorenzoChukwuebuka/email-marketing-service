import React from 'react';
import EmptyState from '../../../../components/emptyStateComponent';

const RecentCampaigns = () => {
    return (

        <>
            <div className="font-sans p-6">
                <h2 className="text-2xl font-bold mb-4">Recent campaigns</h2>
                <EmptyState className='shadow-md'
                    title="You have not created any Template"
                    description="Create and easily send marketing emails to your audience"
                    icon={<i className="bi bi-emoji-frown text-xl"></i>}
                    buttonText="Create Template"
                  //  onButtonClick={() => navigate("/user/dash/marketing")}
                />
            </div>
        </>
    );
};

export default RecentCampaigns;