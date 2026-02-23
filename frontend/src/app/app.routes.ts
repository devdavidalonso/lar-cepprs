import { Routes } from '@angular/router';
import { AuthGuard } from './core/guards/auth.guard';
import { LayoutComponent } from './layout/layout.component';

export const routes: Routes = [

  {
    path: '',
    redirectTo: 'dashboard',
    pathMatch: 'full'
  },
  {
    path: 'auth',
    loadChildren: () => import('./features/auth/auth.routes').then(m => m.AUTH_ROUTES)
  },
  {
    path: 'access-denied',
    loadComponent: () => import('./shared/components/access-denied/access-denied.component').then(m => m.AccessDeniedComponent)
  },
  {
    path: 'acesso-negado',
    redirectTo: 'access-denied',
    pathMatch: 'full'
  },
  {
    path: '',
    component: LayoutComponent, // Use LayoutComponent directly in the route
    canActivate: [AuthGuard],
    children: [
      {
        path: 'dashboard',
        loadComponent: () => import('./features/dashboard/dashboard.component').then(m => m.DashboardComponent)
      },
      {
        path: 'students',
        loadChildren: () => import('./features/students/students.routes').then(m => m.STUDENTS_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'administrador'] }
      },
      {
        path: 'courses',
        loadChildren: () => import('./features/courses/courses.routes').then(m => m.COURSES_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'administrador'] }
      },
      {
        path: 'enrollments',
        loadChildren: () => import('./features/enrollments/enrollments.routes').then(m => m.ENROLLMENTS_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'administrador'] }
      },
      {
        path: 'attendance',
        loadChildren: () => import('./features/attendance/attendance.routes').then(m => m.ATTENDANCE_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'administrador'] }
      },
      {
        path: 'reports',
        loadChildren: () => import('./features/reports/reports.routes').then(m => m.REPORTS_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'administrador'] }
      },
      {
        path: 'teachers',
        loadChildren: () => import('./features/teachers/teachers.routes').then(m => m.TEACHERS_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'administrador'] } 
      },
      {
        path: 'teacher',
        loadChildren: () => import('./features/teacher-portal/teacher-portal.routes').then(m => m.TEACHER_PORTAL_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['professor', 'admin', 'administrador'] }
      },
      {
        path: 'student',
        loadChildren: () => import('./features/student-portal/student-portal.routes').then(m => m.STUDENT_PORTAL_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['aluno', 'responsavel', 'responsável', 'admin', 'administrador'] }
      },
      {
        path: 'administration',
        loadChildren: () => import('./features/administration/administration.routes').then(m => m.ADMINISTRATION_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'administrador'] }
      },
      {
        path: 'interviews',
        loadChildren: () => import('./features/interviews/interviews.routes').then(m => m.INTERVIEWS_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'administrador'] }
      },
      {
        path: 'profile',
        loadComponent: () => import('./features/profile/profile.component').then(m => m.ProfileComponent),
        canActivate: [AuthGuard]
      },
      {
        path: 'volunteering',
        loadChildren: () => import('./features/volunteering/volunteering.routes').then(m => m.VOLUNTEERING_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'administrador'] }
      },
      {
        path: 'locations',
        loadChildren: () => import('./features/locations/locations.routes').then(m => m.LOCATIONS_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'administrador'] }
      },
      {
        path: 'class-sessions',
        loadChildren: () => import('./features/class-sessions/class-sessions.routes').then(m => m.CLASS_SESSIONS_ROUTES),
        canActivate: [AuthGuard],
        data: { roles: ['admin', 'professor'] }
      },
    ]
  },
  {
    path: '**',
    loadComponent: () => import('./shared/components/not-found/not-found.component').then(m => m.NotFoundComponent)
  }
];
