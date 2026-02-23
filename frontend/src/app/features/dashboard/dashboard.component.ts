import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatGridListModule } from '@angular/material/grid-list';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

import { AuthService } from '../../core/services/auth.service';

interface DashboardCard {
  title: string;
  subtitle: string;
  icon: string;
  value: string | number;
  route: string;
  color: string;
}

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatIconModule,
    MatButtonModule,
    MatGridListModule,
    TranslateModule
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="dashboard-container">
      <h1 class="dashboard-title">{{ 'NAV.DASHBOARD' | translate }}</h1>
      
      <div class="welcome-card">
        <mat-card>
          <mat-card-content>
            <h2>{{ 'HOME.WELCOME_TITLE' | translate:{name: (authService.currentUser$ | async)?.name || 'User'} }}!</h2>
            <p>{{ 'HOME.WELCOME_MESSAGE' | translate }}</p>
          </mat-card-content>
        </mat-card>
      </div>

      <!-- Logo CECOR centralizado entre os cards -->
      <div class="cecor-logo-divider">
        <img src="assets/images/cecor-logo.png" alt="LAR-CEPPRS" class="cecor-center-logo" />
      </div>
      
      <div class="dashboard-grid">
        <div class="dashboard-card" *ngFor="let card of dashboardCards">
          <mat-card [routerLink]="card.route" class="clickable">
            <mat-card-content>
              <div class="card-content">
                <div class="card-info">
                  <h3>{{ card.title | translate }}</h3>
                  <p>{{ card.subtitle | translate }}</p>
                  <div class="card-value" [style.color]="card.color">{{ card.value }}</div>
                </div>
                <div class="card-icon" [style.background-color]="card.color">
                  <mat-icon>{{ card.icon }}</mat-icon>
                </div>
              </div>
            </mat-card-content>
          </mat-card>
        </div>
      </div>
      
      <div class="actions-container">
        <h2>{{ 'COMMON.QUICK_ACTIONS' | translate }}</h2>
        <div class="quick-actions">
          <button mat-raised-button color="primary" routerLink="/students/new">
            <mat-icon>person_add</mat-icon> {{ 'NAV.STUDENTS_NEW' | translate }}
          </button>
          
          <button mat-raised-button color="accent" routerLink="/enrollments/new">
            <mat-icon>how_to_reg</mat-icon> {{ 'NAV.ENROLLMENTS_NEW' | translate }}
          </button>
          
          <button mat-raised-button color="primary" routerLink="/attendance">
            <mat-icon>fact_check</mat-icon> {{ 'NAV.ATTENDANCE' | translate }}
          </button>
          
          <button mat-raised-button color="accent" routerLink="/reports">
            <mat-icon>assessment</mat-icon> {{ 'NAV.REPORTS' | translate }}
          </button>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .dashboard-container {
      padding: 20px;
      margin-top: 64px; /* Altura do header */
    }
    
    .dashboard-title {
      margin-bottom: 20px;
      color: #3f51b5;
    }
    
    .welcome-card {
      margin-bottom: 24px;
    }
    
    .dashboard-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 20px;
      margin-bottom: 24px;
    }

    .cecor-logo-divider {
      display: flex;
      justify-content: center;
      align-items: center;
      padding: 16px 0 20px;
    }

    .cecor-center-logo {
      width: 220px;
      height: 220px;
      object-fit: contain;
      opacity: 0.9;
    }
    
    .card-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    
    .card-info {
      flex: 1;
    }
    
    .card-info h3 {
      margin: 0;
      font-size: 18px;
    }
    
    .card-info p {
      margin: 4px 0 12px;
      color: rgba(0, 0, 0, 0.6);
    }
    
    .card-value {
      font-size: 24px;
      font-weight: bold;
    }
    
    .card-icon {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 50px;
      height: 50px;
      border-radius: 50%;
      color: white;
    }
    
    .clickable {
      cursor: pointer;
      transition: transform 0.2s, box-shadow 0.2s;
    }
    
    .clickable:hover {
      transform: translateY(-5px);
      box-shadow: 0 6px 10px rgba(0, 0, 0, 0.15);
    }
    
    .actions-container {
      margin-top: 30px;
    }
    
    .quick-actions {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
    }
    
    @media (max-width: 599px) {
      .dashboard-grid {
        grid-template-columns: 1fr;
      }
    }
  `]
})
export class DashboardComponent implements OnInit {
  dashboardCards: DashboardCard[] = [
    {
      title: 'NAV.STUDENTS',
      subtitle: 'DASHBOARD.ACTIVE_STUDENTS',
      icon: 'people',
      value: '---',
      route: '/students',
      color: '#4caf50'
    },
    {
      title: 'NAV.COURSES',
      subtitle: 'DASHBOARD.COURSES_IN_PROGRESS',
      icon: 'school',
      value: '---',
      route: '/courses',
      color: '#2196f3'
    },
    {
      title: 'NAV.ENROLLMENTS',
      subtitle: 'DASHBOARD.ACTIVE_ENROLLMENTS',
      icon: 'how_to_reg',
      value: '---',
      route: '/enrollments',
      color: '#ff9800'
    },
    {
      title: 'NAV.ATTENDANCE',
      subtitle: 'DASHBOARD.CURRENT_ATTENDANCE',
      icon: 'fact_check',
      value: '---',
      route: '/attendance',
      color: '#9c27b0'
    }
  ];
  
  constructor(
    public authService: AuthService,
    private router: Router
  ) {}
  
  ngOnInit(): void {
    // Evita que perfil aluno/professor fique no dashboard administrativo genérico.
    const target = this.authService.getDefaultRouteByRole();
    if (target !== '/dashboard') {
      this.router.navigate([target]);
      return;
    }

    // Aqui você carregaria dados reais da API
    // Este é apenas um exemplo com dados fictícios
    this.loadDashboardData();
  }
  
  loadDashboardData(): void {
    // Simular carregamento de dados
    setTimeout(() => {
      this.dashboardCards = [
        {
          title: 'NAV.STUDENTS',
          subtitle: 'DASHBOARD.ACTIVE_STUDENTS',
          icon: 'people',
          value: '243',
          route: '/students',
          color: '#4caf50'
        },
        {
          title: 'NAV.COURSES',
          subtitle: 'DASHBOARD.COURSES_IN_PROGRESS',
          icon: 'school',
          value: '18',
          route: '/courses',
          color: '#2196f3'
        },
        {
          title: 'NAV.ENROLLMENTS',
          subtitle: 'DASHBOARD.ACTIVE_ENROLLMENTS',
          icon: 'how_to_reg',
          value: '312',
          route: '/enrollments',
          color: '#ff9800'
        },
        {
          title: 'NAV.ATTENDANCE',
          subtitle: 'DASHBOARD.CURRENT_ATTENDANCE',
          icon: 'fact_check',
          value: '87%',
          route: '/attendance',
          color: '#9c27b0'
        }
      ];
    }, 1000);
  }
}
