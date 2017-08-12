import { Component, Input} from '@angular/core';
import { PodsService } from '../pods.service';

@Component({
  selector: '[app-pod]',
  templateUrl: './pod.component.html',
  styleUrls: ['./pod.component.css']
})
export class PodComponent {
  @Input() pod: any;
  constructor(private podsService: PodsService) { }
}
