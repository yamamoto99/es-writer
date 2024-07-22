export default function getCookieValue(name: string): string | null {
  const nameString = name + "=";
  const value = document.cookie.split('; ').find(row => row.startsWith(nameString));
  return value ? value.split('=')[1] : null;
}