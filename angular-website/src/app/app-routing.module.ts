import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { AddCourseComponent } from '../app/components/add-course/add-course.component';
import { LoginComponent } from './components/login/login.component';
import { RegisterComponent } from './components/register/register.component';
import { AddProjectComponent } from './components/add-project/add-project.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { UnverifiedBlocksComponent } from './components/unverified-blocks/unverified-blocks.component';
import { GenerateCVComponent } from './components/generate-cv/generate-cv.component';
import { HomeComponent } from './components/home/home.component';
import { SearchCVComponent } from './components/search-cv/search-cv.component';
import { VerifyContentComponent } from './components/verify-content/verify-content.component';
import { RunServerComponent } from './components/run-server/run-server.component';
import { GetBlockComponent } from './components/get-block/get-block.component';
import { ShowBlocksComponent } from './components/show-blocks/show-blocks.component';
 
const routes: Routes = [
  {path: 'addCourse', component: AddCourseComponent},
  {path: 'addProject', component: AddProjectComponent},
  {path: 'unverifiedBlocks', component: UnverifiedBlocksComponent},
  {path: 'generateCV', component: GenerateCVComponent},
  {path: 'register', component: RegisterComponent},
  {path: 'login', component: LoginComponent},
  {path: '', component: HomeComponent},
  {path: 'runServer', component: RunServerComponent},
  {path: 'mineBlock/:hash', component: VerifyContentComponent},
  {path: 'searchCV', component: SearchCVComponent},
  {path: 'dashboard', component: DashboardComponent},
  {path: 'getBlock', component: GetBlockComponent},
  {path: 'showBlocks', component: ShowBlocksComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
