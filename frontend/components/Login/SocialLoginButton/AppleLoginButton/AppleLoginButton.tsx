import SocialLoginButton from "../SocialLoginButton";
import Logo from "./apple.svg";

interface AppleLoginButtonProps {
  onClick: React.MouseEventHandler<HTMLButtonElement>;
  disabled: boolean;
}

export default function AppleLoginButton(props: AppleLoginButtonProps) {
  const { onClick, disabled } = props;
  return (
    <SocialLoginButton
      onClick={onClick}
      label="Sign in with Apple"
      logo={<Logo />}
      style={{ backgroundColor: "#000000", color: "#ffffff" }}
      disabled={disabled}
    />
  );
}
