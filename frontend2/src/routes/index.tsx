import { createBrowserRouter, Navigate, RouterProvider } from 'react-router-dom';
import React, { Suspense } from 'react';
import LoadingSpinnerComponent from '../components/loadingSpinnerComponent';
import Nopage from '../pages/Nopage/Nopage';
import landingPageRoute from './landingpage.route';
import authRoute from '../features/auth/routes/index.router';
import { ProtectedRoute } from './protectedRoute';
import UserDashBoardPage from '../pages/Dashboard/UserDashBoard';
import UserDashIndexTemplate from '../templates/UserDashIndexTemplate';
import contactRoute from '../features/contacts/routes/contacts.route';


const router = createBrowserRouter([
    {
        path: "*",
        element: <Nopage />
    },

    {
        path: '/',  // Root path
        element: <Navigate to="/home" replace />,  // Redirect from root to login
    },
    {
        path: '/home',
        children: [
            ...landingPageRoute,
        ],
    },

    {
        path: "/auth",
        children: [
            ...authRoute
        ]
    },
    {
        path: "app",
        element: <ProtectedRoute />,
        children: [
            {
                element: <UserDashBoardPage />,
                children: [
                    {
                        index: true,
                        element: <UserDashIndexTemplate />
                    },

                    {
                        path: "contacts",
                        children: [
                            ...contactRoute
                        ]
                    }
                ]
            },


        ]

    }


])

const AppRouter: React.FC = () => {
    return (
        <Suspense fallback={<LoadingSpinnerComponent />}>
            <RouterProvider router={router} />
        </Suspense>
    );

};

export default AppRouter;