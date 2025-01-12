import { useParams } from 'react-router-dom';
import useCampaignStore, { Campaign, CampaignData } from '../../../../store/userstore/campaignStore';
import { useEffect, useState } from 'react';
import CampaignInfo from './campaignInfoComponent';
import { BaseEntity } from '../../../../interface/baseentity.interface';
import CampaignRecipientComponent from './campaignRecipientComponent';

type StatProp = { value: string; label: string }

const StatItem = ({ value, label }: StatProp) => (
    <div className="bg-white p-4 rounded-lg shadow-sm flex flex-col items-center justify-center">
        <span className="text-3xl font-bold text-gray-800">{value}</span>
        <span className="text-sm text-gray-500 mt-2">{label}</span>
    </div>
);

const CampaignReport: React.FC = () => {

    const { getSingleCampaign, campaignData, getCampaignStats, campaignStatData } = useCampaignStore();
    const { id } = useParams<{ id: string }>() as { id: string };
    const [isLoading, setIsLoading] = useState<boolean>(false)

    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true)
            await getSingleCampaign(id);
            await getCampaignStats(id)
            await new Promise((resolve) => setTimeout(resolve, 500))
            setIsLoading(false)
        };
        fetchData();
    }, [id, getSingleCampaign, getCampaignStats]);


    const stats = [
        { value: `${campaignStatData.total_emails_sent ?? 0} `, label: 'Emails Sent' },
        { value: `${campaignStatData.total_deliveries ?? 0}`, label: 'Delivered' },
        { value: `${campaignStatData.total_bounces ?? 0}`, label: 'Bounce' },
        { value: `0`, label: 'Complaints' },
        { value: `${campaignStatData.hard_bounces ?? 0}`, label: 'Rejected' },
        { value: `${campaignStatData.total_opens ?? 0}`, label: 'Opens' },
        { value: `${campaignStatData.unique_opens ?? 0}`, label: 'Unique Opens' },
        { value: `${campaignStatData.open_rate ?? 0}%`, label: 'Open rate' },
        { value: `${campaignStatData.total_clicks ?? 0}`, label: 'Total Clicks' },
        { value: `${campaignStatData.unique_clicks ?? 0}`, label: 'Unique Clicks' },
    ];

    return (
        <>
            <div className="bg-gray-100 mt-10 mb-5 p-6">

                {isLoading ? <div className="flex items-center justify-center mt-20"><span className="loading loading-spinner loading-lg"></span></div> : (
                    <>
                        <div className="flex items-center mb-5">
                            <button className="text-blue-600 mr-2" onClick={() => window.history.back()}>
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                                </svg>
                            </button>

                            <h1 className='mt-2 mb-4 font-medium text-lg'> {(campaignData as (BaseEntity & Campaign))?.name} </h1>

                        </div>
                        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                            {stats.map((stat, index) => (
                                <StatItem key={index} value={stat.value} label={stat.label} />
                            ))}
                        </div>

                        <CampaignInfo campaignData={campaignData as (BaseEntity & Campaign)} />

                        <CampaignRecipientComponent campaignId={(campaignData as (BaseEntity & CampaignData))?.uuid as string} />
                    </>
                )}
            </div>
        </>
    )
}


export default CampaignReport