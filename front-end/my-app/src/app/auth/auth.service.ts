// src/app/auth/auth.service.ts
import { Injectable, Inject, PLATFORM_ID } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { tap } from 'rxjs/operators';
import { isPlatformBrowser } from '@angular/common';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private API = 'https://api.ton-backend.com';
  private _token: string | null = null;

  constructor(
    private http: HttpClient,
    @Inject(PLATFORM_ID) private platformId: any
  ) {
    if (isPlatformBrowser(this.platformId)) {
      this._token = localStorage.getItem('token');
    }
  }

  login(creds: { email: string; password: string }) {
    return this.http
      .post<{ token: string }>(`${this.API}/login`, creds)
      .pipe(
        tap(res => {
          if (isPlatformBrowser(this.platformId)) {
            localStorage.setItem('token', res.token);
            this._token = res.token;
          }
        })
      );
  }

  signup(data: any) {
    return this.http.post(`${this.API}/signup`, data);
  }

  logout() {
    if (isPlatformBrowser(this.platformId)) {
      localStorage.removeItem('token');
    }
    this._token = null;
  }

  get token(): string | null {
    return this._token;
  }

  get isLoggedIn(): boolean {
    return !!this._token;
  }
}
