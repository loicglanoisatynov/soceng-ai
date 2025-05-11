import { Component, OnInit, inject } from '@angular/core';
import { CommonModule }               from '@angular/common';
import { FormBuilder, ReactiveFormsModule, FormGroup } from '@angular/forms';
import { TranslateModule }            from '@ngx-translate/core';
import { Subscription }               from 'rxjs';

import { AuthService, UserProfile }   from '../../auth/auth.service';
import { ProfileHeroComponent }       from '../../shared/profile-hero/profile-hero.component';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule,
    ProfileHeroComponent
  ],
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
  private fb   = inject(FormBuilder);
  private auth = inject(AuthService);

  settingsForm!: FormGroup;
  profile?: UserProfile & { avatarUrl: string; score: number; progress: number };
  private sub!: Subscription;

  ngOnInit(): void {
    // Initialiser le formulaire de settings
    this.settingsForm = this.fb.group({
      emailNotifications: [false],
      darkMode:           [false]
    });

    // S'abonner au même profile$ que ProfileHeroComponent
    this.sub = this.auth.profile$.subscribe(p => {
      if (p) {
        this.profile = {
          ...p,
          avatarUrl: p.avatarUrl || '/assets/images/bg-login.jpg',
          score:     p.score     || 0,
          progress:  p.progress  || 0
        };
      }
    });

    // Si jamais on n'a pas encore chargé le profil (rafraîchissement direct)
    if (!this.auth.profile) {
      this.auth.loadProfile().subscribe({ error: () => {/* ignore */} });
    }
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  saveSettings(): void {
    if (this.settingsForm.valid) {
      console.log('Settings saved', this.settingsForm.value);
      // TODO → appel API pour sauvegarder
    }
  }
}
