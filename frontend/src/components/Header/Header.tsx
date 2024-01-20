import { useNavigate } from "react-router-dom";
import classes from "./header.module.scss";
// import { useContext } from "react";
// import { UsernameContext } from "../../context/Username.context";

export interface HeaderProps {}

export const Header = () => {
  const navigate = useNavigate();
  // const { username } = useContext(UsernameContext);
  return (
    <header className={classes.container}>
      <div onClick={() => navigate("/")} className={classes.home}>
        <h1>File Uploader</h1>
      </div>
      <div className={classes.username}>
        <h1>{"לא נבחר שם משתמש"}</h1>
      </div>
    </header>
  );
};
