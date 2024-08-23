import { useEffect, useState } from "react";
import { Modal } from "../../../../components";
import { BaseEntity } from "../../../../interface/baseentity.interface";
import { Campaign, CampaignGroup } from "../../../../store/userstore/campaignStore";
import useCampaignStore from "../../../../store/userstore/campaignStore";
import useContactGroupStore, { ContactGroupData } from "../../../../store/userstore/contactGroupStore";

interface Props {
    isOpen: boolean;
    onClose: () => void;
    campaign: (Campaign & BaseEntity) | null;
}

const AddCampaignRecipients: React.FC<Props> = ({ isOpen, onClose, campaign }) => {
    const { getAllGroups, contactgroupData } = useContactGroupStore();
    const { selectedGroupIds, setSelectedGroupIds, createCampaignGroup } = useCampaignStore();
    const [selectedGroups, setSelectedGroups] = useState<ContactGroupData[]>([]);

    useEffect(() => {
        const fetchGroup = async () => {
            await getAllGroups();
        };

        fetchGroup();
    }, [getAllGroups]);

    useEffect(() => {
        if (contactgroupData && selectedGroupIds) {
            setSelectedGroups(
                (contactgroupData as ContactGroupData[])?.filter(group => selectedGroupIds.includes(group.uuid))
            );
        }
    }, [selectedGroupIds, contactgroupData]);


    useEffect(() => {
        if (contactgroupData && campaign && campaign.campaign_groups) {
            setSelectedGroups((contactgroupData as ContactGroupData[]).filter(group => {
                // Check if the campaign has any associated groups
                if (campaign.campaign_groups.length > 0) {
                    // Filter the contactgroupData to include only the groups associated with the campaign
                    return campaign.campaign_groups.some(cg => cg.group_id === group.id);
                }
                // If the campaign has no associated groups, don't select any groups
                return false;
            }));
        } else {
            setSelectedGroups([]); // Clear selected groups if there are no campaign groups
        }
    }, [contactgroupData, campaign]);


    const handleGroupSelect = (group: ContactGroupData) => {
        setSelectedGroups((prevSelected) => {
            const isSelected = prevSelected.find(g => g.uuid === group.uuid);
            return isSelected
                ? prevSelected.filter((g) => g.uuid !== group.uuid)
                : [...prevSelected, group];
        });
    };

    const handleRemoveGroup = (uuid: string) => {
        setSelectedGroups((prevSelected) => {
            const updatedGroups = prevSelected.filter((g) => g.uuid !== uuid);
            setSelectedGroupIds(updatedGroups.map(g => g.uuid));
            return updatedGroups;
        });
    };

    const handleGroupSubmit = async (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();

        // Assuming selectedGroups is your local state in the component
        const groupIds = selectedGroups.map(group => group.uuid);

        if (campaign?.uuid) {
            await createCampaignGroup(campaign.uuid, groupIds);
        } else {
            console.error("Campaign UUID is not available");
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Add Recipients">
            <>
                <div className="mb-2">
                    {selectedGroups.length > 0 && (
                        <div className="flex flex-wrap gap-2">
                            {selectedGroups.map((group) => (
                                <div key={group.uuid} className="flex items-center space-x-2 bg-gray-200 rounded px-2 py-1">
                                    <span>{group.group_name}</span>
                                    <button
                                        onClick={() => handleRemoveGroup(group.uuid)}
                                        className="text-red-500 font-semibold"
                                    >
                                        &times;
                                    </button>
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                <div className="max-h-60 overflow-y-auto">
                    <h1 className="mt-4 mb-4">Select one or more groups</h1>
                    {contactgroupData && (contactgroupData as ContactGroupData[]).length > 0 ? (
                        (contactgroupData as ContactGroupData[]).map((group: ContactGroupData) => (
                            <div key={group.uuid} className="mb-4">
                                <label className="flex items-center space-x-2">
                                    <input
                                        type="checkbox"
                                        name="group"
                                        className="form-checkbox text-blue-600"
                                        checked={selectedGroups.some(g => g.uuid === group.uuid)}
                                        onChange={() => handleGroupSelect(group)}
                                    />
                                    <span className="font-semibold space-x-5">{group.group_name} ({group.contacts ? group.contacts.length : 0}) contacts</span>
                                </label>
                                <p className="text-sm text-gray-500 ml-6">{group.description}</p>
                            </div>
                        ))
                    ) : (
                        <div className="flex items-center justify-center text-lg font-semibold">No groups found</div>
                    )}
                </div>

                <div className="flex justify-end space-x-2">
                    <button
                        type="button"
                        onClick={onClose}
                        className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
                    >
                        Cancel
                    </button>
                    <button
                        type="button"
                        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                        onClick={handleGroupSubmit}
                        disabled={selectedGroups.length === 0}
                    >
                        Save
                    </button>
                </div>
            </>
        </Modal>
    );
};

export default AddCampaignRecipients;
