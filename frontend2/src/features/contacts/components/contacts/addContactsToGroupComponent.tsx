import {  useState } from "react";
import { Modal, Radio, Button, Empty } from "antd";
import useContactGroupStore from "./../../store/contactgroup.store";
import { useContactGroupQuery } from "../../hooks/useContactGroupQuery";


interface CGProps {
    isOpen: boolean;
    onClose: () => void;
}

const AddContactsToGroupComponent: React.FC<CGProps> = ({ isOpen, onClose }) => {
    const {  setSelectedGroupIds, addContactToGroup } = useContactGroupStore();

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [currentPage, _setCurrentPage] = useState(1);
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [pageSize, _setPageSize] = useState(20);

    const { data: contactgroupData } = useContactGroupQuery(currentPage, pageSize, undefined)

    const cData = contactgroupData?.payload.data

    const [selectedGroup, setSelectedGroup] = useState<string | null>(null);

    const handleGroupSubmit = async () => {
        if (selectedGroup) {
            setSelectedGroupIds([selectedGroup]);
            await addContactToGroup();
            setSelectedGroup(null);
            onClose();
        }
    };

    const handleGroupSelect = (uuid: string) => {
        setSelectedGroup(uuid);
    };

 

    return (
        <Modal
            title="Add selected Contact(s) to Group"
            open={isOpen}
            onCancel={onClose}
            footer={[
                <Button key="cancel" onClick={onClose}>
                    Cancel
                </Button>,
                <Button
                    key="submit"
                    type="primary"
                    onClick={handleGroupSubmit}
                    disabled={!selectedGroup}
                >
                    Submit
                </Button>
            ]}
        >
            <div className="space-y-4">
                <h4 className="mt-4 mb-4">You can only select one group</h4>

                <div className="max-h-60 overflow-y-auto">
                    {cData && cData.length > 0 ? (
                        <Radio.Group
                            onChange={(e) => handleGroupSelect(e.target.value)}
                            value={selectedGroup}
                            className="w-full"
                        >

                            {/*  eslint-disable-next-line @typescript-eslint/no-explicit-any */}
                            {cData.map((group: any) => (
                                <div key={group.uuid} className="mb-4">
                                    <Radio value={group.uuid}>
                                        <div>
                                            <span className="font-semibold">
                                                {group.group_name} ({group.contacts ? group.contacts.length : 0} contacts)
                                            </span>
                                            <p className="text-sm text-gray-500 ml-6">
                                                {group.description}
                                            </p>
                                        </div>
                                    </Radio>
                                </div>
                            ))}
                        </Radio.Group>
                    ) : (
                        <Empty description="No groups found" />
                    )}
                </div>
            </div>
        </Modal>
    );
}

export default AddContactsToGroupComponent;