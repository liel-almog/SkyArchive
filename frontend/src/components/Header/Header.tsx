import { useNavigate } from "react-router-dom";
import classes from "./header.module.scss";
import { useAuthContext } from "../../context/AuthContext/useAuthProvider";

export interface HeaderProps {}

export const Header = () => {
  const navigate = useNavigate();
  const { user } = useAuthContext();

  return (
    <header className={classes.container}>
      <div onClick={() => navigate("/")} className={classes.home}>
        <h1>File Uploader</h1>
      </div>
      <div className={classes.username}>
        <h1>{user?.username}</h1>
      </div>
    </header>
  );
};
