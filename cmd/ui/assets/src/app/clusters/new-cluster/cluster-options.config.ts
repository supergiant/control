export const CLUSTER_OPTIONS = {
  archs: ['amd64'],
  networkProviders: ['Flannel', 'Calico', 'Weave'],
  operatingSystems: ['linux'],
  networkTypes: ['vxlan'],
  ubuntuVersions: ['xenial'],
  helmVersions: ['2.11.0'],
  dockerVersions: ['18.06.3'],
  K8sVersions: ['1.12.10', '1.13.8', '1.14.4', '1.15.1']
};
