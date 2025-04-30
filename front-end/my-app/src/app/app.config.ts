// src/app/app.config.ts
import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideClientHydration, withEventReplay } from '@angular/platform-browser';
import {
  provideHttpClient,
  withInterceptorsFromDi
} from '@angular/common/http';
import { HTTP_INTERCEPTORS } from '@angular/common/http';

import { routes } from './app.routes';
import { TokenInterceptor } from './core/token.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),

    // 1) le router
    provideRouter(routes),

    // 2) HTTP client + intercepteur
    provideHttpClient(
      withInterceptorsFromDi()  // prend en compte tous les HTTP_INTERCEPTORS enregistrés
    ),
    { 
      provide: HTTP_INTERCEPTORS, 
      useClass: TokenInterceptor, 
      multi: true 
    },

    // 3) SSR hydration (si tu utilises le SSR)
    provideClientHydration(withEventReplay())
  ]
};