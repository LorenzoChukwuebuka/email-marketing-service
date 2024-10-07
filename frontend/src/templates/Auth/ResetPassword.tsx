import { useState, ChangeEvent, FormEvent } from "react";
import * as Yup from "yup";
import useAuthStore from "../../store/userstore/AuthStore";
import renderApiName from "../../utils/name";


const ResetPasswordTemplate: React.FC = () => {
    const [errors, setErrors] = useState<{ [key: string]: string }>({});

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
            .oneOf([Yup.ref("password")], "Passwords must match")  
            .required("Confirm Password is required"),
    });

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        try {
            await validationSchema.validate(resetPasswordValues, {
                abortEarly: false,
            });

            const emailFromURL = new URLSearchParams(window.location.search).get("email");
            const tokenFromURL = new URLSearchParams(window.location.search).get("token");

            setResetPasswordValues({
                ...resetPasswordValues,
                token: tokenFromURL || "",
                email: emailFromURL || "",
            });

            await resetPassword();
        } catch (error) {
            if (error instanceof Yup.ValidationError) {
                const validationErrors: { [key: string]: string } = {};
                error.inner.forEach((err) => {
                    if (err.path) {
                        validationErrors[err.path] = err.message;
                    }
                });
                setErrors(validationErrors);
            }

        }
    };

    return (
        <div className="container mx-auto px-4">
            <div className="max-w-lg mx-auto mt-5">
                <div className="bg-white shadow-md rounded-lg p-8">
                 {renderApiName()}

                    <h3 className="text-2xl font-bold text-center mb-4">
                        Reset Password
                    </h3>

                    <p className="text-center mb-2 mt-2 text-gray-400"></p>

                    <form onSubmit={handleSubmit}>
                        <div className="mb-4">
                            <label
                                htmlFor="password"
                                className="block text-sm font-medium text-gray-700"
                            >
                                <strong>Password </strong>
                                <span className="text-red-500"> *</span>
                            </label>
                            <input
                                type="password"
                                id="password"
                                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                onChange={(event: ChangeEvent<HTMLInputElement>) => {
                                    setResetPasswordValues({
                                        ...resetPasswordValues,
                                        password: event.target.value,
                                    });
                                }}
                                value={resetPasswordValues.password}
                            />
                            {errors.password && (
                                <div style={{ color: "red" }}>{errors.password}</div>
                            )}
                        </div>
                        <div className="mb-4">
                            <label
                                htmlFor="confirmPassword"
                                className="block text-sm font-medium text-gray-700"
                            >
                                <strong>Confirm Password </strong>
                                <span className="text-red-500"> *</span>
                            </label>
                            <input
                                type="password"
                                id="confirmPassword"
                                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                onChange={(event: ChangeEvent<HTMLInputElement>) => {
                                    setResetPasswordValues({
                                        ...resetPasswordValues,
                                        confirmPassword: event.target.value,
                                    });
                                }}
                                value={resetPasswordValues.confirmPassword}
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
                                <button className="w-full bg-gray-800 text-white py-2 px-4 rounded-md hover:bg-gray-700" disabled>
                                    Please wait{" "}
                                    <span className="loading loading-dots loading-sm"></span>
                                </button>
                            )}
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default ResetPasswordTemplate;
