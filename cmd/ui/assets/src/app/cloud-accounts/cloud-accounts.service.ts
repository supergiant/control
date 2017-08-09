import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { NotificationsService } from 'angular2-notifications';

@Injectable()
export class CloudAccountsService {
    newModal = new Subject<any>();
    newEditModal = new Subject<any>();
    cloudAccounts: any;
    selectedItems= new Array();

    constructor(
      private _service: NotificationsService,
    ) {}

    openNewCloudServiceModal(message){
      this.newModal.next(message);
    }
    openNewCloudServiceEditModal(message, object){
      this.newEditModal.next([message, object]);
    }

    showNotification(kind, header, body) {
      switch (kind) {
        case "success": {
          this._service.success(header, body, {})
          break
        }
        case "error": {
          this._service.error(header, body, {})
          break
        }
        case "warn": {
          this._service.warn(header, body, {})
          break
        }
      }
    }

    // return all selected cloud accounts
    returnSelectedCloudAccount(){
      return this.selectedItems
    }

    ngOnInit() {
    }

    // Record/Delete a cloud account selection from the "selected items" array.
    selectItem(val,event){
     if (event) {
       this.selectedItems.push(val);
     } else {
       var index = this.selectedItems.indexOf(val);
         if (index > -1) {
          this.selectedItems.splice(index, 1);
         }
     }
   }
}
