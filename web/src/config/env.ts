// Environment configuration types and utilities
export interface EnvConfig {
  apiUrl: string;
  apiGraphqlUrl: string;
  environment: 'development' | 'staging' | 'production';
  enableDevTools: boolean;
  enableAnalytics: boolean;
  googleClientId?: string;
  appleClientId?: string;
}

// Load and validate environment variables
export const env: EnvConfig = {
  apiUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  apiGraphqlUrl:
    import.meta.env.VITE_API_GRAPHQL_URL || 'http://localhost:8080/graphql',
  environment:
    (import.meta.env.VITE_ENVIRONMENT as EnvConfig['environment']) ||
    'development',
  enableDevTools: import.meta.env.VITE_ENABLE_DEV_TOOLS === 'true',
  enableAnalytics: import.meta.env.VITE_ENABLE_ANALYTICS === 'true',
  googleClientId: import.meta.env.VITE_GOOGLE_CLIENT_ID,
  appleClientId: import.meta.env.VITE_APPLE_CLIENT_ID,
};

export const isDevelopment = env.environment === 'development';
export const isStaging = env.environment === 'staging';
export const isProduction = env.environment === 'production';

// Log configuration in development
if (isDevelopment) {
  console.log('Environment Configuration:', env);
}
