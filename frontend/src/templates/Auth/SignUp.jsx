import * as Yup from "yup";
import useAuthStore from "../../store/AuthStore";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";

const SignUpTemplate = () => {
  const [errors, setErrors] = useState({});

  const {
    formValues,
    isLoading,
    setFormValues,
    registerUser,
    userId,
  } = useAuthStore();
  const navigate = useNavigate();

  const validationSchema = Yup.object().shape({
    fullname: Yup.string()
      .required("Name is required")
      .min(5, "Name must be at least 2 characters"),
    email: Yup.string()
      .email("Invalid email format")
      .required("Email is required"),
    company: Yup.string().required(""),
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
      await validationSchema.validate(formValues, { abortEarly: false });
      let registeredData = await registerUser();

      if (registeredData) {
        navigate(
          `/auth/otp-token?email=${encodeURIComponent(registeredData.email)}
          &username=${encodeURIComponent(registeredData.fullname)}
          &userId=${encodeURIComponent(registeredData.userId)}`
        );
      }

      setErrors({});
    } catch (err) {
      const validationErrors = {};
      err.inner.forEach((error) => {
        validationErrors[error.path] = error.message;
      });
      setErrors(validationErrors);
    }
  };

  useEffect(() => {
    if (userId) {
      navigate(
        `/auth/otp-token?email=${encodeURIComponent(formValues.email)}
          &username=${encodeURIComponent(formValues.fullname)}
          &userId=${encodeURIComponent(userId)}`
      );
    }
  });

  return (
    <main className="min-h-screen">
      <div className="bg-[rgb(4,22,43)] h-[15em] pt-2">
        <h1 className="text-center text-2xl font-semibold text-white mt-8">
          MailCrib
        </h1>
      </div>

      <div className="bg-white w-[60%] min-h-auto md:h-[20em] -mt-[7em] mx-auto rounded-btn">
        <h1 className="text-[rgb(4,22,43)] text-2xl font-semibold text-center mt-10">
          Get Started with Mail Crib
        </h1>

        <div className="mt-8 mb-5">
          <form
            className="mx-auto w-full max-w-xs space-y-4"
            onSubmit={handleSubmit}
          >
            <label className="block">
              <span className="text-medium font-medium">Full Name</span>
              <input
                type="text"
                placeholder=""
                value={formValues.fullname}
                onChange={(event) =>
                  setFormValues({
                    ...formValues,
                    fullname: event.target.value,
                  })
                }
                className="mt-1 block w-full border-2  h-10 rounded-md
                 border-gray-300 shadow-sm focus:border-indigo-300
                  focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
              />
              {errors.fullname && (
                <div style={{ color: "red" }}>{errors.fullname}</div>
              )}
            </label>
            <label className="block">
              <span className="text-medium font-medium">Email</span>
              <input
                type="email"
                placeholder=""
                value={formValues.email}
                onChange={(event) =>
                  setFormValues({
                    ...formValues,
                    email: event.target.value,
                  })
                }
                className="mt-1 block w-full border-2  h-10 rounded-md
                 border-gray-300 shadow-sm focus:border-indigo-300
                  focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
              />
              {errors.email && (
                <div style={{ color: "red" }}>{errors.email}</div>
              )}
            </label>
            <label className="block">
              <span className="text-medium font-medium">Company</span>
              <input
                type="text"
                placeholder=""
                value={formValues.company}
                onChange={(event) =>
                  setFormValues({
                    ...formValues,
                    company: event.target.value,
                  })
                }
                className="mt-1 block w-full border-2  h-10 rounded-md
                 border-gray-300 shadow-sm focus:border-indigo-300
                  focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
              />
              {errors.company && (
                <div style={{ color: "red" }}>{errors.company}</div>
              )}
            </label>
            <label className="block">
              <span className="text-medium font-medium">Password</span>
              <input
                type="password"
                placeholder=""
                value={formValues.password}
                onChange={(event) =>
                  setFormValues({
                    ...formValues,
                    password: event.target.value,
                  })
                }
                className="mt-1 block w-full rounded-md border-2 
                 h-10 border-gray-300 shadow-sm
                 focus:border-indigo-300 focus:ring
                  focus:ring-indigo-200 focus:ring-opacity-50"
              />
              {errors.password && (
                <div style={{ color: "red" }}>{errors.password}</div>
              )}
            </label>
            <label className="block">
              <span className="text-medium font-medium">Confirm Password</span>
              <input
                type="password"
                value={formValues.confirmPassword}
                placeholder=""
                onChange={(event) =>
                  setFormValues({
                    ...formValues,
                    confirmPassword: event.target.value,
                  })
                }
                className="mt-1 block w-full rounded-md border-2  h-10 border-gray-300
                 shadow-sm focus:border-indigo-300 focus:ring
                  focus:ring-indigo-200 focus:ring-opacity-50"
              />
              {errors.confirmPassword && (
                <div style={{ color: "red" }}>{errors.confirmPassword}</div>
              )}
            </label>

            <div className="flex flex-row justify-between">
              {!isLoading ? (
                <button
                  type="submit"
                  className="rounded-full pt-1 bg-red-500 text-white p-2"
                >
                  Create Account
                </button>
              ) : (
                <button className="rounded-full pt-1 bg-red-500 text-white p-2">
                  <span className="flex flex-row items-center">
                    Please wait
                    <span className="loading loading-dots loading-sm"></span>
                  </span>
                </button>
              )}

              <button
                type=""
                className="rounded-full pt-1 bg-gray-500 text-white p-2"
              >
                <Link to="/auth/login"> Login </Link>
              </button>
            </div>
          </form>
        </div>
      </div>
    </main>
  );
};

export default SignUpTemplate;
