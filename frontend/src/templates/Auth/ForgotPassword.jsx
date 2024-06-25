import useAuthStore from "../../store/AuthStore";
import * as Yup from "yup";
import { useState } from "react";
const ForgotPasswordTemplate = () => {
  const {
    isLoading,
    forgetPasswordValues,
    setForgetPasswordValues,
    forgotPass,
  } = useAuthStore();

  const [errors, setErrors] = useState({});

  const validationSchema = Yup.object().shape({
    email: Yup.string()
      .email("Invalid email format")
      .required("Email is required"),
  });

  const handleSubmit = async (e) => {
    try {
      e.preventDefault();

      await validationSchema.validate(forgetPasswordValues, {
        abortEarly: false,
      });

      await forgotPass();
      setErrors({});
    } catch (error) {
      const validationErrors = {};
      error.inner.forEach((error) => {
        validationErrors[error.path] = error.message;
      });
      setErrors(validationErrors);
    }
  };

  const resendEmail = async (e) => {
    e.preventDefault();
  };

  return (
    <>
      <div className="container mx-auto px-4">
        <div className="max-w-lg mx-auto mt-5">
          <div className="bg-white shadow-md rounded-lg p-8">
            <h3 className="text-2xl font-bold  text-center mb-4">MailCrib</h3>

            <h3 className="text-2xl font-bold text-center mb-4">
              Forgot Password
            </h3>

            <p className="text-center mb-2 mt-2 text-gray-400">
              You will receive an email if your mail is registered with us{" "}
            </p>

            <form onSubmit={handleSubmit}>
              <div className="mb-4">
                <label
                  htmlFor="otpInput"
                  className="block text-sm font-medium text-gray-700"
                >
                  <strong>Email</strong>
                  <span className="text-red-500"> *</span>
                </label>
                <input
                  type="text"
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  id="otpInput"
                  placeholder="Enter your registered email"
                  value={forgetPasswordValues.email}
                  onChange={(event) =>
                    setForgetPasswordValues({
                      ...forgetPasswordValues,
                      email: event.target.value,
                    })
                  }
                />
                 {errors.email && (
                <div  className="text-red-500 text-center">{errors.email}</div>
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

            <div className="text-center mt-4">
              <p>
                Didn`t receive the OTP?
                <button
                  className="text-blue-600 hover:underline ml-2"
                  type="submit"
                >
                  Resend Email
                </button>
              </p>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default ForgotPasswordTemplate;
