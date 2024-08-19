import { useEffect, useState } from "react";
import useAPIKeyStore from "../../../../store/userstore/apiKeyStore";
import { convertToNormalTime, maskAPIKey } from "../../../../utils/utils";
import EmptyState from "../../../../components/emptyStateComponent";

interface APIKey {
    uuid: string;
    name: string;
    api_key: string;
    created_at: string;
}

const APIKeysTableComponent: React.FC = () => {
    const { getAPIKey, apiKeyData, deleteAPIKey } = useAPIKeyStore();

    const [deletingId, setDeletingId] = useState<string | null>(null);

    const shouldRenderNoKey = () => {
        return (
            !apiKeyData || !apiKeyData.payload || apiKeyData.payload.length === 0
        );
    };

    const handleDelete = async (id: string) => {
        setDeletingId(id);
        await deleteAPIKey(id);
        getAPIKey();
        setDeletingId(null);
    };

    useEffect(() => {
        getAPIKey();
    }, [getAPIKey]);

    if (shouldRenderNoKey()) {
        return (
            <> <EmptyState title="You  have not generated any API Key"
                description="Kindly Generate an API Key to enjoy our services"
                icon={<i className="bi bi-emoji-frown text-xl"></i>} />
            </>
        );
    }

    return (
        <div className="max-w-4xl mx-auto p-6">
            <h2 className="text-xl font-semibold mb-4">API Keys</h2>

            <div className="bg-white shadow-sm rounded-lg overflow-hidden">
                <table className="w-full">
                    <thead>
                        <tr className="bg-gray-50">
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Name
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                API key value
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Status
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Created on
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {apiKeyData.payload.map((key: APIKey) => (
                            <tr key={key.uuid}>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    {key.name}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="flex items-center">
                                        <span className="text-sm text-gray-500 mr-2">
                                            {maskAPIKey(key.api_key)}
                                        </span>
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                                        Active
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

export default APIKeysTableComponent;
