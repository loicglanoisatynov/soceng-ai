// src/app/settings/settings/settings.component.ts
import { Component, OnInit, inject } from '@angular/core';
import { CommonModule }               from '@angular/common';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { TranslateModule }            from '@ngx-translate/core';

import { AuthService, UserProfile }      from '../../auth/auth.service';
import { ProfileHeroComponent }          from '../../shared/profile-hero/profile-hero.component';

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

  profile!: UserProfile & { avatarUrl: string; score: number; progress: number };
  settingsForm = this.fb.group({
    emailNotifications: [false],
    darkMode:           [false]
  });

  ngOnInit(): void {
    const p = this.auth.profile!;
    this.profile = {
      ...p,
      avatarUrl: p.avatarUrl || '/assets/images/bg-login.jpg',
      score:     p.score     || 0,
      progress:  p.progress  || 0
    };
  }

  saveSettings() { /* … */ }
}
