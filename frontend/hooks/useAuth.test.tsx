import Router from "next/router";
import { render, screen } from "@testing-library/react";
import { OnFailOptions, useAuth } from "./useAuth";
import mockRouter from "next-router-mock";
import { User } from "../interfaces/user.interface";
import { AuthContext } from "../context/auth/AuthContext";

jest.mock("next/router", () => require("next-router-mock"));

interface TestUseAuthProps {
  options?: OnFailOptions;
}

function TestUseAuth(props: TestUseAuthProps) {
  const auth = useAuth(props.options);
  return <div data-testid="auth">{JSON.stringify(auth)}</div>;
}

describe("useAuth hook", () => {
  beforeEach(() => {
    mockRouter.setCurrentUrl("/initial");
  });

  describe("With onFailOptions", () => {
    it("Should redirect page to provided onFailOptions url when auth user is null", () => {
      render(
        <AuthContext.Provider value={null}>
          <TestUseAuth options={{ redirectTo: "/test" }} />
        </AuthContext.Provider>,
      );
      expect(Router.pathname).toEqual("/test");
    });

    it("Should not redirect page to onFailOptionsU url when auth user is valid", () => {
      const user: User = {
        uid: "fake",
        name: "Fake User",
        email: "fake@user.com",
        admin: false,
      };
      render(
        <AuthContext.Provider value={user}>
          <TestUseAuth options={{ redirectTo: "/test" }} />
        </AuthContext.Provider>,
      );
      expect(Router.pathname).toEqual("/initial");
    });

    it("Should render the user object", async () => {
      const user: User = {
        uid: "fake",
        name: "Fake User",
        email: "fake@user.com",
        admin: false,
      };
      render(
        <AuthContext.Provider value={user}>
          <TestUseAuth options={{ redirectTo: "/test" }} />
        </AuthContext.Provider>,
      );
      const renderedHookResult = await screen.findByTestId("auth");
      expect(renderedHookResult.innerHTML).toEqual(JSON.stringify(user));
    });
  });

  describe("Without onFailOptions", () => {
    it("Should render null if useAuth returns null", async () => {
      render(
        <AuthContext.Provider value={null}>
          <TestUseAuth />
        </AuthContext.Provider>,
      );
      const renderedHookResult = await screen.findByTestId("auth");
      expect(renderedHookResult.innerHTML).toEqual(JSON.stringify(null));
    });

    it("Should render the user object if useAuth returns an object", async () => {
      const user: User = {
        uid: "fake",
        name: "Fake User",
        email: "fake@user.com",
        admin: false,
      };
      render(
        <AuthContext.Provider value={user}>
          <TestUseAuth />
        </AuthContext.Provider>,
      );
      const renderedHookResult = await screen.findByTestId("auth");
      expect(renderedHookResult.innerHTML).toEqual(JSON.stringify(user));
    });
  });
});
