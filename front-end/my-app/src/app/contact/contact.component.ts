import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ReactiveFormsModule,
  FormBuilder,
  FormGroup,
  Validators
} from '@angular/forms';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-contact',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule      // ← indispensable pour la pipe | translate
  ],
  templateUrl: './contact.component.html',
  styleUrls: ['./contact.component.scss']
})
export class ContactComponent {
  private fb = inject(FormBuilder);
  private translate = inject(TranslateService);

  loading = false;
  error: string | null = null;

  contactForm: FormGroup = this.fb.group({
    name:    ['', Validators.required],
    email:   ['', [Validators.required, Validators.email]],
    message: ['', Validators.required],
  });

  onSubmit(): void {
    if (this.contactForm.invalid) {
      this.contactForm.markAllAsTouched();
      return;
    }

    this.loading = true;
    this.error = null;

    // Simule un appel HTTP
    setTimeout(() => {
      const success = Math.random() > 0.2;
      if (success) {
        alert(this.translate.instant('CONTACT.SUCCESS_ALERT'));
        this.contactForm.reset();
      } else {
        this.error = this.translate.instant('CONTACT.ERROR.SUBMIT_FAILED');
      }
      this.loading = false;
    }, 1500);
  }
}
