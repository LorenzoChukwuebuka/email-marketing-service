import { useEffect, useState } from "react";
import useCampaignStore from "../../../store/userstore/campaignStore";
import AnalyticsTableComponent from "./analyticsTabletComponent";

type StatProp = { value: string; label: string }

const StatItem = ({ value, label }: StatProp) => (
    <div className="bg-white p-4 rounded-lg shadow-sm flex flex-col items-center justify-center">
        <span className="text-3xl font-bold text-gray-800">{value}</span>
        <span className="text-sm text-gray-500 mt-2">{label}</span>
    </div>
);




const AnalyticsTemplateDash: React.FC = () => {

    const [isLoading, setIsLoading] = useState<boolean>(false)
    const { getCampaignUserStats, campaignUserStatsData } = useCampaignStore()

    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true)
            await getCampaignUserStats()
            await new Promise((resolve) => setTimeout(resolve, 1000))
            setIsLoading(false)
        }

        fetchData()
    }, [getCampaignUserStats])

    const stats = [
        { value: `${campaignUserStatsData.total_emails_sent}`, label: 'Total Emails Sent' },
        { value: `${campaignUserStatsData.total_deliveries}`, label: 'Total Delivered' },
        { value: `${campaignUserStatsData.total_bounces}`, label: 'Total Bounce' },
        { value: `0`, label: 'Total Complaints' },
        { value: `0`, label: 'Total Rejected' },
        { value: `${campaignUserStatsData.total_opens}`, label: 'Total Opens' },
        { value: `${campaignUserStatsData.unique_opens}`, label: 'Total Unique Opens' },
        { value: `${campaignUserStatsData.open_rate}%`, label: 'Total Open Rate' },
        { value: `${campaignUserStatsData.total_clicks}`, label: 'Total Clicks' },
        { value: `${campaignUserStatsData.unique_clicks}`, label: 'Total Unique Clicks' },
    ];


    return <>
        <div className="bg-gray-100 mt-10 mb-5 p-6">

            <h1 className="font-semibold text-lg   mb-4"> Analytics </h1>

            {isLoading ? <div className="flex items-center justify-center mt-20"><span className="loading loading-spinner loading-lg"></span></div> : (
                <>
                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                        {stats.map((stat, index) => (
                            <StatItem key={index} value={stat.value} label={stat.label} />
                        ))}
                    </div>

                    <AnalyticsTableComponent />
                </>
            )}


        </div>
    </>
}

export default AnalyticsTemplateDash