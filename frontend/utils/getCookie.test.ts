import { getCookie } from "./getCookie";

describe("getCookie", () => {
  const beforeCookies = document.cookie;

  afterEach(() => {
    setCookieDev(beforeCookies);
  });

  it("Should return undefined if the cookie does not exist", () => {
    setCookieDev("");
    const cookie = getCookie("fake-cookie");
    expect(cookie).toEqual(undefined);
  });

  it("Should return the value 'hello' for the cookie 'test'", () => {
    setCookieDev("fake=no;test=hello;session=1");
    const cookie = getCookie("test");
    expect(cookie).toEqual("hello");
  });
});

function setCookieDev(cookie: string) {
  Object.defineProperty(window.document, "cookie", {
    writable: true,
    value: cookie,
  });
}
