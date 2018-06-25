import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DownloadingItemComponent } from './downloading-item.component';

describe('DownloadingItemComponent', () => {
  let component: DownloadingItemComponent;
  let fixture: ComponentFixture<DownloadingItemComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DownloadingItemComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DownloadingItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
