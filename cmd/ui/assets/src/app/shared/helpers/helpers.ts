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

export function convertIsoToHumanReadable(isoString) {
	const date = new Date(isoString);
	const month = date.getMonth() + 1;
	const day = date.getDay() + 1;
	const fullYear = date.getFullYear();
	const hours = date.getHours();
	const minutes = date.getMinutes();
	const seconds = date.getSeconds();
	return `${month}-${day}-${fullYear}, ${hours}:${minutes}:${seconds}`;
}
