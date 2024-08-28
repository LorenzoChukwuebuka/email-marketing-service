import React, { useEffect, useState } from 'react';
import { CampaignData, CampaignGroup } from '../../../../store/userstore/campaignStore';
import Cookies from 'js-cookie'
import useContactGroupStore, { ContactGroupData } from '../../../../store/userstore/contactGroupStore';

type InfoItemProp = { label: string; value: string | React.ReactNode }

const InfoItem = ({ label, value }: InfoItemProp) => (
    <div className="mb-4">
        <span className="text-sm font-medium text-gray-500">{label}</span>
        <div className="mt-1 text-sm text-gray-900">{value}</div>
    </div>
);


type GProps = { campaignGroups: CampaignGroup[], contactGroups: ContactGroupData[] }

const CampaignGroupsList = ({ campaignGroups, contactGroups }: GProps) => {
    const matchedGroups = Array.isArray(campaignGroups) && Array.isArray(contactGroups)
        ? campaignGroups.map(campaignGroup => {
            const matchedGroup = contactGroups.find(contactGroup => contactGroup.id === campaignGroup.group_id);
            return matchedGroup ? matchedGroup.group_name : 'Unknown Group';
        })
        : [];

    return (
        <ul className="mt-1 text-sm text-gray-900">
            {matchedGroups.length > 0 ? (
                matchedGroups.map((groupName, index) => (
                    <li key={index}>{groupName}</li>
                ))
            ) : (
                <li>No groups available</li>
            )}
        </ul>
    );
};



type CampaignInfoProps = { campaignData: CampaignData }

const CampaignInfo: React.FC<CampaignInfoProps> = ({ campaignData }) => {

    let cookie: any = Cookies.get("Cookies");
    let user = JSON.parse(cookie)?.details?.email;
    const { contactgroupData, getAllGroups } = useContactGroupStore()
    const [templatePreview, setTemplatePreview] = useState<string | null>(null);

    useEffect(() => {
        getAllGroups()
    }, [getAllGroups])


    const handleClick = () => {
        setTemplatePreview(campaignData.template?.email_html as string)
    }


    return (
        <div className="bg-white mt-5 shadow overflow-hidden sm:rounded-lg">
            <div className="px-4 py-5 sm:px-6">
                <h3 className="text-lg leading-6 font-medium text-indigo-900">Campaign Info</h3>
            </div>
            <div className="border-t border-gray-200 px-4 py-5 sm:p-0">
                <dl className="sm:divide-y sm:divide-gray-200">
                    <div className="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                        <InfoItem label="Email Subject" value={campaignData?.subject as string} />
                        <InfoItem label="Sender Email" value={user} />
                        <InfoItem label="Sender From" value={campaignData?.sender_from_name as string} />
                    </div>
                    <div className="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                        <InfoItem
                            label="Group/Segment Sent To"
                            value={<CampaignGroupsList campaignGroups={campaignData.campaign_groups} contactGroups={contactgroupData as ContactGroupData[]} />}
                        />
                        <InfoItem
                            label="Sent On"
                            value={new Date(campaignData.sent_at as string).toLocaleDateString('en-GB', {
                                day: '2-digit',
                                month: '2-digit',
                                year: 'numeric',
                                hour: '2-digit',
                                minute: '2-digit',
                                second: '2-digit',
                            })}
                        />
                        {/* <InfoItem label="Audience" value={campaignData.audience} /> */}
                    </div>
                    <div className="py-4 sm:py-5 sm:px-6">
                        <h4 className="text-md font-medium text-gray-900 mb-2">Content</h4>
                        <InfoItem label="Template Selected" value={campaignData.template?.template_name as string} />
                        <div className="mt-4 space-x-2">
                            <button className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500" onClick={() => handleClick()}>
                                Preview template
                            </button>
                            <button className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                                Save as new template
                            </button>
                        </div>

                        {templatePreview && (
                            <div className="mt-4">
                                <h3 className="text-lg font-semibold">Template Preview</h3>
                                <div className="border p-4 mt-2" dangerouslySetInnerHTML={{ __html: templatePreview }} />
                            </div>
                        )}
                    </div>
                </dl>
            </div>
        </div>
    );
};

export default CampaignInfo;