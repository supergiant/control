import { Component, OnDestroy } from '@angular/core';
import { VolumesService } from '../volumes.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service';
import { SystemModalService } from '../../shared/system-modal/system-modal.service';
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service';
import { EditModalService } from '../../shared/edit-modal/edit-modal.service';
import { LoginComponent } from '../../login/login.component';

@Component({
  selector: 'app-volumes-header',
  templateUrl: './volumes-header.component.html',
  styleUrls: ['./volumes-header.component.css']
})
export class VolumesHeaderComponent implements OnDestroy {
  subscriptions = new Subscription();
  constructor(
    private volumesService: VolumesService,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private dropdownModalService: DropdownModalService,
    private editModalService: EditModalService,
    private loginComponent: LoginComponent,
  ) { }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  openSystemModal(message) {
    this.systemModalService.openSystemModal(message);
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteVolume() {
    const selectedItems = this.volumesService.returnSelected();
    if (selectedItems.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Volume Selected.');
    } else {
      for (const provider of selectedItems) {
        this.subscriptions.add(this.supergiant.KubeResources.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Volume: ' + provider.name, 'Deleted...');
          },
          (err) => {
            this.notifications.display('error', 'Volume: ' + provider.name, 'Error:' + err);
          },
        ));
      }
    }
  }
}
