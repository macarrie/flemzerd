import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NotificationsCenterComponent } from './notifications-center.component';

describe('NotificationsCenterComponent', () => {
  let component: NotificationsCenterComponent;
  let fixture: ComponentFixture<NotificationsCenterComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NotificationsCenterComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NotificationsCenterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
