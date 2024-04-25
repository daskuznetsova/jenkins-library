import groovy.transform.Field

@Field String STEP_NAME = getClass().getName()
@Field String METADATA_FILE = 'metadata/dasterExecuteScan.yaml'

void call(Map parameters = [:]) {
    List credentials = [
    [type: 'usernamePassword', id: 'oAuthCredentialsId', env: ['PIPER_clientId', 'PIPER_clientSecret']]
    ]
    piperExecuteBin(parameters, STEP_NAME, METADATA_FILE, credentials)
}
