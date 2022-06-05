import { CSSProperties } from "react";
import styles from "./SocialLoginButton.module.scss";

interface SocialLoginButtonProps {
  onClick: React.MouseEventHandler<HTMLButtonElement>;
  label: string;
  logo: JSX.Element;
  style?: CSSProperties;
  disabled: boolean;
}

export default function SocialLoginButton(props: SocialLoginButtonProps) {
  const { onClick, logo, label, style, disabled } = props;
  return (
    <button
      onClick={onClick}
      className={styles.button}
      style={style}
      disabled={disabled}>
      {logo}
      <div className={styles.label}>{label}</div>
    </button>
  );
}
