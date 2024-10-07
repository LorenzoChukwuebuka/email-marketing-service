import React, { useState, ChangeEvent, FormEvent } from "react";
import * as Yup from "yup";
import useAuthStore from "../../store/userstore/AuthStore";
import { Link } from "react-router-dom";

const ForgotPasswordTemplate: React.FC = () => {
    const {
        isLoading,
        forgetPasswordValues,
        setForgetPasswordValues,
        forgotPass,
    } = useAuthStore();

    const [errors, setErrors] = useState<{ [key: string]: string }>({});

    const validationSchema = Yup.object().shape({
        email: Yup.string()
            .email("Invalid email format")
            .required("Email is required"),
    });

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        try {
            await validationSchema.validate(forgetPasswordValues, {
                abortEarly: false,
            });

            await forgotPass();
            setErrors({});
        } catch (error) {
            if (error instanceof Yup.ValidationError) {
                const validationErrors: { [key: string]: string } = {};
                error.inner.forEach((err) => {
                    validationErrors[err.path || ""] = err.message;
                });
                setErrors(validationErrors);
            }
        }
    };

    const apiName = import.meta.env.VITE_API_NAME;
    const firstFourLetters = apiName.slice(0, 4);
    const remainingLetters = apiName.slice(4);

    return (
        <>
            <div className="container mx-auto px-4">
                <div className="max-w-lg mx-auto mt-[10em]">
                    <h3 className="text-2xl font-bold text-center mb-4">
                        <span>{firstFourLetters}</span>
                        <span className="text-blue-500">{remainingLetters}</span> <i className="bi bi-mailbox2-flag text-blue-500"></i>
                    </h3>
                    <div className="bg-white shadow-md rounded-lg p-8">
                        <h3 className="text-2xl font-bold text-center mb-4">
                            Forgot Password
                        </h3>
                        <p className="text-center mb-2 mt-2 text-gray-400">
                            You will receive an email if your mail is registered with us
                        </p>

                        <form onSubmit={handleSubmit}>
                            <div className="mb-4">
                                <label
                                    htmlFor="emailInput"
                                    className="block text-sm font-medium text-gray-700"
                                >
                                    <strong>Email</strong>
                                    <span className="text-red-500"> *</span>
                                </label>
                                <input
                                    type="text"
                                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                    id="emailInput"
                                    placeholder="Enter your registered email"
                                    value={forgetPasswordValues.email}
                                    onChange={(event: ChangeEvent<HTMLInputElement>) =>
                                        setForgetPasswordValues({
                                            ...forgetPasswordValues,
                                            email: event.target.value,
                                        })
                                    }
                                />
                                {errors.email && (
                                    <div className="text-red-500 text-center">{errors.email}</div>
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
                                    <button className="w-full bg-gray-800 text-white py-2 px-4 rounded-md hover:bg-gray-700" disabled>
                                        Please wait{" "}
                                        <span className="loading loading-dots loading-sm"></span>
                                    </button>
                                )}
                            </div>
                        </form>
                        <div className="text-center mt-4">
                            <p>
                                Remember your password?{" "}
                                <Link to="/auth/login" className="text-gray-700 hover:underline ml-4">
                                    Login
                                </Link>
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
};

export default ForgotPasswordTemplate;
