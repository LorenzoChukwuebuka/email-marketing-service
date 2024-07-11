import { Modal } from "../../../../../components";
import usePlanStore from "../../../../../store/admin/planStore";

import PropTypes from "prop-types";

const CreatePlan = ({ isOpen, onClose }) => {
  const { setPlanValues, planValues, createPlan, getPlans } = usePlanStore();
  const handleSubmit = async (e) => {
    e.preventDefault();
    await createPlan();
    await getPlans();
    onClose();
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Create Plan"
      width="w-[30em]"
    >
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
            value={planValues.planname}
            onChange={(event) =>
              setPlanValues({ ...planValues, planname: event.target.value })
            }
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            required
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
            value={planValues.details}
            onChange={(event) =>
              setPlanValues({
                ...planValues,
                details: event.target.value,
              })
            }
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
            value={planValues.duration}
            onChange={(event) =>
              setPlanValues({
                ...planValues,
                duration: event.target.value,
              })
            }
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            required
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
            value={planValues.price}
            onChange={(event) =>
              setPlanValues({ ...planValues, price: event.target.value })
            }
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            required
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
            value={planValues.number_of_mails_per_day}
            onChange={(event) => {
              setPlanValues({
                ...planValues,
                number_of_mails_per_day: event.target.value,
              });
            }}
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            required
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

CreatePlan.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  children: PropTypes.any,
};

export default CreatePlan;
