import { useEffect } from "react";
import { Modal } from "../../../../../components";
import usePlanStore from "../../../../../store/admin/planStore";

interface Plan {
    uuid: string;
    planname: string;
    details: string;
    duration: string;
    price: string;
    number_of_mails_per_day: string;
}

interface EditPlansProps {
    isOpen: boolean;
    onClose: () => void;
    plan: Plan | null | any
}

const EditPlans: React.FC<EditPlansProps> = ({ isOpen, onClose, plan }) => {
    const {
        setEditPlanValues,
        editPlanValues,
        updatePlan,
        getPlans,
    } = usePlanStore();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        await updatePlan();
        await getPlans();
        onClose();
    };

    useEffect(() => {
        if (plan) {
            setEditPlanValues({
                uuid: plan.uuid,
                planname: plan.planname,
                details: plan.details,
                duration: plan.duration,
                price: plan.price,
                number_of_mails_per_day: plan.number_of_mails_per_day,
            });
        }
    }, [plan, setEditPlanValues]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { id, value } = e.target;
        setEditPlanValues({ ...editPlanValues, [id]: value });
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Edit Plan">
            <form onSubmit={handleSubmit}>
                <div className="mb-4">
                    <label
                        htmlFor="planName"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Plan Name
                    </label>
                    <input
                        type="text"
                        id="planName"
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        onChange={handleChange}
                        value={editPlanValues.planname}
                    />
                </div>

                <div className="mb-4">
                    <label
                        htmlFor="planDescription"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Description
                    </label>
                    <input
                        type="text"
                        id="planDescription"
                        onChange={handleChange}
                        value={editPlanValues.details}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                    />
                </div>

                <div className="mb-4">
                    <label
                        htmlFor="planDuration"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Duration (days)
                    </label>
                    <input
                        type="text"
                        id="planDuration"
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        onChange={handleChange}
                        value={editPlanValues.duration}
                    />
                </div>

                <div className="mb-4">
                    <label
                        htmlFor="planPrice"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Price
                    </label>
                    <input
                        type="text"
                        id="planPrice"
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        onChange={handleChange}
                        value={editPlanValues.price}
                    />
                </div>

                <div className="mb-4">
                    <label
                        htmlFor="number_of_mails_per_day"
                        className="block text-sm font-medium text-gray-700"
                    >
                        Number of Mails Per Day
                    </label>
                    <input
                        type="text"
                        id="number_of_mails_per_day"
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        onChange={handleChange}
                        value={editPlanValues.number_of_mails_per_day}
                    />
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
                        Submit
                    </button>
                </div>
            </form>
        </Modal>
    );
};

export default EditPlans;
