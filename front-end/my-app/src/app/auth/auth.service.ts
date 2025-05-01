import { Injectable, Inject, PLATFORM_ID } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { isPlatformBrowser } from '@angular/common';
import { BehaviorSubject, throwError, Observable } from 'rxjs';
import { tap } from 'rxjs/operators';

interface LoginResponse {
  accessToken: string;
  refreshToken: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private API = 'https://api.ton-backend.com';
  private _token: string | null = null;
  public loggedIn$ = new BehaviorSubject<boolean>(false);

  constructor(
    private http: HttpClient,
    @Inject(PLATFORM_ID) private platformId: any
  ) {
    if (isPlatformBrowser(this.platformId)) {
      this._token = localStorage.getItem('access_token');
      this.loggedIn$.next(!!this._token);
    }
  }

  login(creds: { email: string; password: string }): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(`${this.API}/login`, creds).pipe(
      tap(res => {
        if (isPlatformBrowser(this.platformId)) {
          localStorage.setItem('access_token', res.accessToken);
          localStorage.setItem('refresh_token', res.refreshToken);
          this._token = res.accessToken;
          this.loggedIn$.next(true);
        }
      })
    );
  }

  signup(data: { name: string; email: string; password: string }): Observable<any> {
    return this.http.post<any>(`${this.API}/create-user`, data);
  }

  logout(): void {
    if (isPlatformBrowser(this.platformId)) {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
    this._token = null;
    this.loggedIn$.next(false);
  }

  get token(): string | null {
    return this._token;
  }

  get isLoggedIn(): boolean {
    return !!this._token;
  }

  refreshToken(): Observable<{ accessToken: string }> {
    if (!isPlatformBrowser(this.platformId)) {
      return throwError(() => new Error('Platform not supported'));
    }
    const refresh = localStorage.getItem('refresh_token');
    if (!refresh) {
      return throwError(() => new Error('No refresh token'));
    }
    return this.http
      .post<{ accessToken: string }>(`${this.API}/refresh-token`, { token: refresh })
      .pipe(
        tap(res => {
          localStorage.setItem('access_token', res.accessToken);
          this._token = res.accessToken;
          this.loggedIn$.next(true);
        })
      );
  }
}
