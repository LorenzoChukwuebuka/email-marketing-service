import { Helmet, HelmetProvider } from "react-helmet-async";
import useAdminUserCamapaignStore from "../../../../store/admin/AdminUserCampaignsStore";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

import useCampaignStore, { CampaignData } from "../../../../store/userstore/campaignStore";
import { ContactGroupData } from "../../../../store/userstore/contactGroupStore";
import CampaignReport from "../../../user/components/campaign/campaignReportComponent";
import CampaignInfo from "../../../user/components/campaign/campaignInfoComponent";
import CampaignRecipientComponent from "../../../user/components/campaign/campaignRecipientComponent";


type StatProp = { value: string; label: string }

const StatItem = ({ value, label }: StatProp) => (
    <div className="bg-white p-4 rounded-lg shadow-sm flex flex-col items-center justify-center">
        <span className="text-3xl font-bold text-gray-800">{value}</span>
        <span className="text-sm text-gray-500 mt-2">{label}</span>
    </div>
);

const AdminUserSpecificCampaigns: React.FC = () => {
    const { campaignData, getSingleUserCampaign } = useAdminUserCamapaignStore();
    const { getCampaignStats, campaignStatData } = useCampaignStore()
    const { campaignid } = useParams<{ campaignid: string }>() as { campaignid: string };

    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [campData, setCampData] = useState<CampaignData | null>(null);

    useEffect(() => {
        const fetchCampaignData = async () => {
            setIsLoading(true);
            await getSingleUserCampaign(campaignid);
            await getCampaignStats(campaignid);
            setIsLoading(false);
        };
        fetchCampaignData();
    }, [getSingleUserCampaign, getCampaignStats, campaignid]);

    useEffect(() => {
        if (campaignData) {
            setCampData(campaignData as CampaignData);
        }
    }, [campaignData]);

    if (isLoading) {
        return <div className="flex items-center justify-center mt-20"><span className="loading loading-spinner loading-lg"></span></div>;
    }

    const stats = [
        { value: `${campaignStatData.total_emails_sent ?? 0} `, label: "Emails Sent" },
        { value: `${campaignStatData.total_deliveries ?? 0}`, label: "Delivered" },
        { value: `${campaignStatData.total_bounces ?? 0}`, label: "Bounce" },
        { value: `0`, label: "Complaints" },
        { value: `${campaignStatData.hard_bounces ?? 0}`, label: "Rejected" },
        { value: `${campaignStatData.total_opens ?? 0}`, label: "Opens" },
        { value: `${campaignStatData.unique_opens ?? 0}`, label: "Unique Opens" },
        { value: `${campaignStatData.open_rate ?? 0}%`, label: "Open rate" },
        { value: `${campaignStatData.total_clicks ?? 0}`, label: "Total Clicks" },
        { value: `${campaignStatData.unique_clicks ?? 0}`, label: "Unique Clicks" },
    ];

    return (
        <HelmetProvider>
            <Helmet title={ "Campaign - " + campData?.name} />
            {campData && (
                <>

                    <div className="bg-gray-100 mt-10 mb-5 p-6">
                        <div className="flex items-center mb-5">
                            <button
                                className="text-blue-600 mr-2"
                                onClick={() => window.history.back()}
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    className="h-6 w-6"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M10 19l-7-7m0 0l7-7m-7 7h18"
                                    />
                                </svg>


                            </button>
                            <h1 className="mt-2 mb-4 font-medium text-lg">
                                {(campData as CampaignData)?.name}
                            </h1>
                        </div>
                        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                            {stats.map((stat, index) => (
                                <StatItem key={index} value={stat.value} label={stat.label} />
                            ))}
                        </div>
                    </div>


                    Campaign Info
                    <CampaignInfo campaignData={campData} />

                    {/* Campaign Recipients */}
                    <CampaignRecipientComponent campaignId={campData.uuid} />
                </>
            )}
        </HelmetProvider>
    )
};

export default AdminUserSpecificCampaigns;
