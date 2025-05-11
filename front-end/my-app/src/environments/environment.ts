// src/environments/environment.ts

export const environment = {
  production: false,
  apiBaseUrl: 'http://localhost:8080/api',
  routes: {
    login:     '/auth/login',
    signup:    '/auth/signup',
    dashboard: '/dashboard'
  }
};
