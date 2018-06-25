import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TvShowMiniatureComponent } from './tv-show-miniature.component';

describe('TvShowMiniatureComponent', () => {
  let component: TvShowMiniatureComponent;
  let fixture: ComponentFixture<TvShowMiniatureComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TvShowMiniatureComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TvShowMiniatureComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
