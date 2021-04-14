import { Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from "@angular/forms";
import * as $ from 'jquery';

@Component({
  selector: 'app-add-course',
  templateUrl: './add-course.component.html',
  styleUrls: ['./add-course.component.css']
})
export class AddCourseComponent implements OnInit {

  User: any;
  submitted = false;
  playerForm: FormGroup;

  constructor(
    public fb: FormBuilder,
    private router: Router,
    private ngZone: NgZone,
    private blockService: BlockService,
  ) { 
    this.mainForm();
  }

  ngOnInit() {
    this.getUser();
    $(document).ready(function(){
      $('.nav_btn').click(function(){
        $('.mobile_nav_items').toggleClass('active');
      });
    });
  }
  
  getUser() {
    if (this.blockService.isLoggedIn == false)
    {
      alert("Access not allowed");
      this.router.navigateByUrl('/login')
    }
    else
    {   
      this.blockService.getUser().subscribe(
      (res) => {
        if (res['status'] == 200)
        {
          this.User = res['body']['result'];
        }
      }, (error) => {
        console.log(error);
        console.exception("ss");
      });
    }
  }

  mainForm() {
    this.playerForm = this.fb.group({
      courseCode: ['', [Validators.required]],
      courseName: ['', [Validators.required]],
      courseCHrs: ['', [Validators.required]],
      courseGrade: ['', [Validators.required]],
      courseEmail: ['', [Validators.required]],
      userEmail: ['', [Validators.required]],
      userPass: ['', [Validators.required]]
    })
  }
  
  get myForm(){
    return this.playerForm.controls;
  }

  onSubmit() {
   
    this.submitted = true;

    if (!this.playerForm.valid) {
      return false;
    } else {
      this.blockService.addCourse(this.playerForm.value).subscribe(
        (res) => {
            this.ngZone.run(() => this.router.navigateByUrl('/dashboard'))
        }, (error) => {
          console.log(error);
          console.exception("ss");
        });
    }
  }

  logout() {
    this.blockService.logout();
    this.ngZone.run(() => this.router.navigateByUrl('/login'))
  }

}
