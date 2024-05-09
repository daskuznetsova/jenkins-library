import groovy.transform.Field

@Field String STEP_NAME = getClass().getName()
@Field String METADATA_FILE = 'metadata/dasterExecuteScan.yaml'

void call(Map parameters = [:]) {
    List credentials = [
    [type: 'usernamePassword', id: 'oAuthCredentialsId', env: ['PIPER_clientId', 'PIPER_clientSecret']],
    [type: 'token', id: 'dasterTokenCredentialsId', env: ['PIPER_dasterToken']],
    [type: 'token', id: 'userCredentialsId', env: ['PIPER_user']],
    [type: 'usernamePassword', id: 'targetAuthCredentialsId', env: ['PIPER_dasterTargetUser', 'PIPER_dasterTargetPassword']]
    ]
    piperExecuteBin(parameters, STEP_NAME, METADATA_FILE, credentials)
}
