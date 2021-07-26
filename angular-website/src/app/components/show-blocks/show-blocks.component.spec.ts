import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShowBlocksComponent } from './show-blocks.component';

describe('ShowBlocksComponent', () => {
  let component: ShowBlocksComponent;
  let fixture: ComponentFixture<ShowBlocksComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShowBlocksComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShowBlocksComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
