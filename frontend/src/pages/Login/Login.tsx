import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "antd";
import { SubmitHandler, useForm } from "react-hook-form";
import { Login, loginSchema } from "../../models/auth.model";
import classes from "./login.module.scss";

export const LoginForm = () => {
  const { handleSubmit, register } = useForm<Login>({
    defaultValues: {
      email: "",
      password: "",
    },
    resolver: zodResolver(loginSchema),
  });

  const onSubmit: SubmitHandler<Login> = (data) => {
    console.log(data);
  };

  return (
    <main className={classes.container}>
      <section className={classes.loginSection}>
        <header className={classes.header}>
          <h2>כניסה</h2>
          <p>אנא הכנס את האימייל והסיסמא על מנת להתחבר לחשבון</p>
        </header>
        <form className={classes.form} onSubmit={handleSubmit(onSubmit)}>
          <Input {...register("email")} />
          <Input type="password" {...register("password")} />
          <button type="submit">שליחה</button>
        </form>
      </section>
    </main>
  );
};
