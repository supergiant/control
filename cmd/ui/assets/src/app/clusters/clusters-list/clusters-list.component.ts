import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Observable } from 'rxjs/Observable';
import { Supergiant } from '../../shared/supergiant/supergiant.service';
import { Notifications } from '../../shared/notifications/notifications.service';

@Component({
  selector: 'app-clusters-list',
  templateUrl: './clusters-list.component.html',
  styleUrls: ['./clusters-list.component.scss']
})
export class ClustersListComponent implements OnInit, OnDestroy {
  rows = [];
  selected = [];
  columns = [
    { prop: 'name' },
    { prop: 'master' },
    { prop: 'status' }
  ];
  private subscriptions = new Subscription();
  public kubes = [];
  constructor(
    private supergiant: Supergiant,
    private notifications: Notifications,
  ) { }

  ngOnInit() {
    this.getKubes();
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  onSelect({ selected }) {
    this.selected.splice(0, this.selected.length);
    this.selected.push(...selected);
  }

  getKubes() {
    this.subscriptions.add(Observable.timer(0, 5000)
      .switchMap(() => this.supergiant.Kubes.get()).subscribe(
      (kubes) => {
        this.rows = kubes.items.map(kube => ({
          id: kube.id, name: kube.name, master: kube.master_node_size, status: kube.status.description
        }));

        // Copy over any kubes that happen to be currently selected.
        this.selected.forEach((kube, index, array) => {
          for (const row of this.rows) {
            if (row.id === kube.id) {
              array[index] = row;
            }
          }
        });
      },
      (err) => { this.notifications.display('warn', 'Connection Issue.', err); }));
  }

  deleteKube() {
    if (this.selected.length === 0) {
      this.notifications.display('warn', 'Warning:', 'No Kube Selected.');
    } else {
      for (const provider of this.selected) {
        this.subscriptions.add(this.supergiant.Kubes.delete(provider.id).subscribe(
          (data) => {
            this.notifications.display('success', 'Kube: ' + provider.name, 'Deleted...');
            this.selected = [];
          },
          (err) => {
            this.notifications.display('error', 'Kube: ' + provider.name, 'Error:' + err);
          },
        ));
      }
    }
  }

}
