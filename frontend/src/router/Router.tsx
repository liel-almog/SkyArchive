import { Navigate, RouterProvider, createBrowserRouter } from "react-router-dom";
import { BaseErrorBoundary } from "../components/common/BaseErrorBoundary";
import { useAuthContext } from "../context/AuthContext/useAuthProvider";
import { Home } from "../pages/Home";
import { LoginForm } from "../pages/Login";
import { SignupForm } from "../pages/Signup";
import { FilesManager } from "../pages/FilesManager";

const authenticatedRouter = createBrowserRouter([
  {
    path: "/",
    errorElement: <BaseErrorBoundary />,
    element: <Home />,
    children: [
      {
        path: "/",
        element: <FilesManager />,
        index: true,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/" />,
  },
]);

const unauthenticatedRouter = createBrowserRouter([
  {
    path: "/login",
    errorElement: <BaseErrorBoundary />,
    element: <LoginForm />,
  },
  {
    path: "/signup",
    errorElement: <BaseErrorBoundary />,
    element: <SignupForm />,
  },
  {
    path: "*",
    element: <Navigate to="/login" />,
  },
]);

export const Router = () => {
  const { authenticated } = useAuthContext();
  const router = authenticated ? authenticatedRouter : unauthenticatedRouter;

  return <RouterProvider router={router} />;
};
