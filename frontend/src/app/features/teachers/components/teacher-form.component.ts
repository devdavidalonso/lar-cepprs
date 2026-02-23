import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, FormArray, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatStepperModule } from '@angular/material/stepper';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MAT_DATE_FORMATS, MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { TeacherService } from '../../../core/services/teacher.service';
import { Teacher, UserContact } from '../../../core/models/teacher.model';
import { ProgramService } from '../../../core/services/program.service';
import { Program } from '../../../core/models/program.model';

export const BRAZILIAN_DATE_FORMATS = {
  parse: {
    dateInput: 'DD/MM/YYYY',
  },
  display: {
    dateInput: 'DD/MM/YYYY',
    monthYearLabel: 'MMM YYYY',
    dateA11yLabel: 'LL',
    monthYearA11yLabel: 'MMMM YYYY',
  },
};

@Component({
  selector: 'app-teacher-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule,
    MatStepperModule,
    MatInputModule,
    MatButtonModule,
    MatSelectModule,
    MatIconModule,
    MatCardModule,
    MatSnackBarModule,
    MatCheckboxModule,
    MatDatepickerModule,
    MatTooltipModule,
    MatNativeDateModule,
    RouterModule
  ],
  providers: [
    { provide: MAT_DATE_LOCALE, useValue: 'pt-BR' },
    { provide: MAT_DATE_FORMATS, useValue: BRAZILIAN_DATE_FORMATS }
  ],
  template: `
    <div class="teacher-form-container">
      <mat-card class="form-card">
        <mat-card-header>
          <div mat-card-avatar>
              <mat-icon class="header-icon">person</mat-icon>
          </div>
          <mat-card-title>{{ isEditMode ? ('TEACHER.EDIT' | translate) : ('TEACHER.NEW' | translate) }}</mat-card-title>
          <mat-card-subtitle>Teacher Registration Wizard</mat-card-subtitle>
        </mat-card-header>
        
        <mat-card-content>
          <mat-stepper [linear]="true" #stepper orientation="vertical">
            
            <!-- Step 1: Dados Pessoais -->
            <mat-step [stepControl]="personalInfoForm">
              <form [formGroup]="personalInfoForm">
                <ng-template matStepLabel>Personal Information</ng-template>
                <p class="step-desc">Basic identification and contact details.</p>
                
                <div class="form-grid">
                  <div class="row">
                      <mat-form-field appearance="outline" class="flex-2">
                        <mat-label>{{ 'TEACHER.NAME' | translate }}</mat-label>
                        <input matInput formControlName="name" placeholder="Full Name">
                        <mat-error *ngIf="personalInfoForm.get('name')?.hasError('required')">{{ 'VALIDATION.REQUIRED' | translate }}</mat-error>
                      </mat-form-field>

                      <mat-form-field appearance="outline" class="flex-1">
                          <mat-label>{{ 'TEACHER.CPF' | translate }}</mat-label>
                          <input matInput formControlName="cpf" placeholder="000.000.000-00" (input)="formatCPF($event, 'cpf')">
                      </mat-form-field>
                  </div>

                  <div class="row">
                      <mat-form-field appearance="outline" class="flex-2">
                        <mat-label>{{ 'TEACHER.EMAIL' | translate }}</mat-label>
                        <input matInput type="email" formControlName="email" placeholder="email@example.com">
                        <mat-error *ngIf="personalInfoForm.get('email')?.hasError('required')">{{ 'VALIDATION.REQUIRED' | translate }}</mat-error>
                        <mat-error *ngIf="personalInfoForm.get('email')?.hasError('email')">{{ 'VALIDATION.INVALID_EMAIL' | translate }}</mat-error>
                      </mat-form-field>

                      <mat-form-field appearance="outline" class="flex-1">
                        <mat-label>Birth Date</mat-label>
                        <input matInput 
                               [matDatepicker]="picker" 
                               formControlName="birthDate" 
                               (input)="formatBirthDate($event)"
                               placeholder="DD/MM/AAAA"
                               maxlength="10"
                               autocomplete="off">
                        <mat-datepicker-toggle matIconSuffix [for]="picker" matTooltip="Abrir calendário"></mat-datepicker-toggle>
                        <mat-datepicker #picker startView="multi-year"></mat-datepicker>
                        <mat-hint>Digite DD/MM/AAAA ou use o calendário</mat-hint>
                        <mat-error *ngIf="personalInfoForm.get('birthDate')?.hasError('matDatepickerParse')">
                            Data inválida. Use DD/MM/AAAA
                        </mat-error>
                      </mat-form-field>

                      <mat-form-field appearance="outline" class="flex-1">
                          <mat-label>{{ 'TEACHER.PHONE' | translate }}</mat-label>
                          <input matInput formControlName="phone" placeholder="(00) 00000-0000" (input)="formatPhone($event, 'phone')">
                      </mat-form-field>
                  </div>

                  <mat-form-field appearance="outline" class="full-width">
                    <mat-label>{{ 'TEACHER.LINKEDIN' | translate }}</mat-label>
                    <input matInput formControlName="linkedinUrl" placeholder="https://linkedin.com/in/...">
                    <mat-icon matSuffix>link</mat-icon>
                  </mat-form-field>
                </div>

                <div class="stepper-actions">
                  <button mat-button color="warn" routerLink="/teachers">Cancel</button>
                  <button mat-raised-button color="primary" matStepperNext>Next: Qualifications</button>
                </div>
              </form>
            </mat-step>

            <!-- Step 2: Qualificações -->
            <mat-step [stepControl]="qualificationsForm">
              <form [formGroup]="qualificationsForm">
                <ng-template matStepLabel>Qualifications</ng-template>
                <p class="step-desc">Academic background and specific skills.</p>
                
                <div class="form-grid">
                  <mat-form-field appearance="outline" class="full-width">
                    <mat-label>Specialization</mat-label>
                    <input matInput formControlName="specialization" placeholder="e.g. Mathematics, Sciences...">
                  </mat-form-field>

                  <mat-form-field appearance="outline" class="full-width">
                    <mat-label>Bio</mat-label>
                    <textarea matInput formControlName="bio" rows="3" placeholder="Brief teacher biography..."></textarea>
                  </mat-form-field>

                  <mat-form-field appearance="outline" class="full-width">
                    <mat-label>Programas</mat-label>
                    <mat-select formControlName="programIds" multiple>
                      <mat-option *ngFor="let program of programs" [value]="program.id">
                        {{ program.name }}
                      </mat-option>
                    </mat-select>
                    <mat-hint>Selecione um ou mais programas de atuação.</mat-hint>
                    <mat-error *ngIf="qualificationsForm.get('programIds')?.hasError('required')">
                      Selecione pelo menos um programa.
                    </mat-error>
                  </mat-form-field>
                </div>

                <div class="stepper-actions">
                  <button mat-button matStepperPrevious>Back</button>
                  <button mat-raised-button color="primary" matStepperNext>Next: Address</button>
                </div>
              </form>
            </mat-step>

            <!-- Step 3: Endereço -->
            <mat-step [stepControl]="addressForm">
              <form [formGroup]="addressForm">
                <ng-template matStepLabel>Address</ng-template>
                <p class="step-desc">Residential address details.</p>
                
                <div class="form-grid">
                  <div class="row">
                    <mat-form-field appearance="outline" class="flex-1">
                      <mat-label>CEP (Zip Code)</mat-label>
                      <input matInput formControlName="cep" placeholder="00000-000" (input)="formatCEP($event, 'cep')">
                      <mat-error *ngIf="addressForm.get('cep')?.hasError('required')">Required</mat-error>
                    </mat-form-field>
                    
                    <mat-form-field appearance="outline" class="flex-2">
                      <mat-label>Street</mat-label>
                      <input matInput formControlName="street" placeholder="Av. Paulista">
                      <mat-error *ngIf="addressForm.get('street')?.hasError('required')">Required</mat-error>
                    </mat-form-field>
                  </div>

                  <div class="row">
                    <mat-form-field appearance="outline" class="flex-1">
                      <mat-label>Number</mat-label>
                      <input matInput formControlName="number" placeholder="1000">
                      <mat-error *ngIf="addressForm.get('number')?.hasError('required')">Required</mat-error>
                    </mat-form-field>

                    <mat-form-field appearance="outline" class="flex-2">
                      <mat-label>Complement</mat-label>
                      <input matInput formControlName="complement" placeholder="Apt 12">
                    </mat-form-field>
                  </div>

                  <div class="row">
                    <mat-form-field appearance="outline" class="flex-1">
                      <mat-label>Neighborhood</mat-label>
                      <input matInput formControlName="neighborhood">
                      <mat-error *ngIf="addressForm.get('neighborhood')?.hasError('required')">Required</mat-error>
                    </mat-form-field>

                    <mat-form-field appearance="outline" class="flex-1">
                      <mat-label>City</mat-label>
                      <input matInput formControlName="city">
                      <mat-error *ngIf="addressForm.get('city')?.hasError('required')">Required</mat-error>
                    </mat-form-field>

                    <mat-form-field appearance="outline" class="flex-1">
                      <mat-label>State (UF)</mat-label>
                      <input matInput formControlName="state" placeholder="SP" maxlength="2">
                      <mat-error *ngIf="addressForm.get('state')?.hasError('required')">Required</mat-error>
                    </mat-form-field>
                  </div>
                </div>

                <div class="stepper-actions">
                  <button mat-button matStepperPrevious>Back</button>
                  <button mat-raised-button color="primary" matStepperNext>Next: Emergency Contacts</button>
                </div>
              </form>
            </mat-step>

            <!-- Step 4: Contatos de Emergência -->
            <mat-step [stepControl]="contactsForm">
              <form [formGroup]="contactsForm">
                <ng-template matStepLabel>{{ 'TEACHER.CONTACTS' | translate }}</ng-template>
                <p class="step-desc">People to contact in case of emergency (Guardians/Relatives).</p>
                
                <div formArrayName="contactsList" class="contacts-array">
                  <div *ngFor="let contact of contactsArray.controls; let i=index" [formGroupName]="i" class="contact-box">
                      <div class="contact-header">
                          <h4>Contact #{{i + 1}}</h4>
                          <button mat-icon-button color="warn" (click)="removeContact(i)" type="button" *ngIf="contactsArray.length > 1">
                              <mat-icon>delete</mat-icon>
                          </button>
                      </div>
                      
                      <div class="row">
                          <mat-form-field appearance="outline" class="flex-2">
                             <mat-label>Name</mat-label>
                             <input matInput formControlName="name">
                          </mat-form-field>

                          <mat-form-field appearance="outline" class="flex-1">
                             <mat-label>Relationship</mat-label>
                             <mat-select formControlName="relationship">
                                 <mat-option value="father">Father</mat-option>
                                 <mat-option value="mother">Mother</mat-option>
                                 <mat-option value="spouse">Spouse</mat-option>
                                 <mat-option value="sibling">Sibling</mat-option>
                                 <mat-option value="friend">Friend</mat-option>
                                 <mat-option value="other">Other</mat-option>
                             </mat-select>
                          </mat-form-field>
                      </div>

                      <div class="row">
                          <mat-form-field appearance="outline" class="flex-1">
                             <mat-label>CPF</mat-label>
                             <input matInput formControlName="cpf" (input)="formatCPF($event, 'cpf', contactsArray.at(i))">
                          </mat-form-field>

                          <mat-form-field appearance="outline" class="flex-1">
                             <mat-label>Phone</mat-label>
                             <input matInput formControlName="phone" (input)="formatPhone($event, 'phone', contactsArray.at(i))">
                          </mat-form-field>

                          <mat-form-field appearance="outline" class="flex-2">
                             <mat-label>Email</mat-label>
                             <input matInput formControlName="email">
                          </mat-form-field>
                      </div>

                      <div class="row">
                         <mat-checkbox formControlName="canPickup">Can Pickup</mat-checkbox>
                         <mat-checkbox formControlName="receiveNotifications">Receive Notifications</mat-checkbox>
                         <mat-checkbox formControlName="authorizeActivities">Authorize Activities</mat-checkbox>
                      </div>
                  </div>
                </div>

                <button mat-stroked-button color="primary" type="button" (click)="addContact()" class="add-btn">
                   <mat-icon>add</mat-icon> Add Contact
                </button>

                <div class="stepper-actions">
                  <button mat-button matStepperPrevious>Back</button>
                  <button mat-raised-button color="primary" matStepperNext>Review</button>
                </div>
              </form>
            </mat-step>

            <!-- Step 5: Review -->
            <mat-step>
              <ng-template matStepLabel>Review & Save</ng-template>
              
              <div class="review-container">
                 <div class="review-card">
                    <h4>{{ personalInfoForm.get('name')?.value }}</h4>
                    <p>{{ personalInfoForm.get('email')?.value }} | {{ personalInfoForm.get('phone')?.value }}</p>
                    <hr>
                    <p><strong>Specialization:</strong> {{ qualificationsForm.get('specialization')?.value || 'Not provided' }}</p>
                    <p><strong>Bio:</strong> {{ qualificationsForm.get('bio')?.value || 'Not provided' }}</p>
                    <p><strong>Programas:</strong> {{ getProgramNames(qualificationsForm.get('programIds')?.value) }}</p>
                    <hr>
                    <p><strong>Address:</strong> {{ addressForm.get('street')?.value }}, {{ addressForm.get('number')?.value }} - {{ addressForm.get('city')?.value }}/{{ addressForm.get('state')?.value | uppercase }}</p>
                    <hr>
                    <p><strong>Emergency Contacts:</strong> {{ contactsArray.length }} registered</p>
                </div>

                <div class="stepper-actions">
                  <button mat-button matStepperPrevious>Back</button>
                  <button mat-raised-button color="accent" (click)="submit()" [disabled]="isSubmitting">
                    {{ isSubmitting ? 'Saving...' : 'Save Teacher' }}
                  </button>
                </div>
              </div>
            </mat-step>
          </mat-stepper>
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .teacher-form-container {
      display: flex;
      justify-content: center;
      padding: 24px;
      min-height: 85vh;
      background-color: #f8f9fa;
    }
    
    .form-card {
      width: 100%;
      max-width: 800px;
      padding: 16px;
    }

    .header-icon {
        font-size: 24px;
        color: #006aac;
    }

    .step-desc {
        color: #666;
        margin-bottom: 24px;
        font-style: italic;
    }

    .form-grid {
      display: flex;
      flex-direction: column;
      gap: 16px;
      margin-top: 10px;
    }

    .row {
      display: flex;
      gap: 16px;
      margin-bottom: 8px;
    }
    
    .flex-1 { flex: 1; }
    .flex-2 { flex: 2; }
    .full-width { width: 100%; }

    .contacts-array {
       display: flex;
       flex-direction: column;
       gap: 16px;
       margin-bottom: 16px;
    }

    .contact-box {
        padding: 16px;
        border: 1px dashed #ccc;
        border-radius: 8px;
        background: #fafafa;
    }

    .contact-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
    }
    .contact-header h4 { margin: 0; color: #555; }

    .add-btn { margin-top: 8px; }

    .stepper-actions {
      display: flex;
      justify-content: flex-end;
      gap: 12px;
      margin-top: 32px;
      padding-top: 16px;
      border-top: 1px solid #eee;
    }
    
    .review-card {
        background: #f0f4f8;
        padding: 16px;
        border-radius: 8px;
        text-align: left;
        margin-top: 16px;
    }
  `]
})
export class TeacherFormComponent implements OnInit {
  personalInfoForm: FormGroup;
  qualificationsForm: FormGroup;
  contactsForm: FormGroup;
  addressForm: FormGroup;
  programs: Program[] = [];
  
  isSubmitting = false;
  isEditMode = false;
  teacherId: number | null = null;
  
  constructor(
    private _formBuilder: FormBuilder,
    private teacherService: TeacherService,
    private programService: ProgramService,
    private router: Router,
    private route: ActivatedRoute,
    private snackBar: MatSnackBar
  ) {
    this.personalInfoForm = this._formBuilder.group({
      name: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      cpf: ['', [Validators.required, Validators.pattern(/^\d{11}$/)]],
      birthDate: ['', Validators.required],
      phone: ['', [Validators.pattern(/^\(\d{2}\)\s\d{5}-\d{4}$/)]],
      linkedinUrl: ['']
    });

    this.qualificationsForm = this._formBuilder.group({
      specialization: [''],
      bio: [''],
      programIds: [[], Validators.required]
    });

    this.addressForm = this._formBuilder.group({
      cep: ['', [Validators.required, Validators.pattern(/^\d{5}-\d{3}$/)]],
      street: ['', Validators.required],
      number: ['', Validators.required],
      complement: [''],
      neighborhood: ['', Validators.required],
      city: ['', Validators.required],
      state: ['', Validators.required]
    });

    this.contactsForm = this._formBuilder.group({
      contactsList: this._formBuilder.array([ this.createContactGroup() ])
    });
  }

  ngOnInit(): void {
    this.loadPrograms();
    this.checkEditMode();
  }

  loadPrograms(): void {
    this.programService.listPrograms(true).subscribe({
      next: (programs) => {
        this.programs = programs;
      },
      error: (err) => {
        console.error('Error loading programs', err);
        this.snackBar.open('Falha ao carregar programas.', 'Fechar', { duration: 3000 });
      }
    });
  }

  get contactsArray(): FormArray {
      return this.contactsForm.get('contactsList') as FormArray;
  }

  createContactGroup(): FormGroup {
      return this._formBuilder.group({
          name: ['', Validators.required],
          relationship: ['other', Validators.required],
          cpf: ['', [Validators.pattern(/^\d{11}$/)]],
          phone: ['', [Validators.pattern(/^\(\d{2}\)\s\d{5}-\d{4}$/)]],
          email: ['', Validators.email],
          canPickup: [false],
          receiveNotifications: [false],
          authorizeActivities: [false]
      });
  }

  addContact() {
      this.contactsArray.push(this.createContactGroup());
  }

  removeContact(index: number) {
      this.contactsArray.removeAt(index);
  }

  checkEditMode() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEditMode = true;
      this.teacherId = +id;
      this.loadTeacherData(this.teacherId);
    }
  }

  loadTeacherData(id: number) {
    this.teacherService.getTeacher(id).subscribe({
      next: (teacher) => {
        this.personalInfoForm.patchValue({
          name: teacher.name,
          email: teacher.email,
          cpf: teacher.cpf,
          birthDate: teacher.birthDate,
          phone: teacher.phone,
          linkedinUrl: teacher.linkedinUrl
        });

        this.qualificationsForm.patchValue({
          specialization: teacher.specialization,
          bio: teacher.bio,
          programIds: teacher.programIds || []
        });
        
        // Ensure email isn't easily changed if coming from Keycloak sync, but leave it as is for MVP unless locked 
        this.personalInfoForm.get('email')?.disable();

        if (teacher.address) {
          this.addressForm.patchValue(teacher.address);
        }

        if (teacher.userContacts && teacher.userContacts.length > 0) {
            this.contactsArray.clear();
            teacher.userContacts.forEach((contact: UserContact) => {
                const group = this.createContactGroup();
                group.patchValue(contact);
                this.contactsArray.push(group);
            });
        }
      },
      error: (err) => {
        this.snackBar.open('Error loading teacher data.', 'Close', { duration: 3000 });
        this.router.navigate(['/teachers']);
      }
    });
  }

  // Mask formatters
  formatCPF(event: any, controlName: string, formGroup: any = null): void {
    const form = (formGroup as FormGroup) || this.personalInfoForm;
    let value = event.target.value.replace(/\D/g, '');
    if (value.length > 11) value = value.slice(0, 11);
    form.get(controlName)?.setValue(value, { emitEvent: false });
  }

  formatPhone(event: any, controlName: string, formGroup: any = null): void {
    const form = (formGroup as FormGroup) || this.personalInfoForm;
    let value = event.target.value.replace(/\D/g, '');
    if (value.length > 11) value = value.slice(0, 11);
    
    if (value.length >= 11) {
      value = `(${value.slice(0, 2)}) ${value.slice(2, 7)}-${value.slice(7, 11)}`;
    } else if (value.length >= 7) {
      value = `(${value.slice(0, 2)}) ${value.slice(2, 7)}-${value.slice(7)}`;
    } else if (value.length >= 2) {
      value = `(${value.slice(0, 2)}) ${value.slice(2)}`;
    }
    
    form.get(controlName)?.setValue(value, { emitEvent: false });
  }

  formatBirthDate(event: any): void {
    let value = event.target.value.replace(/\D/g, '');
    if (value.length > 8) value = value.slice(0, 8);
    
    if (value.length >= 5) {
      value = `${value.slice(0, 2)}/${value.slice(2, 4)}/${value.slice(4)}`;
    } else if (value.length >= 3) {
      value = `${value.slice(0, 2)}/${value.slice(2)}`;
    }
    
    this.personalInfoForm.get('birthDate')?.setValue(value, { emitEvent: false });
  }

  formatCEP(event: any, controlName: string, formGroup: any = null): void {
    const form = (formGroup as FormGroup) || this.addressForm;
    let value = event.target.value.replace(/\D/g, '');
    if (value.length > 8) value = value.slice(0, 8);
    
    if (value.length >= 6) {
      value = `${value.slice(0, 5)}-${value.slice(5)}`;
    }
    
    form.get(controlName)?.setValue(value, { emitEvent: false });
  }

  submit() {
    if (this.personalInfoForm.valid && this.qualificationsForm.valid && this.addressForm.valid && this.contactsForm.valid) {
      this.isSubmitting = true;

      // Reactivate email so it can be read in raw value or just use getRawValue()
      const personalData = this.personalInfoForm.getRawValue();

      // Convert date string from DD/MM/YYYY to YYYY-MM-DD
      let isoDate = null;
      if (typeof personalData.birthDate === 'string' && personalData.birthDate.length === 10) {
        const parts = personalData.birthDate.split('/');
        if (parts.length === 3) {
           isoDate = `${parts[2]}-${parts[1]}-${parts[0]}T00:00:00Z`;
        }
      } else if (personalData.birthDate instanceof Date) {
         isoDate = personalData.birthDate.toISOString();
      }

      const teacherData: Teacher = {
        ...personalData,
        birthDate: isoDate,
        ...this.qualificationsForm.value,
        address: this.addressForm.value,
        userContacts: this.contactsForm.value.contactsList,
        active: true
      };

      const request = this.isEditMode && this.teacherId
        ? this.teacherService.updateTeacher(this.teacherId, teacherData)
        : this.teacherService.createTeacher(teacherData);

      request.subscribe({
        next: (res: Teacher) => {
          this.snackBar.open(this.isEditMode ? 'Teacher updated successfully!' : 'Teacher created successfully!', 'Close', { duration: 3000 });
          this.router.navigate(['/teachers']);
        },
        error: (err: any) => {
          console.error(err);
          this.snackBar.open(`Error ${this.isEditMode ? 'updating' : 'creating'} teacher. Please try again.`, 'Close', { duration: 5000 });
          this.isSubmitting = false;
        }
      });
    } else {
      // Mark all as touched to show errors
      this.personalInfoForm.markAllAsTouched();
      this.qualificationsForm.markAllAsTouched();
      this.addressForm.markAllAsTouched();
      this.contactsForm.markAllAsTouched();
      this.snackBar.open('Please fill all required fields correctly.', 'Close', { duration: 3000 });
    }
  }

  getProgramNames(ids: number[] | undefined): string {
    if (!ids || ids.length === 0) {
      return 'Nenhum programa selecionado';
    }

    const names = this.programs
      .filter((program) => ids.includes(program.id))
      .map((program) => program.name);

    return names.length > 0 ? names.join(', ') : 'Nenhum programa selecionado';
  }
}
