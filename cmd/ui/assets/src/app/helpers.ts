export function getBaseLocation() {
  return window.location.pathname.split('/').splice(1).join('/');
}
