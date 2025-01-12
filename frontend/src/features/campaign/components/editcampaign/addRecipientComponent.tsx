import { useEffect, useMemo, useState } from "react";
import { Modal, Checkbox } from "antd";
import { useContactGroupQuery } from "./../../../contacts/hooks/useContactGroupQuery";
import { ContactGroupData } from "./../../../contacts/interface/contactgroup.interface";
import useCampaignStore from "./../../store/campaign.store";
import { Campaign } from "./../../interface/campaign.interface"
import { BaseEntity } from "./../../../../interface/baseentity.interface";


interface Props {
    isOpen: boolean;
    onClose: () => void;
    campaign: (Campaign & BaseEntity) | null;
}

const AddCampaignRecipients: React.FC<Props> = ({ isOpen, onClose, campaign }) => {

    const { selectedGroupIds, setSelectedGroupIds, createCampaignGroup } = useCampaignStore();
    const [selectedGroups, setSelectedGroups] = useState<ContactGroupData[]>([]);
    /* eslint-disable @typescript-eslint/no-unused-vars */
    const [currentPage, _setCurrentPage] = useState(1);
    const [pageSize, _setPageSize] = useState(2000);

    const { data: contactgroupData } = useContactGroupQuery(currentPage, pageSize, undefined)

    const cgdata = useMemo(() => contactgroupData?.payload.data || [], [contactgroupData]);


    useEffect(() => {
        if (cgdata && selectedGroupIds) {
            setSelectedGroups(
                cgdata?.filter(group => selectedGroupIds.includes(group.uuid))
            );
        }
    }, [selectedGroupIds, cgdata]);

    useEffect(() => {
        if (cgdata && campaign && campaign.campaign_groups) {
            setSelectedGroups((cgdata).filter(group => {
                if (campaign.campaign_groups.length > 0) {
                    return campaign.campaign_groups.some(cg => cg.group_id === group.id);
                }
                return false;
            }));
        } else {
            setSelectedGroups([]);
        }
    }, [cgdata, campaign]);

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

    const handleGroupSubmit = async () => {
        const groupIds = selectedGroups.map(group => group.uuid);

        if (campaign?.uuid) {
            await createCampaignGroup(campaign.uuid, groupIds);
        } else {
            console.error("Campaign UUID is not available");
        }

        onClose();
    };

    return (
        <Modal
            open={isOpen}
            onCancel={onClose}
            onOk={handleGroupSubmit}
            title="Add Recipients"
            okButtonProps={{ disabled: selectedGroups.length === 0 }}
        >
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
                {cgdata && cgdata.length > 0 ? (
                    cgdata.map((group: ContactGroupData) => (
                        <div key={group.uuid} className="mb-4">
                            <Checkbox
                                checked={selectedGroups.some(g => g.uuid === group.uuid)}
                                onChange={() => handleGroupSelect(group)}
                            >
                                <span className="font-semibold space-x-5">{group.group_name} ({group.contacts ? group.contacts.length : 0}) contacts</span>
                            </Checkbox>
                            <p className="text-sm text-gray-500 ml-6">{group.description}</p>
                        </div>
                    ))
                ) : (
                    <div className="flex items-center justify-center text-lg font-semibold">No groups found</div>
                )}
            </div>
        </Modal>
    );
};

export default AddCampaignRecipients;
