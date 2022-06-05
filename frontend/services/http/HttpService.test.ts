import { HttpService } from "./HttpService";
import fetchMock from "jest-fetch-mock";

describe("HttpService", () => {
  beforeAll(() => {
    fetchMock.enableMocks();
  });

  beforeEach(() => {
    fetchMock.resetMocks();
  });

  afterAll(() => {
    fetchMock.disableMocks();
  });

  describe("get", () => {
    it.todo("Should throw an ");
  });

  describe("post", () => {});
});
