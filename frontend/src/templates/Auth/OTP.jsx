import useAuthStore from "../../store/AuthStore";

import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import eventBus from "../../utils/eventBus";

const OTPTemplate = () => {
  const navigate = useNavigate();

  const {
    setOTPValue,
    verifyUser,
    isLoading,
    isVerified,
    resendOTP,
  } = useAuthStore();

  const handleVerify = async (e) => {
    e.preventDefault();

    const tokenFromURL = new URLSearchParams(location.search).get("token");
    setOTPValue({ token: tokenFromURL });
    await verifyUser();
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
      <div className="container mx-auto mt-[10em] px-4">
        <div className="max-w-lg mx-auto mt-5">
          <div className="bg-white shadow-md rounded-lg p-8">
            <h3 className="text-2xl font-bold  text-center mb-4">MailCrib</h3>

            <h3 className="text-2xl font-bold text-center mb-4">
              Verify Email
            </h3>

            <form onSubmit={handleVerify}>
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
