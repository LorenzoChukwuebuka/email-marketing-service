import { useState } from "react";
import CreatePlan from "../components/Plans/createPlan";
import GetAllPlans from "../components/Plans/AllPlans";
import usePlanStore from "../../../../store/admin/planStore";
const PlansDashTemplate = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { selectedId, deletePlan, getPlans } = usePlanStore();

  const handleDelete = (e) => {
    e.preventDefault();

    deletePlan();

    getPlans();
  };

  return (
    <>
      <h1 className="text-2xl font-bold mb-4">Plan</h1>
      <div className="flex justify-between items-center">
        <div className="space-x-2">
          <button
            className="bg-gray-300  px-4 py-2 rounded-md  transition duration-300"
            onClick={() => setIsModalOpen(true)}
          >
            Create Plan
          </button>

          {selectedId.length > 0 && (
            <button
              className="bg-red-200  px-4 py-2 rounded-md  transition duration-300"
              onClick={() => handleDelete()}
            >
              <span className="text-red-500"> Delete Plan </span>{" "}
              <i className="bi bi-trash text-red-500"></i>
            </button>
          )}
        </div>
      </div>

      <CreatePlan isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />

      <GetAllPlans />
    </>
  );
};

export default PlansDashTemplate;
