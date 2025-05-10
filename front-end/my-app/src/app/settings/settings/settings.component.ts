import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { AuthService, UserProfile } from '../../auth/auth.service';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule],
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
  private fb = inject(FormBuilder);
  private auth = inject(AuthService);

  profile!: UserProfile;
  settingsForm = this.fb.group({
    emailNotifications: [false],
    darkMode: [false]
  });

  ngOnInit(): void {
    this.profile = this.auth.profile!;
  }

  saveSettings() {
    console.log('Settings saved', this.settingsForm.value);
  }
}