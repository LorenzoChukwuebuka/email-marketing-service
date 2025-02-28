import React, { useState, useEffect } from "react";
import { Link, Outlet, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { useNotificationQuery } from "./../../hooks/useNotificationQuery";
import LoadingSpinnerComponent from "./../../components/loadingSpinnerComponent";
import eventBus from "../../utils/eventbus";
import { useMailCalcQuery } from "../../hooks/useMailDataQuery";
import useUserNotificationStore from "../../store/notification.store";
import { Modal } from "antd";

interface ContentProps {
    sidebarOpen: boolean;
    setSidebarOpen: (open: boolean) => void;
}

const DashContent: React.FC<ContentProps> = ({ sidebarOpen, setSidebarOpen }) => {
    const [userName, setUserName] = useState<string>("");
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [isNotificationDropdownOpen, setIsNotificationDropdownOpen] = useState(false);
    const { data: notificationsData, isError: notificationError, isLoading: notificationLoader } = useNotificationQuery();
    const { data: mailData } = useMailCalcQuery()
    const navigate = useNavigate();
    const { updateReadStatus } = useUserNotificationStore()

    const hasNotifications = notificationsData && notificationsData.payload?.length > 0;
    const hasUnreadNotifications = hasNotifications && notificationsData.payload.some(notification => notification.read_status === false);

    const toggleDropdown = () => {
        setIsDropdownOpen(!isDropdownOpen);
    };

    const closeDropdown = () => {
        setIsDropdownOpen(false);
    };

    if (notificationError) {
        eventBus.emit('error', 'Failed to fetch notifications');
    }

    const toggleNotificationDropdown = () => {
        setIsNotificationDropdownOpen(!isNotificationDropdownOpen);
    };

    const readNotifications = async () => {
        if (hasUnreadNotifications) {
            await updateReadStatus()
        }
    };

    const Logout = (): void => {
        Modal.confirm({
            title: "Are you sure?",
            content: "Do you want to logout?",
            okText: "Yes",
            cancelText: "No",
            onOk: () => {
                const cookies = Cookies.get("Cookies");
                if (cookies) {
                    Cookies.remove("Cookies");
                    navigate("/auth/login");
                }
            },
        })
    };

    useEffect(() => {
        const cookie = Cookies.get("Cookies");
        const user = cookie ? JSON.parse(cookie)?.details?.fullname : "";
        setUserName(user);
    }, []);

    return (
        <div className="flex-1 flex flex-col">
            <header className="bg-white h-14 p-4 flex justify-between items-center w-full">
                <button
                    onClick={() => setSidebarOpen(!sidebarOpen)}
                    className="text-gray-500 hover:text-gray-700"
                >
                    <span style={{ fontSize: "24px" }}>{sidebarOpen ? "≡" : "☰"}</span>
                </button>

                <div className="space-x-4 flex items-center">
                    <div className="relative">
                        <button className="relative" onClick={toggleNotificationDropdown}>
                            <i className="bi bi-bell font-bold text-xl" onClick={readNotifications}></i>
                            {hasUnreadNotifications && (
                                <span className="absolute top-0 right-0 block h-2 w-2 rounded-full bg-red-500"></span>
                            )}
                        </button>
                        {isNotificationDropdownOpen && (
                            <div className="absolute right-0 w-80 mt-4 bg-white shadow-lg rounded-md p-2 z-50">
                                <h3 className="font-bold text-lg mb-2 px-3">Notifications</h3>
                                <div className="max-h-80 overflow-y-auto">
                                    {notificationLoader ? (
                                        <LoadingSpinnerComponent />
                                    ) : hasNotifications ? (
                                        notificationsData?.payload?.map((notification, index) => (
                                            <div key={index} className="p-3 border-b border-gray-200 last:border-b-0">
                                                <p className="text-sm font-semibold">{notification.title}</p>
                                                <p className="text-xs text-gray-500 mt-1">{notification.created_at}</p>
                                            </div>
                                        ))
                                    ) : (
                                        <p className="text-sm p-3">No new notifications</p>
                                    )}
                                </div>
                                {hasNotifications && (
                                    <Link to="/app/notifications" className="text-blue-500 text-sm font-semibold p-3 block text-center hover:bg-gray-100 rounded-b-md">
                                        View all activity
                                    </Link>
                                )}
                            </div>
                        )}
                    </div>

                    <div className="relative">
                        <button className="flex items-center" onClick={toggleDropdown}>
                            {userName} <i className={`bi ${isDropdownOpen ? 'bi-chevron-up' : 'bi-chevron-down'}`}></i>
                        </button>
                        {isDropdownOpen && (
                            <ul className="absolute right-0 w-60 mt-4 bg-white shadow-lg rounded-md p-2 z-50">
                                <li className="p-2 border-b">
                                    <div className="text-center">
                                        <span>{userName}</span>
                                        <span className="block text-xs text-gray-500">Emails sent: {mailData?.payload?.remainingMails}/{mailData?.payload?.mailsPerDay}</span>
                                        <span className="block text-sm bg-gray-300 px-2 py-1 rounded-md">Plan: {mailData?.payload?.plan}</span>
                                    </div>
                                </li>
                                <li className="p-2">
                                    <Link to="/app/settings/account-management" className="block">My Profile</Link>
                                    <Link to="/app/settings/api" className="block">API & SMTP</Link>
                                    <button onClick={() => { Logout(); closeDropdown(); }} className="block w-full text-left">Logout</button>
                                </li>
                            </ul>
                        )}
                    </div>
                </div>
            </header>

            <main className="flex-1 p-2 w-full overflow-auto">
                <Outlet />
            </main>
        </div>
    );
};

export default DashContent;
