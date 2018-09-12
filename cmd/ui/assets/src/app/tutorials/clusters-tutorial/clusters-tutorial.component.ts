import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-clusters-tutorial',
  templateUrl: './clusters-tutorial.component.html',
  styleUrls: ['./clusters-tutorial.component.scss']
})
export class ClustersTutorialComponent implements OnInit {

  constructor() { }

  ngOnInit() {
    console.log('Loading Tutorial Component');
  }

}
