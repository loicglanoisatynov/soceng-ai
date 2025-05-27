import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ReactiveFormsModule,
  FormBuilder,
  FormGroup,
  Validators
} from '@angular/forms';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-contact',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule
  ],
  templateUrl: './contact.component.html',
  styleUrls: ['./contact.component.scss']
})
export class ContactComponent {
  private fb = inject(FormBuilder);
  private translate = inject(TranslateService);
  private http = inject(HttpClient);

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

    const formData = this.contactForm.value;

    const headers = new HttpHeaders({ 'Accept': 'application/json' });

    this.http.post('https://formspree.io/f/xldbqrzz', formData, { headers }).subscribe({
      next: () => {
        alert(this.translate.instant('CONTACT.SUCCESS_ALERT'));
        this.contactForm.reset();
        this.loading = false;
      },
      error: () => {
        this.error = this.translate.instant('CONTACT.ERROR.SUBMIT_FAILED');
        this.loading = false;
      }
    });
  }
}
