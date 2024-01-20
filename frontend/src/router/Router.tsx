import {
  Navigate,
  RouterProvider,
  createBrowserRouter,
} from "react-router-dom";
import { BaseErrorBoundary } from "../components/common/BaseErrorBoundary";
import { useAuthContext } from "../context/AuthContext/useAuthProvider";
import { Home } from "../pages/Home";
import { LoginForm } from "../pages/Login";
import { UploadFiles } from "../pages/UploadFiles";

const authenticatedRouter = createBrowserRouter([
  {
    path: "/",
    errorElement: <BaseErrorBoundary />,
    element: <Home />,
    children: [
      {
        path: "/",
        element: <UploadFiles />,
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
    path: "/",
    errorElement: <BaseErrorBoundary />,
    element: <LoginForm />,
  },
  {
    path: "*",
    element: <Navigate to="/" />,
  },
]);

export const Router = () => {
  const { authenticated } = useAuthContext();
  const router = authenticated ? authenticatedRouter : unauthenticatedRouter;

  return <RouterProvider router={router} />;
};
