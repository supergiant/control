import { Injectable } from '@angular/core';
import { Http, Response, Headers } from '@angular/http';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class CloudAccountsService {
    newModal = new Subject<any>();
    newEditModal = new Subject<any>();
    cloudAccounts: any;
    selectedItems= new Array();

    constructor(private http: Http) {}

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
      console.log("Deleting these items:" + this.selectedItems)
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
