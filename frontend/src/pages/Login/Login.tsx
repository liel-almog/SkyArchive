import { FormProvider } from "react-hook-form";
import { Link } from "react-router-dom";
import { Input, InputPassword } from "../../components/common/FormElements/Input";
import classes from "./login.module.scss";
import { useLogin } from "./useLogin";

export const LoginForm = () => {
  const { methods, onSubmit } = useLogin();

  return (
    <main className={classes.container}>
      <section className={classes.loginSection}>
        <header className={classes.header}>
          <h2>כניסה</h2>
          <p>אנא הכנס את האימייל והסיסמא על מנת להתחבר לחשבון</p>
        </header>
        <FormProvider {...methods}>
          <form className={classes.form} onSubmit={methods.handleSubmit(onSubmit)}>
            <Input name="email" placeholder="m@example.com" label="אימייל" />
            <InputPassword name="password" label="סיסמא" type="password" />
            <button type="submit">שליחה</button>
            <p>
              אין חשבון עדיין? <Link to="/signup">הרשמה</Link>
            </p>
          </form>
        </FormProvider>
      </section>
    </main>
  );
};
