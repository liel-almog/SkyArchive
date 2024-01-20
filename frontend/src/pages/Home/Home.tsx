import { Outlet } from "react-router-dom";
import classes from "./home.module.scss";
import { Footer } from "../../components/Footer";
import { Header } from "../../components/Header";

export interface HomeProps {}

export const Home = () => {
  return (
    <div className={classes.container}>
      <Header />
      <main>
        <Outlet />
      </main>
      <Footer />
    </div>
  );
};
