import { Component, OnDestroy, OnInit } from '@angular/core';
import { NodesService } from './nodes.service';
import { Subscription } from 'rxjs/Subscription';
import { Supergiant } from '../shared/supergiant/supergiant.service'
import { Notifications } from '../shared/notifications/notifications.service'
import { Observable } from 'rxjs/Rx';


@Component({
  selector: 'app-nodes',
  templateUrl: './nodes.component.html',
  styleUrls: ['./nodes.component.css']
})
export class NodesComponent implements OnInit {
  private nodes = [];
  private subscription: Subscription;

  constructor(
    private nodesService: NodesService,
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  //get accouts when page loads
  ngOnInit() {
    this.getAccounts()
  }
  //get accounts
  getAccounts() {
    this.subscription = Observable.timer(0, 5000)
    .switchMap(() => this.supergiant.Nodes.get()).subscribe(
      (nodesObj) => {
        // Because of the check boxes we must reconsile the array.
        // If does not exist locally push it locally.
        for(let node of nodesObj.json().items) {
          var present = false
          for(let uiNodes of this.nodes) {
            if ( node.id === uiNodes.id ) {present = true}
          }
          if (!present) {this.nodes.push(node)}
         }

         // If does not exist on the API remove it locally.
         for(let uiNode of this.nodes) {
           var present = false
           for(let node of nodesObj.json().items) {
             if ( node.id === uiNode.id ) {present = true}
           }
           if (!present) {
             var index = this.nodes.indexOf(uiNode)
             this.nodes.splice(index, 1)}
          }
      },
      (err) => { this.notifications.display("warn", "Connection Issue.", err)});
  }

  ngOnDestroy(){
    this.subscription.unsubscribe();
  }
}
