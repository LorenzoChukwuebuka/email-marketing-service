import Cookies from "js-cookie";
import { useEffect, useState } from "react";
const UserDashboardTemplate = () => {
  const [userName, setUserName] = useState("");

  useEffect(() => {
    let cookie = Cookies.get("Cookies");
    let user = JSON.parse(cookie)?.details.fullname;
    setUserName(user);
  }, []);

  return (
    <>
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-bold mb-4">Welcome {userName}</h2>
        <p>This is the main content area. You can add your components here.</p>
      </div>
    </>
  );
};

export default UserDashboardTemplate;
