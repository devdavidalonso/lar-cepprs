import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

// Canonical guard for management/admin area.
// Keeps compatibility with both "admin/administrador" and "gestor" roles.
export const AdminGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  const isAuthorized =
    authService.checkAuth() &&
    (authService.hasRole('admin') || authService.hasRole('administrator') || authService.hasRole('administrador') || authService.hasRole('gestor'));

  if (isAuthorized) {
    return true;
  }

  if (!authService.checkAuth()) {
    router.navigate(['/auth/login'], { queryParams: { returnUrl: state.url } });
  } else {
    router.navigate(['/access-denied']);
  }

  return false;
};
