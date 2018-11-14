import { Component, OnInit, Input, OnDestroy, ViewChild, ChangeDetectionStrategy } from '@angular/core';
import { Subscription, timer as observableTimer } from 'rxjs';
import { switchMap } from 'rxjs/operators';

import { Supergiant } from '../../shared/supergiant/supergiant.service';

@Component({
  selector: 'usage-orb',
  templateUrl: './usage-orb.component.html',
  styleUrls: ['./usage-orb.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class UsageOrbComponent implements OnInit, OnDestroy {

  constructor( private supergiant: Supergiant ) { }

  @Input() cluster: any;
  @ViewChild("cpu_usage_path") cpu_usage_path;
  @ViewChild("ram_usage_path") ram_usage_path;



  private subscriptions = new Subscription;

  updateMetrics(metrics) {

      const cpu_length = this.cpu_usage_path.nativeElement.getTotalLength();
      const ram_length = this.ram_usage_path.nativeElement.getTotalLength();

      // minus 7 for a strange calibration issue
      const new_cpu_usage = (cpu_length * metrics.cpu) - 7;
      const new_ram_usage = ram_length * metrics.memory;

      this.cpu_usage_path.nativeElement.style.strokeDashoffset = cpu_length - new_cpu_usage;
      // add 1 because it's physically a mirror of cpu usage
      this.ram_usage_path.nativeElement.style.strokeDashoffset = (ram_length + new_ram_usage);
    }

  getMetrics() {
    this.subscriptions.add(observableTimer(0, 10000).pipe(
      switchMap(() => this.supergiant.Kubes.getClusterMetrics(this.cluster.id))).subscribe(
        res => this.updateMetrics(res),
        err => console.error(err)
      )
    )
  }

  ngOnInit() {
    if (this.cluster.state == "operational") {
      this.getMetrics();
    }
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe()
  }

}
