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
import campaignRoute from '../features/campaign/routes/campaign.route';
import NotificationTemplate from '../templates/NotificationTemplate';
import emailTemplateRoute from '../features/email-templates/routes/email-template.route';
import EditorRouter from '../features/editors/routes/editor.route';
import analyticsRoute from '../features/analytics/routes/analytics.route';
import settingsRoute from '../features/settings/routes/settings.route';
import VerifySenderTemplate from '../features/settings/templates/verifySenderTemplate';
import billingRoutes from '../features/billing/routes/billing.route';
import adminAuthRoute from '../features/auth/routes/admin.auth.route';
import supportRoute from '../features/support/routes/support.route';
import AdminDashTemplate from '../templates/AdminDashTemplate';
import AdminDashIndexTemplate from '../templates/AdminDashIndex';
import planRoute from '../features/plans/routes/plan.route';

const router = createBrowserRouter([
    {
        path: "*",
        element: <Nopage />
    },

    {
        path: "/404",
        element: <Nopage />
    },
    {
        path: '/',  // Root path
        element: <Navigate to="/home" replace />,  // Redirect from root to  landingpage...
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

    //admin auth route
    {
        path: "/next",
        children: [
            ...adminAuthRoute
        ]
    },

    //for admin routes

    {
        path: "/zen",
        element: <ProtectedRoute />,
        children: [
            {
                element: <AdminDashTemplate />,
                children: [
                    {
                        index: true,
                        element: <AdminDashIndexTemplate />
                    },
                    {
                        path: "plan",
                        children: [
                            ...planRoute
                        ]
                    }
                ]
            }
        ]
    },

    //editor route

    {
        path: "/editor/:editorType",
        element: <ProtectedRoute />,
        children: [
            {
                index: true,
                element: <EditorRouter />
            }
        ]
    },
    {
        path: "/verifysender",
        element: <VerifySenderTemplate />
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
                    },
                    {
                        path: "campaign",
                        children: [
                            ...campaignRoute
                        ]
                    },
                    {
                        path: 'notifications',
                        element: <NotificationTemplate />
                    },
                    {
                        path: "templates",
                        children: [
                            ...emailTemplateRoute
                        ]
                    },
                    {
                        path: "billing",
                        children: [
                            ...billingRoutes
                        ]
                    },
                    {
                        path: "analytics",
                        children: [
                            ...analyticsRoute
                        ]
                    },
                    {
                        path: "settings",
                        children: [
                            ...settingsRoute
                        ]
                    },
                    {
                        path: "support",
                        children: [
                            ...supportRoute
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