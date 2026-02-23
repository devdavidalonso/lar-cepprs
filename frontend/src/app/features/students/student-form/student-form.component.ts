import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, FormArray, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router, ActivatedRoute, RouterModule } from '@angular/router';
import { MatStepperModule } from '@angular/material/stepper';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDividerModule } from '@angular/material/divider';
import { MatTooltipModule } from '@angular/material/tooltip';
import { StudentService } from '../../../core/services/student.service';
import { CreateStudentRequest, UpdateStudentRequest, Student, StudentStatus } from '../../../core/models/student.model';
import { USER_PROFILES } from '../../../core/models/user.model';
import { ProgramService } from '@core/services/program.service';
import { Program } from '@core/models/program.model';

// Configuração do formato de data brasileiro
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
  selector: 'app-student-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    MatStepperModule,
    MatInputModule,
    MatButtonModule,
    MatSelectModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatIconModule,
    MatCardModule,
    MatSnackBarModule,
    MatCheckboxModule,
    MatDividerModule,
    MatTooltipModule
  ],
  templateUrl: './student-form.component.html',
  styleUrl: './student-form.component.scss',
  providers: [
    { provide: MAT_DATE_LOCALE, useValue: 'pt-BR' },
    { provide: MAT_DATE_FORMATS, useValue: BRAZILIAN_DATE_FORMATS }
  ]
})
export class StudentFormComponent implements OnInit {
  personalDataForm!: FormGroup;
  studentDataForm!: FormGroup;
  guardiansForm!: FormGroup;
  
  isSubmitting = false;
  isEditMode = false;
  studentId: number | null = null;
  student?: Student;
  programs: Program[] = [];

  // ✅ Status options - only for edit mode
  statusOptions: { value: StudentStatus; label: string }[] = [
    { value: 'active', label: 'Ativo' },
    { value: 'inactive', label: 'Inativo' },
    { value: 'suspended', label: 'Suspenso' }
  ];

  relationshipOptions = [
    { value: 'father', label: 'Pai' },
    { value: 'mother', label: 'Mãe' },
    { value: 'grandfather', label: 'Avô' },
    { value: 'grandmother', label: 'Avó' },
    { value: 'uncle', label: 'Tio' },
    { value: 'aunt', label: 'Tia' },
    { value: 'brother', label: 'Irmão' },
    { value: 'sister', label: 'Irmã' },
    { value: 'other', label: 'Outro' }
  ];

  constructor(
    private fb: FormBuilder,
    private studentService: StudentService,
    private programService: ProgramService,
    private router: Router,
    private route: ActivatedRoute,
    private snackBar: MatSnackBar
  ) {
    this.initForms();
  }

  ngOnInit(): void {
    this.loadPrograms();
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEditMode = true;
      this.studentId = +id;
      if (!this.studentDataForm.contains('status')) {
        this.studentDataForm.addControl('status', this.fb.control('active', Validators.required));
      }
      this.loadStudent();
    }
  }

  private loadPrograms(): void {
    this.programService.listPrograms(true).subscribe({
      next: (programs) => {
        this.programs = programs;
      },
      error: (err) => {
        console.error('Error loading programs:', err);
        this.snackBar.open('Erro ao carregar programas.', 'Fechar', { duration: 3000 });
      }
    });
  }

  initForms(): void {
    // Step 1: Personal Data
    this.personalDataForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      cpf: ['', [Validators.required, Validators.pattern(/^\d{11}$/)]],
      email: ['', [Validators.required, Validators.email]],
      birthDate: ['', Validators.required],
      phone: ['', [Validators.required, Validators.pattern(/^\(\d{2}\)\s\d{5}-\d{4}$/)]],
      address: this.fb.group({
        cep: ['', [Validators.required, Validators.pattern(/^\d{5}-\d{3}$/)]],
        street: ['', Validators.required],
        number: ['', Validators.required],
        complement: [''],
        neighborhood: ['', Validators.required],
        city: ['', Validators.required],
        state: ['', [Validators.required, Validators.maxLength(2)]]
      })
    });

    // Step 2: Student Data
    // ✅ REMOVED: registrationNumber (auto-generated by backend)
    // ✅ REMOVED: status from creation form (defaults to 'active')
    this.studentDataForm = this.fb.group({
      programIds: [[], Validators.required],
      medicalInfo: [''],
      specialNeeds: [''],
      notes: ['']
    });

    // ✅ Add status field ONLY in edit mode
    if (this.isEditMode) {
      this.studentDataForm.addControl('status', this.fb.control('active', Validators.required));
    }

    // Step 3: Guardians (FormArray)
    this.guardiansForm = this.fb.group({
      guardians: this.fb.array([])
    });
  }

  /**
   * Load student data for editing
   */
  private loadStudent(): void {
    if (!this.studentId) return;

    this.isSubmitting = true;
    this.studentService.getStudent(this.studentId).subscribe({
      next: (student) => {
        this.student = student;
        
        // Populate personal data form
        this.personalDataForm.patchValue({
          name: student.user.name,
          cpf: student.user.cpf,
          email: student.user.email,
          birthDate: student.user.birthDate,
          phone: student.user.phone,
          address: student.user.address
        });

        // Populate student data form
        this.studentDataForm.patchValue({
          programIds: student.programIds || [],
          medicalInfo: student.medicalInfo || '',
          specialNeeds: student.specialNeeds || '',
          notes: student.notes || '',
          status: student.status // ✅ Load status in edit mode
        });

        this.isSubmitting = false;
      },
      error: (err) => {
        console.error('Error loading student:', err);
        this.snackBar.open('Erro ao carregar dados do aluno.', 'Fechar', { duration: 3000 });
        this.isSubmitting = false;
        this.router.navigate(['/students']);
      }
    });
  }

  get guardians(): FormArray {
    return this.guardiansForm.get('guardians') as FormArray;
  }

  addGuardian(): void {
    const guardianGroup = this.fb.group({
      name: ['', Validators.required],
      relationship: ['', Validators.required],
      cpf: ['', [Validators.required, Validators.pattern(/^\d{11}$/)]],
      phone: ['', [Validators.required, Validators.pattern(/^\(\d{2}\)\s\d{5}-\d{4}$/)]],
      email: ['', Validators.email],
      address: [''],
      permissions: this.fb.group({
        canPickup: [true],
        canAuthorizeLeave: [false],
        receivesNotifications: [true],
        portalAccess: [false]
      })
    });
    this.guardians.push(guardianGroup);
  }

  removeGuardian(index: number): void {
    if (confirm('Tem certeza que deseja remover este responsável?')) {
      this.guardians.removeAt(index);
    }
  }

  copyStudentAddress(guardianIndex: number): void {
      // Address copy logic needs update for structured address
      // For now, simpler to just copy fields if needed, or disable copy for structured address
      // this.guardians.at(guardianIndex).get('address')?.setValue(studentAddress);
  }

  // Mask formatters
  formatCPF(event: any, controlName: string, formGroup: any = null): void {
    const form = (formGroup as FormGroup) || this.personalDataForm;
    let value = event.target.value.replace(/\D/g, '');
    if (value.length > 11) value = value.slice(0, 11);
    form.get(controlName)?.setValue(value, { emitEvent: false });
  }

  formatPhone(event: any, controlName: string, formGroup: any = null): void {
    const form = (formGroup as FormGroup) || this.personalDataForm;
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
    
    // Aplica a máscara visual
    let maskedValue = value;
    if (value.length >= 5) {
      maskedValue = `${value.slice(0, 2)}/${value.slice(2, 4)}/${value.slice(4)}`;
    } else if (value.length >= 3) {
      maskedValue = `${value.slice(0, 2)}/${value.slice(2)}`;
    }
    
    // Atualiza o valor do input visualmente
    event.target.value = maskedValue;
    
    // Se tiver data completa (8 dígitos), converte para objeto Date
    if (value.length === 8) {
      const day = parseInt(value.slice(0, 2), 10);
      const month = parseInt(value.slice(2, 4), 10) - 1; // Meses em JS são 0-11
      const year = parseInt(value.slice(4), 10);
      
      const dateObj = new Date(year, month, day);
      
      // Valida se a data é válida
      if (dateObj.getDate() === day && 
          dateObj.getMonth() === month && 
          dateObj.getFullYear() === year) {
        this.personalDataForm.get('birthDate')?.setValue(dateObj, { emitEvent: true });
      }
    } else {
      // Para valores incompletos, guarda a string temporariamente
      this.personalDataForm.get('birthDate')?.setValue(maskedValue, { emitEvent: false });
    }
  }
  
  // Método auxiliar para converter string DD/MM/YYYY em Date
  parseBrazilianDate(dateString: string): Date | null {
    if (!dateString || typeof dateString !== 'string') return null;
    
    const match = dateString.match(/^(\d{2})\/(\d{2})\/(\d{4})$/);
    if (!match) return null;
    
    const day = parseInt(match[1], 10);
    const month = parseInt(match[2], 10) - 1;
    const year = parseInt(match[3], 10);
    
    const date = new Date(year, month, day);
    
    // Validação de data válida
    if (date.getDate() !== day || date.getMonth() !== month || date.getFullYear() !== year) {
      return null;
    }
    
    return date;
  }

  displayCPF(cpf: string): string {
    if (!cpf) return '';
    const cleaned = cpf.replace(/\D/g, '');
    if (cleaned.length === 11) {
      return cleaned.replace(/(\d{3})(\d{3})(\d{3})(\d{2})/, '$1.$2.$3-$4');
    }
    return cpf;
  }

  submit(): void {
    if (this.personalDataForm.valid && this.studentDataForm.valid) {
      this.isSubmitting = true;

      // Garante que a data de nascimento seja um objeto Date válido
      let birthDateValue = this.personalDataForm.value.birthDate;
      let birthDateObj: Date;
      
      if (birthDateValue instanceof Date) {
        birthDateObj = birthDateValue;
      } else if (typeof birthDateValue === 'string') {
        // Tenta fazer parse da string DD/MM/YYYY
        const parsed = this.parseBrazilianDate(birthDateValue);
        if (parsed) {
          birthDateObj = parsed;
        } else {
          this.snackBar.open('Data de nascimento inválida. Use o formato DD/MM/AAAA.', 'Fechar', { duration: 3000 });
          this.isSubmitting = false;
          return;
        }
      } else {
        this.snackBar.open('Data de nascimento é obrigatória.', 'Fechar', { duration: 3000 });
        this.isSubmitting = false;
        return;
      }
      
      if (this.isEditMode && this.studentId) {
        // UPDATE - CAN change status
        const updateData: UpdateStudentRequest = {
          user: {
            name: this.personalDataForm.value.name,
            email: this.personalDataForm.value.email,
            cpf: this.personalDataForm.value.cpf,
            birthDate: birthDateObj.toISOString(),
            phone: this.personalDataForm.value.phone,
            address: this.personalDataForm.value.address
          },
          programIds: this.studentDataForm.value.programIds || [],
          status: this.studentDataForm.value.status, // ✅ Allowed in UPDATE
          medicalInfo: this.studentDataForm.value.medicalInfo,
          specialNeeds: this.studentDataForm.value.specialNeeds,
          notes: this.studentDataForm.value.notes
        };

        this.studentService.updateStudent(this.studentId, updateData).subscribe({
          next: () => {
            this.snackBar.open('Aluno atualizado com sucesso!', 'Fechar', { duration: 3000 });
            this.router.navigate(['/students']);
          },
          error: (err) => {
            console.error('Error updating student:', err);
            this.snackBar.open('Erro ao atualizar aluno. ' + (err.error?.message || err.message), 'Fechar', { duration: 5000 });
            this.isSubmitting = false;
          }
        });
      } else {
        // CREATE - NO registrationNumber, NO status (backend generates/defaults)
        const createData: CreateStudentRequest = {
          user: {
            name: this.personalDataForm.value.name,
            email: this.personalDataForm.value.email,
            cpf: this.personalDataForm.value.cpf,
            birthDate: birthDateObj.toISOString(),
            phone: this.personalDataForm.value.phone,
            address: this.personalDataForm.value.address,
            profileId: USER_PROFILES.STUDENT, // ✅ Changed from 'profile: student' to 'profileId: 3'
            active: true,
            password: 'temp123' // Temporary password, will be changed by Keycloak
          },
          programIds: this.studentDataForm.value.programIds || [],
          // ❌ registrationNumber removed - auto-generated by backend
          // ❌ status removed - defaults to 'active'
          medicalInfo: this.studentDataForm.value.medicalInfo,
          specialNeeds: this.studentDataForm.value.specialNeeds,
          notes: this.studentDataForm.value.notes
        };

        this.studentService.createStudent(createData).subscribe({
          next: (response) => {
            // ✅ Backend returns student with auto-generated registrationNumber
            this.snackBar.open(
              `Aluno cadastrado com sucesso! Matrícula: ${response.registrationNumber}`,
              'Fechar',
              { duration: 5000 }
            );
            this.router.navigate(['/students']);
          },
          error: (err) => {
            console.error('Error creating student:', err);
            this.snackBar.open('Erro ao cadastrar aluno. ' + (err.error?.message || err.message), 'Fechar', { duration: 5000 });
            this.isSubmitting = false;
          }
        });
      }
    } else {
      this.snackBar.open('Por favor, preencha todos os campos obrigatórios.', 'Fechar', { duration: 3000 });
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
