import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-mychallenge',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './mychallenge.component.html',
  styleUrls: ['./mychallenge.component.scss']
})
export class MyChallengeComponent implements OnInit {
  challenges = [
    { name: 'Challenge 1', info: 'Lorem ipsum…' },
    { name: 'Challenge 2', info: 'Dolor sit amet…' },
    { name: 'Challenge 3', info: 'Consectetur…' }
  ];
  ngOnInit(): void {}
}