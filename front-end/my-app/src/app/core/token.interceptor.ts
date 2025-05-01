import { Injectable } from '@angular/core';
import {
  HttpInterceptor,
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpErrorResponse
} from '@angular/common/http';
import { Observable, throwError, BehaviorSubject } from 'rxjs';
import { catchError, switchMap, filter, take } from 'rxjs/operators';
import { AuthService } from '../auth/auth.service';

@Injectable()
export class TokenInterceptor implements HttpInterceptor {
  private refreshInProgress = false;
  private refreshSubject = new BehaviorSubject<string | null>(null);

  constructor(private auth: AuthService) {}

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const token = this.auth.token;
    const authReq = token
      ? req.clone({ setHeaders: { Authorization: `Bearer ${token}` } })
      : req;

    return next.handle(authReq).pipe(
      catchError(err => {
        if (err instanceof HttpErrorResponse && err.status === 401) {
          return this.handle401(authReq, next);
        }
        return throwError(() => err);
      })
    );
  }

  private handle401(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    if (!this.refreshInProgress) {
      this.refreshInProgress = true;
      this.refreshSubject.next(null);

      return this.auth.refreshToken().pipe(
        switchMap(res => {
          this.refreshInProgress = false;
          this.refreshSubject.next(res.accessToken);
          return next.handle(
            req.clone({ setHeaders: { Authorization: `Bearer ${res.accessToken}` } })
          );
        }),
        catchError(err => {
          this.refreshInProgress = false;
          this.auth.logout();
          return throwError(() => err);
        })
      );
    } else {
      return this.refreshSubject.pipe(
        filter(tok => tok != null),
        take(1),
        switchMap(tok =>
          next.handle(req.clone({ setHeaders: { Authorization: `Bearer ${tok}` } }))
        )
      );
    }
  }
}
