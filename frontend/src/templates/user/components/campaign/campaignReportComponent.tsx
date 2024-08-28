import { useParams } from 'react-router-dom';
import useCampaignStore, { Campaign, CampaignData } from '../../../../store/userstore/campaignStore';
import { useEffect } from 'react';
import CampaignInfo from './campaignInfoComponent';
import { BaseEntity } from '../../../../interface/baseentity.interface';

type StatProp = { value: string; label: string }

const StatItem = ({ value, label }: StatProp) => (
    <div className="bg-white p-4 rounded-lg shadow-sm flex flex-col items-center justify-center">
        <span className="text-3xl font-bold text-gray-800">{value}</span>
        <span className="text-sm text-gray-500 mt-2">{label}</span>
    </div>
);


const CampaignReport: React.FC = () => {

    const { getSingleCampaign, campaignData } = useCampaignStore();
    const { id } = useParams<{ id: string }>() as { id: string };

    useEffect(() => {
        const fetchData = async () => {
            await getSingleCampaign(id);
        };
        fetchData();
    }, [id, getSingleCampaign]);


    const stats = [
        { value: '2', label: 'Emails Sent' },
        { value: '0', label: 'Delivered' },
        { value: '0', label: 'Bounce' },
        { value: '0', label: 'Complaints' },
        { value: '0', label: 'Rejected' },
        { value: '0', label: 'Opens' },
        { value: '0', label: 'Unique Opens' },
        { value: '0%', label: 'Open rate' },
        { value: '0', label: 'Total Clicks' },
        { value: '0', label: 'Unique Clicks' },
    ];

    return (
        <>
            <div className="bg-gray-100 mt-10 mb-5 p-6">
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                    {stats.map((stat, index) => (
                        <StatItem key={index} value={stat.value} label={stat.label} />
                    ))}
                </div>

                <CampaignInfo campaignData={campaignData as (BaseEntity & Campaign)} />

            </div>


        </>

    )
}







export default CampaignReport