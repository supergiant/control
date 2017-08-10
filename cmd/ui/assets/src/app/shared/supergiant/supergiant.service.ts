import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import { CloudAccounts } from './cloud-accounts/cloud-accounts.service'
import { Sessions } from './sessions/sessions.service'
import { Users } from './users/users.service'
import { Kubes } from './kubes/kubes.service'
import { KubeResources } from './kube-resources/kube-resources.service'
import { Nodes } from './nodes/nodes.service'
import { LoadBalancers } from './load-balancers/load-balancers.service'
import { HelmRepos } from './helm-repos/helm-repos.service'
import { HelmCharts } from './helm-charts/helm-charts.service'
import { HelmReleases } from './helm-releases/helm-releases.service'
import { Logs } from './logs/logs.service'
import { UtilService } from './util/util.service'

@Injectable()
export class Supergiant {
constructor(
  public CloudAccounts: CloudAccounts,
  public Sessions: Sessions,
  public Users: Users,
  public Kubes: Kubes,
  public KubeResources: KubeResources,
  public Nodes: Nodes,
  public LoadBalancers: LoadBalancers,
  public HelmRepos: HelmRepos,
  public HelmCharts: HelmCharts,
  public HelmReleases: HelmReleases,
  public Logs: Logs,
) {}
}
