import { Outlet } from "react-router-dom";
import classes from "./home.module.scss";

export interface HomeProps {}

export const Home = () => {
  return (
    <div className={classes.container}>
      <Outlet />
    </div>
  );
};
