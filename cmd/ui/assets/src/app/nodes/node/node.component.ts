import { Component, Input } from '@angular/core';
import { NodesService } from '../nodes.service';

@Component({
  selector: '[app-node]', // tslint:disable-line
  templateUrl: './node.component.html',
  styleUrls: ['./node.component.scss']
})
export class NodeComponent {
  @Input() node: any;
  constructor(
    public nodesService: NodesService,
  ) { }

  round(value, precision?) {
    const multiplier = Math.pow(10, precision || 1);
    return Math.round(value * multiplier) / multiplier;
  }

  status(node) {
    if (node.status && node.status.error && node.status.retries === node.status.max_retries) {
      return 'status status-danger';
    } else if (node.status) {
      return 'status status-transitioning';
    } else if (node.passive_status && !node.passive_status_okay) {
      return 'status status-warning';
    } else {
      return 'status status-ok';
    }
  }
}
