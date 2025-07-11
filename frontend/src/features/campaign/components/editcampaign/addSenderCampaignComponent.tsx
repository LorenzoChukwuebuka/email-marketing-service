import React, { useEffect, useMemo, useState } from "react";
import { Form, Input, Select, Button, Modal } from "antd";
import { Campaign } from "./../../interface/campaign.interface"
import { BaseEntity } from "./../../../../interface/baseentity.interface";
import useCampaignStore from "./../../store/campaign.store";
import { useSenderQuery } from "./../../../settings/hooks/useSenderQuery";

const { Option } = Select;

interface Props {
    isOpen: boolean;
    onClose: () => void;
    campaign: (Campaign & BaseEntity) | null;
}

const AddSenderComponent: React.FC<Props> = ({ isOpen, onClose, campaign }) => {
    const { createCampaignValues, setCreateCampaignValues, updateCampaign } = useCampaignStore();
    /* eslint-disable @typescript-eslint/no-unused-vars */
    const [currentPage, _setCurrentPage] = useState(1);
    const [pageSize, _setPageSize] = useState(2000);

    const { data: senderData } = useSenderQuery(currentPage, pageSize, undefined)
    const [form] = Form.useForm();

    const sData = useMemo(() => senderData?.payload.data || [], [senderData])

    useEffect(() => {
        if (isOpen && campaign) {
            // Set form values from campaign data
            form.setFieldsValue({
                sender_email: campaign.sender || "",
                sender_from_name: campaign.sender_from_name || "",
            });
            
            // Also update the store with current campaign values
            setCreateCampaignValues({
                sender: campaign.sender || "",
                sender_from_name: campaign.sender_from_name || "",
            });
        }
    }, [isOpen, campaign, form, setCreateCampaignValues]);

    // Handle when sender email is selected from dropdown
    const handleSenderChange = (selectedEmail: string) => {
        const sender = sData.find(s => s.email === selectedEmail);
        if (sender) {
            form.setFieldsValue({
                sender_email: selectedEmail,
                sender_from_name: sender.name,
            });
        }
    };

    const handleFormSubmit = async (values: { sender_email: string; sender_from_name: string }) => {
        setCreateCampaignValues({
            sender: values.sender_email,
            sender_from_name: values.sender_from_name,
        });
        await updateCampaign(campaign?.id as string);
        onClose();
    };

    return (
        <Modal
            title="Sender Details"
            open={isOpen}
            onCancel={onClose}
            footer={null}
            destroyOnClose
        >
            <h1 className="mt-4 text-lg font-semibold mb-4">Who is sending this email campaign?</h1>

            <Form
                form={form}
                layout="vertical"
                onFinish={handleFormSubmit}
            >
                <Form.Item
                    label="Email Address"
                    name="sender_email"
                    rules={[{ required: true, message: "Please select an email address" }]}
                >
                    <Select 
                        placeholder="Select an email..."
                        onChange={handleSenderChange}
                    >
                        {Array.isArray(sData) &&
                            sData.map((sender) => (
                                <Option key={sender.id} value={sender.email}>
                                    {sender.email}
                                </Option>
                            ))}
                    </Select>
                </Form.Item>

                <Form.Item
                    label="Name"
                    name="sender_from_name"
                    rules={[{ required: true, message: "Please enter a name" }]}
                >
                    <Input placeholder="Enter sender name..." />
                </Form.Item>

                <div className="flex justify-end space-x-2">
                    <Button onClick={onClose}>Cancel</Button>
                    <Button type="primary" htmlType="submit">
                        Save
                    </Button>
                </div>
            </Form>
        </Modal>
    );
};

export default AddSenderComponent;