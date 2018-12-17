// TODO: remove this file once
// "new cluster wizard" is separated by components

export function sortDigitalOceanMachineTypes(machineTypes) {

  return machineTypes
    .sort(sortByCategory)
    .sort(sortTypeWithinCategory);
}

function sortByCategory(a, b) {
  const typeA = getDigitalOceanMachineType(a);
  const typeB = getDigitalOceanMachineType(b);

  return typeA - typeB;
}

const typeSorters = [
  (a, b) => { // 512mb
    const separator = 'mb';
    const A = +a.split(separator)[0];
    const B = +b.split(separator)[0];
    return A - B;
  },
  (a, b) => { // 1gb
    const separator = 'gb';
    const A = +a.split(separator)[0];
    const B = +b.split(separator)[0];
    return A - B;
  },
  (a, b) => { // m-4
    const separator = '-';
    const A = +a.split(separator)[1];
    const B = +b.split(separator)[1];
    return A - B;
  },
  (a, b) => { // m-16gb
    const separator = '-';
    const A = +a.split(separator)[1].split('gb')[0];
    const B = +b.split(separator)[1].split('gb')[0];
    return A - B;
  },
  (a, b) => { // s-1cpu-1gb
    const separator = '-';
    const A = +a.split(separator)[2].split('gb')[0];
    const B = +b.split(separator)[2].split('gb')[0];
    return A - B;
  },
];

function sortTypeWithinCategory(a, b) {
  const typeA = getDigitalOceanMachineType(a);
  const typeB = getDigitalOceanMachineType(b);

  if (typeA === typeB) {
    return typeSorters[typeA](a, b);
  }
  return 0;
}

export function getDigitalOceanMachineType(machineSize): number {
  // NOTE: type index is the order of types sorting
  const sizeTypes = [
    /^\d+mb$/, // 1mb
    /^\d+gb$/, // 1gb
    /^[a-z]-\d+$/, // m-4
    /^[a-z]-\d+gb$/, // m-16gb
    /^[a-z]-\d+[a-z]+-\d+gb$/, // s-12vcpu-48gb
  ];

  return sizeTypes
    .findIndex(
      typeRegEx => machineSize.match(typeRegEx),
    );
}
