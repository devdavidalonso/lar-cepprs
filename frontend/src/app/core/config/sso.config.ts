import { AuthConfig } from 'angular-oauth2-oidc';

export const authCodeFlowConfig: AuthConfig = {
  // URL do Keycloak (via Proxy para evitar CORS)
  issuer: window.location.origin + '/realms/lar-sso',

  // URL do CECO para onde voltar após login
  redirectUri: window.location.origin,

  // ID do Cliente
  clientId: 'lar-cepprs-frontend',

  responseType: 'code',
  scope: 'openid profile email offline_access',
  showDebugInformation: true,
  requireHttps: false,
  strictDiscoveryDocumentValidation: false,
  skipIssuerCheck: true,
  disableAtHashCheck: true
};
