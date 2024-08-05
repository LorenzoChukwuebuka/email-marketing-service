import React, { useState, FormEvent } from "react";
import * as Yup from "yup";
import useAuthStore from "../../../../store/AuthStore";

const ChangePasswordComponent: React.FC = () => {
    const {
        isLoading,
        changePasswordValues,
        changePassword,
        setChangePasswordValues,
    } = useAuthStore();
    const [errors, setErrors] = useState<{ [key: string]: string }>({});

    const validationSchema = Yup.object().shape({
        old_password: Yup.string().required("Current password is required"),
        new_password: Yup.string().required("New password is required"),
        confirm_password: Yup.string()
            .oneOf([Yup.ref("new_password")], "Passwords must match")
            .required("Confirm Password is required"),
    });

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        try {
            await validationSchema.validate(changePasswordValues, {
                abortEarly: false,
            });
            await changePassword();
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

    return (
        <div className="mt-8 mb-5">
            <form className="w-full max-w-xs space-y-4" onSubmit={handleSubmit}>
                <label className="block">
                    <span className="text-medium font-medium">Current Password</span>
                    <input
                        type="password"
                        value={changePasswordValues.old_password}
                        onChange={(event) =>
                            setChangePasswordValues({
                                ...changePasswordValues,
                                old_password: event.target.value,
                            })
                        }
                        className="mt-1 block w-full rounded-md border-2 h-10 border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
                    />
                    {errors.old_password && (
                        <div style={{ color: "red" }}>{errors.old_password}</div>
                    )}
                </label>

                <label className="block">
                    <span className="text-medium font-medium">New Password</span>
                    <input
                        type="password"
                        value={changePasswordValues.new_password}
                        onChange={(event) =>
                            setChangePasswordValues({
                                ...changePasswordValues,
                                new_password: event.target.value,
                            })
                        }
                        className="mt-1 block w-full rounded-md border-2 h-10 border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
                    />
                    {errors.new_password && (
                        <div style={{ color: "red" }}>{errors.new_password}</div>
                    )}
                </label>

                <label className="block">
                    <span className="text-medium font-medium">Confirm New Password</span>
                    <input
                        type="password"
                        value={changePasswordValues.confirm_password}
                        onChange={(event) =>
                            setChangePasswordValues({
                                ...changePasswordValues,
                                confirm_password: event.target.value,
                            })
                        }
                        className="mt-1 block w-full rounded-md border-2 h-10 border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
                    />
                    {errors.confirm_password && (
                        <div style={{ color: "red" }}>{errors.confirm_password}</div>
                    )}
                </label>

                <div className="flex flex-row justify-between">
                    {!isLoading ? (
                        <button
                            type="submit"
                            className="mt-6 bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold py-2 px-4 rounded-md"
                        >
                            Change Password
                        </button>
                    ) : (
                        <button className="mt-6 bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold py-2 px-4 rounded-md" disabled>
                            <span className="flex flex-row items-center">
                                Please wait
                                <span className="loading loading-dots loading-sm"></span>
                            </span>
                        </button>
                    )}
                </div>
            </form>
        </div>
    );
};

export default ChangePasswordComponent;
