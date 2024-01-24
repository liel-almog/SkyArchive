import { FormProvider } from "react-hook-form";
import {
  Input,
  InputPassword,
} from "../../components/common/FormElements/Input";
import classes from "./signup.module.scss";
import { useSignup } from "./useSignup";
import { Link } from "react-router-dom";

export const SignupForm = () => {
  const { methods, onSubmit } = useSignup();

  return (
    <main className={classes.container}>
      <section className={classes.signupSection}>
        <header className={classes.header}>
          <h2>הרשמה</h2>
          <p>אנא הכנס את הפרטים שלך על מנת לפתוח חשבון</p>
        </header>
        <FormProvider {...methods}>
          <form
            className={classes.form}
            onSubmit={methods.handleSubmit(onSubmit)}
          >
            <Input name="username" label="שם משתמש" />
            <Input name="email" placeholder="m@example.com" label="אימייל" />
            <InputPassword name="password" label="סיסמא" />
            <InputPassword
              name="confirmPassword"
              label="אימות סיסמא"
              type="password"
            />
            <button type="submit">שליחה</button>
            <p>
              כבר יש חשבון? <Link to="/login">כניסה</Link>
            </p>
          </form>
        </FormProvider>
      </section>
    </main>
  );
};
