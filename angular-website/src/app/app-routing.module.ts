import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { AddCourseComponent } from '../app/components/add-course/add-course.component';
import { LoginComponent } from './components/login/login.component';
import { RegisterComponent } from './components/register/register.component';
import { AddProjectComponent } from './components/add-project/add-project.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { UnverifiedBlocksComponent } from './components/unverified-blocks/unverified-blocks.component';
import { GenerateCVComponent } from './components/generate-cv/generate-cv.component';
 
const routes: Routes = [
  {path: 'addCourse', component: AddCourseComponent},
  {path: 'addProject', component: AddProjectComponent},
  {path: 'unverifiedBlocks', component: UnverifiedBlocksComponent},
  {path: 'generateCV', component: GenerateCVComponent},
  {path: 'register', component: RegisterComponent},
  {path: 'login', component: LoginComponent},
  {path: 'dashboard', component: DashboardComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
