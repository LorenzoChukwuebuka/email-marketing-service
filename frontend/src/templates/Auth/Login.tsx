import { Link } from "react-router-dom";
import useAuthStore from "../../store/userstore/AuthStore";
import * as Yup from "yup";
import { useEffect, useState, ChangeEvent, FormEvent } from "react";

const LoginTemplate: React.FC = () => {
    const {
        isLoading,
        loginValues,
        setLoginValues,
        loginUser,
        isLoggedIn,
    } = useAuthStore();

    const [errors, setErrors] = useState<{ [key: string]: string }>({});

    const validationSchema = Yup.object().shape({
        email: Yup.string()
            .required("Email is required")
            .email("Invalid email format"),
        password: Yup.string().required("Password is required"),
    });

    const handleLogin = async (e: FormEvent) => {
        e.preventDefault();

        try {
            await validationSchema.validate(loginValues, { abortEarly: false });
            await loginUser();
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
            console.log(error);
        }
    };

    useEffect(() => {
        if (isLoggedIn) {
            window.location.href = "/user/dash"; // Use `window.location.href` for redirecting
        }
    }, [isLoggedIn]);

    return (
        <>
            <div className="flex justify-center items-center h-screen bg-gray-100">
                <div className="container mx-auto">
                    <h3 className="text-2xl font-bold text-center mb-4">
                        {import.meta.env.VITE_API_NAME}
                    </h3>
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
                                    value={loginValues.email}
                                    onChange={(event: ChangeEvent<HTMLInputElement>) => {
                                        setLoginValues({
                                            ...loginValues,
                                            email: event.target.value,
                                        });
                                    }}
                                    required
                                />
                                {errors.email && (
                                    <div className="text-red-500">{errors.email}</div>
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
                                        type="password"
                                        id="password"
                                        className="block w-full p-2 border border-gray-300 rounded-md"
                                        placeholder=""
                                        value={loginValues.password}
                                        onChange={(event: ChangeEvent<HTMLInputElement>) => {
                                            setLoginValues({
                                                ...loginValues,
                                                password: event.target.value,
                                            });
                                        }}
                                    />
                                    {errors.password && (
                                        <div className="text-red-500">{errors.password}</div>
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
                                    <button className="bg-black text-white py-2 px-4 rounded-md mt-3 hover:bg-gray-800" disabled>
                                        Please wait{" "}
                                        <span className="loading loading-dots loading-sm"></span>
                                    </button>
                                )}
                            </div>
                        </form>
                    </div>
                    <div className="text-center mt-4">
                        <p>
                            <Link
                                to="/auth/forgot-password"
                                className="text-gray-700 hover:underline"
                            >
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
