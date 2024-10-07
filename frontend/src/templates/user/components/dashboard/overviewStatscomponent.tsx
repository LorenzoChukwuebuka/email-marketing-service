import { useEffect } from "react";
import useContactStore from "../../../../store/userstore/contactStore";

const OverviewStats: React.FC = () => {
    const { getContactSubEngagement, engagementCount } = useContactStore()

    useEffect(() => {
        const fetchEngagement = async () => {
            await getContactSubEngagement()
        }

        fetchEngagement()
    }, [getContactSubEngagement])
    return (
        <div className="bg-gray-100 p-6 rounded-lg">
            <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-bold text-gray-800">Overview Stats</h2>

            </div>

            <div>
                <h3 className="text-sm font-medium text-gray-500 mb-4">Subscription & Audience</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                    <StatCard title="Total active subscribers" value={engagementCount.total} />
                    <StatCard title="New subscribers" value={engagementCount.new} />
                    <StatCard title="Unsubscribed" value={engagementCount.unsubscribed} />
                    <StatCard title="Engaged Subscribers" value={engagementCount.engaged} />
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
    <div className="bg-white p-4 rounded-lg shadow transition-transform transform hover:scale-105 hover:shadow-md hover:rounded-md hover:bg-gray-50">
        <p className="text-3xl font-bold text-gray-800">{value}</p>
        <h3 className="text-sm text-gray-500 mb-2">{title}</h3>
    </div>
);


export default OverviewStats;