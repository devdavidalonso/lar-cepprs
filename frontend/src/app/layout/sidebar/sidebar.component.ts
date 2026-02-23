import { Component, Input, ChangeDetectionStrategy, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatSidenav } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { Observable, Subject } from 'rxjs';
import { map, takeUntil } from 'rxjs/operators';

import { AuthService } from '../../core/services/auth.service';

interface MenuItem {
  text: string;
  icon: string;
  route: string;
  roles?: string[];
  dividerBefore?: boolean;
}

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatListModule,
    MatIconModule,
    MatDividerModule,
    MatTooltipModule,
    TranslateModule,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="sidebar-container">

      <!-- Cabeçalho / Logo -->
      <div class="sidebar-header">
        <img src="assets/images/cecor-logo.png" alt="LAR-CEPPRS" class="logo-cecor" />
      </div>

      <!-- Usuário Logado -->
      <div class="user-info" *ngIf="currentUser$ | async as user">
        <div class="user-avatar">{{ getInitials(user.name) }}</div>
        <div class="user-details">
          <span class="user-name">{{ user.name }}</span>
          <span class="user-role">{{ getRoleLabel(user.roles) }}</span>
        </div>
      </div>

      <mat-divider></mat-divider>

      <!-- Menu de Navegação -->
      <nav class="sidebar-nav">
        <mat-nav-list>
          <ng-container *ngFor="let item of menuItems">
            <mat-divider *ngIf="item.dividerBefore" class="section-divider"></mat-divider>
            <ng-container *ngIf="hasRequiredRole(item.roles) | async">
              <a
                mat-list-item
                [routerLink]="item.route"
                routerLinkActive="active-link"
                [routerLinkActiveOptions]="{exact: item.route === '/dashboard'}"
                (click)="closeIfMobile()"
                class="nav-item"
                [matTooltip]="item.text | translate"
                matTooltipPosition="right">
                <mat-icon matListItemIcon class="nav-icon">{{ item.icon }}</mat-icon>
                <span matListItemTitle class="nav-label">{{ item.text | translate }}</span>
              </a>
            </ng-container>
          </ng-container>
        </mat-nav-list>
      </nav>

      <!-- Rodapé Sidebar -->
      <div class="sidebar-footer">
        <img src="assets/images/paulo-rossi.png" alt="GECS - Paulo Rossi" class="footer-paulo-rossi" />
      </div>

    </div>
  `,
  styles: [`
    :host {
      display: block;
      height: 100%;
    }

    .sidebar-container {
      display: flex;
      flex-direction: column;
      height: 100%;
      background: #ffffff;
      border-right: 1px solid rgba(0, 106, 172, 0.12);
      overflow: hidden;
    }

    /* ---- Cabeçalho ---- */
    /* ---- Cabeçalho ---- */
    .sidebar-header {
      background: linear-gradient(135deg, #006aac 0%, #0083c0 100%);
      padding: 16px 12px;
      flex-shrink: 0;
      display: flex;
      justify-content: center;
      align-items: center;
    }

    .logo-cecor {
      width: 80px;
      height: 80px;
      object-fit: contain;
    }

    /* ---- Usuário ---- */
    .user-info {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 14px 16px;
      background: #f8fbff;
      flex-shrink: 0;
    }

    .user-avatar {
      width: 36px;
      height: 36px;
      border-radius: 50%;
      background: linear-gradient(135deg, #006aac, #0083c0);
      color: #ffffff;
      display: flex;
      align-items: center;
      justify-content: center;
      font-family: 'Manrope', sans-serif;
      font-weight: 700;
      font-size: 14px;
      flex-shrink: 0;
    }

    .user-details {
      display: flex;
      flex-direction: column;
      min-width: 0;
    }

    .user-name {
      font-family: 'Manrope', sans-serif;
      font-weight: 600;
      font-size: 13px;
      color: #302424;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .user-role {
      font-family: 'Manrope', sans-serif;
      font-weight: 400;
      font-size: 11px;
      color: #0083c0;
    }

    /* ---- Navegação ---- */
    .sidebar-nav {
      flex: 1;
      overflow-y: auto;
      overflow-x: hidden;
      padding: 8px 0;
    }

    .section-divider {
      margin: 8px 16px !important;
      border-color: rgba(0, 106, 172, 0.1) !important;
    }

    mat-nav-list {
      padding-top: 0 !important;
    }

    .nav-item {
      height: 44px !important;
      margin: 2px 8px !important;
      border-radius: 8px !important;
      transition: background-color 0.15s ease, color 0.15s ease !important;

      &:hover {
        background-color: rgba(0, 106, 172, 0.06) !important;
      }

      &.active-link {
        background-color: rgba(0, 106, 172, 0.12) !important;
        border-left: 3px solid #006aac;
        padding-left: 13px !important;

        .nav-icon {
          color: #006aac !important;
        }

        .nav-label {
          color: #006aac !important;
          font-weight: 700 !important;
        }
      }
    }

    .nav-icon {
      color: #8b8d94 !important;
      font-size: 20px !important;
      width: 20px !important;
      height: 20px !important;
      transition: color 0.15s ease;
    }

    .nav-label {
      font-family: 'Manrope', sans-serif !important;
      font-size: 14px !important;
      font-weight: 500 !important;
      color: #302424 !important;
      transition: color 0.15s ease, font-weight 0.15s ease;
    }

    /* ---- Rodapé ---- */
    .sidebar-footer {
      padding: 16px;
      border-top: 1px solid rgba(0, 0, 0, 0.06);
      flex-shrink: 0;
      display: flex;
      justify-content: center;
      align-items: center;
      background: #f8fbff;
    }

    .footer-paulo-rossi {
      width: 90px;
      height: auto;
      object-fit: contain;
      display: block;
    }
  `]
})
export class SidebarComponent implements OnInit, OnDestroy {
  @Input() sidenav!: MatSidenav;

  private destroy$ = new Subject<void>();

  currentUser$: Observable<any>;

  readonly menuItems: MenuItem[] = [
    {
      text: 'NAV.DASHBOARD',
      icon: 'dashboard',
      route: '/dashboard',
      roles: ['admin', 'administrador', 'gestor'],
    },
    {
      text: 'NAV.DASHBOARD',
      icon: 'dashboard',
      route: '/teacher/dashboard',
      roles: ['professor'],
    },
    {
      text: 'NAV.DASHBOARD',
      icon: 'dashboard',
      route: '/student/dashboard',
      roles: ['aluno', 'responsavel', 'responsável'],
    },
    {
      text: 'NAV.STUDENTS',
      icon: 'people',
      route: '/students',
      roles: ['admin', 'administrador'],
      dividerBefore: true,
    },
    {
      text: 'NAV.STUDENTS_NEW',
      icon: 'person_add',
      route: '/students/new',
      roles: ['admin', 'administrador'],
    },
    {
      text: 'NAV.COURSES',
      icon: 'school',
      route: '/courses',
      roles: ['admin', 'administrador'],
    },
    {
      text: 'NAV.ENROLLMENTS',
      icon: 'how_to_reg',
      route: '/enrollments',
      roles: ['admin', 'administrador'],
    },
    {
      text: 'NAV.ATTENDANCE',
      icon: 'fact_check',
      route: '/attendance',
      roles: ['admin', 'administrador'],
    },
    {
      text: 'NAV.REPORTS',
      icon: 'assessment',
      route: '/reports',
      roles: ['admin', 'administrador'],
      dividerBefore: true,
    },
    {
      text: 'NAV.INTERVIEWS',
      icon: 'question_answer',
      route: '/interviews',
      roles: ['admin', 'administrador'],
    },
    {
      text: 'NAV.INTERVIEWS_DASHBOARD',
      icon: 'analytics',
      route: '/interviews/dashboard',
      roles: ['admin', 'administrador'],
    },
    {
      text: 'NAV.INTERVIEWS_REPORTS',
      icon: 'bar_chart',
      route: '/interviews/reports',
      roles: ['admin', 'administrador'],
    },
    {
      text: 'NAV.TEACHERS',
      icon: 'supervisor_account',
      route: '/teachers',
      roles: ['admin', 'administrador'],
      dividerBefore: true,
    },
    {
      text: 'NAV.ADMINISTRATION',
      icon: 'admin_panel_settings',
      route: '/administration',
      roles: ['admin', 'administrador'],
    },
    {
      text: 'NAV.LOCATIONS',
      icon: 'meeting_room',
      route: '/locations',
      roles: ['admin', 'administrador'],
    },
    {
      text: 'NAV.VOLUNTEERING',
      icon: 'volunteer_activism',
      route: '/volunteering',
      roles: ['admin', 'administrador'],
    },
  ];

  constructor(private authService: AuthService) {
    this.currentUser$ = this.authService.currentUser$;
  }

  ngOnInit(): void {}

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  hasRequiredRole(requiredRoles?: string[]): Observable<boolean> {
    if (!requiredRoles || requiredRoles.length === 0) {
      return new Observable(observer => observer.next(true));
    }
    return this.authService.currentUser$.pipe(
      map(user => {
        if (!user?.roles) return false;
        return requiredRoles.some(role => user.roles!.includes(role));
      }),
    );
  }

  closeIfMobile(): void {
    if (window.innerWidth < 960) {
      this.sidenav.close();
    }
  }

  getInitials(name: string): string {
    if (!name) return '?';
    return name
      .split(' ')
      .slice(0, 2)
      .map(n => n[0])
      .join('')
      .toUpperCase();
  }

  getRoleLabel(roles?: string[]): string {
    if (!roles?.length) return '';
    if (roles.includes('admin') || roles.includes('administrador')) return 'Administrador';
    if (roles.includes('professor')) return 'Professor';
    if (roles.includes('aluno')) return 'Aluno';
    return roles[0];
  }
}
