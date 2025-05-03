// src/app/auth/login/login.component.ts
import { Component, OnInit, inject } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule
} from '@angular/forms';
import {
  Router,
  ActivatedRoute,
  RouterModule
} from '@angular/router';
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
  private fb = inject(FormBuilder);
  private auth = inject(AuthService);
  private router = inject(Router);
  private route = inject(ActivatedRoute);

  form!: FormGroup;
  loading = false;
  error   = '';

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
    // On lit le returnUrl (ex: /dashboard) ou on met une valeur par défaut
    const returnUrl = this.route.snapshot.queryParams['returnUrl'] || '/dashboard';

    this.auth.login({ username, password })
      .pipe(finalize(() => (this.loading = false)))
      .subscribe({
        next: (res: LoginResponse) => {
          if (res.status) {
            // redirige vers la page d’origine ou /dashboard
            this.router.navigateByUrl(returnUrl);
          } else {
            this.error = res.message || 'LOGIN.ERROR.LOGIN_FAILED';
          }
        },
        error: err => {
          const msg = typeof err.error === 'string'
                    ? err.error.trim()
                    : 'LOGIN.ERROR.LOGIN_FAILED';
          this.error = msg;
        }
      });
  }
}
