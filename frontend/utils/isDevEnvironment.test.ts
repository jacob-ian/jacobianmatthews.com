import { isDevEnvironment } from "./isDevEnvironment";

describe("isDevEnvironment", () => {
  const envBefore = process.env;

  afterEach(() => {
    process.env = envBefore;
  });

  it("Should return false if NODE_ENV is 'production'", () => {
    (process.env as any).NODE_ENV = "production";
    expect(isDevEnvironment()).toEqual(false);
  });

  it("Should return true if NODE_ENV is 'development'", () => {
    (process.env as any).NODE_ENV = "development";
    expect(isDevEnvironment()).toEqual(true);
  });
});
