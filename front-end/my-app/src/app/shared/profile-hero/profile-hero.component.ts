// profile-hero.component.ts
import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { UserProfile } from '../../auth/auth.service';

@Component({
  selector: 'app-profile-hero',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './profile-hero.component.html',
  styleUrls: ['./profile-hero.component.scss']
})
export class ProfileHeroComponent {
  /**
   * Profil complet, incluant avatarUrl, score, progress et biography
   */
  @Input() profile!: UserProfile & { avatarUrl: string; score: number; progress: number; biography?: string };
}