// src/app/core/language.service.ts
import { Injectable, Inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { TranslateService } from '@ngx-translate/core';

const STORAGE_KEY = 'app_language';

@Injectable({ providedIn: 'root' })
export class LanguageService {
  constructor(
    private translate: TranslateService,
    @Inject(PLATFORM_ID) private platformId: object
  ) {
    this.translate.addLangs(['fr', 'en']);
    this.translate.setDefaultLang('fr');

    let initial = 'fr';
    const browser = this.translate.getBrowserLang() ?? 'fr';
    const fallback = ['fr', 'en'].includes(browser) ? browser : 'fr';

    if (isPlatformBrowser(this.platformId)) {
      const saved = localStorage.getItem(STORAGE_KEY);
      initial = saved || fallback;
    } else {
      initial = fallback;
    }

    this.use(initial);
  }

  get currentLang(): string {
    return this.translate.currentLang;
  }

  use(lang: string) {
    this.translate.use(lang);
    if (isPlatformBrowser(this.platformId)) {
      localStorage.setItem(STORAGE_KEY, lang);
    }
  }
}
