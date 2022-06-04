export function getCookie(name: string): string | undefined {
  const cookies = decodeURIComponent(document.cookie).split(";");
  const cookie = cookies.find((c) => c.startsWith(`${name}=`));
  if (!cookie) {
    return undefined;
  }
  return cookie.split("=").pop();
}
