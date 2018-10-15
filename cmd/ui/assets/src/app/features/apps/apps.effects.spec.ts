import { TestBed, inject } from '@angular/core/testing';
import { provideMockActions } from '@ngrx/effects/testing';
import { Observable } from 'rxjs';

import { AppsEffects } from './apps.effects';

describe('AppsEffects', () => {
  let actions$: Observable<any>;
  let effects: AppsEffects;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        AppsEffects,
        provideMockActions(() => actions$)
      ]
    });

    effects = TestBed.get(AppsEffects);
  });

  it('should be created', () => {
    expect(effects).toBeTruthy();
  });
});
