export function getBaseLocation() {

  if (window.location.pathname.split('/')[0] === '') {
    const paths: string[] = location.pathname.split('/').splice(1, 1);
    const basePath: string = (paths && paths[0]);
    return '/';
  } else {
    const paths: string[] = location.pathname.split('/').splice(1, 1);
    const basePath: string = (paths && paths[0]);
    return '/' + basePath;
  }
}
