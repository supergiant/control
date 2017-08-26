import { Component, Input } from '@angular/core';
import { PodsService } from '../pods.service';

@Component({
  selector: '[app-pod]',  // tslint:disable-line
  templateUrl: './pod.component.html',
  styleUrls: ['./pod.component.css']
})
export class PodComponent {
  @Input() pod: any;
  constructor(private podsService: PodsService) { }

  status(pod) {
    if (pod.status && pod.status.error && pod.status.retries === pod.status.max_retries) {
      return 'status status-danger';
    } else if (pod.status) {
      return 'status status-transitioning';
    } else if (pod.passive_status && !pod.passive_status_okay) {
      return 'status status-warning';
    } else {
      return 'status status-ok';
    }
  }
}
