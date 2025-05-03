// src/app/auth/auth.guard.ts
import { Injectable } from '@angular/core';
import {
  CanActivate, Router,
  ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree
} from '@angular/router';
import { Observable } from 'rxjs';
import { map, take } from 'rxjs/operators';
import { AuthService } from './auth.service';

@Injectable({ providedIn: 'root' })
export class AuthGuard implements CanActivate {
  constructor(
    private auth: AuthService,
    private router: Router
  ) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<boolean|UrlTree> {
    return this.auth.loggedIn$.pipe(
      take(1),
      map(isIn => {
        if (isIn) return true;
        return this.router.createUrlTree(
          ['/auth/login'],
          { queryParams: { returnUrl: state.url } }
        );
      })
    );
  }
}
