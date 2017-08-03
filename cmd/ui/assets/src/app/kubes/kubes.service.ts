import { Injectable } from '@angular/core';
import { Kube } from './kubes.model'

@Injectable()
export class KubesService {
  private kubes: Kube[] = [
    {id: 1, name: "kube1", masterNodeSize: "desc1", status: "true"},
    {id: 2, name: "kube2", masterNodeSize: "desc2", status: "true"},
    {id: 3, name: "kube3", masterNodeSize: "desc3", status: "true"},
  ]

  getKubes() {
    return this.kubes;
  }

  getKube(id: number) {
    return this.kubes[id];
  }

  constructor() { }

}
