export const CLUSTER_OPTIONS = {
  archs: ['amd64'],
  networkProviders: ['Flannel', 'Calico', 'Weave'],
  operatingSystems: ['linux'],
  networkTypes: ['vxlan'],
  ubuntuVersions: ['xenial'],
  helmVersions: ['2.11.0'],
  dockerVersions: ['18.06.3'],
  K8sVersions: ['1.12.10', '1.13.9', '1.14.5', '1.15.2']
};
