import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { tap, map } from 'rxjs/operators';
import { environment } from '../../environments/environment';

export interface UserProfile {
  id: number;
  username: string;
  email: string;
  avatarUrl?: string;
  score?: number;
  progress?: number;
}

export interface LoginResponse {
  status: boolean;
  message: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  // Utilisation de l'URL de base définie dans environment
  private readonly API = environment.apiBaseUrl;

  public loggedIn$: BehaviorSubject<boolean>;
  private _profile: UserProfile | null = null;

  constructor(private http: HttpClient) {
    const saved = typeof window !== 'undefined' &&
      localStorage.getItem('isLoggedIn') === 'true';
    this.loggedIn$ = new BehaviorSubject<boolean>(saved);
  }

  /** SIGN UP */
  signup(data: { name: string; email: string; password: string }): Observable<void> {
    return this.http.post<void>(
      `${this.API}/create-user`,
      { username: data.name, email: data.email, password: data.password },
      { withCredentials: true }
    );
  }

  /** LOGIN : stocke état, profil minimal, persistance */
  login(creds: { username: string; password: string }): Observable<boolean> {
    return this.http.post<LoginResponse>(
      `${this.API}/login`,
      creds,
      { withCredentials: true }
    ).pipe(
      tap(res => {
        this.loggedIn$.next(res.status);
        if (res.status) {
          this._profile = { id: 0, username: creds.username, email: '' };
          if (typeof window !== 'undefined') {
            localStorage.setItem('isLoggedIn', 'true');
            localStorage.setItem('username', creds.username);
          }
        }
      }),
      map(res => res.status)
    );
  }

  /** LOGOUT : récupère username, envoie body, reset état */
  logout(): Observable<void> {
    const username =
      typeof window !== 'undefined'
        ? localStorage.getItem('username') || ''
        : '';
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.request<void>('DELETE', `${this.API}/logout`, {
      headers,
      withCredentials: true,
      body: { username }
    }).pipe(
      tap(() => {
        this.loggedIn$.next(false);
        this._profile = null;
        if (typeof window !== 'undefined') {
          localStorage.removeItem('isLoggedIn');
          localStorage.removeItem('username');
        }
      })
    );
  }

  /** Getter profil */
  get profile(): UserProfile | null {
    return this._profile;
  }
}
