// src/app/dashboard/dashboard/dashboard.component.ts
import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { ReactiveFormsModule, FormBuilder, FormGroup } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { AuthService, UserProfile } from '../../auth/auth.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,
    TranslateModule
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  private auth = inject(AuthService);
  private router = inject(Router);
  private fb     = inject(FormBuilder);

  // On ne stocke QUE ces trois champs pour l'instant
  user = {
    name:     '',
    photoUrl: '/assets/images/default-avatar.png',
    score:    0
  };

  selectedTab: 'details' | 'settings' = 'details';
  profileForm!: FormGroup;

  // Exemple statique ; vous remplacerez ça par un vrai fetch de vos challenges
  challenges = [
    { name: 'Challenge 1', info: 'Information', message: 'Message' },
    { name: 'Challenge 2', info: 'Information', message: 'Message' },
    { name: 'Challenge saisonnier', info: 'Information', message: 'Message' }
  ];

  ngOnInit() {
    // 1) On valide la session via le back (cookie HTTP)
    this.auth.checkAuth().subscribe(isAuth => {
      if (!isAuth) {
        // pas auth -> retour au login
        this.router.navigate(['/auth/login'], { queryParams: { returnUrl: '/dashboard' } });
        return;
      }

      // 2) Si OK, on récupère le profil stocké
      const profile: UserProfile | null = this.auth.profile;
      if (profile) {
        this.user.name     = profile.username;                // backend renvoie username
        this.user.photoUrl = (profile as any).avatarUrl  || this.user.photoUrl;
        this.user.score    = (profile as any).score      ?? this.user.score;
      }

      // 3) On construit le formulaire
      this.profileForm = this.fb.group({
        fullName: [ this.user.name ],
        email:    [ profile?.email || '' ],
        password: [ '' ]
      });
    });
  }

  logout() {
    this.auth.logout().subscribe(() => {
      this.router.navigate(['/home']);
    });
  }

  switchTab(tab: 'details' | 'settings') {
    this.selectedTab = tab;
  }

  saveDetails() {
    if (!this.profileForm.valid) return;
    // TODO → PUT /api/edit-profile ou /api/edit-user
  }
}
