import { Link } from "react-router-dom";
import useAuthStore from "../../store/AuthStore";
import * as Yup from "yup";
import { useState } from "react";

const LoginTemplate = () => {
  const { isLoading, loginValues, setLoginValues, loginUser } = useAuthStore();
  const [errors, setErrors] = useState({});

  const validationSchema = Yup.object().shape({
    password: Yup.string().required("token is required"),
    email: Yup.string()
      .required("password is required")
      .email("invalid email format"),
  });

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      await validationSchema.validate(loginValues, { abortEarly: false });
      await loginUser();
    } catch (error) {
      const validationErrors = {};
      error.inner.forEach((error) => {
        validationErrors[error.path] = error.message;
      });
      setErrors(validationErrors);
    }
  };

  return (
    <>
      <div className="flex justify-center items-center h-screen bg-gray-100">
        <div className="container mx-auto">
          <h3 className="text-2xl font-bold  text-center mb-4">MailCrib</h3>
          <div className="bg-white shadow-lg rounded-lg max-w-lg mx-auto mt-2 p-6">
            <h3 className="text-2xl font-semibold text-center mb-4">Log in</h3>
            <form onSubmit={handleLogin}>
              <div className="mb-4">
                <label
                  htmlFor="email"
                  className="block text-gray-700 font-bold mb-2"
                >
                  Email <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  id="email"
                  className="block w-full p-2 border border-gray-300 rounded-md"
                  placeholder=""
                  onChange={(event) => {
                    setLoginValues({
                      ...loginValues,
                      email: event.target.value,
                    });
                  }}
                  required
                />
                {errors.email && (
                  <div style={{ color: "red" }}>{errors.email}</div>
                )}
              </div>
              <div className="mb-4">
                <label
                  htmlFor="password"
                  className="block text-gray-700 font-bold mb-2"
                >
                  Password <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                  <input
                    id="password"
                    className="block w-full p-2 border border-gray-300 rounded-md"
                    placeholder=""
                    onChange={(event) => {
                      setLoginValues({
                        ...loginValues,
                        password: event.target.value,
                      });
                    }}
                  />
                  {errors.password && (
                    <div style={{ color: "red" }}>{errors.password}</div>
                  )}
                </div>
              </div>
              <div className="text-center">
                {!isLoading ? (
                  <button
                    className="bg-black text-white py-2 px-4 rounded-md mt-3 hover:bg-gray-800"
                    type="submit"
                  >
                    Login
                  </button>
                ) : (
                  <button className="bg-black text-white py-2 px-4 rounded-md mt-3 hover:bg-gray-800">
                    Please wait{" "}
                    <span className="loading loading-dots loading-sm"></span>
                  </button>
                )}
              </div>
            </form>
          </div>
          <div className="text-center mt-4">
            <p>
              <Link to="/auth/forgot-password" className="text-gray-700 hover:underline">
                Forgot Password
              </Link>
              <Link
                to="/auth/sign-up"
                className="text-gray-700 hover:underline ml-4"
              >
                Create Account
              </Link>
            </p>
          </div>
        </div>
      </div>
    </>
  );
};

export default LoginTemplate;
