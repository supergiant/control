import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class CloudAccountsService {
    newModal = new Subject<any>();
    newEditModal = new Subject<any>();
    cloudAccounts: any;
    selectedItems= new Array();

    constructor() {}

    // msg to Cloud Service Modeal Dropdown to popup
    openNewCloudServiceDropdownModal(message){
      this.newModal.next(message);
    }

    // msg to Cloud Service New/Edit modal
    openNewCloudServiceEditModal(type, message, object){
      this.newEditModal.next([type, message, object]);
    }

    // return all selected cloud accounts
    returnSelectedCloudAccount(){
      return this.selectedItems
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
