import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VerifyContentComponent } from './verify-content.component';

describe('VerifyContentComponent', () => {
  let component: VerifyContentComponent;
  let fixture: ComponentFixture<VerifyContentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VerifyContentComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VerifyContentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
