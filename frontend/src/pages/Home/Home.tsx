import { Outlet } from "react-router-dom";
import { Header } from "../../components/Header";
import classes from "./home.module.scss";

export interface HomeProps {}

export const Home = () => {
  return (
    <div className={classes.container}>
      <Header />
      <main>
        <Outlet />
      </main>
    </div>
  );
};
