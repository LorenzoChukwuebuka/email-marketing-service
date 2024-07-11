import { useEffect, useState } from "react";
import usePlanStore from "../../../../../store/admin/planStore";
import EditPlans from "./EditPlans";

const GetAllPlans = () => {
  const { getPlans, planData, selectedId, setSelectedId } = usePlanStore();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedPlan, setSelectedPlan] = useState(null);

  useEffect(() => {
    getPlans();
  }, [getPlans]);

  const handleSelectAll = (e) => {
    if (e.target.checked) {
      const allIds = planData.map((plan) => plan.uuid);
      setSelectedId(allIds);
    } else {
      setSelectedId([]);
    }
  };

  const handleSelect = (uuid) => {
    if (selectedId.includes(uuid)) {
      setSelectedId(selectedId.filter((id) => id !== uuid));
    } else {
      setSelectedId([...selectedId, uuid]);
    }
  };

  const openEditModal = (plan) => {
    setSelectedPlan(plan);
    setIsModalOpen(true);
  };

  const closeEditModal = () => {
    setIsModalOpen(false);
    setSelectedPlan(null);
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
              <th className="py-3 px-4"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {planData.map((plan) => (
              <tr key={plan.uuid}>
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
                <td className="py-4 px-4">
                  <button
                    className="text-gray-400 hover:text-gray-600"
                    onClick={() => openEditModal(plan)}
                  >
                    ✏️
                  </button>
                </td>
              </tr>
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
