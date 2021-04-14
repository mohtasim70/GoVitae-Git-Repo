import { Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from "@angular/forms";
import { flatMap } from 'rxjs/operators';
declare const addcl: any;
declare const remcl: any;

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  submitted = false;
  registerForm: FormGroup;

  constructor(
    public fb: FormBuilder,
    private router: Router,
    private ngZone: NgZone,
    private blockService: BlockService,
  ) { 
    this.mainForm();
  }

  ngOnInit() {
    this.getAllUsers();
  }

  getAllUsers() {
    this.blockService.getAllUsers().subscribe(
      (res) => {
        if (res['status'] == 200)
        {            
          console.log(res);
        }
      }, (error) => {
        console.log(error);
        console.exception("ss");
      });
  }

  mainForm() {
    this.registerForm = this.fb.group({
      username: ['', [Validators.required]],
      password: ['', [Validators.required]],
    })
  }
  
  get myForm(){
    return this.registerForm.controls;
  }

  onSubmit() {
   
    this.submitted = true;

    if (!this.registerForm.valid) {
      return false;
    } else {
      this.blockService.login(this.registerForm.value).subscribe(
        (res) => {
          if (res['status'] == 200)
          {            
            if (res['body']['error'])
            {
              if (res['body']['error'] == "Invalid username")
              {
                alert("Invalid Username!");
              }
              if (res['body']['error'] == "Invalid password")
              {
                alert("Invalid Password!");
              }
            }
            else
            {
              console.log(res['body']['token'])
              localStorage.setItem("TOKEN", res['body']['token']);
              this.router.navigateByUrl('/dashboard');
            }
          }
        }, (error) => {
          alert("Invalid Username or Password!");          
          console.log(error);
          console.exception("ss");
        });
    }
  }
}
