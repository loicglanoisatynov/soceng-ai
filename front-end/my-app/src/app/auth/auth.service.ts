import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, of } from 'rxjs';
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
  biography?: string;
}

interface ApiResponse {
  status: boolean;
  message: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  // chemin absolu pour le proxy Angular
  private readonly API = 'http://localhost:8080';

  public loggedIn$ = new BehaviorSubject<boolean>(false);
  public profile$  = new BehaviorSubject<UserProfile | null>(null);

  constructor(private http: HttpClient) {}

  /** Accès synchrone au profil courant */
  get profile(): UserProfile | null {
    return this.profile$.value;
  }

  /** Inscription */
  signup(data: { name: string; email: string; password: string }): Observable<string> {
    return this.http.post<string>(
      `${this.API}/create-user`,
      { username: data.name, email: data.email, password: data.password },
      { responseType: 'text' as 'json', withCredentials: true }
    );
  }

  /**
   * Connexion :
   * - responseType:'text' pour ne pas forcer JSON.parse sur un 500 HTML
   * - observe:'response' pour capter le status
   */
  login(creds: { username: string; password: string }): Observable<LoginResponse> {
    console.log('🔐 Tentative login avec :', creds); // ✅ DEBUG FRONT
  
    return this.http.post<LoginResponse>(
      `${this.API}/login`,
      creds,
      { withCredentials: true }
    ).pipe(
      tap(res => {
        this.loggedIn$.next(res.status);
        if (res.status) {
          this.loadProfile().subscribe();
        }
      })
    );
  }
  

  /** Déconnexion */
  logout(): Observable<void> {
    return this.http
      .delete<void>(`${this.API}/logout`, { withCredentials: true })
      .pipe(tap(() => {
        this.loggedIn$.next(false);
        this.profile$.next(null);
      }));
  }

  /**
   * Chargement de profil.
   * Comme votre back n’a pas de /profile, on stub à partir du cookie.
   */
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
    return of(stub).pipe(tap(p => this.profile$.next(p)));
  }

  /** MAJ du profil public */
  updateProfile(data: { username: string; biography: string; avatar: string }): Observable<ApiResponse> {
    return this.http.put<ApiResponse>(
      `${this.API}/edit-profile`,
      data,
      { withCredentials: true }
    ).pipe(tap(res => {
      if (res.status && this.profile) {
        this.profile$.next({
          ...this.profile,
          username:  data.username,
          biography: data.biography,
          avatarUrl: data.avatar
        });
      }
    }));
  }

  /** MAJ des identifiants */
  updateUser(data: { email?: string; password: string; newpassword?: string }): Observable<ApiResponse> {
    return this.http.put<ApiResponse>(
      `${this.API}/edit-user`,
      data,
      { withCredentials: true }
    ).pipe(tap(res => {
      if (res.status && this.profile && data.email) {
        this.profile$.next({ ...this.profile, email: data.email });
      }
    }));
  }
}