export interface EducationalCenter {
  id: number;
  name: string;
  code: string;
  isActive: boolean;
}

export interface Program {
  id: number;
  centerId: number;
  code: string;
  name: string;
  isActive: boolean;
  center?: EducationalCenter;
}
