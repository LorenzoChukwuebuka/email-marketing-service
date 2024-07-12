import { useState } from "react";
import * as Yup from "yup";
import useAuthStore from "../../store/AuthStore";



const ResetPasswordTemplate = () => {
  const [errors, setErrors] = useState({});

  const {
    resetPasswordValues,
    setResetPasswordValues,
    isLoading,
    resetPassword,
  } = useAuthStore();

  const validationSchema = Yup.object().shape({
    password: Yup.string()
      .required("Password is required")
      .min(8, "Password must be at least 8 characters"),
    confirmPassword: Yup.string()
      .oneOf([Yup.ref("password"), null], "Passwords must match")
      .required("Confirm Password is required"),
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await validationSchema.validate(resetPasswordValues, {
        abortEarly: false,
      });

      const emailFromURL = new URLSearchParams(location.search).get("email");
      const tokenFromURL = new URLSearchParams(location.search).get("token");

      setResetPasswordValues({
        ...resetPasswordValues,
        token: tokenFromURL,
        email: emailFromURL,
      });

      await resetPassword();
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
      <div className="container mx-auto px-4">
        <div className="max-w-lg mx-auto mt-5">
          <div className="bg-white shadow-md rounded-lg p-8">
            <h3 className="text-2xl font-bold text-center mb-4">{import.meta.env.VITE_API_NAME}</h3>

            <h3 className="text-2xl font-bold text-center mb-4">
              Reset Password
            </h3>

            <p className="text-center mb-2 mt-2 text-gray-400"></p>

            <form onSubmit={handleSubmit}>
              <div className="mb-4">
                <label
                  htmlFor="otpInput"
                  className="block text-sm font-medium text-gray-700"
                >
                  <strong>Password </strong>
                  <span className="text-red-500"> *</span>
                </label>
                <input
                  type="password"
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  onChange={(event) => {
                    setResetPasswordValues({
                      ...resetPasswordValues,
                      password: event.target.value,
                    });
                  }}
                />
                {errors.password && (
                  <div style={{ color: "red" }}>{errors.password}</div>
                )}
              </div>
              <div className="mb-4">
                <label
                  htmlFor="otpInput"
                  className="block text-sm font-medium text-gray-700"
                >
                  <strong>Confirm Password </strong>
                  <span className="text-red-500"> *</span>
                </label>
                <input
                  type="password"
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  onChange={(event) => {
                    setResetPasswordValues({
                      ...resetPasswordValues,
                      confirmPassword: event.target.value,
                    });
                  }}
                />
                {errors.confirmPassword && (
                  <div style={{ color: "red" }}>{errors.confirmPassword}</div>
                )}
              </div>
              <div className="text-center">
                {!isLoading ? (
                  <button
                    className="w-full bg-gray-800 text-white py-2 px-4 rounded-md hover:bg-gray-700"
                    type="submit"
                  >
                    Submit
                  </button>
                ) : (
                  <button className="w-full bg-gray-800 text-white py-2 px-4 rounded-md hover:bg-gray-700">
                    Please wait{" "}
                    <span className="loading loading-dots loading-sm"></span>
                  </button>
                )}
              </div>
            </form>
          </div>
        </div>
      </div>
    </>
  );
};

export default ResetPasswordTemplate;
