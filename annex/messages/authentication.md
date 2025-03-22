# Authentication Messages
OO -> PO: Authorization Request
- client_id
- client_secret

PO -> Keycloak: Credential verification
- client_id
- client_secrect
- grant_type

Keycloak -> PO: Issue Token
- access_token
- expires_in 

PO -> OO: Authorization Response
- access_token
- expires_at
