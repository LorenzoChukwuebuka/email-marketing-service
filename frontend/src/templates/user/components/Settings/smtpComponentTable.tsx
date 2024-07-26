import { useEffect } from "react";
import useSMTPKeyStore from "../../../../store/userstore/smtpkeyStore";

const SMTPKeysTableComponent: React.FC = () => {
    const { getSMTPKeys, smtpKeyData } = useSMTPKeyStore()

    useEffect(() => {
        getSMTPKeys()
    }, [])

    return (
        <div className="max-w-4xl mx-auto p-6">
            <h2 className="text-2xl font-bold mb-6">Your SMTP Settings</h2>
            <div className="mb-6">
                <div className="mb-2">
                    <span className="font-semibold">SMTP Server:</span>{" "}
                    {smtpKeyData.smtp_server}
                </div>
                <div className="mb-2">
                    <span className="font-semibold">Port:</span> {smtpKeyData.smtp_port}
                </div>
                <div className="mb-2">
                    <span className="font-semibold">Login:</span>  {smtpKeyData.smtp_login}
                </div>
            </div>
            <button className="text-blue-600 hover:text-blue-800 mb-6">
                Regenerate SMTP Login and Master password
            </button>

            <h2 className="text-xl font-semibold mb-4">Your SMTP Keys</h2>

            <div className="bg-white shadow-sm rounded-lg overflow-hidden">
                <table className="w-full">
                    <thead>
                        <tr className="bg-gray-50">
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                SMTP key name
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                SMTP key value
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Status
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Created on
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        <tr>
                            <td className="px-6 py-4 whitespace-nowrap">
                                <div className="flex items-center">
                                    <input
                                        type="checkbox"
                                        className="mr-3 h-4 w-4 text-blue-600"
                                    />
                                    <div className="text-sm font-medium text-gray-900">
                                        {smtpKeyData.smtp_master}
                                    </div>
                                </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                                <div className="flex items-center">
                                    <span className="text-sm text-gray-500 mr-2">
                                        {smtpKeyData.smtp_master_password}
                                    </span>
                                    <button className="p-1 rounded-full bg-gray-200 hover:bg-gray-300">
                                        <svg
                                            className="h-4 w-4 text-gray-600"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                strokeLinecap="round"
                                                strokeLinejoin="round"
                                                strokeWidth={2}
                                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                                            />
                                            <path
                                                strokeLinecap="round"
                                                strokeLinejoin="round"
                                                strokeWidth={2}
                                                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                                            />
                                        </svg>
                                    </button>
                                </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                                <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                                    {smtpKeyData.smtp_master_status}
                                </span>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {/* You might want to add a created_at field for the master key */}
                                N/A
                            </td>
                        </tr>
                        {smtpKeyData.keys && smtpKeyData.keys.map((key) => (
                            <tr key={key.uuid}>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="flex items-center">
                                        <input
                                            type="checkbox"
                                            className="mr-3 h-4 w-4 text-blue-600"
                                        />
                                        <div className="text-sm font-medium text-gray-900">
                                            {key.key_name}
                                        </div>
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="flex items-center">
                                        <span className="text-sm text-gray-500 mr-2">
                                            {key.password}
                                        </span>
                                        <button className="p-1 rounded-full bg-gray-200 hover:bg-gray-300">
                                            <svg
                                                className="h-4 w-4 text-gray-600"
                                                fill="none"
                                                viewBox="0 0 24 24"
                                                stroke="currentColor"
                                            >
                                                <path
                                                    strokeLinecap="round"
                                                    strokeLinejoin="round"
                                                    strokeWidth={2}
                                                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                                                />
                                                <path
                                                    strokeLinecap="round"
                                                    strokeLinejoin="round"
                                                    strokeWidth={2}
                                                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                                                />
                                            </svg>
                                        </button>
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                                        {key.status}
                                    </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {new Date(key.created_at).toLocaleString()}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default SMTPKeysTableComponent;