import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SearchCVComponent } from './search-cv.component';

describe('SearchCVComponent', () => {
  let component: SearchCVComponent;
  let fixture: ComponentFixture<SearchCVComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SearchCVComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SearchCVComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
