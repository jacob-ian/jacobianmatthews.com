import { render, screen } from "@testing-library/react";
import { AuthServiceContext } from "../context/auth/AuthServiceContext";
import { useAuthService } from "./useAuthService";

function TestUseAuthService() {
  const auth = useAuthService();
  return <div data-testid="auth">{!!auth ? "true" : "false"}</div>;
}

describe("useAuthService hook", () => {
  it("Should render 'false' if AuthService.Provider has a null value", async () => {
    render(
      <AuthServiceContext.Provider value={null}>
        <TestUseAuthService />
      </AuthServiceContext.Provider>,
    );

    const renderedHookResult = await screen.findByTestId("auth");
    expect(renderedHookResult.innerHTML).toEqual("false");
  });

  it("Should render 'true' if AuthService.Provider has the value of an AuthService instance", async () => {
    const fakeService: any = { test: "hello" };
    render(
      <AuthServiceContext.Provider value={fakeService}>
        <TestUseAuthService />
      </AuthServiceContext.Provider>,
    );

    const renderedHookResult = await screen.findByTestId("auth");
    expect(renderedHookResult.innerHTML).toEqual("true");
  });
});
