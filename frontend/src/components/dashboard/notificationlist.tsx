import { Helmet, HelmetProvider } from 'react-helmet-async';
import { useNotificationQuery } from '../../hooks/useNotificationQuery';

const NotificationList = () => {
    const { data: notificationsData } = useNotificationQuery()
    const ndata = notificationsData?.payload || []

    // Group notifications by date
    const groupedNotifications = ndata?.reduce((acc, notification) => {
        const date = new Date(notification.created_at).toLocaleDateString('en-GB', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric'
        });
        if (!acc[date]) {
            acc[date] = [];
        }
        acc[date].push(notification);
        return acc;
    }, {});


    return (
        <HelmetProvider>
            <Helmet title="Notifications - CrabMailer" />
            <div className="p-6">
                <h2 className="text-2xl font-bold mb-4">Notifications</h2>

                {Object.entries(groupedNotifications || {})?.map(([date, notifications]) => (
                    <div key={date} className="mb-6">
                        {/* Date Header */}
                        <h3 className="text-lg font-semibold mb-2">{date}</h3>

                        {/* Notification Items */}
                          {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
                        {(notifications as any).map((notification) => (  // Changed from ndata to notifications
                            <div
                                key={notification.id}
                                className="flex items-start mb-4 p-4 bg-white shadow-sm rounded-md"
                            >
                                {/* Time */}
                                <div className="w-20 text-sm font-medium text-gray-600">
                                    {new Date(notification.created_at).toLocaleTimeString('en-GB', {
                                        hour: '2-digit',
                                        minute: '2-digit'
                                    })}
                                </div>

                                {/* Notification Message */}
                                <div className="ml-4 text-gray-800">
                                    <p className="font-medium">{notification.title}</p>
                                </div>
                            </div>
                        ))}
                    </div>
                ))}
            </div>
        </HelmetProvider>
    );
};

export default NotificationList;
