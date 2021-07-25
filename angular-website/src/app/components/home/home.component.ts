import { Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from "@angular/forms";

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  
  searchCourse = "";
  searchType = "";
  Courses: any;
  keyword = "courseName";
  grades = ['A+', 'A', 'A-', 'B+', 'B', 'B-', 'C+', 'C', 'D'];
  

  constructor(
    public fb: FormBuilder,
    private router: Router,
    private ngZone: NgZone,
    private blockService: BlockService,
  ) { 
  }

  ngOnInit() {
    this.getCourses();
  }

  gradeType(type) {
    this.searchType = type;
  }

  selectEvent(item) {
    this.searchCourse = item.courseName;
  }
 
  onChangeSearch(val: string) {
    // fetch remote data from here
    // And reassign the 'data' which is binded to 'data' property.
  }
  
  onFocused(e){
    // do something when input is focused
  }

  getCourses() {
    this.blockService.getVerifiedContent().subscribe(
      (res) => {
        if (res['status'] == 200)
        {
          console.log(res)
          this.Courses = res['body']['courses'];
          console.log(this.Courses)
        }
      }, (error) => {
        console.log(error);
        console.exception("ss");
      });
  }

  searchCV() {
    sessionStorage.setItem("SearcherCourse", this.searchCourse);
    sessionStorage.setItem("SearcherGrade", this.searchType);
    this.ngZone.run(() => this.router.navigateByUrl('/searchCV'))
  }

}
