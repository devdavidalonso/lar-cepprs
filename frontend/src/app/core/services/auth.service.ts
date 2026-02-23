import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';

import { User } from '../models/user.model';
import { SsoService } from './sso.service';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly isAuthenticatedSubject = new BehaviorSubject<boolean>(false);
  isAuthenticated$ = this.isAuthenticatedSubject.asObservable();

  private readonly currentUserSubject = new BehaviorSubject<User | null>(null);
  currentUser$ = this.currentUserSubject.asObservable();

  constructor(
    private readonly router: Router,
    private readonly ssoService: SsoService
  ) {
    this.checkAuth();
  }

  checkAuth(): boolean {
    const isAuth = this.ssoService.isAuthenticated;
    this.isAuthenticatedSubject.next(isAuth);
    
    if (isAuth) {
      const user = this.getUserFromClaims();
      this.currentUserSubject.next(user);
    } else {
      this.currentUserSubject.next(null);
    }
    
    return isAuth;
  }

  private getUserFromClaims(): User {
    const claims: any = this.ssoService.identityClaims;
    const roles = this.ssoService.getUserRoles().map((role: string) => role.toLowerCase());
    
    return {
      id: claims?.sub || '',
      name: this.ssoService.getUserName(),
      email: this.ssoService.getUserEmail(),
      roles: roles,
      profileId: this.mapRolesToProfileId(roles),
      locale: claims?.locale || 'pt-BR' // ✅ Extrai locale do token ou usa padrão
    };
  }

  private mapRolesToProfileId(roles: string[]): number {
    if (roles.includes('administrator') || roles.includes('admin') || roles.includes('administrador') || roles.includes('gestor')) return 1; // admin
    if (roles.includes('teacher') || roles.includes('professor')) return 2; // professor
    if (roles.includes('student') || roles.includes('aluno') || roles.includes('responsavel') || roles.includes('responsável')) return 3; // student
    return 3; // default to student
  }

  login(): void {
    this.ssoService.login();
  }

  logout(): void {
    this.ssoService.logout();
    this.currentUserSubject.next(null);
    this.isAuthenticatedSubject.next(false);
  }

  getToken(): string | null {
    return this.ssoService.accessToken;
  }

  getCurrentUser(): User | null {
    if (!this.currentUserSubject.value && this.ssoService.isAuthenticated) {
      // Lazy load user data if not already loaded
      this.checkAuth();
    }
    return this.currentUserSubject.value;
  }

  hasRole(role: string): boolean {
    const user = this.getCurrentUser();
    return user?.roles.includes(role) || false;
  }

  hasAnyRole(roles: string[]): boolean {
    const user = this.getCurrentUser();
    if (!user) return false;
    
    return roles.some(role => user.roles.includes(role));
  }

  getDefaultRouteByRole(): string {
    const user = this.getCurrentUser();
    if (!user) return '/auth/login';

    if (this.hasAnyRole(['admin', 'administrator', 'administrador', 'gestor'])) {
      return '/dashboard';
    }
    if (this.hasAnyRole(['teacher', 'professor'])) {
      return '/teacher/dashboard';
    }
    if (this.hasAnyRole(['student', 'aluno', 'responsavel', 'responsável'])) {
      return '/student/dashboard';
    }

    return '/dashboard';
  }
}
