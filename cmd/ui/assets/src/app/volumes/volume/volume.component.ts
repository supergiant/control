import { Component, Input} from '@angular/core';
import { VolumesService } from '../volumes.service';

@Component({
  selector: '[app-volume]',
  templateUrl: './volume.component.html',
  styleUrls: ['./volume.component.css']
})
export class VolumeComponent {
  @Input() volume: any;
  constructor(private volumesService: VolumesService) { }

  status(kube) {
    // Status logic needs to be added here.
      return "status status-ok"
   }
}
