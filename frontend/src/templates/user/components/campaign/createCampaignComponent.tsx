import { FormEvent, useState } from "react";
import { Modal } from "../../../../components";
import * as Yup from "yup";
import useCampaignStore from "../../../../store/userstore/campaignStore";

interface Props {
    isOpen: boolean;
    onClose: () => void;
}
const CreateCampaignComponent: React.FC<Props> = ({ isOpen, onClose }) => {

    const [errors, setErrors] = useState<{ [key: string]: string }>({});

    const { createCampaignValues, createCampaign, getAllCampaigns, setCreateCampaignValues } = useCampaignStore()

    const validationSchema = Yup.object().shape({
        name: Yup.string()
            .required("campaign name is required"),
    });

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()
        try {
            await validationSchema.validate(createCampaignValues, { abortEarly: false });
            await createCampaign()
            await new Promise(resolve => setTimeout(resolve, 500))
            await getAllCampaigns()
            onClose()
        } catch (err) {
            const validationErrors: { [key: string]: string } = {};
            if (err instanceof Yup.ValidationError) {
                err.inner.forEach((error) => {
                    validationErrors[error.path || ""] = error.message;
                });
                setErrors(validationErrors);
            }
        }
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { id, value } = e.target;
        setCreateCampaignValues({ ...createCampaignValues, [id]: value });
    };


    return <>
        <Modal isOpen={isOpen} onClose={onClose} title="Create Campaigns">

            <p className="mt-2 mb-5"> Keep subscribers engaged by sharing your latest news, promoting your bestselling products, or announcing an upcoming event. </p>
            <form onSubmit={handleSubmit}>
                <div className="mb-4">
                    <label
                        htmlFor="first_name"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Campaign name
                    </label>
                    <input
                        type="text"
                        id="name"
                        value={createCampaignValues.name}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />

                    {errors.group_name && (
                        <div style={{ color: "red" }}>{errors.group_name}</div>
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
                        type="submit"
                        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                    >
                        Create
                    </button>
                </div>
            </form>

        </Modal>
    </>
}


export default CreateCampaignComponent