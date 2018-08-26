import {Component, Inject, OnInit} from '@angular/core';
import {Astilectron, ASTILECTRON_TOKEN} from "../astilectron/astilectron.module";
import {FormBuilder, FormControl, Validators} from "@angular/forms";

@Component({
  selector: 'app-connect',
  templateUrl: './connect.component.html',
  styleUrls: ['./connect.component.scss']
})
export class ConnectComponent implements OnInit {
  form = this.fb.group({
    ssh: this.fb.group({
      useSsh: [false],
      sshUserAndHostname: this.fb.control({value: '', disabled: true})
    }),
    dockerHost: ['']
  });

  constructor(@Inject(ASTILECTRON_TOKEN) private astilectron: Astilectron, private fb: FormBuilder) {
  }

  ngOnInit(): void {
    this.getUseSshFormControl().valueChanges.subscribe((useSsh: boolean) => {
      let userAndHostnameControl = this.getSshUserAndHostnameFormControl();
      if (useSsh) {
        userAndHostnameControl.enable();
        userAndHostnameControl.setValidators([Validators.required]);
      } else {
        userAndHostnameControl.disable();
        if (userAndHostnameControl.value === '') {
          userAndHostnameControl.markAsPristine();
          userAndHostnameControl.markAsUntouched();
        }
        userAndHostnameControl.setValidators([]);
      }
      userAndHostnameControl.updateValueAndValidity();
    });
  }

  connect() {
    if (this.form.invalid) {
      this.getSshUserAndHostnameFormControl().markAsTouched();
      return;
    }

    this.astilectron.sendMessage({
      name: "connection_spec",
      payload: {
        sshUserAndHostname: this.getSshUserAndHostnameFormControl().value,
        dockerHost: this.getDockerHostFormControl().value
      }
    }, () => {
    });
  }

  getUseSshFormControl(): FormControl {
    return this.form.get('ssh.useSsh') as FormControl;
  }

  getSshUserAndHostnameFormControl(): FormControl {
    return this.form.get('ssh.sshUserAndHostname') as FormControl;
  }

  getDockerHostFormControl(): FormControl {
    return this.form.get('dockerHost') as FormControl;
  }
}
