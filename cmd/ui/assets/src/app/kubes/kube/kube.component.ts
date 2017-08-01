import { Component, OnInit, Input } from '@angular/core';
import { Kube } from '../kubes.model';

@Component({
  selector: 'app-kube',
  templateUrl: './kube.component.html',
  styleUrls: ['./kube.component.css']
})
export class KubeComponent implements OnInit {
  @Input() kube: Kube[];

  constructor() { }

  ngOnInit() {
  }

}
