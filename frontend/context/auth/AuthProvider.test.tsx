import {
  act,
  cleanup,
  render,
  renderHook,
  screen,
} from "@testing-library/react";
import { useContext } from "react";
import { User } from "../../interfaces/user.interface";
import { AuthContext } from "./AuthContext";
import AuthProvider from "./AuthProvider";
import { AuthServiceContext } from "./AuthServiceContext";
import * as redirectAuth from "./useRedirectAuth";

jest.mock("firebase/app", () => ({
  initializeApp: jest.fn().mockReturnValue("app"),
}));

jest.mock("../../services/auth/AuthService", () => ({
  AuthService: function () {
    return { test: "ok" };
  },
}));

function TestAuthProviderConsumer() {
  const auth = useContext(AuthContext);
  return <div data-testid="test">{JSON.stringify(auth)}</div>;
}

function TestAuthServiceProvider() {
  const authService = useContext(AuthServiceContext);
  return <div data-testid="test">{JSON.stringify(authService)}</div>;
}

describe("AuthProvider", () => {
  afterEach(cleanup);

  it("Should open an error alert if there is an error when authenticating via redirect", async () => {
    jest.spyOn(redirectAuth, "useRedirectAuth").mockImplementation(() => ({
      error: new Error("Bad thing"),
      user: null,
      loading: false,
    }));

    render(
      <AuthProvider>
        <div></div>
      </AuthProvider>,
    );

    const errorAlert = await screen.findByText("Bad thing");
    expect(errorAlert).toBeInTheDocument();
  });

  it("Should provide a value of null for AuthContext if not authenticated", async () => {
    jest.spyOn(redirectAuth, "useRedirectAuth").mockImplementation(() => ({
      error: null,
      user: null,
      loading: false,
    }));

    render(
      <AuthProvider>
        <TestAuthProviderConsumer />
      </AuthProvider>,
    );

    const authValue = await screen.findByTestId("test");
    expect(authValue.innerHTML).toEqual(JSON.stringify(null));
  });

  it("Should provide a user value for AuthContext if authenticated", async () => {
    const user: User = {
      uid: "fake",
      name: "Fake User",
      admin: false,
      email: "fake",
    };

    jest.spyOn(redirectAuth, "useRedirectAuth").mockImplementation(() => ({
      user,
      error: null,
      loading: false,
    }));

    render(
      <AuthProvider>
        <TestAuthProviderConsumer />
      </AuthProvider>,
    );

    const authValue = await screen.findByTestId("test");
    expect(authValue.innerHTML).toEqual(JSON.stringify(user));
  });

  it("Should provide an instance of AuthService", async () => {
    jest.spyOn(redirectAuth, "useRedirectAuth").mockImplementation(() => ({
      user: null,
      error: null,
      loading: false,
    }));

    render(
      <AuthProvider>
        <TestAuthServiceProvider />
      </AuthProvider>,
    );

    const container = await screen.findByTestId("test");
    expect(container.innerHTML).toEqual(JSON.stringify({ test: "ok" }));
  });
});
