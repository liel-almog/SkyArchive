import { faTriangleExclamation } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useRouteError } from "react-router-dom";
import classes from "./base-error-boundary.module.scss";
import { theme } from "antd";
const { useToken } = theme;

export const BaseErrorBoundary = () => {
  const error = useRouteError() as Error;
  const {
    token: { colorError },
  } = useToken();

  return (
    <div className={classes.container}>
      <FontAwesomeIcon color={colorError} icon={faTriangleExclamation} />
      <h1>אירעה שגיאה</h1>
      <p>{error.message}</p>
    </div>
  );
};
