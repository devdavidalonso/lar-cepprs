// src/app/features/prototype-controls/prototype-controls.component.ts
import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatCardModule } from '@angular/material/card';
import { MatSelectModule } from '@angular/material/select';
import { MatDividerModule } from '@angular/material/divider';
import { MatSliderModule } from '@angular/material/slider';

import { PrototypeService } from '../../core/services/prototype/prototype.service';

@Component({
  selector: 'app-prototype-controls',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatButtonModule,
    MatIconModule,
    MatSlideToggleModule,
    MatCardModule,
    MatSelectModule,
    MatDividerModule,
    MatSliderModule
  ],
  template: `
    <div class="prototype-controls" [class.expanded]="isPanelExpanded">
      <button mat-fab color="warn" (click)="togglePanel()" [attr.aria-label]="isPanelExpanded ? 'Esconder controles' : 'Mostrar controles'">
        <mat-icon>{{ isPanelExpanded ? 'close' : 'construction' }}</mat-icon>
      </button>
      
      <mat-card class="control-panel" *ngIf="isPanelExpanded">
        <mat-card-header>
          <mat-card-title>Controles de Protótipo</mat-card-title>
          <mat-card-subtitle>Gerencie o modo de protótipo</mat-card-subtitle>
        </mat-card-header>
        
        <mat-card-content>
          <div class="control-section">
            <h3>Configurações Gerais</h3>
            <div class="control-item">
              <mat-slide-toggle 
                [checked]="isPrototypeMode" 
                (change)="togglePrototypeMode($event.checked)">
                Modo Protótipo
              </mat-slide-toggle>
              <p class="control-description">Ativa/desativa o uso de dados mockados</p>
            </div>
            
            <mat-divider></mat-divider>
            
            <div class="control-item">
              <h3>Simular Delay de Rede</h3>
              <mat-slider min="0" max="3000" step="100" [discrete]="true">
                <input matSliderThumb [(ngModel)]="networkDelay">
              </mat-slider>
              <p class="control-description">{{ networkDelay }}ms</p>
            </div>
          </div>
          
          <div class="control-section">
            <h3>Simular Perfil</h3>
            <mat-select [(ngModel)]="selectedUserProfile" placeholder="Selecione um perfil">
              <mat-option value="admin">Administrator</mat-option>
              <mat-option value="teacher">Teacher</mat-option>
              <mat-option value="student">Student</mat-option>
              <mat-option value="visitante">Visitante</mat-option>
            </mat-select>
            <p class="control-description">Simula diferentes níveis de acesso</p>
          </div>
          
          <div class="control-section">
            <h3>Simulação de Erros</h3>
            <button mat-raised-button color="accent" (click)="toggleErrorSimulation()">
              {{ simulateErrors ? 'Desativar Erros' : 'Ativar Erros' }}
            </button>
            <p class="control-description">Simula erros de API aleatórios</p>
          </div>
        </mat-card-content>
        
        <mat-card-actions>
          <button mat-button color="warn" (click)="resetSettings()">
            Resetar
          </button>
        </mat-card-actions>
      </mat-card>
    </div>
  `,
  styles: [`
    .prototype-controls {
      position: fixed;
      bottom: 20px;
      right: 20px;
      z-index: 1000;
    }
    
    .control-panel {
      position: absolute;
      bottom: 70px;
      right: 0;
      width: 280px;
      max-width: 90vw;
      max-height: 80vh;
      overflow-y: auto;
      background-color: white;
      border-radius: 4px;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
      transition: all 0.3s ease;
    }
    
    .control-section {
      margin-bottom: 16px;
    }
    
    .control-item {
      margin-bottom: 12px;
    }
    
    .control-description {
      font-size: 12px;
      color: #666;
      margin-top: 4px;
      margin-bottom: 8px;
    }
    
    h3 {
      margin: 8px 0;
      font-size: 14px;
      font-weight: 500;
    }
    
    mat-divider {
      margin: 12px 0;
    }
    
    mat-slider {
      width: 100%;
    }
    
    mat-select {
      width: 100%;
    }
  `]
})
export class PrototypeControlsComponent implements OnInit {
  isPanelExpanded = false;
  isPrototypeMode = false;
  networkDelay = 500;
  selectedUserProfile = 'visitante';
  simulateErrors = false;
  
  constructor(private prototypeService: PrototypeService) {}
  
  ngOnInit(): void {
    // Carregar o estado atual do modo de protótipo
    this.isPrototypeMode = this.prototypeService.isPrototypeEnabled();
    
    // Carregar outras configurações do localStorage, se existirem
    const networkDelay = localStorage.getItem('prototype_network_delay');
    if (networkDelay) {
      this.networkDelay = parseInt(networkDelay, 10);
    }
    
    const userProfile = localStorage.getItem('prototype_user_profile');
    if (userProfile) {
      this.selectedUserProfile = userProfile;
    }
    
    const simulateErrors = localStorage.getItem('prototype_simulate_errors');
    if (simulateErrors) {
      this.simulateErrors = simulateErrors === 'true';
    }
  }
  
  togglePanel(): void {
    this.isPanelExpanded = !this.isPanelExpanded;
  }
  
  togglePrototypeMode(enabled: boolean): void {
    this.isPrototypeMode = enabled;
    
    if (enabled) {
      this.prototypeService.enablePrototypeMode();
    } else {
      this.prototypeService.disablePrototypeMode();
    }
    
    // Recarrega a página para aplicar as alterações
    // (em uma implementação real, você provavelmente usaria eventos para notificar os componentes)
    window.location.reload();
  }
  
  toggleErrorSimulation(): void {
    this.simulateErrors = !this.simulateErrors;
    localStorage.setItem('prototype_simulate_errors', this.simulateErrors.toString());
  }
  
  resetSettings(): void {
    // Resetar todas as configurações
    this.networkDelay = 500;
    this.selectedUserProfile = 'visitante';
    this.simulateErrors = false;
    
    // Atualizar localStorage
    localStorage.setItem('prototype_network_delay', '500');
    localStorage.setItem('prototype_user_profile', 'visitante');
    localStorage.setItem('prototype_simulate_errors', 'false');
    
    // Mantem o estado do modo protótipo
  }
}