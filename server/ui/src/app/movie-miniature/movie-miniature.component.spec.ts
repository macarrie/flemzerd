import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MovieMiniatureComponent } from './movie-miniature.component';

describe('MovieMiniatureComponent', () => {
  let component: MovieMiniatureComponent;
  let fixture: ComponentFixture<MovieMiniatureComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MovieMiniatureComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MovieMiniatureComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
