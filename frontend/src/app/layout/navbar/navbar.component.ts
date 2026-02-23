// src/app/layout/navbar/navbar.component.ts
import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatSidenav } from '@angular/material/sidenav';
import { Observable } from 'rxjs';

import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatMenuModule
  ],
  template: `
    <mat-toolbar color="primary" class="navbar">
      <button 
        mat-icon-button 
        class="menu-button" 
        (click)="toggleSidenav()"
        aria-label="Menu button">
        <mat-icon>menu</mat-icon>
      </button>
      
      <a routerLink="/" class="logo-container">
        <img 
          src="assets/images/logo.png" 
          alt="CECOR Logo" 
          class="logo-image"
          height="40">
        <span class="logo-text">CECOR</span>
      </a>
      
      <span class="spacer"></span>
      
      <!-- Menu de navegação para desktop -->
      <div class="nav-links">
        <a mat-button routerLink="/dashboard">Dashboard</a>
        <a mat-button routerLink="/students">Alunos</a>
        <a mat-button routerLink="/courses">Cursos</a>
      </div>
      
      <!-- Notificações -->
      <button mat-icon-button aria-label="Ver notificações">
        <mat-icon>notifications</mat-icon>
      </button>
      
      <!-- Botão de perfil -->
      <button mat-icon-button [matMenuTriggerFor]="userMenu" aria-label="Menu do usuário">
        <mat-icon>person</mat-icon>
      </button>
      <mat-menu #userMenu="matMenu">
        <div class="user-info" *ngIf="currentUser$ | async as user">
          <div class="user-name">{{user.name}}</div>
          <div class="user-email">{{user.email}}</div>
        </div>
        <a mat-menu-item routerLink="/profile">
          <mat-icon>account_circle</mat-icon>
          <span>Meu Perfil</span>
        </a>
        <a mat-menu-item routerLink="/administration/settings">
          <mat-icon>settings</mat-icon>
          <span>Configurações</span>
        </a>
        <button mat-menu-item (click)="logout()">
          <mat-icon>exit_to_app</mat-icon>
          <span>Sair</span>
        </button>
      </mat-menu>
    </mat-toolbar>
  `,
  styles: [`
    .navbar {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      z-index: 999;
      height: 64px;
      padding: 0 16px;
    }
    
    .menu-button {
      display: block;
    }
    
    .logo-container {
      display: flex;
      align-items: center;
      text-decoration: none;
      color: white;
      margin-left: 8px;
    }
    
    .logo-image {
      height: 40px;
      margin-right: 8px;
    }
    
    .logo-text {
      font-size: 20px;
      font-weight: 500;
    }
    
    .spacer {
      flex: 1 1 auto;
    }
    
    .nav-links {
      display: flex;
      gap: 8px;
    }
    
    .user-info {
      padding: 16px;
      min-width: 180px;
      border-bottom: 1px solid rgba(0, 0, 0, 0.12);
    }
    
    .user-name {
      font-weight: 500;
    }
    
    .user-email {
      font-size: 12px;
      opacity: 0.7;
    }
    
    @media (max-width: 768px) {
      .nav-links {
        display: none;
      }
      
      .logo-text {
        display: none;
      }
    }
  `]
})
export class NavbarComponent {
  @Input() sidenav!: MatSidenav;
  
  currentUser$: Observable<any>;
  
  constructor(private authService: AuthService) {
    this.currentUser$ = this.authService.currentUser$;
  }
  
  toggleSidenav(): void {
    this.sidenav.toggle();
  }
  
  logout(): void {
    this.authService.logout();
  }
}
