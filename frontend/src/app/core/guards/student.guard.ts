// src/app/core/guards/student.guard.ts
import { Injectable, inject } from '@angular/core';
import { CanActivate, Router, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { USER_PROFILES } from '../models/user.model';
import { Observable, map, take } from 'rxjs';

/**
 * Guard para proteger rotas do Portal do Aluno
 * Verifica se o usuário está autenticado e tem perfil de aluno
 */
@Injectable({
  providedIn: 'root'
})
export class StudentGuard implements CanActivate {
  private authService = inject(AuthService);
  private router = inject(Router);

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<boolean> | Promise<boolean> | boolean {
    return this.authService.currentUser$.pipe(
      take(1),
      map(user => {
        // Verificar se está autenticado
        if (!user) {
          console.log('[StudentGuard] Usuário não autenticado, redirecionando para login');
          this.router.navigate(['/login'], { queryParams: { returnUrl: state.url } });
          return false;
        }

        // Verificar se tem perfil de aluno
        const isStudent = user.profileId === USER_PROFILES.STUDENT || 
                         user.profile?.name === 'student';
        
        if (!isStudent) {
          console.log('[StudentGuard] Usuário não é aluno, redirecionando');
          
          // Redirecionar para o portal correto baseado no perfil
          if (user.profileId === USER_PROFILES.ADMIN || user.profile?.name === 'admin' || user.profile?.name === 'administrator') {
            this.router.navigate(['/admin/dashboard']);
          } else if (user.profileId === USER_PROFILES.PROFESSOR || user.profile?.name === 'professor' || user.profile?.name === 'teacher') {
            this.router.navigate(['/teacher/dashboard']);
          } else {
            this.router.navigate(['/']);
          }
          return false;
        }

        console.log('[StudentGuard] Acesso permitido para aluno:', user.email);
        return true;
      })
    );
  }
}

/**
 * Guard funcional (Angular 15+) - alternativa moderna
 */
export const studentGuard = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return authService.currentUser$.pipe(
    take(1),
    map(user => {
      if (!user) {
        router.navigate(['/login']);
        return false;
      }

      const isStudent = user.profileId === USER_PROFILES.STUDENT || 
                       user.profile?.name === 'student';
      
      if (!isStudent) {
        if (user.profileId === USER_PROFILES.ADMIN || user.profile?.name === 'admin' || user.profile?.name === 'administrator') {
          router.navigate(['/admin/dashboard']);
        } else if (user.profileId === USER_PROFILES.PROFESSOR || user.profile?.name === 'professor' || user.profile?.name === 'teacher') {
          router.navigate(['/teacher/dashboard']);
        } else {
          router.navigate(['/']);
        }
        return false;
      }

      return true;
    })
  );
};
