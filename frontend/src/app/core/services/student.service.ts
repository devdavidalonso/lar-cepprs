// src/app/core/services/student.service.ts

import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '@environments/environment';
import { Student, Guardian, CreateStudentRequest, UpdateStudentRequest } from '@core/models/student.model';
import { PaginatedResponse } from '@core/models/common.model';

export interface StudentFilters {
  name?: string;
  email?: string;
  cpf?: string;
  status?: string;
  courseId?: number;
  programId?: number;
  ageMin?: number;
  ageMax?: number;
}

@Injectable({
  providedIn: 'root'
})
export class StudentService {
  private readonly API_URL = `${environment.apiUrl}/students`;

  constructor(private http: HttpClient) { }

  getStudents(page: number = 1, pageSize: number = 20, filters?: StudentFilters): Observable<PaginatedResponse<Student>> {
    let params = new HttpParams()
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());

    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          params = params.set(key, value.toString());
        }
      });
    }

    return this.http.get<PaginatedResponse<Student>>(this.API_URL, { params });
  }

  getStudent(id: number): Observable<Student> {
    return this.http.get<Student>(`${this.API_URL}/${id}`);
  }

  createStudent(student: CreateStudentRequest): Observable<Student> {
    return this.http.post<Student>(this.API_URL, student);
  }

  updateStudent(id: number, student: UpdateStudentRequest): Observable<Student> {
    return this.http.put<Student>(`${this.API_URL}/${id}`, student);
  }

  deleteStudent(id: number): Observable<void> {
    return this.http.delete<void>(`${this.API_URL}/${id}`);
  }

  // Métodos para Responsáveis (Guardians)
  getGuardians(studentId: number): Observable<Guardian[]> {
    return this.http.get<Guardian[]>(`${this.API_URL}/${studentId}/guardians`);
  }

  addGuardian(studentId: number, guardian: Partial<Guardian>): Observable<Guardian> {
    return this.http.post<Guardian>(`${this.API_URL}/${studentId}/guardians`, guardian);
  }

  updateGuardian(guardianId: number, guardian: Partial<Guardian>): Observable<Guardian> {
    return this.http.put<Guardian>(`${environment.apiUrl}/guardians/${guardianId}`, guardian);
  }

  deleteGuardian(guardianId: number): Observable<void> {
    return this.http.delete<void>(`${environment.apiUrl}/guardians/${guardianId}`);
  }

  // Métodos para Documentos
  getDocuments(studentId: number): Observable<any[]> {
    return this.http.get<any[]>(`${this.API_URL}/${studentId}/documents`);
  }

  // Métodos para Notas
  getNotes(studentId: number, includeConfidential: boolean = false): Observable<any[]> {
    let params = new HttpParams()
      .set('includeConfidential', includeConfidential.toString());
    
    return this.http.get<any[]>(`${this.API_URL}/${studentId}/notes`, { params });
  }
}
