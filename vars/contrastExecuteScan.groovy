import groovy.transform.Field

@Field String STEP_NAME = getClass().getName()
@Field String METADATA_FILE = 'metadata/contrastExecuteScan.yaml'

void call(Map parameters = [:]) {
    List credentials = [
    [type: 'token', id: 'userApiKeyCredentialsId', env: ['PIPER_userApiKey']],
    [type: 'token', id: 'serviceKeyCredentialsId', env: ['PIPER_serviceKey']]
    ]
    piperExecuteBin(parameters, STEP_NAME, METADATA_FILE, credentials)
}
