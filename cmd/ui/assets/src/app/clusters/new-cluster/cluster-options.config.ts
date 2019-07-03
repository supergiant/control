export const CLUSTER_OPTIONS = {
  archs: ['amd64'],
  networkProviders: ['Flannel', 'Calico', 'Weave'],
  operatingSystems: ['linux'],
  networkTypes: ['vxlan'],
  ubuntuVersions: ['xenial'],
  helmVersions: ['2.11.0'],
  dockerVersions: ['18.06.3'],
  K8sVersions: ['1.12.7', '1.13.7', '1.14.3', '1.15.0']
};
