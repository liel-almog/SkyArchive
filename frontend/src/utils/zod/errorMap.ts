import { z } from "zod";
import validator from "validator";
import { isEmpty } from "remeda";

export const errorMessages = {
  url: "כתובת ה-URL מכילה נתיב לא חוקי",
  mongoId: "צריך להיות מזהה יחודי תקין",
  date: "פורמט תאריך לא תקין",
  minNumber: "חייב להיות מעל או שווה ל",
  maxNumber: "חייב להיות עד",
  numberString: "חייב להיות מספר חיובי/שלילי/שלם/עשרוני",
  fixedNumber: "חייב להיות מספר שלם",
  phone: "מספר טלפון לא תקין",
  positive: "חייב להיות מספר חיובי",
  idNumber: "תעודת זהות לא תקינה",
  required: "חובה להזין ערך בשדה זה",
  enumOption: "יש לבחור ערך מהרשימה",
  alphabet: "צריך להיות אותיות באנגלית או בעברית בלבד",
  odata: "odata not valid",
  english: "צריך להיות אותיות באנגלית בלבד",
  hebrew: "צריך להיות אותיות בעברית בלבד",
  israeliID: "מספר תעודת זהות אינו תקין",
  localPhone: "מספר טלפון אינו תקין",
  aadId: "חייב להיות מספר זיהוי תקין",
  email: "האימייל אינו תקין",
  time: "צריך להיות שעה תקינה",
  alphanumericHebEn: "חייב להיות טקסט בעברית או אנגלית, יכול להכיל גם מספרים",
  text: "חייב להיות טקסט תקין, ללא סימנים מיוחדים (@,#,%,&,$,*,^,+,=,<,>,~,`,|,',/,_)",
} as const;

const stringSchema = z.string().trim().max(255);
const textAreaSchema = z.string().trim().max(2000);
const dateLike = z.union([z.number(), z.string(), z.date()]);

// Backslash does not need to be escaped in character sets.
const validRegex =
  /^[a-zA-Z0-9\u0590-\u05FF- \u0020.,?!:'+ /_ \\\- –';"(){}[\] \n\t\r\f\v]*/gu;
const odateRegex =
  /^[a-zA-Z0-9\u0590-\u05FF- \u0020.,?!:'+ _/$#@ \\\- –';"(){}[\] \n\t\r\f\v]*/gu;

export const customValidation = {
  alphabet: stringSchema.regex(
    /^[a-zA-Z\u0590-\u05FF() ,.'-]*$/i,
    errorMessages.alphabet
  ),
  english: stringSchema.regex(/^[a-zA-z]*$/i, errorMessages.english),
  hebrew: stringSchema.regex(/^[\u0590-\u05FF]*$/i, errorMessages.hebrew),
  phone: {
    localPhone: stringSchema.refine(
      (str) => validator.isMobilePhone(str, ["he-IL"]),
      {
        message: errorMessages.localPhone,
      }
    ),
    optionalLocalPhone: stringSchema
      .refine(
        (str) => validator.isMobilePhone(str, ["he-IL"]) || isEmpty(str),
        {
          message: errorMessages.localPhone,
        }
      )
      .nullish(),
  },
  url: z.string().trim().url({ message: errorMessages.url }).max(1000),
  email: stringSchema.email(errorMessages.email),
  dateLikeToDate: dateLike.pipe(
    z.coerce.date({
      invalid_type_error: errorMessages.date,
      required_error: errorMessages.required,
    })
  ),
  age: z.coerce
    .number({ invalid_type_error: errorMessages.positive })
    .int()
    .min(0, `${errorMessages.minNumber} 0`)
    .max(120, `${errorMessages.maxNumber} 120`)
    .nullish(),
  number: (options?: { max: number }) =>
    z
      .number({ invalid_type_error: "נא להזין מספר תקין" })
      .int()
      .min(0)
      .max(options?.max ? options.max : 9999999999)
      .transform((value) => Math.round(value)),
  alphanumericHebEn: stringSchema.refine((value: string) => {
    const isHebrew = validator.isAlphanumeric(value, "he", {
      ignore: "-, ",
    });
    const isEnglish = validator.isAlphanumeric(value, "en-US", {
      ignore: "-, ",
    });

    return isHebrew || isEnglish;
  }, errorMessages.alphanumericHebEn),
  odata: stringSchema.regex(odateRegex, errorMessages.odata),
  aadId: z
    .string()
    .max(50)
    .refine(
      (value: string) =>
        validator.isAlphanumeric(value, "en-US", { ignore: "-" }),
      errorMessages.aadId
    ),
  enumOption: stringSchema.refine((value: string) => {
    const isHebrew = validator.isAlphanumeric(value, "he", {
      ignore: "'\\ /",
    });
    const isEnglish = validator.isAlphanumeric(value, "en-US", {
      ignore: "'\\ /",
    });

    return isHebrew || isEnglish;
  }, errorMessages.enumOption),
  uuId: stringSchema.uuid(),
  text: z.string().trim().regex(validRegex, errorMessages.text),
  longText: textAreaSchema.regex(validRegex, errorMessages.text),
  //   ObjectId: z.preprocess(
  //     (arg) => {
  //       if (ObjectId.isValid(arg as any)) {
  //         // We know that arg is a valid ObjectId
  //         return new ObjectId(arg as any);
  //       } else {
  //         // throw new Error("Not a valid ObjectId");
  //         return undefined;
  //       }
  //     },
  //     z.instanceof(ObjectId),
  //     { invalid_type_error: errorMessages.mongoId }
  //   ),
} as const;

const CustomErrorMap: z.ZodErrorMap = (error, ctx) => {
  /*
  This is where you override the various error codes
  */
  switch (error.code) {
    case z.ZodIssueCode.invalid_date:
      return { message: errorMessages.date };
    case z.ZodIssueCode.invalid_type:
      if (error.received !== "undefined" && error.received !== "null") {
        if (error.expected === "integer" && error.received === "float") {
          return { message: errorMessages.fixedNumber };
        }

        return {
          message: `סוג הערך שגוי, הסוג הרצוי הוא ${error.expected}`,
        };
      } else {
        return { message: errorMessages.required };
      }
    case z.ZodIssueCode.invalid_enum_value:
      return { message: errorMessages.enumOption };
    case z.ZodIssueCode.too_big:
      if (error.type === "string") {
        return { message: `מספר התווים המקסימלי הוא ${error.maximum}` };
      }

      if (error.type === "array" || error.type === "set") {
        return {
          message: `מספר הפריטים המקסימלי הוא ${error.maximum}`,
        };
      }

      return { message: `${errorMessages.maxNumber} ${error.maximum}` };

    case z.ZodIssueCode.too_small:
      if (error.type === "string") {
        return { message: `מספר התווים המינימלי הוא ${error.minimum}` };
      }

      if (error.type === "array" || error.type === "set") {
        return {
          message: `מספר הפריטים המינימלי הוא ${error.minimum}`,
        };
      }

      return { message: `${errorMessages.minNumber} ${error.minimum}` };

    case z.ZodIssueCode.invalid_union_discriminator:
      return { message: `יש לבחור את אחד מהשדות` };
  }

  // fall back to default message!
  return { message: ctx.defaultError };
};

z.setErrorMap(CustomErrorMap);
