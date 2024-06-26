import useAuthStore from "../../store/AuthStore";
import * as Yup from "yup";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import eventBus from "../../utils/eventBus";

const OTPTemplate = () => {
  const validationSchema = Yup.object().shape({
    token: Yup.string()
      .required("token is required")
      .min(8, "token must be at least 8 characters"),
  });

  const navigate = useNavigate();

  const {
    otpValue,
    setOTPValue,
    verifyUser,
    isLoading,
    isVerified,
    resendOTP,
  } = useAuthStore();
  const [errors, setErrors] = useState({});

  const handleVerify = async (e) => {
    e.preventDefault();

    try {
      await validationSchema.validate(otpValue, { abortEarly: false });
      await verifyUser();
      setErrors({});
    } catch (err) {
      const validationErrors = {};
      err.inner.forEach((error) => {
        validationErrors[error.path] = error.message;
      });
      setErrors(validationErrors);
    }
  };

  const handleResendOTP = async () => {
    eventBus.emit(
      "message",
      "You have successfully resent the token, kindly check your mail"
    );
    const emailFromURL = new URLSearchParams(location.search).get("email");
    const usernameFromURL = new URLSearchParams(location.search).get(
      "username"
    );
    const userIdFromURL = new URLSearchParams(location.search).get("userId");

    let data = {
      user_id: userIdFromURL,
      username: usernameFromURL,
      email: emailFromURL,
      otp_type: "emailVerify",
    };

    await resendOTP(data);
  };

  useEffect(() => {
    if (isVerified) {
      setTimeout(() => {
        navigate("/auth/login");
      }, 1500);
    }
  });

  return (
    <>
      <div className="container mx-auto px-4">
        <div className="max-w-lg mx-auto mt-5">
          <div className="bg-white shadow-md rounded-lg p-8">
            <h3 className="text-2xl font-bold  text-center mb-4">MailCrib</h3>

            <h3 className="text-2xl font-bold text-center mb-4">
              Verify Email
            </h3>

            <form onSubmit={handleVerify}>
              <div className="mb-4">
                <label
                  htmlFor="otpInput"
                  className="block text-sm font-medium text-gray-700"
                >
                  <strong>OTP</strong>
                  <span className="text-red-500"> *</span>
                </label>
                <input
                  type="text"
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  id="otpInput"
                  placeholder="Enter your OTP"
                  value={otpValue.token}
                  onChange={(event) =>
                    setOTPValue({
                      ...otpValue,
                      token: event.target.value,
                    })
                  }
                  maxLength="8"
                />
                {errors.token && (
                  <div style={{ color: "red" }}>{errors.token}</div>
                )}
              </div>
              <div className="text-center">
                {!isLoading ? (
                  <button
                    className="w-full bg-gray-800 text-white py-2 px-4 rounded-md hover:bg-gray-700"
                    type="submit"
                  >
                    Verify Email
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
                  onClick={handleResendOTP}
                >
                  Resend OTP
                </button>
              </p>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default OTPTemplate;
