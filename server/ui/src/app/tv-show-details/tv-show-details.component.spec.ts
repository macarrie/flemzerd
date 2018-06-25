import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TvShowDetailsComponent } from './tv-show-details.component';

describe('TvShowDetailsComponent', () => {
  let component: TvShowDetailsComponent;
  let fixture: ComponentFixture<TvShowDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TvShowDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TvShowDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
