import React, { useEffect, useState } from "react";
import { Modal } from "../../../../../components";
import usePlanStore from "../../../../../store/admin/planStore";

type Feature = {
    name: string;
    identifier: string;
    count_limit: number;
    size_limit: number;
    is_active: boolean;
    description: string;
}

interface Plan {
    uuid: string;
    planname: string;
    details: string;
    duration: string;
    price: number;
    number_of_mails_per_day: string;
    status: string;
    features: Feature[];
}

interface EditPlansProps {
    isOpen: boolean;
    onClose: () => void;
    plan: Plan | null;
}

const EditPlans: React.FC<EditPlansProps> = ({ isOpen, onClose, plan }) => {
    const {
        setEditPlanValues,
        editPlanValues,
        updatePlan,
        getPlans,
    } = usePlanStore();

    const [localEditPlanValues, setLocalEditPlanValues] = useState<Plan>({
        uuid: '',
        planname: '',
        details: '',
        duration: '',
        price: 0,
        number_of_mails_per_day: '',
        status: 'active',
        features: [],
    });

    useEffect(() => {
        if (plan) {
            setLocalEditPlanValues(plan);
        }
    }, [plan]);

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        //@ts-ignore
        setEditPlanValues(localEditPlanValues);
        await updatePlan();
        await getPlans();
        onClose();
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const { id, value } = e.target;
        setLocalEditPlanValues((prev) => ({
            ...prev,
            [id]: id === "price" ? parseFloat(value) : value,
        }));
    };

    const handleFeatureChange = (index: number, field: keyof Feature, value: string | number | boolean) => {
        setLocalEditPlanValues((prev) => {
            const updatedFeatures = [...prev.features];
            updatedFeatures[index] = {
                ...updatedFeatures[index],
                [field]: field === "is_active" ? value === "true" : value,
            };
            return { ...prev, features: updatedFeatures };
        });
    };

    const addFeature = () => {
        setLocalEditPlanValues((prev) => ({
            ...prev,
            features: [...prev.features, { name: "", identifier: "", count_limit: 0, size_limit: 0, is_active: true, description: "" }],
        }));
    };

    const removeFeature = (index: number) => {
        setLocalEditPlanValues((prev) => ({
            ...prev,
            features: prev.features.filter((_, i) => i !== index),
        }));
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Edit Plan">
            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label htmlFor="planname" className="block text-sm font-medium text-gray-700">
                        Plan Name
                    </label>
                    <input
                        type="text"
                        id="planname"
                        value={localEditPlanValues.planname}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>

                <div>
                    <label htmlFor="details" className="block text-sm font-medium text-gray-700">
                        Description
                    </label>
                    <input
                        type="text"
                        id="details"
                        value={localEditPlanValues.details}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                    />
                </div>

                <div>
                    <label htmlFor="duration" className="block text-sm font-medium text-gray-700">
                        Duration (days)
                    </label>
                    <input
                        type="text"
                        id="duration"
                        value={localEditPlanValues.duration}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>

                <div>
                    <label htmlFor="price" className="block text-sm font-medium text-gray-700">
                        Price
                    </label>
                    <input
                        type="number"
                        id="price"
                        value={localEditPlanValues.price}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>

                <div>
                    <label htmlFor="number_of_mails_per_day" className="block text-sm font-medium text-gray-700">
                        Number of Mails Per Day
                    </label>
                    <input
                        type="text"
                        id="number_of_mails_per_day"
                        value={localEditPlanValues.number_of_mails_per_day}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>

                <div>
                    <label htmlFor="status" className="block text-sm font-medium text-gray-700">
                        Status
                    </label>
                    <select
                        id="status"
                        value={localEditPlanValues.status}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    >
                        <option value="active">Active</option>
                        <option value="inactive">Inactive</option>
                    </select>
                </div>

                <div>
                    <h3 className="text-lg font-medium text-gray-900">Features</h3>
                    {localEditPlanValues.features.map((feature, index) => (
                        <div key={index} className="mt-4 p-4 border border-gray-300 rounded-md">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700">Name</label>
                                    <input
                                        type="text"
                                        value={feature.name}
                                        onChange={(e) => handleFeatureChange(index, "name", e.target.value)}
                                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700">Identifier</label>
                                    <input
                                        type="text"
                                        value={feature.identifier}
                                        onChange={(e) => handleFeatureChange(index, "identifier", e.target.value)}
                                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700">Count Limit</label>
                                    <input
                                        type="number"
                                        value={feature.count_limit}
                                        onChange={(e) => handleFeatureChange(index, "count_limit", parseInt(e.target.value))}
                                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700">Size Limit</label>
                                    <input
                                        type="number"
                                        value={feature.size_limit}
                                        onChange={(e) => handleFeatureChange(index, "size_limit", parseInt(e.target.value))}
                                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700">Is Active</label>
                                    <select
                                        value={feature.is_active.toString()}
                                        onChange={(e) => handleFeatureChange(index, "is_active", e.target.value === "true")}
                                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                    >
                                        <option value="true">Yes</option>
                                        <option value="false">No</option>
                                    </select>
                                </div>
                                <div className="col-span-2">
                                    <label className="block text-sm font-medium text-gray-700">Description</label>
                                    <input
                                        type="text"
                                        value={feature.description}
                                        onChange={(e) => handleFeatureChange(index, "description", e.target.value)}
                                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                    />
                                </div>
                            </div>
                            <button
                                type="button"
                                onClick={() => removeFeature(index)}
                                className="mt-2 px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
                            >
                                Remove Feature
                            </button>
                        </div>
                    ))}
                    <button
                        type="button"
                        onClick={addFeature}
                        className="mt-4 px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
                    >
                        Add Feature
                    </button>
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