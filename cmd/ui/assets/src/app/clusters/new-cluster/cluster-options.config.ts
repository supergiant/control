export const CLUSTER_OPTIONS = {
  archs: ['amd64'],
  networkProviders: [
    'Flannel',
    'Calico',
    'Weave',
  ],
  operatingSystems: ['linux'],
  networkTypes: ['vxlan'],
  ubuntuVersions: ['xenial'],
  helmVersions: ['2.11.0'],
  dockerVersions: ['17.06.0'],
  K8sVersions: [
    '1.12.0',
    '1.13.0',
  ]
};
