import { NgModule } from '@angular/core';
import {MatButtonModule, MatCheckboxModule, MatInputModule} from "@angular/material";

const matModules = [
  MatCheckboxModule, MatInputModule, MatButtonModule
];

@NgModule({
  imports: matModules,
  exports: matModules
})
export class MaterialModule { }
