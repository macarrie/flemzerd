import { TestBed, inject } from '@angular/core/testing';

import { FanartService } from './fanart.service';

describe('FanartService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [FanartService]
    });
  });

  it('should be created', inject([FanartService], (service: FanartService) => {
    expect(service).toBeTruthy();
  }));
});
