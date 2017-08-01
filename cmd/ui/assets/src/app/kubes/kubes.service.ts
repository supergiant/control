import { Injectable } from '@angular/core';
import { Kube } from './kubes.model'

@Injectable()
export class KubesService {
  private kubes: Kube[] = [
    {name: "kube1", description: "desc1"},
    {name: "kube2", description: "desc2"},
    {name: "kube3", description: "desc3"},
  ]

  getKubes() {
    return this.kubes;
  }

  getKube(id: number) {
    return this.kubes[id];
  }

  constructor() { }

}
