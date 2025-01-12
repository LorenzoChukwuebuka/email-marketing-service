import { FormEvent } from "react"
import { Modal } from "antd"
import useTemplateStore from "../store/template.store"

interface Props {
    isOpen: boolean
    onClose: () => void,
    template_id: string
}

const SendTestEmail: React.FC<Props> = ({ isOpen, onClose, template_id }) => {

    const { setEmailTestValues, sendTestMail, sendEmailTestValues } = useTemplateStore()

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()

        setEmailTestValues({ ...sendEmailTestValues, template_id: template_id })
        await sendTestMail()
        new Promise(resolve => setTimeout(resolve, 1000))
        onClose()
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { id, value } = e.target;
        setEmailTestValues({ ...sendEmailTestValues, [id]: value });
    };

    return <>
        <Modal
            title="Test Email"
            open={isOpen}
            onCancel={onClose}
            footer={null}
        >
            <form onSubmit={handleSubmit}>
                <div className="mb-4">
                    <label
                        htmlFor="first_name"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Subject
                    </label>
                    <input
                        type="text"
                        id="subject"
                        value={sendEmailTestValues.subject}
                        onChange={handleChange}
                        placeholder="subject"
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>

                <div className="mb-4">
                    <label
                        htmlFor="first_name"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Send test email to (Max of 10 emails)
                    </label>
                    <textarea
                        id="email_address"
                        value={sendEmailTestValues.email_address}
                        onChange={handleChange}
                        placeholder="use comma to separate the emails"
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>

                <small className="text-red-500 text-center mb-4"> Note: Sending test emails will count against your  daily usage </small>

                <div className="flex justify-end mt-4 space-x-2">
                    <button
                        type="button"
                        onClick={onClose}
                        className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                    >
                        Send Test
                    </button>
                </div>
            </form>
        </Modal>
    </>
}


export default SendTestEmail