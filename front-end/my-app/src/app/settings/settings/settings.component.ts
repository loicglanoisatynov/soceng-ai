import { Component, OnInit, OnDestroy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, FormGroup, Validators } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { AuthService, UserProfile } from '../../auth/auth.service';
import { ProfileHeroComponent } from '../../shared/profile-hero/profile-hero.component';

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
export class SettingsComponent implements OnInit, OnDestroy {
  private fb   = inject(FormBuilder);
  private auth = inject(AuthService);

  profileForm!: FormGroup;
  credentialsForm!: FormGroup;
  profile?: UserProfile & { avatarUrl: string; score: number; progress: number; biography: string };

  private subProfile!: Subscription;

  ngOnInit(): void {
    // 1) Initialise les formulaires
    this.profileForm = this.fb.group({
      username: ['', Validators.required],
      biography: [''],
      avatar: ['']
    });
    this.credentialsForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', Validators.required],
      newpassword: ['']
    });

    // 2) Sync initial store → formulaire
    this.subProfile = this.auth.profile$.subscribe(p => {
      if (!p) return;
      this.profile = {
        ...p,
        avatarUrl: p.avatarUrl || '',
        score:    p.score    || 0,
        progress: p.progress || 0,
        biography:p.biography|| ''
      };
      this.profileForm.patchValue({
        username:  p.username,
        biography: p.biography || '',
        avatar:    p.avatarUrl || ''
      });
      this.credentialsForm.patchValue({ email: p.email });
    });

    // 3) Charge le profil si non déjà présent
    if (!this.auth.profile) {
      this.auth.loadProfile().subscribe({ error: () => {} });
    }
  }

  ngOnDestroy(): void {
    this.subProfile.unsubscribe();
  }

  /** Gestion de l'upload de l'avatar */
  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (!input.files?.length) return;
    const file = input.files[0];
    const reader = new FileReader();
    reader.onload = () => {
      const base64 = reader.result as string;
      this.profileForm.get('avatar')?.setValue(base64);
      if (this.profile) {
        this.profile.avatarUrl = base64;
      }
    };
    reader.readAsDataURL(file);
  }

  onSaveProfile(): void {
    if (this.profileForm.invalid) return;

    const { username, biography, avatar } = this.profileForm.value;

    this.auth.updateProfile({ username, biography, avatar }).subscribe({
      next: res => {
        if (res.status) {
          // Mise à jour du store APRES succès back
          if (this.auth.profile) {
            this.auth.profile$.next({
              ...this.auth.profile,
              username,
              biography,
              avatarUrl: avatar
            });
          }
          alert('Profil mis à jour !');
        } else {
          alert(`Échec : ${res.message}`);
        }
      },
      error: err => {
        // Affiche code + body de l'erreur pour diagnostiquer
        console.error('Erreur HTTP edit-profile', err.status, err.error);
        alert(`Erreur serveur ${err.status} :\n${err.error}`);
      }
    });
  }

  onSaveCredentials(): void {
    if (this.credentialsForm.invalid) return;
    const { email, password, newpassword } = this.credentialsForm.value;
    this.auth.updateUser({ email, password, newpassword }).subscribe(res => {
      alert(res.status ? 'Identifiants mis à jour !' : `Échec : ${res.message}`);
    });
  }
}
