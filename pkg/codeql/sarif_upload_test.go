package codeql

//func TestWaitSarifUploaded(t *testing.T) {
//	t.Parallel()
//	config := codeqlExecuteScanOptions{SarifCheckRetryInterval: 1, SarifCheckMaxRetries: 5}
//	t.Run("Fast complete upload", func(t *testing.T) {
//		codeqlScanAuditMock := CodeqlSarifUploaderMock{counter: 0}
//		timerStart := time.Now()
//		err := waitSarifUploaded(&config, &codeqlScanAuditMock)
//		assert.Less(t, time.Now().Sub(timerStart), time.Second)
//		assert.NoError(t, err)
//	})
//	t.Run("Long completed upload", func(t *testing.T) {
//		codeqlScanAuditMock := CodeqlSarifUploaderMock{counter: 2}
//		timerStart := time.Now()
//		err := waitSarifUploaded(&config, &codeqlScanAuditMock)
//		assert.GreaterOrEqual(t, time.Now().Sub(timerStart), time.Second*2)
//		assert.NoError(t, err)
//	})
//	t.Run("Failed upload", func(t *testing.T) {
//		codeqlScanAuditMock := CodeqlSarifUploaderMock{counter: -1}
//		err := waitSarifUploaded(&config, &codeqlScanAuditMock)
//		assert.Error(t, err)
//		assert.ErrorContains(t, err, "failed to upload sarif file")
//	})
//	t.Run("Error while checking sarif uploading", func(t *testing.T) {
//		codeqlScanAuditErrorMock := CodeqlSarifUploaderErrorMock{counter: -1}
//		err := waitSarifUploaded(&config, &codeqlScanAuditErrorMock)
//		assert.Error(t, err)
//		assert.ErrorContains(t, err, "test error")
//	})
//	t.Run("Completed upload after getting errors from server", func(t *testing.T) {
//		codeqlScanAuditErrorMock := CodeqlSarifUploaderErrorMock{counter: 3}
//		err := waitSarifUploaded(&config, &codeqlScanAuditErrorMock)
//		assert.NoError(t, err)
//	})
//	t.Run("Max retries reached", func(t *testing.T) {
//		codeqlScanAuditErrorMock := CodeqlSarifUploaderErrorMock{counter: 6}
//		err := waitSarifUploaded(&config, &codeqlScanAuditErrorMock)
//		assert.Error(t, err)
//		assert.ErrorContains(t, err, "max retries reached")
//	})
//}
