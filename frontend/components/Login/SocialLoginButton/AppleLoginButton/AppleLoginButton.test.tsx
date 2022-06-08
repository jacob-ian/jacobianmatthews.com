import { cleanup, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import AppleLoginButton from "./AppleLoginButton";

jest.mock("./apple.svg", () => "svg");

describe("AppleLoginButton", () => {
  afterEach(() => {
    cleanup();
  });
  it("Should render a button", () => {
    const onClick = () => {};
    render(<AppleLoginButton onClick={onClick} disabled={false} />);
    const button = screen.getByRole("button");
    expect(button).toBeInTheDocument();
  });

  it("Should render a button with the text 'Sign in with Apple'", async () => {
    const onClick = () => {};
    render(<AppleLoginButton onClick={onClick} disabled={false} />);
    const element = await screen.findByText(/Sign in with Apple/);
    expect(element).toBeInTheDocument();
    expect(element.parentElement?.tagName).toEqual("BUTTON");
  });

  it("Should render a button with a background color of #000000 and font color of #ffffff", async () => {
    const onClick = () => {};
    render(<AppleLoginButton onClick={onClick} disabled={false} />);
    const button = await screen.findByRole("button");
    expect(button).toHaveStyle({
      backgroundColor: "#000000",
      color: "#ffffff",
    });
  });

  it("Should call onClick when button is clicked", async () => {
    const onClick = jest.fn();
    render(<AppleLoginButton onClick={onClick} disabled={false} />);
    const button = await screen.findByRole("button");
    await userEvent.click(button);
    expect(onClick).toHaveBeenCalled();
  });

  it("Should not call onClick when button is clicked but is disabled", async () => {
    const onClick = jest.fn();
    render(<AppleLoginButton onClick={onClick} disabled={true} />);
    const button = await screen.findByRole("button");
    await userEvent.click(button);
    expect(onClick).not.toHaveBeenCalled();
  });
});
