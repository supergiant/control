import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class AppsService {
    newModal = new Subject<any>();
    newEditModal = new Subject<any>();
    apps: any;
    selectedItems= new Array();

    constructor() {}

    // return all selected cloud accounts
    returnSelected(){
      return this.selectedItems
    }

    isChecked(item) {
      for (let obj of this.selectedItems) {
        if (item.id === obj.id) {return true}
      }
      return false
    }

    // Record/Delete a selection from the "selected items" array.
    selectItem(item,event){
     if (event) {
       this.selectedItems.push(item);
     } else {
       for (let obj of this.selectedItems) {
         if (item.id === obj.id) {
           this.selectedItems.splice(
             this.selectedItems.indexOf(obj), 1);
         }
       }
     }
   }
}
