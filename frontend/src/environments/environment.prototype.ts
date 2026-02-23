// src/environments/environment.prototype.ts
export const environment = {
  production: false,
  apiUrl: '/api/v1',
  ssoApiUrl: 'http://localhost:8081',
  ssoIssuer: 'http://localhost:8081/realms/cepprs',
  ssoClientId: 'lar-cepprs-frontend',
  ssoRedirectUri: 'http://localhost:4201',
  prototype: true,
  version: '1.0.0-prototype',
  useMocks: true,
};
