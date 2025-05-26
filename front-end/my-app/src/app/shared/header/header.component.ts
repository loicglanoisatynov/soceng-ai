import { Component, OnInit, Inject, PLATFORM_ID } from '@angular/core';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { LanguageService } from '../../core/language.service';
import { AuthService } from '../../auth/auth.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule],
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {
  public defaultAvatar = 'assets/images/bg-login.jpg';
  public isLoggedIn   = false;
  public menuOpen     = false;

  constructor(
    public auth: AuthService,
    public lang: LanguageService,
    private router: Router,
    @Inject(PLATFORM_ID) private platformId: any
  ) {}

  ngOnInit(): void {
    if (isPlatformBrowser(this.platformId)) {
      this.auth.loggedIn$.subscribe(status => this.isLoggedIn = status);
    }
  }

  get showAuthButtons(): boolean {
    if (!isPlatformBrowser(this.platformId)) return false;
    const url = this.router.url;
    return !this.isLoggedIn && !url.startsWith('/auth');
  }

  toggleMenu(): void {
    this.menuOpen = !this.menuOpen;
  }

  onLogout(): void {
    this.auth.logout().subscribe({
      next: () => this.navigateToLogin(),
      error: () => this.navigateToLogin()
    });
  }

  private navigateToLogin(): void {
    this.menuOpen = false;
    this.router.navigate(['/auth/login']);
  }

  switchLang(lang: string): void {
    this.lang.use(lang);
  }
}