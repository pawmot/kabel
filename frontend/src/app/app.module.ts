import {BrowserModule} from '@angular/platform-browser';
import {Inject, NgModule} from '@angular/core';

import {AppComponent} from './app.component';
import {Astilectron, ASTILECTRON_TOKEN, AstilectronModule} from "./astilectron/astilectron.module";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {MaterialModule} from "./material/material.module";
import {Hotkey, HotkeyModule, HotkeysService} from "angular2-hotkeys";
import {ReactiveFormsModule} from "@angular/forms";
import {ConnectComponent} from './connect/connect.component';
import {RouterModule, Routes} from "@angular/router";

const routes: Routes = [
  {path: '', redirectTo: '/connect', pathMatch: 'full'},
  {path: 'connect', component: ConnectComponent}
];

@NgModule({
  declarations: [
    AppComponent,
    ConnectComponent
  ],
  imports: [
    BrowserModule,
    ReactiveFormsModule,
    AstilectronModule,
    BrowserAnimationsModule,
    MaterialModule,
    HotkeyModule.forRoot(),
    RouterModule.forRoot(routes)
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
  constructor(
    @Inject(ASTILECTRON_TOKEN) astilectron: Astilectron,
    hotkeysSvc: HotkeysService
  ) {
    hotkeysSvc.add(new Hotkey('f12', (ev): boolean => {
      astilectron.sendMessage({name: "devtools", payload: {}}, () => {
      });
      return false;
    }))
  }
}
