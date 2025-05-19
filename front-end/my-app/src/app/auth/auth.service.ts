// src/app/auth/auth.service.ts
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams, HttpResponse } from '@angular/common/http';
import { BehaviorSubject, Observable, of } from 'rxjs';
import { map, tap } from 'rxjs/operators';

/** Réponses communes de l'API */
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
  biography?: string;
}

interface ApiResponse {
  status: boolean;
  message: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  /** Préfixe API pour le proxy Angular */
  private readonly API = '/api';

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
    console.log('🔐 Tentative login avec :', creds);

    const body = new HttpParams()
      .set('username', creds.username)
      .set('password', creds.password);

    return this.http.post(
      `/login`,               // chemin relatif pour le proxy Angular
      body.toString(),
      {
        headers: new HttpHeaders({
          'Content-Type': 'application/x-www-form-urlencoded'
        }),
        withCredentials: true,
        responseType: 'text',   // on attend du texte (HTML ou JSON)
        observe: 'response'     // on veut HttpResponse<string>
      }
    ).pipe(
      map((resp: HttpResponse<string>) => {
        // 200 = succès, 302 + suivi redirection OR 200 index.html => on considère que c'est OK
        const ok = resp.status >= 200 && resp.status < 300;
        return {
          status: ok,
          message: ok ? 'Connexion réussie' : 'Échec de la connexion'
        };
      }),
      tap(res => {
        this.loggedIn$.next(res.status);
        if (res.status) {
          this.loadProfile().subscribe();
        }
      })
    );
  }

  logout(): Observable<ApiResponse> {
    return this.http.delete<ApiResponse>(
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
    const username = document.cookie
      .split('; ')
      .find(r => r.startsWith('socengai-username='))
      ?.split('=')[1] ?? 'Utilisateur';

    const stub: UserProfile = {
      id:        0,
      username,
      email:     '',
      avatarUrl: '',
      score:     0,
      progress:  0,
      biography: ''
    };

    return of(stub).pipe(
      tap(p => this.profile$.next(p))
    );
  }

  /** Mise à jour du profil public */
  updateProfile(data: { username: string; biography: string; avatar: string }) {
    return this.http.post<ApiResponse>(
      `${this.API}/edit-profiles`,
      data,
      { withCredentials: true }
    );
  }

  /** Mise à jour des identifiants */
  updateUser(data: { email?: string; password: string; newpassword?: string }): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(
      `${this.API}/edit-user`,
      data,
      { withCredentials: true }
    ).pipe(
      tap(res => {
        if (res.status && this.profile && data.email) {
          this.profile$.next({ ...this.profile, email: data.email });
        }
      })
    );
  }
}
