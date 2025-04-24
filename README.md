# poc-fiber

Starting with authentik 2024.2, applications only receive an access token. To receive a refresh token, both applications and authentik must be configured to request the offline_access scope

Select 'offline_access' in go authentik provider
Append 'offline_access' in oauth2.Config

Changing go module go version

```bash
go mod edit -go=1.24.2
```