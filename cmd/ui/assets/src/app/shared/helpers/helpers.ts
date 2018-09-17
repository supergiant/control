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
	let date = new Date(isoString);
	let month = date.getMonth() + 1;
	let day = date.getDay() + 1;
	let fullYear = date.getFullYear();
	let hours = date.getHours();
	let minutes = date.getMinutes();
	let seconds = date.getSeconds();
	return `${month}-${day}-${fullYear}, ${hours}:${minutes}:${seconds}`;
}
