import clsx from "clsx";
import classes from "./helper-text.module.scss";

interface HelperTextProps {
  message: string | undefined;
  error: boolean;
}

export const HelperText = ({ message: text, error }: HelperTextProps) => {
  if (!text) {
    return <></>;
  }

  if (error) {
    return (
      <span className={clsx(classes.helperText, classes.error)}>{text}</span>
    );
  }

  return <span className={clsx(classes.helperText)}>{text}</span>;
};
