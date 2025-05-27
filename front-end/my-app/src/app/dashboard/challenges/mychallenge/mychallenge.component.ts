// src/app/dashboard/challenges/mychallenge/mychallenge.component.ts
import { Component, OnInit } from '@angular/core';
import { CommonModule }       from '@angular/common';
import { TranslateModule }    from '@ngx-translate/core';
import { ProfileHeroComponent } from '../../../shared/profile-hero/profile-hero.component';

@Component({
  selector: 'app-mychallenge',
  standalone: true,
  imports: [CommonModule, TranslateModule, ProfileHeroComponent],
  templateUrl: './mychallenge.component.html',
  styleUrls: ['./mychallenge.component.scss']
})
export class MyChallengeComponent implements OnInit {
  ngOnInit(): void {
    // Initialization logic here
  }
}
