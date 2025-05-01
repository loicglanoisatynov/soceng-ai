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
  isLoggedIn = false;

  constructor(
    public auth: AuthService,
    private router: Router,
    public lang: LanguageService,
    @Inject(PLATFORM_ID) private platformId: any
  ) {}

  ngOnInit(): void {
    if (isPlatformBrowser(this.platformId)) {
      this.auth.loggedIn$.subscribe(status => (this.isLoggedIn = status));
    }
  }

  get showAuthButtons(): boolean {
    if (!isPlatformBrowser(this.platformId)) return false;
    const url = this.router.url;
    return !this.isLoggedIn && !url.startsWith('/auth');
  }

  logout(): void {
    this.auth.logout();
    this.router.navigate(['/auth/login']);
  }

  switchLang(lang: string) {
    this.lang.use(lang);
  }
}
