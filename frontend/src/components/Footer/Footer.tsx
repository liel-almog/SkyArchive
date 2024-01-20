import packageJson from "../../../package.json";
import classes from "./footer.module.scss";

export const Footer = () => {
  return (
    <footer className={classes.container}>
      <EnvironmentIndicator />
      <span>פותח ע"י ליאל אלמוג</span>
      <VersionIndicator />
    </footer>
  );
};

const VersionIndicator = () => (
  <span role="version">{packageJson.version}</span>
);

const EnvironmentIndicator = () => {
  const BlinkingDot = () => <span className={classes.dot} />;

  const Body = ({ label }: { label: string }) => {
    return (
      <span role="environment">
        <BlinkingDot />
        {"סביבת " + label}
      </span>
    );
  };

  return <Body label={"פיתוח מקומית"} />;
};
