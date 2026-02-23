import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';

import { AuthService } from '../../../core/services/auth.service';

@Component({
  selector: 'app-access-denied',
  standalone: true,
  imports: [CommonModule, RouterModule, MatButtonModule, MatIconModule, MatCardModule],
  template: `
    <div class="access-denied-container">
      <mat-card>
        <mat-card-content>
          <div class="error-icon">
            <mat-icon>block</mat-icon>
          </div>
          <h1>Acesso Negado</h1>
          <p>Você não tem permissão para acessar esta página.</p>
          <p class="user-info" *ngIf="authService.currentUser$ | async as user">
            Olá, <strong>{{ user.name }}</strong>. Seus papéis atuais são
            <strong>{{ formatRoles(user.roles) }}</strong>, o que não permite acesso a este recurso.
          </p>
          <div class="actions">
            <button mat-raised-button color="primary" routerLink="/dashboard">
              <mat-icon>dashboard</mat-icon>
              Ir para o Dashboard
            </button>
            <button mat-button (click)="logout()">
              <mat-icon>exit_to_app</mat-icon>
              Sair
            </button>
          </div>
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [
    `
      .access-denied-container {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        background-color: #f5f5f5;
      }

      mat-card {
        max-width: 500px;
        width: 90%;
        text-align: center;
        padding: 40px 20px;
      }

      .error-icon {
        display: flex;
        justify-content: center;
        margin-bottom: 20px;
      }

      .error-icon mat-icon {
        font-size: 80px;
        height: 80px;
        width: 80px;
        color: #f44336;
      }

      h1 {
        font-size: 28px;
        margin: 0 0 16px;
      }

      p {
        margin-bottom: 16px;
        font-size: 16px;
      }

      .user-info {
        background-color: #f0f0f0;
        padding: 16px;
        border-radius: 4px;
        margin: 20px 0;
      }

      .actions {
        display: flex;
        flex-direction: column;
        gap: 10px;
        margin-top: 30px;
      }

      @media (min-width: 600px) {
        .actions {
          flex-direction: row;
          justify-content: center;
        }
      }
    `,
  ],
})
export class AccessDeniedComponent {
  constructor(public authService: AuthService) {}

  formatRoles(roles?: string[]): string {
    if (!roles || roles.length === 0) return 'indefinido';
    return roles.join(', ');
  }

  logout(): void {
    this.authService.logout();
  }
}
