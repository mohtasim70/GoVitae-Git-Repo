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
  grades = ['A', 'A-', 'B+', 'B', 'B-', 'C+', 'C', 'D'];
  

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

  changeType(type) {
    console.log(type)
    if (type == "cName")
    {
      this.keyword = "courseName"
      document.getElementById("showGrade").style.display = "none";
      document.getElementById("showName").style.display = "block";
      document.getElementById("typer").style.marginLeft = "-340px";
    }
    if (type == "cGrade")
    {
      this.keyword = "courseGrade"
      document.getElementById("showGrade").style.display = "block";
      document.getElementById("showName").style.display = "none";
      document.getElementById("typer").style.marginLeft = "-150px";
      this.searchCourse = 'A+';
    }
  }

  gradeType(type) {
    this.searchCourse = type;
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
    sessionStorage.setItem("Searcher", this.searchCourse);
    this.ngZone.run(() => this.router.navigateByUrl('/searchCV'))
  }

}
