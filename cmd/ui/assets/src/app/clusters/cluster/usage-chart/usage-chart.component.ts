import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-usage-chart',
  templateUrl: './usage-chart.component.html',
  styleUrls: ['./usage-chart.component.scss']
})
export class UsageChartComponent {

  @Input() name: string;
}
