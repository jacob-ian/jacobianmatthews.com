import React from "react";
import SocialLoginButton from "../SocialLoginButton";
import Logo from "./google.svg";

interface GoogleLoginButtonProps {
  onClick: React.MouseEventHandler<HTMLButtonElement>;
  disabled: boolean;
}

export default function GoogleLoginButton(props: GoogleLoginButtonProps) {
  const { onClick, disabled } = props;
  return (
    <SocialLoginButton
      onClick={onClick}
      logo={<Logo />}
      label="Sign in with Google"
      disabled={disabled}
    />
  );
}
