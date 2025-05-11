// src/app/auth/signup/signup.component.ts
import { Component, OnInit } from '@angular/core';
import { CommonModule }      from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule
} from '@angular/forms';
import { Router, RouterModule }    from '@angular/router';
import { TranslateModule }         from '@ngx-translate/core';
import { finalize }                from 'rxjs/operators';
import { AuthService }             from '../auth.service';
import { environment }             from '../../../environments/environment';

@Component({
  selector: 'app-signup',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    TranslateModule
  ],
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})
export class SignupComponent implements OnInit {
  form!: FormGroup;
  loading = false;
  error   = '';

  constructor(
    private fb: FormBuilder,
    private auth: AuthService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.form = this.fb.group(
      {
        name:     ['', Validators.required],
        email:    ['', [Validators.required, Validators.email]],
        password: ['', [Validators.required, Validators.minLength(6)]],
        confirm:  ['', Validators.required]
      },
      { validators: this.passwordsMatch }
    );
  }

  private passwordsMatch(group: FormGroup) {
    const p = group.get('password')!.value;
    const c = group.get('confirm')!.value;
    return p === c ? null : { mismatch: true };
  }

  submit(): void {
    if (this.form.invalid) {
      this.error = 'SIGNUP.ERROR.FILL_FIELDS';
      return;
    }
    this.error   = '';
    this.loading = true;

    const { name, email, password } = this.form.value;
    this.auth
      .signup({ name, email, password })
      .pipe(finalize(() => (this.loading = false)))
      .subscribe({
        next: () => {
          // on success, navigate to login via environment.routes.login
          this.router.navigate(
            [environment.routes.login],
            { queryParams: { registered: 1 } }
          );
        },
        error: (err: any) => {
          this.error =
            typeof err.error === 'string'
              ? err.error.trim()
              : 'SIGNUP.ERROR.CREATE_FAILED';
        }
      });
  }
}
