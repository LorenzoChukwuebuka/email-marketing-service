import { Link, useParams } from 'react-router-dom';
import { Helmet, HelmetProvider } from 'react-helmet-async';
import { useSingleUserQuery, useUserStatsQuery } from '../hooks/useAdminUsersQueryHook';
import { useMemo } from 'react';

const AdminUserDetailComponent = () => {
    const { id } = useParams<{ id: string }>() as { id: string };
    const { data: userData } = useSingleUserQuery(id);
    const { data: statsData } = useUserStatsQuery(id)

    const userdetailsData = useMemo(() => userData?.payload, [userData])
    const userStatsData = useMemo(() => statsData?.payload, [statsData])
    return (
        <HelmetProvider>
            <Helmet title={`Details for ${userdetailsData?.fullname}`} />
            <div className="container mx-auto p-4">
                <h1 className="text-2xl font-bold mb-4">User Detail</h1>
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
                    <div className="bg-blue-600 text-white p-4 rounded-lg shadow">
                        <div>
                            <p className="text-sm">Total Campaigns Created</p>
                            <p className="text-2xl font-bold">{userStatsData?.total_campaigns}</p>
                        </div>
                        <Link to={{ pathname: "/zen/campaigns/details/" + id, search: `?username=${userdetailsData?.fullname}` }} className="text-xs mt-2 bg-white text-blue-600 px-2 py-1 rounded hover:bg-blue-100 transition-colors duration-200 text-center">
                            View Campaigns
                        </Link>
                    </div>

                    <div className="bg-purple-600 text-white p-4 rounded-lg shadow">
                        <div>
                            <p className="text-sm">Total Templates Created</p>
                            <p className="text-2xl font-bold">{userStatsData?.total_templates}</p>
                        </div>

                    </div>

                    {/* <div className="bg-blue-500 text-white p-4 rounded-lg shadow">
                        <div>
                            <p className="text-sm">Total Campaigns Sent</p>
                            <p className="text-2xl font-bold">{userStatsData.total_campaigns_sent}</p>
                        </div>
                        <a href={"/zen/dash/users/campaigns/"} className="text-xs mt-2 bg-white text-blue-500 px-2 py-1 rounded hover:bg-blue-100 transition-colors duration-200 text-center">
                            View Sent Campaigns
                        </a>
                    </div> */}

                    {/* <div className="bg-blue-700 text-white p-4 rounded-lg shadow">
                    <div>
                        <p className="text-sm">Total Subscriptions</p>
                        <p className="text-2xl font-bold">{userStatsData.total_subscriptions}</p>
                    </div>
                    <a href="#" className="text-xs mt-2 bg-white text-blue-700 px-2 py-1 rounded hover:bg-blue-100 transition-colors duration-200 text-center">
                        View Subscriptions
                    </a>
                </div> */}

                    <div className="bg-blue-600 text-white p-4 rounded-lg shadow">
                        <div>
                            <p className="text-sm">Total Groups Created</p>
                            <p className="text-2xl font-bold">{userStatsData?.total_groups}</p>
                        </div>

                    </div>

                    <div className="bg-purple-600 text-white p-4 rounded-lg shadow">
                        <div>
                            <p className="text-sm">Total Contacts</p>
                            <p className="text-2xl font-bold">{userStatsData?.total_contacts}</p>
                        </div>

                    </div>
                </div>


                <div className="bg-white shadow rounded-lg p-6 mb-6">
                    <h2 className="text-xl font-semibold mb-4">Information of {userdetailsData?.fullname || ""}</h2>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Full Name</label>
                            <input type="text" value={userdetailsData?.fullname || ""} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" readOnly />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Email</label>
                            <input type="email" value={userdetailsData?.email || ""} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" readOnly />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Mobile Number</label>
                            <input type="tel" value={userdetailsData?.phonenumber || "N/A"} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" readOnly />
                        </div>
                        <div className="md:col-span-2">
                            <label className="block text-sm font-medium text-gray-700">Company </label>
                            <input type="text" value={userdetailsData?.company || ""} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" readOnly />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Verified</label>
                            <input type="text" value={userdetailsData?.verified ? "Yes" : "No"} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" readOnly />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Verified At</label>
                            <input type="text" value={userdetailsData?.verified_at && new Date(userdetailsData?.verified_at).toLocaleString('en-US', {
                                timeZone: 'UTC',
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric',
                                hour: 'numeric',
                                minute: 'numeric',
                                second: 'numeric'
                            }) || "Not Verified Yet"} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" readOnly />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700"> Blocked </label>
                            <input type="text" value={userdetailsData?.blocked ? "Yes" : "No"} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" readOnly />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Joined on</label>
                            <input type="text" value={userdetailsData?.created_at && new Date(userdetailsData?.created_at).toLocaleString('en-US', {
                                timeZone: 'UTC',
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric',
                                hour: 'numeric',
                                minute: 'numeric',
                                second: 'numeric'
                            })} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" readOnly />
                        </div>
                    </div>
                </div>

                {/* <div className="bg-white shadow rounded-lg p-6">
                    <h2 className="text-xl font-semibold mb-4">User information</h2>
                    {/* <div className="space-y-2">
                        <p><span className="font-medium">Username:</span> johndoe123</p>
                        <p><span className="font-medium">Status:</span> <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">Active</span></p>
                        <p><span className="font-medium">Current Plan:</span> 2.50000000 BTC</p>
                   
                    </div> 
                </div> */}
            </div>
        </HelmetProvider>
    );
};

export default AdminUserDetailComponent;