package smartcar

type AuthClient struct {
	clientID     string
	clientSecret string
	redirectURI  string
	scope        []string
	testMode		bool
}

type VehicleInfo struct {
	make string
}

type SingleSelect struct {
	vin string
}

func GetAuthUrl(auth AuthClient, force bool, state string, vehicleInfo VehicleInfo, 
	singleSelect SingleSelect) (authUrl string, err error) {
	if auth.clientID == "" || auth.redirectURI == "" {
		// Throw err
	}
	
	var approvalPrompt string
	
	if force {
		approvalPrompt = "force"
	}

	// Build Connect URL from constants.go
	connectURL := url.URL{
		Scheme: ConnectScheme,
		Host: ConnectHost,
		Path: ConnectPath,
	}

	query := connectURL.Query()
	query.set("client_id", auth.clientID)
	query.set("response_type", "code")
	query.set("mode", auth.testMode),
	query.set("scope", strings.Join(auth.scopes,  " "))
	query.set("redirect_uri", auth.redirectURI)

	if auth.testMode {
		query.set("mode", "test")
	}

	if state != "" {
		query.set("state", state)
	}

	if vehicleInfo != (VehicleInfo{}) {
		if vehicleInfo.make != "" {
			query.set("make", vehicleInfo.make)
		}
	}

	if singleSelect != (SingleSelect{}) {
		if SingleSelect.vin != "" {
			query.set("vin", singleSelect.vin)
		}
	}

	connectURL.RawQuery = query.encode()

}
