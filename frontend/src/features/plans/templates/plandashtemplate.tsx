import { FormEvent, useState } from "react";
import { Helmet, HelmetProvider } from "react-helmet-async";
import usePlanStore from "../store/plan.store";
import GetAllPlans from "../components/AllPlans";
import CreatePlan from "../components/createPlan";



const PlansDashTemplate: React.FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { selectedId, deletePlan } = usePlanStore();

    const handleDelete = async (e: FormEvent<HTMLButtonElement>) => {
        try {
            e.preventDefault();

            deletePlan();
        } catch (error) {
            console.log(error)
        }



    };

    return (
        <>

            <HelmetProvider>
                <Helmet title="Plans" />

                <h1 className="text-2xl font-bold mb-4">Plan</h1>
                <div className="flex justify-between items-center">
                    <div className="space-x-2">
                        <button
                            className="bg-gray-300 px-4 py-2 rounded-md transition duration-300"
                            onClick={() => setIsModalOpen(true)}
                        >
                            Create Plan
                        </button>

                        {selectedId.length > 0 && (
                            <button
                                className="bg-red-200 px-4 py-2 rounded-md transition duration-300"
                                onClick={(e) => handleDelete(e)}
                            >
                                <span className="text-red-500"> Delete Plan </span>
                                <i className="bi bi-trash text-red-500"></i>
                            </button>
                        )}
                    </div>
                </div>

                <CreatePlan isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />

                <GetAllPlans />

            </HelmetProvider>
        </>
    );
};

export default PlansDashTemplate;
