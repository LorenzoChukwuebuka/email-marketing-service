import React, { useEffect, useState } from "react"
import { Modal, Form, Input, Button } from "antd"
import {Campaign} from "./../../interface/campaign.interface"
import { BaseEntity } from "./../../../../interface/baseentity.interface";
import useCampaignStore from "./../../store/campaign.store";


interface Props {
    isOpen: boolean
    onClose: () => void
    campaign: (Campaign & BaseEntity) | null
}

interface FormValues {
    subject: string
    preview_text: string
}

const AddCampaignSubjectComponent: React.FC<Props> = ({ isOpen, onClose, campaign }) => {
    const [form] = Form.useForm<FormValues>()
    const { updateCampaign, setCreateCampaignValues } = useCampaignStore()
    const [isLoading, setIsLoading] = useState<boolean>(false)

    useEffect(() => {
        if (campaign) {
            form.setFieldsValue({
                subject: campaign.subject ?? "",
                preview_text: campaign.preview_text ?? ""
            })
        }
    }, [campaign, form])

    const handleSubmit = async (values: FormValues) => {
        try {
            setIsLoading(true)
            if (campaign?.id) {
                setCreateCampaignValues({
                    subject: values.subject,
                    preview_text: values.preview_text
                })
                await updateCampaign(campaign.id)
                new Promise(resolve => setTimeout(resolve, 1000))
                onClose()
            }
        } catch (error) {
            console.log(error)
        } finally {
            setIsLoading(false)
        }

    }

    return (
        <Modal
            open={isOpen}
            onCancel={onClose}
            title="Add Campaign Subject"
            footer={null}
        >
            <Form
                form={form}
                onFinish={handleSubmit}
                layout="vertical"
            >
                <Form.Item
                    label={<>Subject <span style={{ color: '#ff4d4f' }}>*</span></>}
                    name="subject"
                    help="Subject is what your audience sees in the title of your email"
                    rules={[{ required: true, message: 'Please input the subject!' }]}
                >
                    <Input
                        placeholder="Add a subject..."
                    />
                </Form.Item>

                <Form.Item
                    label="Preview text"
                    name="preview_text"
                    help="Preview Text tells your audience more about the mail"
                >
                    <Input
                        placeholder="Add a preview text..."
                    />
                </Form.Item>

                <Form.Item className="flex justify-end">
                    <Button.Group>
                        <Button onClick={onClose}>
                            Cancel
                        </Button>
                        <Button
                            type="primary"
                            htmlType="submit"
                            disabled={isLoading}
                            loading={isLoading}
                        >
                            {isLoading ? "Please wait ..." : "Save"}
                        </Button>
                    </Button.Group>
                </Form.Item>
            </Form>
        </Modal>
    )
}

export default AddCampaignSubjectComponent