export interface Course {
  id?: number;
  programId?: number;
  name: string;
  category?: string;
  coverImage?: string;
  shortDescription?: string;
  detailedDescription?: string;
  difficultyLevel?: 'Beginner' | 'Intermediate' | 'Advanced';
  targetAudience?: string;
  prerequisites?: string;
  workload: number;
  maxStudents: number;
  duration?: number;
  weekDays?: string;
  startTime?: string;
  endTime?: string;
  startDate?: Date;
  endDate?: Date;
  locationId?: number;
  status: string;
  googleClassroomUrl?: string;
  teacherId?: number;
  createdAt?: string;
  updatedAt?: string;
}
