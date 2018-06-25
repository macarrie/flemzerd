import { TestBed, inject } from '@angular/core/testing';

import { TvshowsService } from './tvshows.service';

describe('TvshowsService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [TvshowsService]
    });
  });

  it('should be created', inject([TvshowsService], (service: TvshowsService) => {
    expect(service).toBeTruthy();
  }));
});
