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
import { CommonModule }    from '@angular/common';
import { TranslateModule } from '@ngx-translate/core';
import { finalize }        from 'rxjs/operators';
import { AuthService }     from '../auth.service';
import { environment }     from '../../../environments/environment';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    TranslateModule
  ],
  templateUrl: './login.component.html'
})
export class LoginComponent implements OnInit {
  private fb     = inject(FormBuilder);
  private auth   = inject(AuthService);
  private router = inject(Router);
  private route  = inject(ActivatedRoute);

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

    this.auth.login({ username, password })
      .pipe(finalize(() => (this.loading = false)))
      .subscribe({
        next: (isAuth: boolean) => {
          if (isAuth) {
            // on récupère la route du dashboard depuis l'env
            const returnUrl =
              this.route.snapshot.queryParamMap.get('returnUrl')
                || environment.routes.dashboard;
            this.router.navigateByUrl(returnUrl);
          } else {
            this.error = 'LOGIN.ERROR.INVALID_CREDENTIALS';
          }
        },
        error: err => {
          this.error = typeof err.error === 'string'
            ? err.error.trim()
            : 'LOGIN.ERROR.LOGIN_FAILED';
        }
      });
  }
}
