import * as Yup from "yup";
import { useEffect, useState, FormEvent } from "react";
import useAdminAuthStore from "../../../store/admin/AdminAuthStore";



interface ValidationErrors {
    [key: string]: string;
}

const AdminLoginTemplate: React.FC = () => {
    const {
        isLoading,
        loginValues,
        setLoginValues,
        loginAdmin,
        isLoggedIn,
    } = useAdminAuthStore();

    const [errors, setErrors] = useState<ValidationErrors>({});

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
            await loginAdmin();
        } catch (error: any) {
            const validationErrors: ValidationErrors = {};
            if (error.name === "ValidationError") {
                console.log("Validation errors:", error.errors);
            }
            if (error.inner) {
                error.inner.forEach((err: any) => {
                    validationErrors[err.path] = err.message;
                });
            }
            setErrors(validationErrors);
            console.log(error);
        }
    };

    useEffect(() => {
        if (isLoggedIn) {
            location.href = "/zen/dash";
        }
    }, [isLoggedIn]);

    return (
        <div className="flex justify-center items-center h-screen bg-gray-100">
            <div className="container mx-auto">
                <h3 className="text-2xl font-bold text-center mb-4">MailCrib</h3>
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
                                    type="password"
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
            </div>
        </div>
    );
};

export default AdminLoginTemplate;
