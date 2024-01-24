import {
  Input as AntInput,
  InputNumber as AntInputNumber,
  InputNumberProps as AntInputNumberProps,
  InputProps as AntInputProp,
} from "antd";
import { TextAreaProps } from "antd/lib/input";
import clsx from "clsx";
import { useController, useFormContext } from "react-hook-form";
import { HelperText } from "../HelperText";
import classes from "./input.module.scss";

const { TextArea } = AntInput;
const { Password } = AntInput;
export interface ExtraInputProps {
  name: string;
  label?: string;
  description?: string | false;
}

type InputProps = AntInputProp & ExtraInputProps;

export const Input = ({
  name,
  label,
  description = false,
  ...rest
}: InputProps) => {
  const {
    formState: { isSubmitting },
  } = useFormContext();

  const {
    fieldState: { invalid, error },
    field,
  } = useController({
    name,
  });

  const handleChange: InputProps["onChange"] = (e) => {
    if (!e.target.value) {
      field.onChange(null);
      return;
    }

    field.onChange(e.target.value);
  };

  return (
    <label className={classes.container} htmlFor={name}>
      <h4>
        {label} {rest.required && "*"}
      </h4>
      {description && <p>{description}</p>}
      <AntInput
        {...field}
        onChange={handleChange}
        id={name}
        disabled={isSubmitting}
        aria-errormessage={error?.message}
        className={classes.input}
        status={invalid ? "error" : undefined}
        {...rest}
      />
      <HelperText error={invalid} message={error?.message} />
    </label>
  );
};

export const InputPassword = ({
  name,
  label,
  description = false,
  ...rest
}: InputProps) => {
  const {
    formState: { isSubmitting },
  } = useFormContext();

  const {
    fieldState: { invalid, error },
    field,
  } = useController({
    name,
  });

  const handleChange: InputProps["onChange"] = (e) => {
    if (!e.target.value) {
      field.onChange(null);
      return;
    }

    field.onChange(e.target.value);
  };

  return (
    <label className={classes.container} htmlFor={name}>
      <h4>
        {label} {rest.required && "*"}
      </h4>
      {description && <p>{description}</p>}
      <Password
        {...field}
        onChange={handleChange}
        id={name}
        disabled={isSubmitting}
        aria-errormessage={error?.message}
        className={classes.input}
        status={invalid ? "error" : undefined}
        {...rest}
      />
      <HelperText error={invalid} message={error?.message} />
    </label>
  );
};

type InputNumberProps = AntInputNumberProps & ExtraInputProps;

export const InputNumber = ({
  name,
  label,
  description = false,
  ...rest
}: InputNumberProps) => {
  const {
    formState: { isSubmitting },
  } = useFormContext();

  const {
    fieldState: { error, invalid },
    field,
  } = useController({
    name,
  });

  const handleChange: InputNumberProps["onChange"] = (val) => {
    if ((!val || isNaN(Number(val))) && val !== 0) {
      field.onChange(null);
      return;
    }

    if (rest.type === "text") {
      field.onChange(val.toString());
      return;
    } else if (rest.type === "number") {
      field.onChange(+val);
      return;
    }

    field.onChange(val);
  };

  return (
    <label className={classes.container} htmlFor={name}>
      <h4>
        {label} {rest.required && "*"}
      </h4>
      {description && <p>{description}</p>}
      <AntInputNumber
        {...field}
        dir="rtl"
        inputMode="numeric"
        type="number"
        id={name}
        pattern="[0-9]*"
        onChange={handleChange}
        disabled={isSubmitting}
        className={classes.inputNumber}
        status={invalid ? "error" : undefined}
        {...rest}
      />
      <HelperText error={invalid} message={error?.message} />
    </label>
  );
};

export const TextInput = ({
  name,
  label,
  description = false,
  ...rest
}: TextAreaProps & ExtraInputProps) => {
  const {
    formState: { isSubmitting },
  } = useFormContext();

  const {
    fieldState: { invalid, error },
    field,
  } = useController({ name });

  const handleChange: TextAreaProps["onChange"] = (e) => {
    if (!e.target.value) {
      field.onChange(null);
      return;
    }

    field.onChange(e.target.value);
  };

  return (
    <label className={classes.container} htmlFor={name}>
      <h4>
        {label} {rest.required && "*"}
      </h4>
      {description && <p>{description}</p>}
      <TextArea
        {...field}
        id={name}
        disabled={isSubmitting}
        className={classes.textArea}
        autoSize={{ maxRows: 3 }}
        onChange={handleChange}
        status={invalid ? "error" : undefined}
        {...rest}
        rows={5}
      />
      <HelperText error={invalid} message={error?.message} />
    </label>
  );
};

export const DumbInput = ({ className, ...rest }: AntInputProp) => (
  <AntInput className={clsx(classes.input, className)} {...rest} />
);

export const DumbInputNumber = ({
  className,
  ...rest
}: AntInputNumberProps) => (
  <AntInputNumber
    inputMode="numeric"
    type="number"
    className={clsx(classes.inputNumber, className)}
    {...rest}
  />
);
