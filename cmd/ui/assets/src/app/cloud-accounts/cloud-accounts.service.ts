import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { NotificationsService } from 'angular2-notifications';

@Injectable()
export class CloudAccountsService {
    newModal = new Subject<any>();
    newEditModal = new Subject<any>();
    cloudAccounts: any;
    selectedItems= new Array();

    constructor(private http: Http, private _service: NotificationsService) {}

    openNewCloudServiceModal(message){
      this.newModal.next(message);
    }
    openNewCloudServiceEditModal(message){
      this.newEditModal.next(message);
    }

    getCloudAccounts(){
      // TODO: Specify mech for hostname and auth
      let headers = new Headers();
      headers.append('Authorization', 'SGAPI token="iPIrhzlIHfRtHkRXW8NoAou0pUGFUfmo"');
      return this.http.get("http://localhost:8080/api/v0/cloud_accounts", { headers: headers })
    }

    // Delete all selected cloud accounts
    deleteCloudAccount(){

      if (this.selectedItems.length === 0) {
        this._service.error(
          'No Accounts Selected...',
          'You must select an account before it can be deleted.',
        {}
       )
       return;
      }

      console.log("Deleting these items:" + this.selectedItems);
      var error = false; // remove when api is connected.
      if (!error) {
        this._service.success(
          'Success...',
          'Deleted: ' + this.selectedItems,
        {}
       )
      }
      this.getCloudAccounts();
    }

    // Create a Cloud Account
    createCloudAccount(account){
      console.log(account);
      var error = true; // remove when api is connected.
      if (!error) {
        this._service.success(
          'Success...',
          'New Cloud Account Created.',
        {}
       )
     } else {
       this._service.error(
         'Error...',
         'You jacked something up.',
       {}
      )
      return error
     }
      this.getCloudAccounts();
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
