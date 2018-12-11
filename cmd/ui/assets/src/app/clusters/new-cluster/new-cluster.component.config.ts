import { IMachineType } from 'app/clusters/new-cluster/new-cluster.component.interface';

export const BLANK_MACHINE_TEMPLATE: IMachineType = {
  machineType: null,
  role: '',
  qty: 1,
  availabilityZone: '',
  availableMachineTypes: null,
};

// TODO: an interface type need to be defined
export const DEFAULT_MACHINE_SET = [
  {
    ...BLANK_MACHINE_TEMPLATE,
    role: 'Master',
  },
  {
    ...BLANK_MACHINE_TEMPLATE,
    role: 'Node',
  },
];
