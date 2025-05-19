// src/app/about/about/about.component.ts
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TranslateModule } from '@ngx-translate/core';

interface TeamMember {
  image: string;
  textKey: string;
}

@Component({
  selector: 'app-about',
  standalone: true,
  imports: [
    CommonModule,
    TranslateModule 
  ],
  templateUrl: './about.component.html',
  styleUrls: ['./about.component.scss']
})
export class AboutComponent {
  team: TeamMember[] = [
    { image: 'corentin.png', textKey: 'ABOUT.TEAM.CARD.TEXT' },
    { image: 'lisa.png',      textKey: 'ABOUT.TEAM.CARD.TEXT' },
    { image: 'loic.png',      textKey: 'ABOUT.TEAM.CARD.TEXT' },
    { image: 'corentin.png',  textKey: 'ABOUT.TEAM.CARD.TEXT' }
  ];
}
