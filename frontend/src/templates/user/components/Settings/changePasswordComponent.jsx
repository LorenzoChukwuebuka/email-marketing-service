import useAuthStore from "../../../../store/AuthStore";

const ChangePasswordComponent = () => {
  const { isLoading } = useAuthStore();
  return (
    <>
      <div className="mt-8 mb-5">
        <form className=" w-full max-w-xs space-y-4">
          <label className="block">
            <span className="text-medium font-medium">Current Password</span>
            <input
              type="password"
              placeholder=""
              className="mt-1 block w-full rounded-md border-2 
                 h-10 border-gray-300 shadow-sm
                 focus:border-indigo-300 focus:ring
                  focus:ring-indigo-200 focus:ring-opacity-50"
            />
          </label>

          <label className="block">
            <span className="text-medium font-medium">New Password</span>
            <input
              type="password"
              placeholder=""
              className="mt-1 block w-full rounded-md border-2 
                 h-10 border-gray-300 shadow-sm
                 focus:border-indigo-300 focus:ring
                  focus:ring-indigo-200 focus:ring-opacity-50"
            />
          </label>
          <label className="block">
            <span className="text-medium font-medium">
              Confirm New Password
            </span>
            <input
              type="password"
              placeholder=""
              className="mt-1 block w-full rounded-md border-2  h-10 border-gray-300
                 shadow-sm focus:border-indigo-300 focus:ring
                  focus:ring-indigo-200 focus:ring-opacity-50"
            />
          </label>

          <div className="flex flex-row justify-between">
            {!isLoading ? (
              <button
                type="submit"
                className="mt-6 bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold py-2 px-4 rounded-md"
              >
                Change Password
              </button>
            ) : (
              <button className="mt-6 bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold py-2 px-4 rounded-md">
                <span className="flex flex-row items-center">
                  Please wait
                  <span className="loading loading-dots loading-sm"></span>
                </span>
              </button>
            )}
          </div>
        </form>
      </div>
    </>
  );
};

export default ChangePasswordComponent;
