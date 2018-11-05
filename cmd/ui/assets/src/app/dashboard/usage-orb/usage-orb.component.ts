import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'usage-orb',
  templateUrl: './usage-orb.component.html',
  styleUrls: ['./usage-orb.component.scss']
})
export class UsageOrbComponent implements OnInit {

  constructor() { }

  @Input() cluster: any;

  ngOnInit() {
  }

}
