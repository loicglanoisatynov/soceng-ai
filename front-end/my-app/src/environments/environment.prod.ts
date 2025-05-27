// src/environments/environment.prod.ts

export const environment = {
    production: true,
    apiBaseUrl: 'http://localhost:8080/api',
    routes: {
      login:     '/auth/login',
      signup:    '/auth/signup',
      dashboard: '/dashboard'
    }
  };
  