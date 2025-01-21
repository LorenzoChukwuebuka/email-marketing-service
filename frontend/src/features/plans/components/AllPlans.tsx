import React, {  useMemo, useState } from "react"
import EditPlans from "./EditPlans";
import { usePlansQuery } from "../hooks/usePlanQuery";
import { PlanData } from '../interface/plan.interface';
import usePlanStore from "../store/plan.store";



const GetAllPlans: React.FC = () => { 
    const { selectedId, setSelectedId } = usePlanStore();
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const [selectedPlan, setSelectedPlan] = useState<PlanData | null>(null);
    const [expandedRows, setExpandedRows] = useState<string[]>([]);

    const { data: planAPIData } = usePlansQuery()

    const planData = useMemo(() => planAPIData?.payload || [], [planAPIData])

    const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.checked) {
            const allIds = planData.map((plan) => plan.uuid);
            setSelectedId(allIds);
        } else {
            setSelectedId([]);
        }
    };

    const handleSelect = (uuid: string) => {
        if (selectedId.includes(uuid)) {
            setSelectedId(selectedId.filter((id) => id !== uuid));
        } else {
            setSelectedId([...selectedId, uuid]);
        }
    };

    const openEditModal = (plan: PlanData) => {
        setSelectedPlan(plan);
        setIsModalOpen(true);
    };

    const closeEditModal = () => {
        setIsModalOpen(false);
        setSelectedPlan(null);
    };

    const toggleRowExpansion = (uuid: string) => {
        setExpandedRows(prev =>
            prev.includes(uuid)
                ? prev.filter(id => id !== uuid)
                : [...prev, uuid]
        );
    };

    return (
        <>
            <div className="overflow-x-auto mt-8">
                <table className="min-w-full bg-white">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="py-3 px-4 text-left">
                                <input
                                    type="checkbox"
                                    className="form-checkbox h-4 w-4 text-blue-600"
                                    onChange={handleSelectAll}
                                    checked={selectedId.length === planData.length}
                                />
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Plan Name
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Price
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Duration
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Details
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Number of Mails per day
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Status
                            </th>
                            <th className="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Features
                            </th>
                            <th className="py-3 px-4"></th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                        {planData.map((plan) => (
                            <React.Fragment key={plan.uuid}>
                                <tr>
                                    <td className="py-4 px-4">
                                        <input
                                            type="checkbox"
                                            className="form-checkbox h-4 w-4 text-blue-600"
                                            checked={selectedId.includes(plan.uuid)}
                                            onChange={() => handleSelect(plan.uuid)}
                                        />
                                    </td>
                                    <td className="py-4 px-4">{plan.planname}</td>
                                    <td className="py-4 px-4">{plan.price}</td>
                                    <td className="py-4 px-4">{plan.duration}</td>
                                    <td className="py-4 px-4">{plan.details}</td>
                                    <td className="py-4 px-4">{plan.number_of_mails_per_day}</td>
                                    <td className="py-4 px-4">{plan.status}</td>
                                    <td className="py-4 px-4">
                                        <button
                                            onClick={() => toggleRowExpansion(plan.uuid)}
                                            className="text-blue-600 hover:text-blue-800"
                                        >
                                            {expandedRows.includes(plan.uuid) ? 'Hide' : 'Show'} Features
                                        </button>
                                    </td>
                                    <td className="py-4 px-4">
                                        <button
                                            className="text-gray-400 hover:text-gray-600"
                                            onClick={() => openEditModal(plan)}
                                        >
                                            ✏️
                                        </button>
                                    </td>
                                </tr>
                                {expandedRows.includes(plan.uuid) && (
                                    <tr>
                                        <td colSpan={9}>
                                            <div className="p-4 bg-gray-100">
                                                <h4 className="font-bold mb-2">Features:</h4>
                                                <ul>
                                                    {plan.features.map((feature) => (
                                                        <li key={feature.uuid} className="mb-2">
                                                            <strong>{feature.name}</strong> ({feature.identifier})<br />
                                                            Count Limit: {feature.count_limit}, Size Limit: {feature.size_limit}<br />
                                                            Active: {feature.is_active ? 'Yes' : 'No'}<br />
                                                            Description: {feature.description}
                                                        </li>
                                                    ))}
                                                </ul>
                                            </div>
                                        </td>
                                    </tr>
                                )}
                            </React.Fragment>
                        ))}
                    </tbody>
                </table>
                <div className="mt-4">Total Plans: {planData.length}</div>
            </div>

            <EditPlans
                isOpen={isModalOpen}
                onClose={closeEditModal}
                plan={selectedPlan}
            />
        </>
    );
};

export default GetAllPlans;