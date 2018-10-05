import { Component, OnInit } from '@angular/core';
import { State } from '../../reducers';
import { Store } from '@ngrx/store';
import { LoadAppss } from '../apps/apps.actions';

@Component({
  selector: 'app-app-store',
  templateUrl: './app-store.component.html',
  styleUrls: ['./app-store.component.scss']
})
export class AppStoreComponent implements OnInit {

  constructor(private store: Store<State>) { }

  ngOnInit() {
    this.store.dispatch(new LoadAppss())
  }

}
