// TODO: remove this file once
// "new cluster wizard" is separated by components

export function sortDigitaloceanMachineTypes(machineTypes) {
  return machineTypes
    .sort((a, b) => { // for types like 's-1vcpu-1gb'
      const SEPARATOR = 'gb';
      const [sizeA] = a.split(SEPARATOR);
      const [sizeB] = b.split(SEPARATOR);

      if (sizeA === sizeB) return 0;

      return (isNaN(+sizeA) || isNaN(+sizeB)) ? -1 : 1;
    })
    .sort((a, b) => { // for types like '1gb'
      const SEPARATOR = 'gb';
      const [sizeA] = a.split(SEPARATOR);
      const [sizeB] = b.split(SEPARATOR);

      return +sizeA > +sizeB ? 1 : -1;
    });
}
