import { Routes } from '@angular/router';
import { AuthGuard } from './auth/auth.guard';
import { SettingsComponent } from './../app/settings/settings/settings.component';
import { MyChallengeComponent } from './../app/dashboard/challenges/mychallenge/mychallenge.component';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },

  // alias pour /login et /signup
  { path: 'login', redirectTo: 'auth/login', pathMatch: 'full' },
  { path: 'signup', redirectTo: 'auth/signup', pathMatch: 'full' },

  // Public pages
  {
    path: 'home',
    loadComponent: () => import('./home/home/home.component').then(m => m.HomeComponent)
  },
  {
    path: 'about',
    loadComponent: () => import('./about/about/about.component').then(m => m.AboutComponent)
  },
  {
    path: 'contact',
    loadComponent: () => import('./contact/contact.component').then(m => m.ContactComponent)
  },
  {
    path: 'challenge',
    loadComponent: () => import('./challenge/challenge/challenge.component').then(m => m.ChallengeComponent)
  },

  // Auth routes
  {
    path: 'auth',
    children: [
      {
        path: 'login',
        loadComponent: () => import('./auth/login/login.component').then(m => m.LoginComponent)
      },
      {
        path: 'signup',
        loadComponent: () => import('./auth/signup/signup.component').then(m => m.SignupComponent)
      }
    ]
  },

  // Protected Dashboard with nested tab routes
  {
    path: 'dashboard',
    canActivate: [AuthGuard],
    loadComponent: () => import('./dashboard/dashboard/dashboard.component').then(m => m.DashboardComponent),
    children: [
      { path: '', redirectTo: 'details', pathMatch: 'full' },
      // Tab views rendered inside DashboardComponent
      { path: 'details', redirectTo: '/', pathMatch: 'full' },
      { path: 'settings', component: SettingsComponent },
      { path: 'challenges', component: MyChallengeComponent },
    ]
  },

  // Fallback
  { path: '**', redirectTo: 'home' }
];