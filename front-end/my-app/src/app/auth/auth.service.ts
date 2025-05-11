import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { tap } from 'rxjs/operators';

export interface LoginResponse {
  status: boolean;
  message: string;
}

export interface UserProfile {
  id: number;
  username: string;
  email: string;
  avatarUrl?: string;
  score?: number;
  progress?: number;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly API = 'http://localhost:8080';

  public loggedIn$ = new BehaviorSubject<boolean>(false);
  public profile$  = new BehaviorSubject<UserProfile | null>(null);

  constructor(private http: HttpClient) {}

  get profile(): UserProfile | null {
    return this.profile$.value;
  }

  signup(data: { name: string; email: string; password: string }): Observable<string> {
    return this.http.post<string>(
      `${this.API}/create-user`,
      { username: data.name, email: data.email, password: data.password },
      { responseType: 'text' as 'json', withCredentials: true }
    );
  }

  login(creds: { username: string; password: string }): Observable<LoginResponse> {
    const headers = new HttpHeaders().set('Content-Type', 'application/json');

    return this.http.post<LoginResponse>(
      `${this.API}/login`,
      creds,
      { headers, withCredentials: true }
    ).pipe(
      tap(res => {
        this.loggedIn$.next(res.status);
        if (res.status) {
          this.loadProfile().subscribe();
        }
      })
    );
  }

  logout(): Observable<void> {
    return this.http.delete<void>(
      `${this.API}/logout`,
      { withCredentials: true }
    ).pipe(
      tap(() => {
        this.loggedIn$.next(false);
        this.profile$.next(null);
      })
    );
  }

  loadProfile(): Observable<UserProfile> {
    return this.http.get<UserProfile>(
      `${this.API}/profile`,
      { withCredentials: true }
    ).pipe(
      tap(profile => this.profile$.next(profile))
    );
  }
}