import { useState, useMemo } from 'react';
import { convertToNormalTime, copyToClipboard, maskAPIKey } from "../../../../utils/utils";
import useSMTPKeyStore from "../../store/smtpkey.store";
import { useSMTPKeyQuery } from "../../hooks/useSmtpkeyQuery";
import { Modal } from 'antd';

const SMTPKeysTableComponent: React.FC = () => {
    const { deleteSMTPKey, generateSMTPKey } = useSMTPKeyStore()
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [deletingId, _setDeletingId] = useState<string | null>(null);
    const [copyingKey, setCopyingKey] = useState<string | null>(null);

    const { data: smtpKeyData } = useSMTPKeyQuery()

    const smkdata = useMemo(() => smtpKeyData?.payload, [smtpKeyData])

    const handleDelete = async (id: string) => {
        Modal.confirm({
            title: "Are you sure?",
            content: "Do you want to delete this smtp key?",
            okText: "Yes",
            cancelText: "No",
            onOk: async () => {
                await deleteSMTPKey(id)
                await new Promise(resolve => setTimeout(resolve, 3000));
                location.reload()
            },
        });
    }

    const handleCopy = (key: string) => {
        copyToClipboard(key);
        setCopyingKey(key);
        setTimeout(() => {
            setCopyingKey(null);
        }, 2000);
    };


    return (
        <div className="max-w-4xl mx-auto p-6">
            <h2 className="text-2xl font-bold mb-6">Your SMTP Settings</h2>
            <div className="mb-6">
                <div className="mb-2">
                    <span className="font-semibold">SMTP Server:</span>{" "}
                    {smkdata?.smtp_server}
                </div>
                <div className="mb-2">
                    <span className="font-semibold">Port:</span> {smkdata?.smtp_port}
                </div>
                <div className="mb-2">
                    <span className="font-semibold">Login:</span>  {smkdata?.smtp_login}
                </div>
            </div>
            <button className="text-blue-600 hover:text-blue-800 mb-6" onClick={async () => await generateSMTPKey()}>
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
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"> Delete </th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        <tr>
                            <td className="px-6 py-4 whitespace-nowrap">
                                <div className="flex items-center">

                                    <div className="text-sm font-medium text-gray-900">
                                        {smkdata?.smtp_master}
                                    </div>
                                </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                                <div className="flex space-x-2 items-center">
                                    <span className="text-sm text-gray-500 mr-2">
                                        {maskAPIKey(smkdata?.smtp_master_password as string)}
                                    </span>
                                    <button className="p-1 rounded-md bg-gray-200 hover:bg-gray-300" onClick={() => handleCopy(smkdata?.smtp_master_password as string)} title="click here to copy">
                                        <i className={`bi ${copyingKey === smkdata?.smtp_master_password ? 'bi-clipboard2-check' : 'bi-clipboard'}`}></i>
                                    </button>
                                </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                                <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                                    {smkdata?.smtp_master_status}
                                </span>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {new Date(smkdata?.smtp_created_at as string).toLocaleString('en-US', { timeZone: 'UTC', year: 'numeric', month: 'long', day: 'numeric', hour: 'numeric', minute: 'numeric', second: 'numeric' })}

                            </td>
                        </tr>
                        {Array.isArray(smkdata?.keys) && smkdata?.keys?.length > 0 && smkdata?.keys.map((key) => (
                            <tr key={key.uuid}>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="flex items-center">
                                        <div className="text-sm font-medium text-gray-900">
                                            {key.key_name}
                                        </div>
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="flex items-center">
                                        <span className="text-sm text-gray-500 mr-2">
                                            {maskAPIKey(key.password)}
                                        </span>
                                        <button className="p-1 rounded-full bg-gray-200 hover:bg-gray-300" onClick={() => copyToClipboard(key.password)}>
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
                                    {convertToNormalTime(key.created_at)}
                                </td>

                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    <button
                                        onClick={() => handleDelete(key.uuid)}
                                        className="text-red-600 hover:text-red-900"
                                        disabled={deletingId === key.uuid}
                                    >
                                        {deletingId === key.uuid ? (
                                            <span className="loading loading-spinner loading-sm"></span>
                                        ) : (
                                            <i className="bi bi-trash"></i>
                                        )}
                                    </button>
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