import { Component, OnInit, OnDestroy, inject } from '@angular/core';
import { CommonModule }    from '@angular/common';
import {
  FormBuilder,
  ReactiveFormsModule,
  FormGroup,
  Validators
} from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { Subscription, skip }    from 'rxjs';

import { AuthService, UserProfile } from '../../auth/auth.service';
import { ProfileHeroComponent }     from '../../shared/profile-hero/profile-hero.component';

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
  profile?: UserProfile & { avatarUrl: string; score: number; progress: number; biography?: string };

  // Subscriptions pour le profil et pour les changements de formulaire
  private subProfile!: Subscription;
  private subForm!: Subscription;

  ngOnInit(): void {
    // Initialisation des formulaires
    this.profileForm = this.fb.group({
      username:  ['', Validators.required],
      biography: [''],
      avatar:    ['']
    });
    this.credentialsForm = this.fb.group({
      email:       ['', [Validators.required, Validators.email]],
      password:    ['', Validators.required],
      newpassword: ['']
    });

    // Abonnement au BehaviorSubject du profil
    this.subProfile = this.auth.profile$.subscribe(p => {
      if (p) {
        this.profile = {
          ...p,
          avatarUrl: p.avatarUrl || '',
          score:     p.score     || 0,
          progress:  p.progress  || 0,
          biography: p.biography || ''
        };
        // On patch le formulaire avec les valeurs reçues
        this.profileForm.patchValue({
          username:  p.username,
          biography: p.biography || '',
          avatar:    p.avatarUrl || ''
        });
        this.credentialsForm.patchValue({ email: p.email });
      }
    });

    // Si pas encore chargé, on charge une première fois
    if (!this.auth.profile) {
      this.auth.loadProfile().subscribe({ error: () => {} });
    }

    // Abonnement aux changements de formulaire pour mise à jour instantanée du Hero
    this.subForm = this.profileForm.valueChanges
      .pipe(skip(1)) // on ignore la première émission du patchValue
      .subscribe(({ username, biography, avatar }) => {
        if (this.auth.profile) {
          this.auth.profile$.next({
            ...this.auth.profile,
            username,
            biography,
            avatarUrl: avatar
          });
        }
      });
  }

  ngOnDestroy(): void {
    this.subProfile.unsubscribe();
    this.subForm.unsubscribe();
  }

  onSaveProfile(): void {
    if (this.profileForm.invalid) {
      return;
    }
    const { username, biography, avatar } = this.profileForm.value;

    // Mise à jour optimiste du profil local
    if (this.auth.profile) {
      this.auth.profile$.next({
        ...this.auth.profile,
        username,
        biography,
        avatarUrl: avatar
      });
    }

    // Appel HTTP pour persister sur le serveur
    this.auth.updateProfile({ username, biography, avatar }).subscribe({
      next: res => {
        alert(res.status ? 'Profil mis à jour !' : `Échec : ${res.message}`);
      },
      error: () => {
        alert('Erreur serveur, le profil local peut être désynchronisé.');
      }
    });
  }

  onSaveCredentials(): void {
    if (this.credentialsForm.invalid) {
      return;
    }
    const { email, password, newpassword } = this.credentialsForm.value;
    this.auth.updateUser({ email, password, newpassword }).subscribe(res => {
      alert(res.status ? 'Identifiants mis à jour !' : `Échec : ${res.message}`);
    });
  }
}
