import { Navigate, RouterProvider, createBrowserRouter } from "react-router-dom";
import { BaseErrorBoundary } from "../components/common/BaseErrorBoundary";
import { Home } from "../pages/Home";
import { UploadFiles } from "../pages/UploadFiles";

const router = createBrowserRouter([
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

export const Router = () => {
  return <RouterProvider router={router} />;
};
