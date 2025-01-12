import React from 'react';
import { Outlet, Navigate } from "react-router-dom";
import useAuthStore from '../features/auth/store/auth.store';

export const ProtectedRoute: React.FC = () => {
    const token = useAuthStore(((state) => state.token));
    return token ? <Outlet /> : <Navigate to="/auth/login" />;
};


export const AdminProctedRoute: React.FC = () => {
    return <> hello world </>
}