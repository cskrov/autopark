package main

type AuthRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type ApiError struct {
	ResultCode   string `json:"resultCode"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type AuthResponse struct {
	ApiError
	Token string `json:"token"`
}

type ProfileResponse struct {
	ApiError
	Profile struct {
		KioskMode             string `json:"kioskMode"`
		WelcomeMessage        string `json:"welcomeMessage"`
		KioskNotificationMode string `json:"kioskNotificationMode"`
		TermsAndConditionsId  string `json:"termsAndConditionsId"`
		Agents                []struct {
			Id            string `json:"id"`
			Name          string `json:"name"`
			KioskProducts []struct {
				Id                  string   `json:"id"`
				ProductName         string   `json:"productName"`
				QualificationFields []string `json:"qualificationFields"`
			} `json:"kioskProducts"`
		} `json:"agents"`
	} `json:"profile"`
}

type PlateResponse struct {
	ApiError
	Message string `json:"message"`
}

type StatusResponse struct {
	ApiError
	StartTime          string `json:"startTime"`
	EndTime            string `json:"endTime"`
	ExpiryTime         string `json:"expiryTime"`
	EndOfWaitingPeriod string `json:"endOfWaitingPeriod"`
}

type RegisterRequest struct {
	PlateNumber                  string   `json:"plateNumber"`
	KioskAgentId                 string   `json:"kioskAgentId"`
	ProductId                    string   `json:"productId"`
	TermsAccepted                bool     `json:"termsAccepted"`
	AcceptedTermsAndConditionsId *string  `json:"acceptedTermsAndConditionsId"`
	QualificationForm            []string `json:"qualificationForm"`
}

type RegisterResponse struct {
	ApiError
	Parking struct {
		Id                        string      `json:"id"`
		LicencePlateNumber        string      `json:"licencePlateNumber"`
		NotifyPhone               string      `json:"notifyPhone"`
		PayerUserId               string      `json:"payerUserId"`
		PayerMsisdn               string      `json:"payerMsisdn"`
		AgentId                   string      `json:"agentId"`
		AgentName                 string      `json:"agentName"`
		ProductId                 string      `json:"productId"`
		OperatorId                string      `json:"operatorId"`
		ZoneId                    string      `json:"zoneId"`
		OrgNo                     string      `json:"orgNo"`
		ZoneName                  string      `json:"zoneName"`
		StartTime                 string      `json:"startTime"`
		EndTime                   string      `json:"endTime"`
		ExpiryTime                string      `json:"expiryTime"`
		ValidNow                  bool        `json:"validNow"`
		Active                    bool        `json:"active"`
		DiscountAmount            string      `json:"discountAmount"`
		ParkingFee                string      `json:"parkingFee"`
		VatAmount                 string      `json:"vatAmount"`
		AgentCommission           string      `json:"agentCommission"`
		PaymentMethod             string      `json:"paymentMethod"`
		PaymentReference          string      `json:"paymentReference"`
		Note                      string      `json:"note"`
		RefundDate                string      `json:"refundDate"`
		CanEditLicencePlateNumber bool        `json:"canEditLicencePlateNumber"`
		Gates                     interface{} `json:"gates"`
	} `json:"parking"`
}
