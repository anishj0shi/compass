package testkit

import "fmt"

type queryProvider struct{}

func (qp queryProvider) generateToken(id string) string {
	return fmt.Sprintf(`mutation {
	result: generateApplicationToken(appID: "%s") {
		token
	}
}`, id)
}

func (qp queryProvider) configuration() string {
	return fmt.Sprintf(`query {
	result: configuration() {
		%s
	}
}`, configurationResult())
}

func (qp queryProvider) generateCert(csr string) string {
	return fmt.Sprintf(`mutation {
	result: signCertificateSigningRequest(csr: "%s") {
		%s
	}
}`, csr, certificationResult())
}

func (qp queryProvider) revokeCert() string {
	return fmt.Sprintf(`mutation {
	result: revokeCertificate() {
		boolean
	}
}`)
}

func configurationResult() string {
	return `token { token }
	certificateSigningRequestInfo { subject keyAlgorithm }
	managementPlaneInfo { directorURL }`
}

func certificationResult() string {
	return `certificateChain
			caCertificate
			clientCertificate`
}