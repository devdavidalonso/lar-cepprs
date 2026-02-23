export interface Address {
  id?: number;
  cep: string;
  street: string;
  number: string;
  complement?: string;
  neighborhood: string;
  city: string;
  state: string;
}

export interface UserContact {
  id?: number;
  name: string;
  email?: string;
  cpf?: string;
  phone: string;
  relationship: string;
  canPickup?: boolean;
  receiveNotifications?: boolean;
  authorizeActivities?: boolean;
}

export interface Teacher {
  id?: number;
  name: string;
  email: string;
  cpf?: string;
  phone?: string;
  birthDate?: string;
  specialization?: string;
  bio?: string;
  linkedinUrl?: string; // from Phase 2 database changes
  programIds?: number[];
  active?: boolean;
  address?: Address;
  userContacts?: UserContact[]; // Emergency contacts
  createdAt?: string;
  updatedAt?: string;
}
