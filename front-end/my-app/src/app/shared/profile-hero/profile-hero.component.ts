// src/app/shared/profile-hero/profile-hero.component.ts
import { Component, Input } from '@angular/core';
import { CommonModule }     from '@angular/common';
import { RouterModule }     from '@angular/router';
import { UserProfile }      from '../../auth/auth.service';

@Component({
  selector: 'app-profile-hero',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './profile-hero.component.html',
  styleUrls: ['./profile-hero.component.scss']
})
export class ProfileHeroComponent {
  @Input() user!: UserProfile & { avatarUrl: string; score: number; progress: number };
}
