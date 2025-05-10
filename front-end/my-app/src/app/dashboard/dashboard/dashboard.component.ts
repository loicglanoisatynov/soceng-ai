// src/app/dashboard/dashboard/dashboard.component.ts
import { Component, OnInit, inject } from '@angular/core';
import { CommonModule }                                        from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup }         from '@angular/forms';
import { TranslateModule }                                     from '@ngx-translate/core';
import { take }                                                from 'rxjs/operators';

import { AuthService, UserProfile }     from '../../auth/auth.service';
import { ProfileHeroComponent }         from '../../shared/profile-hero/profile-hero.component';
import { SettingsComponent }            from '../../settings/settings/settings.component';
import { MyChallengeComponent }         from '../challenges/mychallenge/mychallenge.component';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule,
    ProfileHeroComponent,
    SettingsComponent,
    MyChallengeComponent
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  private auth = inject(AuthService);
  private fb   = inject(FormBuilder);

  user: UserProfile & { avatarUrl: string; score: number; progress: number } = {
    id:        0,
    username:  'John Doe',
    email:     '',
    avatarUrl: '/assets/images/bg-login.jpg',
    score:     0,
    progress:  0
  };

  selectedTab: 'details' | 'settings' | 'challenges' | 'help' = 'details';
  profileForm!: FormGroup;
  currentChallenge: { title: string } | null = { title: 'Mon premier challenge' };

  ngOnInit() {
    this.auth.loggedIn$.pipe(take(1)).subscribe(is => {
      if (!is) return;
      const p = this.auth.profile!;
      this.user = {
        ...p,
        avatarUrl: p.avatarUrl || this.user.avatarUrl,
        score:     p.score     || 0,
        progress:  p.progress  || 0
      };

      this.profileForm = this.fb.group({
        fullName: [this.user.username],
        email:    [this.user.email],
        password: ['']
      });
    });
  }

  // …
}
