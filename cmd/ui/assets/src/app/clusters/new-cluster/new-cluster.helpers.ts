// TODO: remove this file once
// "new cluster wizard" is separated by components

export function sortDigitaloceanMachineTypes(machineTypes) {
  const SEPARATOR_GB = 'gb';
  const SEPARATOR_MB = 'mb';

  return machineTypes

    .sort((a, b) => { // for types like 's-1vcpu-1gb'
      const [sizeA] = a.split(SEPARATOR_GB);
      const [sizeB] = b.split(SEPARATOR_GB);

      if (sizeA === sizeB) {
        return 0;
      }

      return (isNaN(+sizeA) || isNaN(+sizeB)) ? -1 : 1;
    })

    .sort((a, b) => { // for types like '1gb'
      const [sizeA] = a.split(SEPARATOR_GB);
      const [sizeB] = b.split(SEPARATOR_GB);

      if (sizeA === sizeB) {
        return 0;
      }

      return +sizeA > +sizeB ? 1 : -1;
    })

    .sort((a, b) => { // for types like '1gb'
      const SEPARATOR_DASH = '-';
      const sizeA = a.split(SEPARATOR_DASH);
      const sizeB = b.split(SEPARATOR_DASH);

      const isSorted = [sizeA, sizeB].every(el => el.length < 3);

      if (isSorted) {
        return 0;
      }

      if (sizeA.length === 2 && sizeB.length === 2) { //

        const A = sizeA[0];
        const B = sizeB[0];

        if (A === B) {
          return 0;
        }

        return A > B ? 1 : -1;
      }

      if (sizeA.length === 3 && sizeB.length === 3) { //
        const A = sizeA[2].split(SEPARATOR_GB)[0];
        const B = sizeB[2].split(SEPARATOR_GB)[0];

        if (A === B) {
          return 0;
        }

        return +A > +B ? 1 : -1;
      }
    })

    .sort((a, b) => {

      if (a.indexOf(SEPARATOR_MB) > -1) {
        return -1;
      }

      if (a === b) {
        return 0;
      }

      return 1;

    });
}
