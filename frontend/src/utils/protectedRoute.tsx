import React from 'react';
import { Outlet, Navigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export const ProtectedRoute: React.FC = () => {
    const { token } = useAuth();
    return token ? <Outlet /> : <Navigate to="/auth/login" />;
};