// src/app/core/guards/auth.guard.ts
import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

import { SsoService } from '../services/sso.service';

export const AuthGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);
  const ssoService = inject(SsoService);

  if (authService.checkAuth()) {
    const requiredRoles = (route.data?.['roles'] as string[] | undefined) ?? [];
    if (requiredRoles.length > 0) {
      const user = authService.getCurrentUser();
      const userRoles = (user?.roles ?? []).map(role => role.toLowerCase());
      const normalizedRequired = requiredRoles.map(role => role.toLowerCase());
      const authorized = normalizedRequired.some(role => userRoles.includes(role));

      if (!authorized) {
        router.navigate(['/access-denied']);
        return false;
      }
    }
    return true;
  }

  // Redirect to the login page
  // Redirect to the login page (SSO)
  ssoService.login();

  return false;
};
