import { Injectable } from '@angular/core';
import { Actions, Effect, ofType } from '@ngrx/effects';
import { AppsActionTypes } from './apps.actions';

@Injectable()
export class AppsEffects {

  @Effect()
  loadFoos$ = this.actions$.pipe(ofType(AppsActionTypes.LoadAppss));

  constructor(private actions$: Actions) {}
}
