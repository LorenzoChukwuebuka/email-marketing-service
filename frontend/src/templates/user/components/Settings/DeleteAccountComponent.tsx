const DeleteAccountComponent:React.FC = () => {
  const handleDeleteAccount = () => {
    // Logic for deleting the account
    alert("Account deletion initiated!");
  };

  return (
    <div className="max-w-2xl  p-6 bg-white mx-auto rounded-md">
      <h1 className="text-2xl font-bold text-red-600 mb-4">Deleting Account</h1>
      <p className="text-gray-600 mb-6">
        Pay special attention when entering this area!
      </p>
      <div className="border border-red-400 p-4 rounded-md bg-red-50">
        <div className="flex items-center mb-4">
          {/* <ExclamationTriangleIcon className="h-8 w-8 text-red-600 mr-2" /> */}

          <i className="bi bi-info-circle text-red-600 mr-2 font-semibold text-xl"></i>
          <p className="text-gray-700">
            Deleting your account removes all projects, inboxes, domains, and
            data associated with the account. If you have multiple accounts, you
            can still access other accounts and their data. Once you click
            &quot; Delete Account &quot;, we send you a confirmation email.
          </p>
        </div>
        <button
          onClick={handleDeleteAccount}
          className="bg-transparent hover:bg-red-600 text-red-600 font-semibold hover:text-white py-2 px-4 border border-red-600 hover:border-transparent rounded"
        >
          Delete Account
        </button>
        <p className="text-gray-600 mt-4">
          <strong>Note:</strong> To delete your profile completely, go to the My
          Profile menu in the upper right corner.
        </p>
      </div>
    </div>
  );
};

export default DeleteAccountComponent;
