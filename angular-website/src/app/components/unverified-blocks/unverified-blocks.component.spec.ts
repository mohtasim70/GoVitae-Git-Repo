import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UnverifiedBlocksComponent } from './unverified-blocks.component';

describe('UnverifiedBlocksComponent', () => {
  let component: UnverifiedBlocksComponent;
  let fixture: ComponentFixture<UnverifiedBlocksComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UnverifiedBlocksComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UnverifiedBlocksComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
