import { Component, Input} from '@angular/core';
import { NodesService } from '../nodes.service';

@Component({
  selector: '[app-node]',
  templateUrl: './node.component.html',
  styleUrls: ['./node.component.css']
})
export class NodeComponent {
  @Input() node: any;
  constructor(private nodesService: NodesService) { }

  status(kube) {
    // Status logic needs to be added here.
      return "status status-ok"
   }
}
