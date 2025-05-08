// src/app/app.routes.ts
import { Routes } from '@angular/router';
import { AuthGuard } from './auth/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },

  // alias pour /login et /signup
  { path: 'login', redirectTo: 'auth/login', pathMatch: 'full' },
  { path: 'signup', redirectTo: 'auth/signup', pathMatch: 'full' },

  // Public pages
  {
    path: 'home',
    loadComponent: () =>
      import('./home/home/home.component').then(m => m.HomeComponent)
  },
  {
    path: 'about',
    loadComponent: () =>
      import('./about/about/about.component').then(m => m.AboutComponent)
  },
  {
    path: 'contact',
    loadComponent: () =>
      import('./contact/contact.component').then(m => m.ContactComponent)
  },
  {
    path: 'challenge',
    loadComponent: () =>
      import('./challenge/challenge/challenge.component').then(m => m.ChallengeComponent)
  },

  // Auth
  {
    path: 'auth',
    children: [
      {
        path: 'login',
        loadComponent: () =>
          import('./auth/login/login.component').then(m => m.LoginComponent)
      },
      {
        path: 'signup',
        loadComponent: () =>
          import('./auth/signup/signup.component').then(m => m.SignupComponent)
      }
    ]
  },

  // Protected pages
  {
    path: 'dashboard',
    canActivate: [AuthGuard],
    loadComponent: () =>
      import('./dashboard/dashboard/dashboard.component').then(m => m.DashboardComponent)
  },
  {
    path: 'settings',
    canActivate: [AuthGuard],
    loadComponent: () =>
      import('./settings/settings/settings.component').then(m => m.SettingsComponent)
  },

  // Wildcard
  { path: '**', redirectTo: 'home' }
];
