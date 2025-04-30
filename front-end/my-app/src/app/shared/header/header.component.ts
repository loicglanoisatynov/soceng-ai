// src/app/shared/header/header.component.ts
import { Component, Inject, PLATFORM_ID } from '@angular/core';
import { RouterModule, Router } from '@angular/router';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { AuthService } from '../../auth/auth.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent {
  constructor(
    public auth: AuthService,
    private router: Router,
    @Inject(PLATFORM_ID) private platformId: any
  ) {}

  /** Affiche Login/SignUp quand on est en navigateur, non connecté, et hors des pages /auth */
  get showAuthButtons(): boolean {
    if (!isPlatformBrowser(this.platformId)) return false;
    const url = this.router.url;
    return !this.auth.isLoggedIn && !url.startsWith('/auth');
  }

  logout() {
    this.auth.logout();
    this.router.navigate(['/login']);
  }
}
