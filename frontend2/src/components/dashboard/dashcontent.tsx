import React, { useState, useEffect } from "react";
import { Link, Outlet, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { useNotificationQuery } from "./../../hooks/useNotificationQuery";
import LoadingSpinnerComponent from "./../../components/loadingSpinnerComponent";
import eventBus from "../../utils/eventbus";
import { useMailCalcQuery } from "../../hooks/useMailDataQuery";

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
            console.log('reading notifications');
        }
    };

    const Logout = (): void => {
        const confirmResult = confirm("Do you want to logout?");

        if (confirmResult) {
            const cookies = Cookies.get("Cookies");
            if (cookies) {
                Cookies.remove("Cookies");
                navigate("/auth/login");
            }
        }

        const confirmDialog = document.querySelector("div[role='dialog']");
        if (confirmDialog) {
            (confirmDialog as HTMLElement).style.backgroundColor = "white";
            (confirmDialog as HTMLElement).style.color = "black";
            (confirmDialog as HTMLElement).style.padding = "20px";
            (confirmDialog as HTMLElement).style.borderRadius = "5px";
            (confirmDialog as HTMLElement).style.boxShadow = "0 0 10px rgba(0, 0, 0, 0.3)";
        }
    };

    useEffect(() => {
        const cookie = Cookies.get("Cookies");
        const user = cookie ? JSON.parse(cookie)?.details?.fullname : "";
        setUserName(user);
    }, []);

    return (
        <div className="flex-1 flex flex-col">
            <header className="bg-white h-14 p-4 flex justify-between items-center">
                <button
                    onClick={() => setSidebarOpen(!sidebarOpen)}
                    className="text-gray-500 hover:text-gray-700"
                >
                    <span style={{ fontSize: "24px" }}>{sidebarOpen ? "≡" : "☰"}</span>
                </button>

                <div className="space-x-4 flex items-center">
                    <div className="dropdown dropdown-end relative">
                        <div
                            tabIndex={0}
                            role="button"
                            className="m-1 cursor-pointer relative"
                            onClick={toggleNotificationDropdown}
                        >
                            <i className="bi bi-bell font-bold text-xl" onClick={readNotifications}></i>
                            {hasUnreadNotifications && (
                                <span className="absolute top-0 right-0 block h-2 w-2 rounded-full bg-red-500"></span>
                            )}
                        </div>
                        {isNotificationDropdownOpen && (
                            <div className="dropdown-content menu bg-white rounded-box z-[50] mt-4 w-80 p-2 shadow absolute right-0">
                                <h3 className="font-bold text-lg mb-2 px-3">Notifications</h3>
                                <div className="max-h-80 overflow-y-auto">
                                    {notificationLoader ? (
                                        <LoadingSpinnerComponent />
                                    ) : (
                                        <>
                                            {hasNotifications ? (
                                                notificationsData?.payload?.map((notification, index) => (
                                                    <div key={index} className="p-3 border-b border-gray-200 last:border-b-0">
                                                        <p className="text-sm font-semibold">{notification.title}</p>
                                                        <p className="text-xs text-gray-500 mt-1">{notification.created_at}</p>
                                                    </div>
                                                ))
                                            ) : (
                                                <p className="text-sm p-3">No new notifications</p>
                                            )}
                                        </>
                                    )}
                                </div>
                                {hasNotifications && (
                                    <Link to="/user/dash/notifications" className="text-blue-500 text-sm font-semibold p-3 block text-center hover:bg-gray-100 rounded-b-lg">
                                        View all activity
                                    </Link>
                                )}
                            </div>
                        )}
                    </div>

                    <div className="dropdown dropdown-end relative">
                        <div
                            tabIndex={0}
                            role="button"
                            className="m-1"
                            onClick={toggleDropdown}
                        >
                            {userName}{' '}
                            <i className={`bi ${isDropdownOpen ? 'bi-chevron-up' : 'bi-chevron-down'}`}></i>
                        </div>
                        {isDropdownOpen && (
                            <ul tabIndex={0} className="dropdown-content menu bg-white rounded-box z-[50] mt-4 w-60 p-2 shadow">
                                <li>
                                    <div className="flex flex-col items-center">
                                        {userName}
                                        <span className="text-base-300 border-b-2 border-black mb-4">
                                            Emails sent: {mailData?.payload?.remainingMails}/
                                            {mailData?.payload?.mailsPerDay}
                                        </span>
                                        <span className="text-black bg-gray-300 px-1 py-1 rounded-md">
                                            Plan: {mailData?.payload?.plan}
                                        </span>
                                        <span className="text-blue-500">
                                            <i className="bi bi-person-fill"></i>
                                            <Link to="/user/dash/settings/account-management"> My Profile</Link>
                                        </span>
                                        <span className="text-black">
                                            <i className="bi bi-gear-wide-connected"></i>
                                            <Link to="/user/dash/settings/api"> API & SMTP</Link>
                                        </span>
                                        <a onClick={() => { Logout(); closeDropdown(); }}>
                                            <i className="bi bi-box-arrow-in-left"></i> Logout
                                        </a>
                                    </div>
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