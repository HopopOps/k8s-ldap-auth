package types

import (
	"encoding/json"
	"os"

	client "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

type ExecCredential struct {
	Kind       string                       `json:"kind"`
	APIVersion string                       `json:"apiVersion,omitempty"`
	Spec       client.ExecCredentialSpec    `json:"spec"`
	Status     *client.ExecCredentialStatus `json:"status,omitempty"`
}

func NewExecCredential(APIVersion string) ExecCredential {
	ec := ExecCredential{
		Kind: "ExecCredential",
	}

	if APIVersion != "" {
		ec.APIVersion = APIVersion
	} else {
		ec.APIVersion = os.Getenv("AUTH_API_VERSION")
	}

	return ec
}

func (ec *ExecCredential) Marshal(APIVersion string) ([]byte, error) {
	ec.Kind = "ExecCredential"

	if APIVersion != "" {
		ec.APIVersion = APIVersion
	} else {
		ec.APIVersion = os.Getenv("AUTH_API_VERSION")
	}

	data, err := json.Marshal(ec)
	if err != nil {
		return nil, err
	}

	return data, nil
}
