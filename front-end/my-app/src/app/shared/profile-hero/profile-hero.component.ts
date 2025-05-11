import { Component, OnInit, OnDestroy, inject } from '@angular/core';
import { CommonModule }                         from '@angular/common';
import { RouterModule }                         from '@angular/router';
import { AuthService, UserProfile }             from '../../auth/auth.service';
import { Subscription }                         from 'rxjs';

@Component({
  selector: 'app-profile-hero',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './profile-hero.component.html',
  styleUrls: ['./profile-hero.component.scss']
})
export class ProfileHeroComponent implements OnInit, OnDestroy {
  private auth = inject(AuthService);
  private sub!: Subscription;

  // Valeurs par défaut pour affichage immédiat
  user: UserProfile & { avatarUrl: string; score: number; progress: number } = {
    id:        0,
    username:  'John Doe',
    email:     'john@example.com',
    avatarUrl: '/assets/images/bg-login.jpg',
    score:     0,
    progress:  0
  };

  ngOnInit(): void {
    // Souscription au BehaviorSubject partagé
    this.sub = this.auth.profile$.subscribe(p => {
      if (p) {
        this.user = {
          ...p,
          avatarUrl: p.avatarUrl || this.user.avatarUrl,
          score:     p.score     ?? this.user.score,
          progress:  p.progress  ?? this.user.progress
        };
      }
    });

    // Si profile n'a jamais été chargé (refresh direct), on force
    if (!this.auth.profile) {
      this.auth.loadProfile().subscribe({ error: () => {/* ignore */} });
    }
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }
}
