// src/app/auth/auth.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, of } from 'rxjs';
import { tap } from 'rxjs/operators';

export interface LoginResponse {
  status: boolean;
  message: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly API = '/api';
  public loggedIn$ = new BehaviorSubject<boolean>(false);

  constructor(private http: HttpClient) {}

  signup(data: {
    name: string; email: string; password: string;
  }): Observable<string> {
    const payload = {
      username: data.name,
      email:    data.email,
      password: data.password
    };
    return this.http.post(
      `${this.API}/create-user`,
      payload,
      { responseType: 'text', withCredentials: true }
    );
  }

  login(creds: { username: string; password: string }): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(
      `${this.API}/login`,
      creds,
      { withCredentials: true }
    ).pipe(
      tap(res => this.loggedIn$.next(res.status))
    );
  }

  logout(): Observable<void> {
    return this.http.delete<void>(
      `${this.API}/logout`,
      { withCredentials: true }
    ).pipe(
      tap(() => this.loggedIn$.next(false))
    );
  }
}
