// src/app/core/guards/teacher.guard.ts
import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

export const TeacherGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  const isTeacherOrAdmin =
    authService.checkAuth() &&
    (authService.hasRole('teacher') || authService.hasRole('professor') || authService.hasRole('admin') || authService.hasRole('administrator') || authService.hasRole('administrador'));

  if (isTeacherOrAdmin) {
    return true;
  }

  if (!authService.checkAuth()) {
    router.navigate(['/auth/login'], { queryParams: { returnUrl: state.url } });
  } else {
    router.navigate(['/access-denied']);
  }

  return false;
};
