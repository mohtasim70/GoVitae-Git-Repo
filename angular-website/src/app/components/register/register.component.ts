import { Router } from '@angular/router';
import { BlockService } from './../../service/block.service';
import { Component, OnInit, NgZone } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from "@angular/forms";

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {

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
  }

  mainForm() {
    this.registerForm = this.fb.group({
      username: ['', [Validators.required]],
      firstname: ['', [Validators.required]],
      lastname: ['', [Validators.required]],
      email: ['', [Validators.required, Validators.email]],
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
      this.blockService.register(this.registerForm.value).subscribe(
        (res) => {
          this.router.navigateByUrl('/login')
        }, (error) => {
          console.log(error);
          console.exception("ss");
        });
    }
  }
}
