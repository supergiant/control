import { Component, OnInit } from '@angular/core';
import { KubesService } from './kubes.service';
import {Kube} from './kubes.model'

@Component({
  selector: 'app-kubes',
  templateUrl: './kubes.component.html',
  styleUrls: ['./kubes.component.css']
})
export class KubesComponent implements OnInit {
  kubes: Kube[]

  constructor( private kubesService: KubesService) { }

  ngOnInit() {
    this.kubes = this.kubesService.getKubes();
  }

}
