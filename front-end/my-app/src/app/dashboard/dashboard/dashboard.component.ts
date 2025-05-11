import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { ReactiveFormsModule, FormBuilder, FormGroup } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { take } from 'rxjs/operators';
import { AuthService, UserProfile } from '../../auth/auth.service';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,
    TranslateModule
  ],
  templateUrl: './dashboard.component.html'
})
export class DashboardComponent implements OnInit {
  private auth   = inject(AuthService);
  private router = inject(Router);
  private fb     = inject(FormBuilder);

  user: UserProfile & { avatarUrl: string; score: number; progress: number } = {
    id:        0,
    username:  'John Doe',
    email:     '',
    avatarUrl: '/assets/images/default-avatar.png',
    score:     0,
    progress:  0
  };

  selectedTab: 'details'|'settings' = 'details';
  profileForm!: FormGroup;
  challenges: Array<{ name: string; info: string }> = [];

  ngOnInit() {
    this.auth.loggedIn$.pipe(take(1)).subscribe(isLoggedIn => {
      if (!isLoggedIn) {
        this.router.navigate(
          [environment.routes.login],
          { queryParams: { returnUrl: environment.routes.dashboard } }
        );
        return;
      }

      const p = this.auth.profile!;
      this.user = {
        ...p,
        avatarUrl: p.avatarUrl || this.user.avatarUrl,
        score:     p.score     || this.user.score,
        progress:  p.progress  || this.user.progress
      };

      this.profileForm = this.fb.group({
        fullName: [ this.user.username ],
        email:    [ this.user.email ],
        password: [ '' ]
      });

      this.challenges = [
        { name: 'Challenge 1', info: 'Lorem ipsum…' },
        { name: 'Challenge 2', info: 'Dolor sit amet…' },
        { name: 'Challenge 3', info: 'Consectetur…' }
      ];
    });
  }

  logout() {
    this.auth.logout().subscribe(() =>
      this.router.navigate([environment.routes.login])
    );
  }

  switchTab(tab: 'details'|'settings') {
    this.selectedTab = tab;
  }

  saveDetails() {
    if (!this.profileForm.valid) return;
    // TODO → PUT `${environment.apiBaseUrl}/edit-profile`
  }
}
