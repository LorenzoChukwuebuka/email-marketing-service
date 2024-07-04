import { useEffect, useState } from "react";
import useAPIKeyStore from "../../../../store/userstore/apiKeyStore";
import {
  convertToNormalTime,
  copyToClipboard,
  maskAPIKey,
} from "../../../../utils/utils";

const APIKeysTableComponent = () => {
  const { getAPIKey, apiKeyData } = useAPIKeyStore();
  const [copied, setCopied] = useState(false);

  const shouldRenderNoKey = () => {
    return !apiKeyData || apiKeyData.payload === null;
  };

  const copyText = (text) => {
    copyToClipboard(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 1000);
  };

  useEffect(() => {
    getAPIKey();
  }, [getAPIKey]);

  if (shouldRenderNoKey()) {
    return (
      <div className="max-w-4xl text-center text-lg font-semibold mx-auto p-6 bg-white shadow-sm rounded-lg overflow-hidden">
        You have not generated an API Key yet
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h2 className="text-xl font-semibold mb-4">API Key</h2>

      <div className="bg-white shadow-sm rounded-lg overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="bg-gray-50">
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                API key value
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
                  <span className="text-sm text-gray-500 mr-2">
                    {maskAPIKey(apiKeyData?.payload?.api_key)}
                  </span>
                  <button
                    className="p-1 px-2 rounded-full bg-gray-200 hover:bg-gray-300"
                    onClick={() => copyText(apiKeyData?.payload?.api_key)}
                  >
                    {copied ? (
                      <i className="bi bi-check-circle"></i>
                    ) : (
                      <i className="bi bi-clipboard2"></i>
                    )}
                  </button>
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                  Active
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {convertToNormalTime(apiKeyData?.payload?.created_at)}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default APIKeysTableComponent;
