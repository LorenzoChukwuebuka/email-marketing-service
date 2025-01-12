import { useMemo } from "react";
import useMetadata from "../../../hooks/useMetaData";
import { HelmetProvider, Helmet } from "react-helmet-async";
import AnalyticsTableComponent from "../components/analyticsTableComponent";
import { useCampaignUserStatsQuery } from "../hooks/useAnalyticsQuery";

type StatProp = { value: string; label: string }

const StatItem = ({ value, label }: StatProp) => (
    <div className="bg-white p-4 rounded-lg shadow-sm flex flex-col items-center justify-center">
        <span className="text-3xl font-bold text-gray-800">{value}</span>
        <span className="text-sm text-gray-500 mt-2">{label}</span>
    </div>
);


const AnalyticsTemplateDash: React.FC = () => {
    const { data: campaignUserStatsData, isLoading } = useCampaignUserStatsQuery()
    const cusdata = useMemo(() => campaignUserStatsData?.payload, [campaignUserStatsData])
    const stats = [
        { value: `${cusdata?.total_emails_sent ?? 0}`, label: 'Total Emails Sent' },
        { value: `${cusdata?.total_deliveries ?? 0}`, label: 'Total Delivered' },
        { value: `${cusdata?.total_bounces ?? 0}`, label: 'Total Bounce' },
        // { value: `0`, label: 'Total Complaints' },
        // { value: `0`, label: 'Total Rejected' },
        { value: `${cusdata?.total_opens ?? 0}`, label: 'Total Opens' },
        { value: `${cusdata?.unique_opens ?? 0}`, label: 'Total Unique Opens' },
        { value: `${cusdata?.open_rate ?? 0}%`, label: 'Total Open Rate' },
        { value: `${cusdata?.total_clicks ?? 0}`, label: 'Total Clicks' },
        { value: `${cusdata?.unique_clicks ?? 0}`, label: 'Total Unique Clicks' },
    ];

    const metaData = useMetadata("Dashboard")

    return <>
        <HelmetProvider>
            <Helmet {...metaData} title="Analytics" />
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
        </HelmetProvider>
    </>
}

export default AnalyticsTemplateDash