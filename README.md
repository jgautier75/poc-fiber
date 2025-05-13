# poc-fiber

## Go

Changing go module go version

```bash
go mod edit -go=1.24.2
```

## Application

Application relies on fiber (expressjs like written in go), postgreSQL for persistance, 

https://localhost:8101/poc-fiber/home

Once authenticated, open [Bruno](https://www.usebruno.com/)  collection (docs directory) for interacting with REST API

## Docker & Telemetry Stack

| Service             | Version | Port           | Description                                                                 |
|---------------------|---------|----------------|-----------------------------------------------------------------------------|
| mailpit             | 1.24    | 1025 & 8025    | Smtp mock server (smtp: 1025, 8025 for web app)Spring app storage           |
| postgreSQL          | 17.4    | 5432           | Application storage                                                         |
| postgreSQL          | 17.4    | 5433           | Authentik storage                                                           |
| redis               | alpine  | 66379          |                                                                             |
| authentik           | 2025.4.0| 9000           | Authentik server oidc provider                                              | 
| authentik           | 2025.4.0| -              | Authentik worker (scheduled tasks)                                          |


Telemetry stack relies on grafana (loki & tempo)

To start telemetry stack, execut docker/run _lgtm.sh script
Once started, open [grafana](http://localhost:3000)

## Authentik

Authentik is a OpenId provider written in python and relying by default on postgreSQL & redis

Starting with authentik oidc provider, applications only receive an access token. To receive a refresh token, both applications and authentik must be configured to request the offline_access scope

Select 'offline_access' in go authentik provider
Append 'offline_access' in oauth2.Config

Default login: akadmin
Default password: Authentik01234567890!

Open Start url to : http://localhost:9000/if/flow/initial-setup/

If, for any reason,you cannot access anymore to admin interface, run the following command to create a recovery token:

```bash
docker-compose -f docker-compose.yml run --rm server create_recovery_key 10 akadmin
```

Output:

```bash
recovery/use-token/YCk1Xednn1Y3YQy84CyfaKHBsxSOE7gXazB4KqNCDWyDT9c0uhs8HaGO4li7/
```

Open url with recovery url above:

```
http://localhost:9000/recovery/use-token/YCk1Xednn1Y3YQy84CyfaKHBsxSOE7gXazB4KqNCDWyDT9c0uhs8HaGO4li7/
```

## Configuration

Main configuration is store in config.config.yml file. To override values, provide environment variables prefixed with "EV_".

Example: To override the "app.name" configuration in config.yml file, provider an EV_APP.NAME environment variable.

# Ansible 

ansible-vault encrypt secrets_file.enc

ansible-playbook -i inventory/hosts.ini docker-setup.yml --connection=local --ask-vault-pass