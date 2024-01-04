## How tests the changes
1. `make docker`
1. ```shell
   docker tag europe-west4-docker.pkg.dev/sbat-gcr-develop/sapig-docker-artifact/securebanking/secureopenbanking-uk-fidc-initializer:latest europe-west4-docker.pkg.dev/sbat-gcr-release/sapig-docker-artifact/securebanking/secureopenbanking-uk-fidc-initializer:latest
   ```
1. ```shell
   docker push !$[+TAB]
   ```
   Or
   ```shell
   docker push europe-west4-docker.pkg.dev/sbat-gcr-release/sapig-docker-artifact/securebanking/secureopenbanking-uk-fidc-initializer:latest
   ```
1. Change kubernetes context to `sbat-master-dev`
1. Set the namespace to `ig`
1. Run `docker delete pod ig-xxxxxxxx`

>The new ig pod will run the latest image pushed in the step 3.

>Check your changes on the [platform](https://iam.dev.forgerock.financial/platform)