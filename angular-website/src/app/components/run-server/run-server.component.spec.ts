import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RunServerComponent } from './run-server.component';

describe('RunServerComponent', () => {
  let component: RunServerComponent;
  let fixture: ComponentFixture<RunServerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RunServerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RunServerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
