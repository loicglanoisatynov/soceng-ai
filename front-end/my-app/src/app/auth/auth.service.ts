// src/app/auth/auth.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { tap } from 'rxjs/operators';

export interface LoginResponse {
  status: boolean;
  message: string;
}

export interface UserProfile {
  id: number;              // Modifié : l'ID est de type number pour correspondre au dashboard
  username: string;
  email: string;
  avatarUrl?: string;
  score?: number;
  progress?: number;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly API = '/api';
  /** Statut de connexion */
  public loggedIn$ = new BehaviorSubject<boolean>(false);
  /** Profil utilisateur chargé */
  public profile$ = new BehaviorSubject<UserProfile | null>(null);

  constructor(private http: HttpClient) {}

  /** Getter pour faciliter l'accès au profil synchroniquement */
  get profile(): UserProfile | null {
    return this.profile$.value;
  }

  /** Inscription et retour du message du serveur en texte */
  signup(data: { name: string; email: string; password: string }): Observable<string> {
    return this.http.post<string>(
      `${this.API}/create-user`,
      {
        username: data.name,
        email: data.email,
        password: data.password
      },
      {
        responseType: 'text' as 'json',
        withCredentials: true
      }
    );
  }

  /** Authentification et mise à jour du statut */
  login(creds: { username: string; password: string }): Observable<LoginResponse> {
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

  /** Déconnexion et remise à zéro du statut et du profil */
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

  /** Chargement du profil utilisateur après connexion */
  loadProfile(): Observable<UserProfile> {
    return this.http.get<UserProfile>(
      `${this.API}/profile`,
      { withCredentials: true }
    ).pipe(
      tap(profile => this.profile$.next(profile))
    );
  }
}
