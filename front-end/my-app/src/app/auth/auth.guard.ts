// src/app/auth/auth.guard.ts
import { Injectable } from '@angular/core';
import {
  CanActivate, Router, ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree
} from '@angular/router';
import { Observable, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { AuthService } from './auth.service';

@Injectable({ providedIn: 'root' })
export class AuthGuard implements CanActivate {
  constructor(private auth: AuthService, private router: Router) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<boolean|UrlTree> {
    return this.auth.checkAuth().pipe(
      map(isAuth => isAuth
        ? true
        : this.router.createUrlTree(['/auth/login'], { queryParams: { returnUrl: state.url } })
      ),
      catchError(() =>
        of(this.router.createUrlTree(['/auth/login'], { queryParams: { returnUrl: state.url } }))
      )
    );
  }
}