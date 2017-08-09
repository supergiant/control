package model

//CloudAccountSchema returns a default model and schema. For use with api.
func CloudAccountSchema() map[string]interface{} {
	return map[string]interface{}{
		"digitalocean": map[string]interface{}{
			// Default UI Object
			"model": map[string]interface{}{
				"name":     "",
				"provider": "digitalocean",
				"credentials": map[string]interface{}{
					"token": "",
				},
			},
			// UI Object rules
			"schema": map[string]interface{}{
				"properties": map[string]interface{}{
					"credentials": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"token": map[string]interface{}{
								"type":        "string",
								"description": "API Token",
							},
						},
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Provider Name",
					},
					"provider": map[string]interface{}{
						"type":        "string",
						"default":     "digitalocean",
						"description": "Provider",
						"widget":      "hidden",
					},
				},
			},
		},
		"openstack": map[string]interface{}{
			"model": map[string]interface{}{
				"name":     "",
				"provider": "openstack",
				"credentials": map[string]interface{}{
					"identity_endpoint": "",
					"username":          "",
					"password":          "",
					"tenant_id":         "",
					"domain_id":         "",
					"domain_name":       "",
				},
			},
			// UI Object rules
			"schema": map[string]interface{}{
				"properties": map[string]interface{}{
					"credentials": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"identity_endpoint": map[string]interface{}{
								"type":        "string",
								"description": "Identity Endpoint",
							},
							"username": map[string]interface{}{
								"type":        "string",
								"description": "User Name",
							},
							"password": map[string]interface{}{
								"type":        "string",
								"description": "Password",
							},
							"tenant_id": map[string]interface{}{
								"type":        "string",
								"description": "Tenant ID",
							},
							"domain_id": map[string]interface{}{
								"type":        "string",
								"description": "Domain ID",
							},
							"domain_name": map[string]interface{}{
								"type":        "string",
								"description": "Domain Name",
							},
						},
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Provider Name",
					},
					"provider": map[string]interface{}{
						"type":        "string",
						"default":     "openstack",
						"description": "Provider",
						"widget":      "hidden",
					},
				},
			},
		},
		"gce": map[string]interface{}{
			"model": map[string]interface{}{
				"name":     "",
				"provider": "gce",
				"credentials": map[string]interface{}{
					"type":                        "",
					"project_id":                  "",
					"private_key_id":              "",
					"private_key":                 "",
					"client_email":                "",
					"client_id":                   "",
					"auth_uri":                    "",
					"token_uri":                   "",
					"auth_provider_x509_cert_url": "",
					"client_x509_cert_url":        "",
				},
			},
			"schema": map[string]interface{}{
				"properties": map[string]interface{}{
					"credentials": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"type": map[string]interface{}{
								"type":        "string",
								"description": "Type",
							},
							"project_id": map[string]interface{}{
								"type":        "string",
								"description": "Project ID",
							},
							"private_key_id": map[string]interface{}{
								"type":        "string",
								"description": "Private Key ID",
							},
							"private_key": map[string]interface{}{
								"type":        "string",
								"description": "Private Key",
							},
							"client_email": map[string]interface{}{
								"type":        "string",
								"description": "Client Email",
							},
							"client_id": map[string]interface{}{
								"type":        "string",
								"description": "Client ID",
							},
							"auth_uri": map[string]interface{}{
								"type":        "string",
								"description": "Auth URI",
							},
							"token_uri": map[string]interface{}{
								"type":        "string",
								"description": "Token URI",
							},
							"auth_provider_x509_cert_url": map[string]interface{}{
								"type":        "string",
								"description": "auth_provider_x509_cert_url",
							},
							"client_x509_cert_url": map[string]interface{}{
								"type":        "string",
								"description": "client_x509_cert_url",
							},
						},
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Provider Name",
					},
					"provider": map[string]interface{}{
						"type":        "string",
						"default":     "gce",
						"description": "Provider",
						"widget":      "hidden",
					},
				},
			},
		},
		"packet": map[string]interface{}{
			"model": map[string]interface{}{
				"name":     "",
				"provider": "packet",
				"credentials": map[string]interface{}{
					"api_token": "",
				},
			},
			"schema": map[string]interface{}{
				"properties": map[string]interface{}{
					"credentials": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"api_token": map[string]interface{}{
								"type":        "string",
								"description": "API Token",
							},
						},
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Provider Name",
					},
					"provider": map[string]interface{}{
						"type":        "string",
						"default":     "packet",
						"description": "Provider",
						"widget":      "hidden",
					},
				},
			},
		},
		"aws": map[string]interface{}{
			"model": map[string]interface{}{
				"name":     "",
				"provider": "aws",
				"credentials": map[string]interface{}{
					"access_key": "",
					"secret_key": "",
				},
			},
			// UI Object rules
			"schema": map[string]interface{}{
				"properties": map[string]interface{}{
					"credentials": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"access_key": map[string]interface{}{
								"type":        "string",
								"description": "Access Key",
							},
							"secret_key": map[string]interface{}{
								"type":        "string",
								"description": "Secret Access Key",
							},
						},
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Provider Name",
					},
					"provider": map[string]interface{}{
						"type":        "string",
						"default":     "aws",
						"description": "Provider",
						"widget":      "hidden",
					},
				},
			},
		},
	}
}
