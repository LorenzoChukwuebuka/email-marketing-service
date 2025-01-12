import React, { useEffect } from 'react';
import useUserNotificationStore, { UserNotification } from '../../../store/userstore/notifications.store';
import { Helmet, HelmetProvider } from 'react-helmet-async';

const NotificationList = () => {
    const { notificationsData, getUserNotifications } = useUserNotificationStore();

    // Group notifications by date
    const groupedNotifications = notificationsData.reduce((acc: any, notification) => {
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

    useEffect(() => {
        getUserNotifications();
    }, []);

    return (
        <HelmetProvider>
            <Helmet title="Notifications - CrabMailer" />
            <div className="p-6">
                <h2 className="text-2xl font-bold mb-4">Notifications</h2>

                {Object.entries(groupedNotifications).map(([date, notifications]) => (
                    <div key={date} className="mb-6">
                        {/* Date Header */}
                        <h3 className="text-lg font-semibold mb-2">{date}</h3>

                        {/* Notification Items */}
                        {(notifications as UserNotification[]).map((notification) => (
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
                                    <p
                                        className="font-medium"

                                    > {notification.title}</p>
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
