import { Component, OnInit, OnDestroy, inject } from '@angular/core';
import { CommonModule }    from '@angular/common';
import {
  FormBuilder,
  ReactiveFormsModule,
  FormGroup,
  Validators
} from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { Subscription }    from 'rxjs';

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
  profile?: UserProfile & { avatarUrl:string; score:number; progress:number; biography?:string };
  private sub!: Subscription;

  ngOnInit(): void {
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

    this.sub = this.auth.profile$.subscribe(p => {
      if (p) {
        this.profile = {
          ...p,
          avatarUrl: p.avatarUrl || '',
          score:     p.score     || 0,
          progress:  p.progress  || 0,
          biography: p.biography || ''
        };
        this.profileForm.patchValue({
          username:  p.username,
          biography: p.biography || '',
          avatar:    p.avatarUrl || ''
        });
        this.credentialsForm.patchValue({ email: p.email });
      }
    });

    if (!this.auth.profile) {
      this.auth.loadProfile().subscribe({ error: () => {} });
    }
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  onSaveProfile(): void {
    if (this.profileForm.invalid) return;
    const { username, biography, avatar } = this.profileForm.value;
    this.auth.updateProfile({ username, biography, avatar }).subscribe(res => {
      alert(res.status ? 'Profil mis à jour !' : `Échec : ${res.message}`);
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
