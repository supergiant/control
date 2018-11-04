import { TestBed } from '@angular/core/testing';

import { NodeProfileService } from './node-profile.service';

describe('NodeProfileService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: NodeProfileService = TestBed.get(NodeProfileService);
    expect(service).toBeTruthy();
  });
});
