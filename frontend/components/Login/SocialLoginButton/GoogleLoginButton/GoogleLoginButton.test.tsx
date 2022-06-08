import { cleanup, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import GoogleLoginButton from "./GoogleLoginButton";

jest.mock("./google.svg", () => "svg");

describe("GoogleLoginButton", () => {
  afterEach(() => {
    cleanup();
  });

  it("Should render a button", () => {
    const onClick = () => {};
    render(<GoogleLoginButton onClick={onClick} disabled={false} />);
    const button = screen.getByRole("button");
    expect(button).toBeInTheDocument();
  });

  it("Should render a button with the text 'Sign in with Google'", async () => {
    const onClick = () => {};
    render(<GoogleLoginButton onClick={onClick} disabled={false} />);
    const element = await screen.findByText(/Sign in with Google/);
    expect(element).toBeInTheDocument();
    expect(element.parentElement?.tagName).toEqual("BUTTON");
  });

  it("Should call onClick when button is clicked", async () => {
    const onClick = jest.fn();
    render(<GoogleLoginButton onClick={onClick} disabled={false} />);
    const button = await screen.findByRole("button");
    await userEvent.click(button);
    expect(onClick).toHaveBeenCalled();
  });

  it("Should not call onClick when button is clicked but is disabled", async () => {
    const onClick = jest.fn();
    render(<GoogleLoginButton onClick={onClick} disabled={true} />);
    const button = await screen.findByRole("button");
    await userEvent.click(button);
    expect(onClick).not.toHaveBeenCalled();
  });
});
