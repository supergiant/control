export const CLOUD_ACCOUNTS_MOCK = [{
  name: 'nick-DO-acc',
  provider: 'digitalocean',
  credentials: {
    accessToken: 'xxx',
    publicKey: 'xxx',
  },
}, {
  name: 'nick-aws-CA',
  provider: 'aws',
  credentials: { 'access_key': 'xxx', 'secret_key': 'xxx' },
}, {
  'name': 'nick-googl-CA', 'provider': 'gce', 'credentials': {
    auth_provider_x509_cert_url: 'https://www.googleapis.com/oauth2/v1/certs',
    auth_uri: 'https://accounts.google.com/o/oauth2/auth',
    client_email: 'test@test.iam.gserviceaccount.com',
    client_id: '1234',
    client_x509_cert_url: 'https://www.googleapis.com/robot/v1/metadata/x509/test-268%40ordinal-case-222023.iam.gserviceaccount.com',
    private_key: 'xxx',
    private_key_id: 'xxx',
    project_id: 'ordinal-case-test',
    token_uri: 'https://oauth2.googleapis.com/token',
    type: 'service_account',
  },
}];
