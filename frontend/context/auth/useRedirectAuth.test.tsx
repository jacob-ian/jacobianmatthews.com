import { act, render, screen } from "@testing-library/react";
import { NoAuthRedirectException } from "../../services/auth/NoAuthRedirectException";
import { UnauthenticatedException } from "../../services/http/UnauthenticatedException";
import { useRedirectAuth } from "./useRedirectAuth";

function TestUseRedirectAuth(props: { authService: any }) {
  const res = useRedirectAuth(props.authService);
  return <div data-testid="test">{JSON.stringify(res)}</div>;
}

describe("useRedirectAuth hook", () => {
  it("Should return a null user, null error and false loading if user is not authenticated and is not an auth redirect", async () => {
    const authService: any = {
      handleAuthRedirect: jest
        .fn()
        .mockRejectedValue(new NoAuthRedirectException()),
      getSignedInUser: jest
        .fn()
        .mockRejectedValue(new UnauthenticatedException()),
    };

    await act(async () => {
      render(<TestUseRedirectAuth authService={authService} />);
    });

    const container = await screen.findByTestId("test");
    expect(container.innerHTML).toEqual(
      JSON.stringify({
        user: null,
        error: null,
        loading: false,
      }),
    );
  });
});
