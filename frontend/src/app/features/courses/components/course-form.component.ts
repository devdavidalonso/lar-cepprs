import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule, FormControl } from '@angular/forms';
import { MatStepperModule } from '@angular/material/stepper';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MAT_DATE_FORMATS, MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';
import { BRAZILIAN_DATE_FORMATS } from '../../../core/utils/date-formats';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatCardModule } from '@angular/material/card';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { CourseService, Course } from '../../../core/services/course.service';
import { LocationService, Location } from '../../../core/services/location.service';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { ProgramService } from '@core/services/program.service';
import { Program } from '@core/models/program.model';

@Component({
  selector: 'app-course-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule,
    MatStepperModule,
    MatInputModule,
    MatButtonModule,
    MatSelectModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatIconModule,
    MatChipsModule,
    MatCardModule,
    MatSnackBarModule,
    MatAutocompleteModule,
    MatTooltipModule,
    MatCheckboxModule,
    RouterModule
  ],
  providers: [
    { provide: MAT_DATE_LOCALE, useValue: 'pt-BR' },
    { provide: MAT_DATE_FORMATS, useValue: BRAZILIAN_DATE_FORMATS }
  ],
  template: `
    <div class="course-form-container">
      <mat-card class="form-card">
        <mat-card-header>
          <div mat-card-avatar>
              <mat-icon class="header-icon">school</mat-icon>
          </div>
          <mat-card-title>{{ isEditMode ? 'Edit Course' : 'Create New Course' }}</mat-card-title>
          <mat-card-subtitle>Course Creation Wizard</mat-card-subtitle>
        </mat-card-header>
        
        <mat-card-content>
          <mat-stepper [linear]="true" #stepper orientation="vertical">
            
            <!-- Step 1: Basic Info -->
            <mat-step [stepControl]="basicInfoForm">
              <form [formGroup]="basicInfoForm">
                <ng-template matStepLabel>Basic Information</ng-template>
                <p class="step-desc">Define the core identity of the course.</p>
                
                <div class="form-grid">
                  <!-- Program, Name & Category -->
                  <mat-form-field appearance="outline" class="full-width">
                    <mat-label>Programa</mat-label>
                    <mat-select formControlName="programId">
                      <mat-option *ngFor="let program of programs" [value]="program.id">
                        {{ program.name }}
                      </mat-option>
                    </mat-select>
                    <mat-hint>Cada curso pertence a um programa.</mat-hint>
                    <mat-error *ngIf="basicInfoForm.get('programId')?.hasError('required')">
                      Programa é obrigatório.
                    </mat-error>
                  </mat-form-field>

                  <div class="row">
                      <mat-form-field appearance="outline" class="flex-2">
                        <mat-label>{{ 'COURSE.NAME' | translate }}</mat-label>
                        <input matInput formControlName="name" placeholder="Ex: Introduction to Python">
                        <mat-error *ngIf="basicInfoForm.get('name')?.hasError('required')">{{ 'VALIDATION.REQUIRED' | translate }}</mat-error>
                      </mat-form-field>

                      <mat-form-field appearance="outline" class="flex-1">
                          <mat-label>{{ 'COURSE.CATEGORY' | translate }}</mat-label>
                          <mat-select formControlName="category">
                              <mat-option value="Technology">Technology</mat-option>
                              <mat-option value="Arts">Arts & Crafts</mat-option>
                              <mat-option value="Languages">Languages</mat-option>
                              <mat-option value="Health">Health & Wellness</mat-option>
                              <mat-option value="Culinary">Culinary</mat-option>
                          </mat-select>
                      </mat-form-field>
                  </div>

                  <!-- Short Description -->
                  <mat-form-field appearance="outline" class="full-width">
                    <mat-label>{{ 'COURSE.SHORT_DESCRIPTION' | translate }}</mat-label>
                    <input matInput formControlName="shortDescription" placeholder="Brief summary for course cards">
                  </mat-form-field>

                  <!-- Cover Image URL -->
                  <div class="row">
                    <mat-form-field appearance="outline" class="flex-3">
                        <mat-label>Cover Image URL</mat-label>
                        <input matInput formControlName="coverImage" placeholder="https://example.com/image.jpg">
                        <mat-icon matSuffix>image</mat-icon>
                    </mat-form-field>
                    
                    <div class="image-preview-mini" *ngIf="basicInfoForm.get('coverImage')?.value">
                        <img [src]="basicInfoForm.get('coverImage')?.value" (error)="onImageError($event)">
                    </div>
                  </div>

                  <!-- Google Classroom -->
                   <mat-form-field appearance="outline" class="full-width">
                    <mat-label>{{ 'COURSE.GOOGLE_CLASSROOM_URL' | translate }}</mat-label>
                    <input matInput formControlName="googleClassroomUrl" placeholder="https://classroom.google.com/c/...">
                    <mat-icon matSuffix>class</mat-icon>
                    <mat-hint>Optional link to Google Classroom</mat-hint>
                  </mat-form-field>
                </div>

                <div class="stepper-actions">
                  <button mat-button color="warn" routerLink="/courses">Cancel</button>
                  <button mat-raised-button color="primary" matStepperNext>Next: Team</button>
                </div>
              </form>
            </mat-step>

            <!-- Step 2: Team (Teachers) -->
            <mat-step [stepControl]="teamForm">
              <form [formGroup]="teamForm">
                <ng-template matStepLabel>Team Allocation</ng-template>
                <p class="step-desc">Assign professors responsible for this course.</p>
                
                <div class="form-grid">
                  <mat-form-field appearance="outline" class="full-width">
                    <mat-label>Main Professor (Titular)</mat-label>
                    <mat-select formControlName="teacherId">
                      <mat-option *ngFor="let prof of professors" [value]="prof.id">
                        {{ prof.name }}
                      </mat-option>
                    </mat-select>
                    <mat-icon matSuffix>person</mat-icon>
                    <mat-hint>Adding a main professor is mandatory</mat-hint>
                  </mat-form-field>

                  <!-- Placeholder for assistants implementation -->
                  <div class="assistants-section">
                      <h4>Assistant Professors (Monitores)</h4>
                      <p class="text-muted">Multi-selection of assistants will be available here.</p>
                  </div>
                </div>

                <div class="stepper-actions">
                  <button mat-button matStepperPrevious>Back</button>
                  <button mat-raised-button color="primary" matStepperNext>Next: Venue</button>
                </div>
              </form>
            </mat-step>

            <!-- Step 3: Venue & Capacity -->
            <mat-step [stepControl]="venueForm">
              <form [formGroup]="venueForm">
                <ng-template matStepLabel>Venue & Capacity</ng-template>
                 <p class="step-desc">Where will the course take place and for how many students?</p>
                
                 <div class="form-grid">
                    <mat-form-field appearance="outline" class="full-width">
                        <mat-label>{{ 'COURSE.LOCATION' | translate }}</mat-label>
                        <mat-select formControlName="locationId">
                            <mat-option [value]="null">-- To be defined --</mat-option>
                             <mat-option *ngFor="let loc of locations" [value]="loc.id" [disabled]="!loc.isAvailable">
                                {{ loc.name }} (Cap: {{ loc.capacity }})
                            </mat-option>
                        </mat-select>
                        <mat-icon matSuffix>place</mat-icon>
                    </mat-form-field>

                    <div class="row">
                        <mat-form-field appearance="outline">
                            <mat-label>Max Students</mat-label>
                            <input matInput type="number" formControlName="maxStudents">
                        </mat-form-field>

                        <mat-form-field appearance="outline">
                            <mat-label>Target Audience</mat-label>
                            <input matInput formControlName="targetAudience">
                        </mat-form-field>
                    </div>
                 </div>

                <div class="stepper-actions">
                  <button mat-button matStepperPrevious>Back</button>
                  <button mat-raised-button color="primary" matStepperNext>Next: Schedule</button>
                </div>
              </form>
            </mat-step>

            <!-- Step 4: Schedule -->
            <mat-step [stepControl]="scheduleForm">
               <form [formGroup]="scheduleForm">
                <ng-template matStepLabel>Schedule Generation</ng-template>
                <p class="step-desc">Define the timeline. We will generate the class sessions automatically.</p>

                <div class="form-grid">
                    <div class="row">
                         <mat-form-field appearance="outline">
                            <mat-label>Start Date</mat-label>
                            <input matInput 
                                   [matDatepicker]="pickerStart" 
                                   formControlName="startDate"
                                   #startDateInput
                                   placeholder="DD/MM/AAAA"
                                   maxlength="10"
                                   autocomplete="off">
                            <mat-datepicker-toggle matSuffix [for]="pickerStart" matTooltip="Abrir calendário"></mat-datepicker-toggle>
                            <mat-datepicker #pickerStart></mat-datepicker>
                            <mat-hint>Digite DD/MM/AAAA ou use o calendário</mat-hint>
                            <mat-error *ngIf="scheduleForm.get('startDate')?.hasError('matDatepickerParse')">
                                Data inválida. Use DD/MM/AAAA
                            </mat-error>
                            <mat-error *ngIf="scheduleForm.get('startDate')?.hasError('required')">
                                Data inicial é obrigatória
                            </mat-error>
                        </mat-form-field>

                        <mat-form-field appearance="outline">
                            <mat-label>End Date</mat-label>
                            <input matInput 
                                   [matDatepicker]="pickerEnd" 
                                   formControlName="endDate"
                                   #endDateInput
                                   placeholder="DD/MM/AAAA"
                                   maxlength="10"
                                   autocomplete="off">
                            <mat-datepicker-toggle matSuffix [for]="pickerEnd" matTooltip="Abrir calendário"></mat-datepicker-toggle>
                            <mat-datepicker #pickerEnd></mat-datepicker>
                            <mat-hint>Digite DD/MM/AAAA ou use o calendário</mat-hint>
                            <mat-error *ngIf="scheduleForm.get('endDate')?.hasError('matDatepickerParse')">
                                Data inválida. Use DD/MM/AAAA
                            </mat-error>
                            <mat-error *ngIf="scheduleForm.get('endDate')?.hasError('required')">
                                Data final é obrigatória
                            </mat-error>
                        </mat-form-field>
                  </div>

                     <div class="row">
                         <mat-form-field appearance="outline">
                            <mat-label>Week Days</mat-label>
                             <mat-select formControlName="weekDays" multiple>
                                <mat-option value="Mon">Monday</mat-option>
                                <mat-option value="Tue">Tuesday</mat-option>
                                <mat-option value="Wed">Wednesday</mat-option>
                                <mat-option value="Thu">Thursday</mat-option>
                                <mat-option value="Fri">Friday</mat-option>
                                <mat-option value="Sat">Saturday</mat-option>
                            </mat-select>
                        </mat-form-field>
                     </div>
                     <div class="row">
                        <mat-form-field appearance="outline">
                            <mat-label>Start Time</mat-label>
                            <input matInput type="time" formControlName="startTime">
                        </mat-form-field>

                        <mat-form-field appearance="outline">
                            <mat-label>End Time</mat-label>
                            <input matInput type="time" formControlName="endTime">
                        </mat-form-field>
                     </div>

                     <div class="row" style="justify-content: flex-end; margin-top: 10px;">
                       <button mat-stroked-button color="accent" (click)="generateSchedule()" type="button">
                           <mat-icon>autorenew</mat-icon> Generate Schedule
                       </button>
                     </div>

                     <div class="schedule-preview" *ngIf="generatedSchedule.length > 0">
                       <h4>Preview Schedule ({{ generatedSchedule.length }} sessions)</h4>
                       <div class="schedule-table-wrapper">
                         <table class="schedule-table">
                           <thead>
                             <tr>
                               <th>Date</th>
                               <th>Week Day</th>
                               <th>Start</th>
                               <th>End</th>
                               <th class="actions-col">Actions</th>
                             </tr>
                           </thead>
                           <tbody>
                             <tr *ngFor="let session of generatedSchedule; let i = index">
                               <td>{{ session.date | date:'dd/MM/yyyy' }}</td>
                               <td>{{ session.date | date:'EEEE' }}</td>
                               <td>{{ session.startTime }}</td>
                               <td>{{ session.endTime }}</td>
                               <td class="actions-col">
                                 <button mat-icon-button color="warn" type="button" (click)="removeSession(i)" matTooltip="Remover sessão">
                                   <mat-icon>delete</mat-icon>
                                 </button>
                               </td>
                             </tr>
                           </tbody>
                         </table>
                       </div>
                     </div>
                </div>

                <div class="stepper-actions">
                  <button mat-button matStepperPrevious>Back</button>
                  <button mat-raised-button color="primary" matStepperNext>Review</button>
                </div>
               </form>
            </mat-step>

            <!-- Step 5: Review -->
            <mat-step>
              <ng-template matStepLabel>Review & Publish</ng-template>
              
              <div class="review-container">
                <h3>Summary</h3>
                <p>Please review all details before publishing.</p>
                
                <div class="review-card">
                    <h4>{{ basicInfoForm.get('name')?.value }}</h4>
                    <p>{{ basicInfoForm.get('shortDescription')?.value }}</p>
                    <p><strong>Category:</strong> {{ basicInfoForm.get('category')?.value }}</p>
                    <p><strong>Programa:</strong> {{ getProgramName(basicInfoForm.get('programId')?.value) }}</p>
                    <p><strong>Professor:</strong> {{ getProfessorName(teamForm.get('teacherId')?.value) }}</p>
                    <p><strong>Location:</strong> {{ getLocationName(venueForm.get('locationId')?.value) }}</p>
                    <hr>
                    <p><strong>Total Workload:</strong> {{ scheduleForm.get('workload')?.value }}h</p>
                    <p><strong>Sessions:</strong> {{ generatedSchedule.length }} classes generated</p>
                </div>

                <div class="stepper-actions">
                  <button mat-button matStepperPrevious>Back</button>
                  <button mat-raised-button color="accent" (click)="submit()" [disabled]="isSubmitting">
                    {{ isSubmitting ? 'Publishing...' : 'Publish Course' }}
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
    .course-form-container {
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
    }
    
    .flex-1 { flex: 1; }
    .flex-2 { flex: 2; }
    .flex-3 { flex: 3; }

    .full-width {
      width: 100%;
    }

    .stepper-actions {
      display: flex;
      justify-content: flex-end;
      gap: 12px;
      margin-top: 32px;
      padding-top: 16px;
      border-top: 1px solid #eee;
    }

    .image-preview-mini {
        width: 60px;
        height: 50px;
        border-radius: 4px;
        overflow: hidden;
        border: 1px solid #ddd;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .image-preview-mini img {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }
    
    .text-muted { color: #888; }
    
    .review-card {
        background: #f0f4f8;
        padding: 16px;
        border-radius: 8px;
        text-align: left;
        margin-top: 16px;
    }

    .schedule-preview {
      margin-top: 8px;
      border: 1px solid #e5e7eb;
      border-radius: 8px;
      padding: 12px;
      background: #fcfdff;
    }

    .schedule-preview h4 {
      margin: 0 0 10px;
    }

    .schedule-table-wrapper {
      max-height: 260px;
      overflow: auto;
      border: 1px solid #eef2f7;
      border-radius: 6px;
      background: #fff;
    }

    .schedule-table {
      width: 100%;
      border-collapse: collapse;
      min-width: 560px;
    }

    .schedule-table th,
    .schedule-table td {
      padding: 8px 10px;
      border-bottom: 1px solid #f1f5f9;
      text-align: left;
      white-space: nowrap;
    }

    .schedule-table thead th {
      position: sticky;
      top: 0;
      background: #f8fafc;
      z-index: 1;
      font-weight: 600;
    }

    .actions-col {
      width: 72px;
      text-align: center !important;
    }
  `]
})
export class CourseFormComponent implements OnInit {
  @ViewChild('startDateInput') startDateInput?: ElementRef<HTMLInputElement>;
  @ViewChild('endDateInput') endDateInput?: ElementRef<HTMLInputElement>;

  // Form Groups for each Step
  basicInfoForm: FormGroup;
  teamForm: FormGroup;
  venueForm: FormGroup;
  scheduleForm: FormGroup;
  
  isSubmitting = false;
  isEditMode = false;
  courseId: number | null = null;
  
  // Data Sources
  professors: any[] = [];
  locations: Location[] = [];
  programs: Program[] = [];
  
  generatedSchedule: any[] = [];

  constructor(
    private _formBuilder: FormBuilder,
    private courseService: CourseService,
    private locationService: LocationService,
    private programService: ProgramService,
    private router: Router,
    private route: ActivatedRoute,
    private snackBar: MatSnackBar
  ) {
    // Step 1: Basic Info
    this.basicInfoForm = this._formBuilder.group({
      programId: [null, Validators.required],
      name: ['', Validators.required],
      category: ['Technology', Validators.required],
      shortDescription: ['', Validators.maxLength(150)],
      coverImage: [''],
      googleClassroomUrl: ['', [Validators.pattern('https?://.*')]]
    });

    // Step 2: Team
    this.teamForm = this._formBuilder.group({
      teacherId: ['', Validators.required], // Titular
      assistantIds: [[]] // Optional
    });

    // Step 3: Venue
    this.venueForm = this._formBuilder.group({
      locationId: [null],
      maxStudents: [20, [Validators.required, Validators.min(1)]],
      targetAudience: ['']
    });

    // Step 4: Schedule
    this.scheduleForm = this._formBuilder.group({
      startDate: [new Date(), Validators.required],
      endDate: [null, Validators.required],
      weekDays: [[], Validators.required],
      startTime: ['', Validators.required],
      endTime: ['', Validators.required],
      workload: [0] // Will be calculated
    });
  }

  ngOnInit(): void {
    this.loadPrograms();
    this.loadProfessors();
    this.loadLocations();
    this.checkEditMode();
  }

  loadPrograms() {
    this.programService.listPrograms(true).subscribe({
      next: (programs) => this.programs = programs,
      error: (err) => console.error('Error loading programs', err)
    });
  }

  loadProfessors() {
    this.courseService.getProfessors().subscribe({
      next: (profs) => this.professors = profs,
      error: (err) => console.error('Error loading professors', err)
    });
  }

  loadLocations() {
    this.locationService.getLocations().subscribe({
      next: (locs) => this.locations = locs,
      error: (err) => console.error('Error loading locations', err)
    });
  }

  checkEditMode() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEditMode = true;
      this.courseId = +id;
      this.loadCourseData(this.courseId);
    }
  }

  loadCourseData(id: number) {
    this.courseService.getCourse(id).subscribe({
      next: (course) => {
        // Patch Basic Info
        this.basicInfoForm.patchValue({
          programId: course.programId ?? null,
          name: course.name,
          category: course.category || 'Technology',
          shortDescription: course.shortDescription,
          coverImage: course.coverImage,
          googleClassroomUrl: course.googleClassroomUrl
        });

        // Patch Team
        this.teamForm.patchValue({
          teacherId: course.teacherId
        });

        // Patch Venue
        this.venueForm.patchValue({
          locationId: course.locationId,
          maxStudents: course.maxStudents,
          targetAudience: course.targetAudience
        });

        // Patch Schedule
        const weekDays = course.weekDays ? course.weekDays.split(',') : [];
        this.scheduleForm.patchValue({
          startDate: course.startDate,
          endDate: course.endDate,
          weekDays: weekDays,
          startTime: course.startTime,
          endTime: course.endTime,
          workload: course.workload
        });
      },
      error: (err) => {
        this.snackBar.open('Error loading course data.', 'Close', { duration: 3000 });
        this.router.navigate(['/courses']);
      }
    });
  }

  onImageError(event: any) {
    event.target.src = 'https://via.placeholder.com/60x50?text=No+Img';
  }

  getProfessorName(id: number): string {
      const prof = this.professors.find(p => p.id === id);
      return prof ? prof.name : 'Unknown';
  }

  getProgramName(id: number | null | undefined): string {
      if (!id) return 'Não definido';
      const program = this.programs.find((p) => p.id === id);
      return program ? program.name : 'Não definido';
  }

  getLocationName(id: number): string {
      const loc = this.locations.find(l => l.id === id);
      return loc ? loc.name : 'Unknown';
  }

  generateSchedule() {
    const { startDate, endDate, weekDays, startTime, endTime } = this.scheduleForm.value;

    let normalizedStartDate = this.parseDateValue(startDate);
    let normalizedEndDate = this.parseDateValue(endDate);
    const normalizedWeekDays = Array.isArray(weekDays) ? weekDays : (weekDays ? [weekDays] : []);
    const normalizedStartTime = this.parseTimeValue(startTime);
    const normalizedEndTime = this.parseTimeValue(endTime);

    // Fallback: quando o MatDatepicker nao parseia o valor digitado,
    // o FormControl pode ficar nulo, mas o texto existe no input.
    if (!normalizedStartDate && this.startDateInput?.nativeElement?.value) {
      normalizedStartDate = this.parseDateValue(this.startDateInput.nativeElement.value);
    }
    if (!normalizedEndDate && this.endDateInput?.nativeElement?.value) {
      normalizedEndDate = this.parseDateValue(this.endDateInput.nativeElement.value);
    }

    if (!normalizedStartDate || !normalizedEndDate || normalizedWeekDays.length === 0 || !normalizedStartTime || !normalizedEndTime) {
      this.snackBar.open('Please select Start Date, End Date, Week Days and start/end times.', 'Close', { duration: 3000 });
      return;
    }

    this.generatedSchedule = [];
    let currentDate = new Date(normalizedStartDate);
    const end = new Date(normalizedEndDate);
    
    // Map week day strings to integers (Sun=0, Mon=1, ...)
    const dayMap: { [key: string]: number } = {
       'Sun': 0, 'Mon': 1, 'Tue': 2, 'Wed': 3, 'Thu': 4, 'Fri': 5, 'Sat': 6
    };
    
    const selectedDays = normalizedWeekDays.map((d: string) => dayMap[d]).filter((d: number) => d !== undefined);

    if (selectedDays.length === 0) {
      this.snackBar.open('Invalid week day selection. Please reselect week days.', 'Close', { duration: 3000 });
      return;
    }

    while (currentDate <= end) {
      if (selectedDays.includes(currentDate.getDay())) {
        this.generatedSchedule.push({
           date: new Date(currentDate),
           startTime: normalizedStartTime,
           endTime: normalizedEndTime,
           topic: 'Class Session' // Default topic
        });
      }
      // Add 1 day
      currentDate.setDate(currentDate.getDate() + 1);
    }
    
    // Auto-calculate workload
    if (normalizedStartTime && normalizedEndTime) {
        const start = parseInt(normalizedStartTime.split(':')[0]);
        const endHour = parseInt(normalizedEndTime.split(':')[0]);
        const hoursPerSession = endHour - start;
        const totalHours = this.generatedSchedule.length * hoursPerSession;
        this.scheduleForm.patchValue({ workload: totalHours });
    }
    
    this.snackBar.open(`Generated ${this.generatedSchedule.length} class sessions.`, 'OK', { duration: 3000 });
  }

  private parseDateValue(value: unknown): Date | null {
    if (!value) {
      return null;
    }

    if (value instanceof Date && !isNaN(value.getTime())) {
      return value;
    }

    if (typeof value === 'string') {
      const trimmed = value.trim();
      const br = /^(\d{2})\/(\d{2})\/(\d{4})$/.exec(trimmed);
      if (br) {
        const day = Number(br[1]);
        const month = Number(br[2]) - 1;
        const year = Number(br[3]);
        const parsed = new Date(year, month, day);
        if (!isNaN(parsed.getTime())) {
          return parsed;
        }
      }

      const iso = new Date(trimmed);
      if (!isNaN(iso.getTime())) {
        return iso;
      }
    }

    return null;
  }

  private parseTimeValue(value: unknown): string {
    if (typeof value !== 'string') {
      return '';
    }

    const match = value.match(/(\d{2}:\d{2})/);
    return match ? match[1] : '';
  }

  removeSession(index: number) {
    if (index < 0 || index >= this.generatedSchedule.length) {
      return;
    }

    this.generatedSchedule.splice(index, 1);
    this.generatedSchedule = [...this.generatedSchedule];

    const { startTime, endTime } = this.scheduleForm.value;
    if (startTime && endTime) {
      const start = parseInt(startTime.split(':')[0], 10);
      const endHour = parseInt(endTime.split(':')[0], 10);
      const hoursPerSession = Math.max(endHour - start, 0);
      this.scheduleForm.patchValue({ workload: this.generatedSchedule.length * hoursPerSession });
    }
  }

  submit() {
    if (this.basicInfoForm.valid && this.teamForm.valid && this.venueForm.valid && this.scheduleForm.valid) {
      this.isSubmitting = true;

      // Combine all forms into one Course object
      const courseData: Course = {
        ...this.basicInfoForm.value,
        ...this.teamForm.value,
        ...this.venueForm.value,
        ...this.scheduleForm.value,
        weekDays: this.scheduleForm.value.weekDays ? this.scheduleForm.value.weekDays.join(',') : '',
        status: 'active'
      };

      const request = this.isEditMode && this.courseId
        ? this.courseService.updateCourse(this.courseId, courseData)
        : this.courseService.createCourse(courseData);

      request.subscribe({
        next: (res: Course) => {
          this.snackBar.open(this.isEditMode ? 'Course updated successfully!' : 'Course created successfully!', 'Close', { duration: 3000 });
          this.router.navigate(['/courses']);
        },
        error: (err: any) => {
          console.error(err);
          this.snackBar.open(`Error ${this.isEditMode ? 'updating' : 'creating'} course. Please try again.`, 'Close', { duration: 5000 });
          this.isSubmitting = false;
        }
      });
    } else {
      this.snackBar.open('Please fill all required fields in the wizard.', 'Close', { duration: 3000 });
    }
  }
}
