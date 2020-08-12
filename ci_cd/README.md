# CI/CD for Git Import

## Setting The Concourse *PLATFORM* Team

1. This pipeline belongs in the `platform` team. If no `platform` concourse team exists run the following:
    - `fly -t prod login --concourse-url https://cicd.zapos.io -u zapatabot -p <PASSWORD IN LASTPASS> -n platform`
    - `fly -t prod sync`
    - `fly -t prod set-team -n platform --local-user=zapatabot`
    - `fly -t prod logout`

## Setting The Pipeline

1. Install the LastPass CLI:
        - `brew install lastpass-cli`

1. The pipeline contains secrets that we do not keep on github, instead we utilize LastPass to store all secrets. Therefore you will need the lpass CLI to have access to them. Run the following to obtain the secrets for the pipeline:
        - `lpass login "<YOUR ZAPATA EMAIL>"`
        - `lpass show zapata-ci-creds.yaml --notes > creds.yaml`

1. Once the secrets are obtained, the pipeline is then set with the following:
       
       - `fly -t prod login --concourse-url https://cicd.zapos.io -u zapatabot -p <ZAPATABOT PASSWORD> -n platform`
       - `fly -t prod set-pipeline -p git-import -c pipeline.yaml --load-vars-from=creds.yaml`

1. In order to enable all zapatistas to see our awesome pipeline we can expose the status and the configuration WITHOUT exposing the logs/secrets. Thus run the following:
    - `fly -t prod expose-pipeline --pipeline git-import`

