import { Component, inject } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-contact',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule
  ],
  templateUrl: './contact.component.html',
})
export class ContactComponent {
  private fb = inject(FormBuilder);

  // **Ajoutez ces deux propriétés** pour lever vos erreurs de template
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

    // Simuler un appel HTTP
    setTimeout(() => {
      // Ici vous remplacerez par votre appel réel, par ex. via HttpClient
      const success = Math.random() > 0.2;
      if (success) {
        alert('Merci pour votre message !');
        this.contactForm.reset();
      } else {
        this.error = 'Une erreur est survenue. Veuillez réessayer plus tard.';
      }
      this.loading = false;
    }, 1500);
  }
}
