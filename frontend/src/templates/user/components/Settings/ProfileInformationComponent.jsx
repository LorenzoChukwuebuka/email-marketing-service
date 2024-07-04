import   { useState } from 'react';

const ProfileInformationComponent = () => {
  const [formData, setFormData] = useState({
    username: 'lawrenceobi2@gmail.com',
    firstName: 'Lawrence',
    lastName: 'Obi',
    email: 'lawrenceobi2@gmail.com',
    phoneNumber: '08134514639',
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleEditInformation = () => {
    // Logic for editing the information
    alert('Edit information functionality!');
  };

  return (
    <div className="max-w-2xl ml-5 p-6 ">
      <h1 className="text-2xl font-bold text-gray-800 mb-4">Profile Information</h1>
      <p className="text-gray-600 mb-6">
        This is the information we have associated with your Seemailer profile, which you can use to access multiple Seemailer accounts. 
        <strong> All contact information is kept strictly confidential.</strong>
      </p>
      <div className="space-y-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="block text-gray-700 font-bold mb-1">Username</label>
            <input
              type="text"
              name="username"
              value={formData.username}
              onChange={handleChange}
              className="w-full bg-gray-100 p-2 rounded-md"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-gray-700 font-bold mb-1">Firstname</label>
              <input
                type="text"
                name="firstName"
                value={formData.firstName}
                onChange={handleChange}
                className="w-full bg-gray-100 p-2 rounded-md"
              />
            </div>
            <div>
              <label className="block text-gray-700 font-bold mb-1">Lastname</label>
              <input
                type="text"
                name="lastName"
                value={formData.lastName}
                onChange={handleChange}
                className="w-full bg-gray-100 p-2 rounded-md"
              />
            </div>
          </div>
        </div>
        <div>
          <label className="block text-gray-700 font-bold mb-1">Email</label>
          <input
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            className="w-full bg-gray-100 p-2 rounded-md"
          />
        </div>
        <div>
          <label className="block text-gray-700 font-bold mb-1">Phone Number</label>
          <input
            type="text"
            name="phoneNumber"
            value={formData.phoneNumber}
            onChange={handleChange}
            className="w-full bg-gray-100 p-2 rounded-md"
          />
        </div>
      </div>
      <button
        onClick={handleEditInformation}
        className="mt-6 bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold py-2 px-4 rounded-md"
      >
        Edit Information
      </button>
    </div>
  );
};

export default ProfileInformationComponent;
