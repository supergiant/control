import { Component, OnInit, Input, ViewChild, OnChanges, ChangeDetectionStrategy } from '@angular/core';

@Component({
  selector: 'app-usage-chart',
  templateUrl: './usage-chart.component.html',
  styleUrls: ['./usage-chart.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class UsageChartComponent implements OnInit, OnChanges {

  @Input() name: string;
  @Input() cpuUsage: number;
  @Input() ramUsage: number;
  @ViewChild("cpu_usage_path") cpu_usage_path;
  @ViewChild("ram_usage_path") ram_usage_path;

  ngOnChanges() {
    const cpu_length = this.cpu_usage_path.nativeElement.getTotalLength();
    const ram_length = this.ram_usage_path.nativeElement.getTotalLength();

    // mins 7 for a strange calibration issue
    const new_cpu_usage = (cpu_length * this.cpuUsage) - 7;
    const new_ram_usage = ram_length * this.ramUsage;


    this.cpu_usage_path.nativeElement.style.strokeDashoffset = cpu_length - new_cpu_usage;
    // add 1 because it's physically a mirror of cpu usage
    this.ram_usage_path.nativeElement.style.strokeDashoffset = (ram_length + new_ram_usage);
  }

  ngOnInit() {
    // for hardcoding the usage chart stroke-dasharray to guard against animation glitches
    // console.log("cpu: ", this.cpu_usage_path.nativeElement.getTotalLength());
    // console.log("ram: ", this.ram_usage_path.nativeElement.getTotalLength());
  }

}
