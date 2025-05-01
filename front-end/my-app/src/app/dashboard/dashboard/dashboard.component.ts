// src/app/dashboard/dashboard/dashboard.component.ts
import { Component, OnInit, inject, PLATFORM_ID } from '@angular/core';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { ReactiveFormsModule, FormBuilder, FormGroup } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { AuthService } from '../../auth/auth.service';

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
  private fb = inject(FormBuilder);
  private platformId = inject(PLATFORM_ID) as any;  // ← cast any

  user = {
    name: '',
    photoUrl: '/assets/images/default-avatar.png',
    score: 0
  };
  selectedTab: 'details' | 'settings' = 'details';
  profileForm!: FormGroup;

  challenges = [
    { name: 'Challenge 1', info: 'Information', message: 'Message' },
    { name: 'Challenge 2', info: 'Information', message: 'Message' },
    { name: 'Challenge saisonnier', info: 'Information', message: 'Message' }
  ];

  ngOnInit() {
    if (isPlatformBrowser(this.platformId) && this.auth.token) {
      try {
        const payloadBase64 = this.auth.token.split('.')[1];
        const payload = JSON.parse(atob(payloadBase64));
        this.user.name = payload.name || '';
        this.user.score = payload.score ?? 0;
      } catch {
        console.warn('Impossible de décoder le token.');
      }
    }

    this.profileForm = this.fb.group({
      fullName: [this.user.name],
      username: [''], // à remplir si besoin
      email: [''],
      password: ['']
    });
  }

  logout() {
    this.auth.logout();
    this.router.navigate(['/home']);
  }

  switchTab(tab: 'details' | 'settings') {
    this.selectedTab = tab;
  }

  saveDetails() {
    if (!this.profileForm.valid) return;
    // TODO: appel API pour mettre à jour le profil
  }
}
