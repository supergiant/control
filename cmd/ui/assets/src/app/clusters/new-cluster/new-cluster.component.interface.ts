export interface IMachineType {
  machineType: any;
  role: string;
  qty: number;
  availabilityZone: string;
  availableMachineTypes: string[];
  recommendedNodesCount?: number;
}
