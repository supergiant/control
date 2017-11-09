import { Component, Input } from '@angular/core';
import { VolumesService } from '../volumes.service';

@Component({
  selector: '[app-volume]', // tslint:disable-line
  templateUrl: './volume.component.html',
  styleUrls: ['./volume.component.scss']
})
export class VolumeComponent {
  @Input() volume: any;
  constructor(public volumesService: VolumesService) { }

  status(volume) {
    if (volume.status && volume.status.error && volume.status.retries === volume.status.max_retries) {
      return 'status status-danger';
    } else if (volume.status) {
      return 'status status-transitioning';
    } else if (volume.passive_status && !volume.passive_status_okay) {
      return 'status status-warning';
    } else {
      return 'status status-ok';
    }
  }
}
