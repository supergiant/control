import { Component } from '@angular/core';
import { VolumesService } from '../volumes.service';
import { Supergiant } from '../../shared/supergiant/supergiant.service'
import { VolumesComponent } from '../volumes.component'
import { Subscription } from 'rxjs/Subscription';
import { Notifications } from '../../shared/notifications/notifications.service'
import { SystemModalService } from '../../shared/system-modal/system-modal.service'
import { DropdownModalService } from '../../shared/dropdown-modal/dropdown-modal.service'
import { EditModalService } from '../../shared/edit-modal/edit-modal.service'
import { LoginComponent } from '../../login/login.component';



@Component({
  selector: 'app-volumes-header',
  templateUrl: './volumes-header.component.html',
  styleUrls: ['./volumes-header.component.css']
})
export class VolumesHeaderComponent {
  providersObj: any;

  constructor(
    private volumesService: VolumesService,
    private volumesComponant: VolumesComponent,
    private supergiant: Supergiant,
    private notifications: Notifications,
    private systemModalService: SystemModalService,
    private dropdownModalService: DropdownModalService,
    private editModalService: EditModalService,
    private loginComponent: LoginComponent,
    ) {}

  // After init, grab the schema
  ngAfterViewInit() {}

  openSystemModal(message){
      this.systemModalService.openSystemModal(message);
  }

  // If the delete button is hit, the seleted accounts are deleted.
  deleteVolume() {
    var selectedItems = this.volumesService.returnSelected()
    if (selectedItems.length === 0) {
      this.notifications.display("warn", "Warning:", "No Provider Selected.")
    } else if (selectedItems.length > 1) {
      this.notifications.display("warn", "Warning:", "You cannot edit more than one provider at a time.")
    } else {
    for(let provider of selectedItems){
      this.supergiant.KubeResources.delete(provider.id).subscribe(
        (data) => {
          if (data.status >= 200 && data.status <= 299) {
            this.notifications.display("success", "Volume: " + provider.name, "Deleted...")
            this.volumesComponant.getAccounts()
           }else{
            this.notifications.display("error", "Volume: " + provider.name, "Error:" + data.statusText)}},
        (err) => {
          if (err) {
            this.notifications.display("error", "Volume: " + provider.name, "Error:" + err)}},
      );
    }
  }
  }
}
