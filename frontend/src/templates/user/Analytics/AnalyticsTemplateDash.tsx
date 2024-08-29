
type StatProp = { value: string; label: string }

const StatItem = ({ value, label }: StatProp) => (
    <div className="bg-white p-4 rounded-lg shadow-sm flex flex-col items-center justify-center">
        <span className="text-3xl font-bold text-gray-800">{value}</span>
        <span className="text-sm text-gray-500 mt-2">{label}</span>
    </div>
);




const AnalyticsTemplateDash: React.FC = () => {

    const stats = [
        { value: `0`, label: 'Total Emails Sent' },
        { value: `0`, label: 'Total Delivered' },
        { value: `0`, label: 'Total Bounce' },
        { value: `0`, label: 'Total Complaints' },
        { value: `0`, label: 'Total Rejected' },
        { value: `0`, label: 'Total Opens' },
        { value: `0`, label: 'Total Unique Opens' },
        { value: `0%`, label: 'Total Open Rate' },
        { value: `0`, label: 'Total Clicks' },
        { value: `0`, label: 'Total Unique Clicks' },
    ];


    return <>
        <div className="bg-gray-100 mt-10 mb-5 p-6">
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                {stats.map((stat, index) => (
                    <StatItem key={index} value={stat.value} label={stat.label} />
                ))}
            </div>
        </div>
    </>
}

export default AnalyticsTemplateDash