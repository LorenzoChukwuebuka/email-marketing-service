import React, { useEffect, useState } from "react";
import { Modal } from "antd";
import usePlanStore from "../store/plan.store";
import { PlanData } from "../interface/plan.interface";

type Feature = {
    name: string;
    description: string;
    value: string;
}

type MailingLimits = {
    daily_limit: number;
    monthly_limit: number;
    max_recipients_per_mail: number;
}

interface EditPlansProps {
    isOpen: boolean;
    onClose: () => void;
    plan: PlanData | null;
}

const EditPlans: React.FC<EditPlansProps> = ({ isOpen, onClose, plan }) => {
    const {
        setEditPlanValues,
        updatePlan,
    } = usePlanStore();

    const [localEditPlanValues, setLocalEditPlanValues] = useState<any>({
        id: '',
        name: '',
        description: '',
        price: 0,
        billing_cycle: 'monthly',
        status: 'active',
        features: [],
        mailing_limits: {
            daily_limit: 0,
            monthly_limit: 0,
            max_recipients_per_mail: 0
        }
    });

    useEffect(() => {
        if (plan) {
            setLocalEditPlanValues({
                id: plan.id,
                name: plan.name,
                description: plan.description,
                price: plan.price,
                billing_cycle: plan.billing_cycle,
                status: plan.status,
                features: plan.features || [],
                mailing_limits: plan.mailing_limits || {
                    daily_limit: 0,
                    monthly_limit: 0,
                    max_recipients_per_mail: 0
                }
            });
        }
    }, [plan]);

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        try {
            setEditPlanValues(localEditPlanValues);
            await updatePlan();
        } catch (error) {
            console.log(error)
        }
        onClose();
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const { id, value } = e.target;
        setLocalEditPlanValues((prev) => ({
            ...prev,
            [id]: id === "price" ? parseFloat(value) : value,
        }));
    };

    const handleMailingLimitChange = (field: keyof MailingLimits, value: string) => {
        setLocalEditPlanValues((prev) => ({
            ...prev,
            mailing_limits: {
                ...prev.mailing_limits,
                [field]: parseInt(value) || 0
            }
        }));
    };

    const handleFeatureChange = (index: number, field: keyof Feature, value: string) => {
        setLocalEditPlanValues((prev) => {
            const updatedFeatures = [...prev.features];
            updatedFeatures[index] = {
                ...updatedFeatures[index],
                [field]: value,
            };
            return { ...prev, features: updatedFeatures };
        });
    };

    const addFeature = () => {
        setLocalEditPlanValues((prev) => ({
            ...prev,
            features: [...prev.features, { name: "", description: "", value: "" }],
        }));
    };

    const removeFeature = (index: number) => {
        setLocalEditPlanValues((prev) => ({
            ...prev,
            features: prev.features.filter((_, i) => i !== index),
        }));
    };

    return (
        <Modal 
            open={isOpen}
            onCancel={onClose}
            footer={null}
            title="Edit Plan"
            width={800}
        >
            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                        Plan Name
                    </label>
                    <input
                        type="text"
                        id="name"
                        value={localEditPlanValues.name}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        required
                    />
                </div>

                <div>
                    <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                        Description
                    </label>
                    <input
                        type="text"
                        id="description"
                        value={localEditPlanValues.description}
                        onChange={handleChange}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                    />
                </div>

                <div className="grid grid-cols-2 gap-4">
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
                        <label htmlFor="billing_cycle" className="block text-sm font-medium text-gray-700">
                            Billing Cycle
                        </label>
                        <select
                            id="billing_cycle"
                            value={localEditPlanValues.billing_cycle}
                            onChange={handleChange}
                            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            required
                        >
                            <option value="monthly">Monthly</option>
                            <option value="yearly">Yearly</option>
                        </select>
                    </div>
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

                {/* Mailing Limits Section */}
                <div>
                    <h3 className="text-lg font-medium text-gray-900 mb-3">Mailing Limits</h3>
                    <div className="grid grid-cols-3 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Daily Limit</label>
                            <input
                                type="number"
                                value={localEditPlanValues.mailing_limits.daily_limit}
                                onChange={(e) => handleMailingLimitChange("daily_limit", e.target.value)}
                                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Monthly Limit</label>
                            <input
                                type="number"
                                value={localEditPlanValues.mailing_limits.monthly_limit}
                                onChange={(e) => handleMailingLimitChange("monthly_limit", e.target.value)}
                                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Max Recipients per Mail</label>
                            <input
                                type="number"
                                value={localEditPlanValues.mailing_limits.max_recipients_per_mail}
                                onChange={(e) => handleMailingLimitChange("max_recipients_per_mail", e.target.value)}
                                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                required
                            />
                        </div>
                    </div>
                </div>

                {/* Features Section */}
                <div>
                    <h3 className="text-lg font-medium text-gray-900">Features</h3>
                    {localEditPlanValues.features.map((feature, index) => (
                        <div key={index} className="mt-4 p-4 border border-gray-300 rounded-md">
                            <div className="grid grid-cols-1 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700">Feature Name</label>
                                    <input
                                        type="text"
                                        value={feature.name}
                                        onChange={(e) => handleFeatureChange(index, "name", e.target.value)}
                                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700">Description</label>
                                    <input
                                        type="text"
                                        value={feature.description}
                                        onChange={(e) => handleFeatureChange(index, "description", e.target.value)}
                                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700">Value</label>
                                    <input
                                        type="text"
                                        value={feature.value}
                                        onChange={(e) => handleFeatureChange(index, "value", e.target.value)}
                                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                        required
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
                        Update Plan
                    </button>
                </div>
            </form>
        </Modal>
    );
};

export default EditPlans;