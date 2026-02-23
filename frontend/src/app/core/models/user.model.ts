/**
 * UserProfile - Catalog of user profiles
 */
export interface UserProfile {
  id: number;
  name: 'administrator' | 'teacher' | 'student' | 'admin' | 'professor';
  description: string;
}

/**
 * User - Updated with profileId as FK to user_profiles
 */
export interface User {
  id: number;
  name: string;
  email: string;
  roles: string[];
  profileId: number; // ✅ Changed from 'profile' (string) to 'profileId' (number)
  profile?: UserProfile; // ✅ New - relationship with UserProfile
  locale?: string; // ✅ New - User locale preference
}

/**
 * User profile ID constants
 */
export const USER_PROFILES = {
  ADMIN: 1,
  PROFESSOR: 2,
  STUDENT: 3,
} as const;

/**
 * Get profile name by ID
 */
export function getProfileName(profileId: number): string {
  switch (profileId) {
    case USER_PROFILES.ADMIN:
      return 'administrator';
    case USER_PROFILES.PROFESSOR:
      return 'teacher';
    case USER_PROFILES.STUDENT:
      return 'student';
    default:
      return 'unknown';
  }
}
