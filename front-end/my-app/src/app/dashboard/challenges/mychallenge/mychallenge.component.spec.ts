import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MychallengeComponent } from './mychallenge.component';

describe('MychallengeComponent', () => {
  let component: MychallengeComponent;
  let fixture: ComponentFixture<MychallengeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MychallengeComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(MychallengeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
