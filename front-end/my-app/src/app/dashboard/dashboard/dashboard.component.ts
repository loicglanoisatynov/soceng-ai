import { Component, OnInit, inject } from '@angular/core';
import { CommonModule }                                    from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup }     from '@angular/forms';
import { TranslateModule }                                 from '@ngx-translate/core';
import { take }                                            from 'rxjs/operators';

import { AuthService, UserProfile }        from '../../auth/auth.service';
import { ProfileHeroComponent }            from '../../shared/profile-hero/profile-hero.component';
import { SettingsComponent }               from '../../settings/settings/settings.component';
import { MyChallengeComponent }            from '../challenges/mychallenge/mychallenge.component';

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

  profileForm!: FormGroup;
  selectedTab: 'details' | 'settings' | 'challenges' | 'help' = 'details';

  ngOnInit(): void {
    this.auth.loadProfile().pipe(take(1)).subscribe({
      next: (p: UserProfile) => {
        this.profileForm = this.fb.group({
          fullName: [p.username],
          email:    [p.email],
          password: ['']
        });
      },
      error: () => {
        this.profileForm = this.fb.group({
          fullName: [''],
          email:    [''],
          password: ['']
        });
      }
    });
  }

  switchTab(tab: 'details' | 'settings' | 'challenges' | 'help'): void {
    this.selectedTab = tab;
  }

  saveDetails(): void {
    if (!this.profileForm.valid) return;
    // …vous pouvez PUT /edit-profile ici si besoin
  }

  logout(): void {
    this.auth.logout().pipe(take(1)).subscribe();
  }
}
