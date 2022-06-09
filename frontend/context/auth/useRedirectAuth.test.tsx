import { act, render, screen } from "@testing-library/react";
import { InvalidAuthException } from "../../services/auth/InvalidAuthException";
import { NoAuthRedirectException } from "../../services/auth/NoAuthRedirectException";
import { InternalErrorException } from "../../services/http/InternalErrorException";
import { UnauthenticatedException } from "../../services/http/UnauthenticatedException";
import { useRedirectAuth } from "./useRedirectAuth";

function TestUseRedirectAuth(props: { authService: any }) {
  const res = useRedirectAuth(props.authService);
  return <div data-testid="test">{JSON.stringify(res)}</div>;
}

describe("useRedirectAuth hook", () => {
  it("Should return a null user, null error and false loading if user is not authenticated and there was no auth redirect", async () => {
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

  it("Should return a null user, a non null error and false loading, if something went wrong during the auth redirect handler", async () => {
    const user: any = { user: true };

    jest.spyOn(console, "error").mockImplementation();

    const authService: any = {
      handleAuthRedirect: jest
        .fn()
        .mockRejectedValue(new InvalidAuthException("Oh no!")),
      getSignedInUser: jest.fn().mockResolvedValue(user),
    };

    await act(async () => {
      render(<TestUseRedirectAuth authService={authService} />);
    });

    const container = await screen.findByTestId("test");
    expect(container.innerHTML).toEqual(
      JSON.stringify({
        user: null,
        error: new Error(
          "An error occurred while signing you in. Please try again.",
        ),
        loading: false,
      }),
    );
  });

  it("Should return a null user, a non null error, and false loading if auth redirect was successful but fetching signed in user failed", async () => {
    jest.spyOn(console, "error").mockImplementation();

    const authService: any = {
      handleAuthRedirect: jest.fn().mockResolvedValue(undefined),
      getSignedInUser: jest
        .fn()
        .mockRejectedValue(new InternalErrorException()),
    };

    await act(async () => {
      render(<TestUseRedirectAuth authService={authService} />);
    });

    const container = await screen.findByTestId("test");
    expect(container.innerHTML).toEqual(
      JSON.stringify({
        user: null,
        error: new Error(
          "An error occurred while signing you in. Please try again.",
        ),
        loading: false,
      }),
    );
  });

  it("Should return a user, a null error, and false loading if auth redirect was successful and fetching signed in user succeeded", async () => {
    const user: any = { test: true };
    const authService: any = {
      handleAuthRedirect: jest.fn().mockResolvedValue(undefined),
      getSignedInUser: jest.fn().mockResolvedValue(user),
    };

    await act(async () => {
      render(<TestUseRedirectAuth authService={authService} />);
    });

    const container = await screen.findByTestId("test");
    expect(container.innerHTML).toEqual(
      JSON.stringify({
        user,
        error: null,
        loading: false,
      }),
    );
  });
});
