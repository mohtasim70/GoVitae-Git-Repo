import { Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from "@angular/forms";

@Component({
  selector: 'app-search-cv',
  templateUrl: './search-cv.component.html',
  styleUrls: ['./search-cv.component.css']
})
export class SearchCVComponent implements OnInit {

  CVs: any;
  searchForm: FormGroup;
  searchCourse = "";

  constructor(
    public fb: FormBuilder,
    private router: Router,
    private ngZone: NgZone,
    private blockService: BlockService,
  ) { 
    this.mainForm();
  }

  ngOnInit() {
    this.findCV();
  }

  mainForm() {
    this.searchForm = this.fb.group({
      courseName: [''],
      courseGrade: ['']
    })
  }

  get myForm(){
    return this.searchForm.controls;
  }

  findCV() {
    this.searchForm.controls.courseName.setValue(sessionStorage.getItem("SearcherCourse"))
    this.searchForm.controls.courseGrade.setValue(sessionStorage.getItem("SearcherGrade"))
    this.blockService.getVerifiedCVs(this.searchForm.value).subscribe(
      (res) => {
        if (res['status'] == 200)
        {
          this.CVs = res['body']['users'];
          console.log(this.CVs)
        }
      }, (error) => {
        console.log(error);
        console.exception("ss");
      });
  }

}
