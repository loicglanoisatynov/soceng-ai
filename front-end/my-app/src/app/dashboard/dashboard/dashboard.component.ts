// src/app/dashboard/dashboard/dashboard.component.ts

import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { take } from 'rxjs/operators';
import { AuthService, UserProfile } from '../../auth/auth.service';
import { SettingsComponent } from '../../settings/settings/settings.component';
import { MyChallengeComponent } from '../challenges/mychallenge/mychallenge.component';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule,
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
    id:         0,
    username:   'John Doe',
    email:      '',
    avatarUrl:  '/assets/images/bg-login.jpg',
    score:      0,
    progress:   0
  };

  // Onglet sélectionné
  selectedTab: 'details' | 'settings' | 'challenges' | 'help' = 'details';
  profileForm!: FormGroup;

  // Propriété ajoutée pour alimenter ta carte
  currentChallenge: { title: string } | null = {
    title: 'Mon premier challenge'
  };

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

      // Si tu souhaites récupérer dynamiquement le challenge en cours :
      // this.currentChallenge = this.challenges?.[0] ? { title: this.challenges[0].name } : null;
    });
  }

  switchTab(tab: 'details' | 'settings' | 'challenges' | 'help') {
    this.selectedTab = tab;
  }

  logout() {
    this.auth.logout().pipe(take(1)).subscribe(() => {
      // TODO → this.router.navigate(['/login']);
    });
  }

  saveDetails() {
    if (!this.profileForm.valid) return;
    // TODO → PUT `${environment.apiBaseUrl}/profile`
  }
}
