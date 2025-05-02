// src/app/auth/auth.service.ts
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { BehaviorSubject, Observable, of } from 'rxjs';
import { tap, map, catchError } from 'rxjs/operators';

export interface UserProfile {
  id: number;
  username: string;
  email: string;
}

export interface LoginResponse {
  status: boolean;
  message: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly API = '/api';

  /** État de connexion */
  public loggedIn$ = new BehaviorSubject<boolean>(false);
  private _profile: UserProfile | null = null;

  constructor(private http: HttpClient) {}

  /**
   * POST /api/create-user
   * – on envoie exactement { username, email, password }
   * – on précise responseType: 'text' pour ne pas parser en JSON
   * – on ajoute withCredentials: true pour que le cookie (même si pas strictement nécessaire ici)
   *   soit correctement géré côté back
   */
  signup(data: { name: string; email: string; password: string }): Observable<string> {
    const payload = {
      username: data.name,
      email:    data.email,
      password: data.password
    };
    return this.http.post(
      `${this.API}/create-user`,
      payload,
      {
        headers: new HttpHeaders({ 'Content-Type': 'application/json' }),
        responseType: 'text',
        withCredentials: true
      }
    );
  }

  /**
   * POST /api/login
   * – on passe { username, password }
   * – on inclut withCredentials pour que Go lise/écrive ses cookies
   */
  login(creds: { username: string; password: string }): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(
      `${this.API}/login`,
      { username: creds.username, password: creds.password },
      { withCredentials: true }
    ).pipe(
      tap(res => {
        // Go renvoie { status: true, message: "Login successful" }
        this.loggedIn$.next(res.status);
      })
    );
  }

  /**
   * GET /api/profile
   * – valide la session via le cookie HTTP
   * – charge le profil
   */
  checkAuth(): Observable<boolean> {
    return this.http.get<UserProfile>(
      `${this.API}/profile`,
      { withCredentials: true }
    ).pipe(
      tap(profile => {
        this._profile = profile;
        this.loggedIn$.next(true);
      }),
      map(() => true),
      catchError(() => {
        this._profile = null;
        this.loggedIn$.next(false);
        return of(false);
      })
    );
  }

  /**
   * DELETE /api/logout
   * – supprime la session côté serveur
   */
  logout(): Observable<void> {
    return this.http.delete<void>(
      `${this.API}/logout`,
      { withCredentials: true }
    ).pipe(
      tap(() => {
        this._profile = null;
        this.loggedIn$.next(false);
      })
    );
  }

  /** Profil chargé après checkAuth */
  get profile(): UserProfile | null {
    return this._profile;
  }
}
