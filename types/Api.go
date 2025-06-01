package types

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Error   error  `json:"error,omitempty"`
}

type ApiAppResponse struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Icon        string            `json:"icon"`
	InstalledAt string            `json:"installedAt"`
	Parts       map[string]string `json:"parts"`
}

type ApiAppInstallRequestBody struct {
	File string `json:"file"`
}
