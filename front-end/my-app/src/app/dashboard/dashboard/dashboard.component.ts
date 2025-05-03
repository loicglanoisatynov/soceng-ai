// src/app/dashboard/dashboard/dashboard.component.ts
import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
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
  private auth   = inject(AuthService);
  private router = inject(Router);
  private fb     = inject(FormBuilder);

  user = {
    name:     'Utilisateur',
    photoUrl: '/assets/images/bg-login.jpg',
    score:    0
  };

  selectedTab: 'details'|'settings' = 'details';
  profileForm!: FormGroup;

  challenges = [
    { name: 'Challenge 1', info: 'Info', message: 'Message' },
    { name: 'Challenge 2', info: 'Info', message: 'Message' },
    { name: 'Challenge 3', info: 'Info', message: 'Message' }
  ];

  ngOnInit() {
    // on ne fait plus de checkAuth() ici
    // on part du principe qu’on est déjà passé par le AuthGuard
    this.profileForm = this.fb.group({
      fullName: [ this.user.name ],
      email:    [''],
      password: ['']
    });
  }

  logout() {
    this.auth.logout().subscribe(() => {
      this.router.navigate(['/home']);
    });
  }

  switchTab(tab: 'details'|'settings') {
    this.selectedTab = tab;
  }

  saveDetails() {
    if (!this.profileForm.valid) return;
    // TODO: PUT /api/edit-profile
  }
}
