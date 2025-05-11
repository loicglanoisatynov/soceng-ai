import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-challenge',
  standalone: true,
  imports: [
    CommonModule,
    TranslateModule   // ← pour la pipe | translate
  ],
  templateUrl: './challenge.component.html',
  styleUrls: ['./challenge.component.scss']
})
export class ChallengeComponent {}
