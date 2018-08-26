import {InjectionToken, NgModule} from '@angular/core';
import { CommonModule } from '@angular/common';

declare var astilectron: Astilectron;

export interface AstilectronMessage {
  name: string;
  payload: any;
}

export interface Astilectron {
  sendMessage(message: AstilectronMessage, callback: (AstilectronMessage) => void)
}

export const ASTILECTRON_TOKEN = new InjectionToken("ASTILECTRON_TOKEN", {
  providedIn: 'root',
  factory: () => astilectron
});

@NgModule({
  imports: [
    CommonModule
  ],
  exports: [],
  declarations: []
})
export class AstilectronModule { }
