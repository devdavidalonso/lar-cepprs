import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '@environments/environment';
import { Program } from '@core/models/program.model';

@Injectable({
  providedIn: 'root'
})
export class ProgramService {
  private readonly API_URL = `${environment.apiUrl}/programs`;

  constructor(private http: HttpClient) {}

  listPrograms(activeOnly: boolean = true): Observable<Program[]> {
    const params = new HttpParams().set('activeOnly', String(activeOnly));
    return this.http.get<Program[]>(this.API_URL, { params });
  }
}
