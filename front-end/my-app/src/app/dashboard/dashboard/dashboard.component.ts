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
  ngOnInit(): void {
    // Initialization logic here
    console.log('DashboardComponent initialized');
  }
}
