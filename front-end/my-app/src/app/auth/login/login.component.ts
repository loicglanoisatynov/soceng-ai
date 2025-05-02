// src/app/auth/login/login.component.ts
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { TranslateModule } from '@ngx-translate/core';
import { finalize } from 'rxjs/operators';
import { AuthService, LoginResponse } from '../auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    TranslateModule
  ],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  form!: FormGroup;
  loading = false;
  error   = '';

  constructor(
    private fb: FormBuilder,
    private auth: AuthService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.form = this.fb.group({
      username: ['', Validators.required],
      password: ['', Validators.required]
    });
  }

  submit(): void {
    if (this.form.invalid) {
      this.error = 'LOGIN.ERROR.FILL_FIELDS';
      return;
    }

    this.error   = '';
    this.loading = true;

    const { username, password } = this.form.value;
    this.auth.login({ username, password })
      .pipe(finalize(() => this.loading = false))
      .subscribe({
        next: (res: LoginResponse) => {
          if (res.status) {
            this.router.navigate(['/dashboard']);
          } else {
            this.error = res.message || 'LOGIN.ERROR.LOGIN_FAILED';
          }
        },
        error: err => {
          // si le back renvoie un texte d’erreur
          const msg = typeof err.error === 'string'
                    ? err.error.trim()
                    : 'LOGIN.ERROR.LOGIN_FAILED';
          this.error = msg;
        }
      });
  }
}
