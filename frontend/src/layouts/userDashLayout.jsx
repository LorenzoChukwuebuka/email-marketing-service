import { useEffect, useState } from "react";
import { Link, Outlet, useLocation } from "react-router-dom";
import Cookies from "js-cookie";
const UserDashLayout = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const [userName, setUserName] = useState("");
  const location = useLocation();

  const getLinkClassName = (path) => {
    const baseClass = "mb-4 text-center text-lg font-semibold";
    const activeClass = "text-white bg-[rgb(56,68,94)] p-2 px-2 rounded-md";
    const inactiveClass =
      "text-gray-300 hover:text-white hover:bg-[rgb(56,68,94)] px-2 p-2 rounded-md";
    return `${baseClass} ${
      location.pathname === path ? activeClass : inactiveClass
    }`;
  };

  useEffect(() => {
    let cookie = Cookies.get("Cookies");
    let user = JSON.parse(cookie)?.details?.fullname;
    setUserName(user);
  }, []);

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <div
        className={`${
          sidebarOpen ? "w-48" : "w-0"
        } transition-all duration-300 bg-[rgb(26,46,68)]`}
      >
        {sidebarOpen && (
          <nav className="p-4 text-white">
            <h2 className="text-xl font-bold mt-4 text-center mb-4">
              Mail Crib
            </h2>
            <ul className="mt-12 w-full">
              <li className={getLinkClassName("/user/dash")}>
                <Link to="/user/dash">
                  <i className="bi bi-house"></i> Home{" "}
                </Link>
              </li>
              <li className={getLinkClassName("")}>
                <Link to="">
                  <i className="bi bi-bar-chart"></i> Analytics{" "}
                </Link>
              </li>
              <li className={getLinkClassName("")}>
                <Link to="">
                  {" "}
                  <i className="bi bi-gear"></i> Settings{" "}
                </Link>
              </li>
            </ul>
          </nav>
        )}
      </div>

      {/* Main content */}
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <header className="bg-white h-12 p-4 flex justify-between items-center">
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="text-gray-500 hover:text-gray-700"
          >
            <span style={{ fontSize: "24px" }}>{sidebarOpen ? "≡" : "☰"}</span>
          </button>
          <h1 className="text-xl font-semibold">Home </h1>
          <button className="hover:bg-blue-200 hover:rounded-btn hover:text-blue-500 font-semibold p-1">
            Usage and Plans
          </button>
          <span className="w-auto h-auto"> {userName} </span>
        </header>

        {/* Content area */}
        <main className="flex-1 p-6 overflow-auto">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default UserDashLayout;
