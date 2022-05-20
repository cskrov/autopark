package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var (
	uFlag  = flag.String("u", "", "Username.")
	pFlag  = flag.String("p", "", "Password.")
	rFlag  = flag.String("r", "", "Reg. nr. of the car.")
	eFlag  = flag.Bool("e", false, "If the parking should be extended, ie. 72 hours, instead of 6 hours.")
	vFalg  = flag.Bool("v", false, "Verbose output.")
	vvFalg = flag.Bool("vv", false, "Extra verbose output. Adds network request logging.")
)

func main() {
	flag.Parse()

	username := *uFlag
	password := *pFlag
	regNr := *rFlag
	extended := *eFlag
	verbose := *vFalg || *vvFalg

	if username == "" || password == "" || regNr == "" {
		log.Fatal("Missing required arguments.")
	}

	if verbose {
		log.Println("Username:", username)
		log.Println("Password:", password)
		log.Println("Reg. nr.:", regNr)
		log.Println("Extended:", extended)
		log.Println("Verbose:", verbose)
	}

	httpClient := NewClient("https://pservice-permit.giantleap.no", "/api/kiosk-v2")

	// Log in
	authRes := AuthResponse{}
	authReq := AuthRequest{
		UserName: username,
		Password: password,
	}
	if err := httpClient.Post("/authenticate.json", &authRes, &authReq); err != nil {
		LogFatalError("Failed to login", err)
	}
	if authRes.ResultCode != "SUCCESS" || authRes.Token == "" {
		LogFatalApiError("Authentication Failed", authRes.ApiError)
	}

	httpClient.SetToken(authRes.Token)

	if verbose {
		log.Println("Authentication Successful")
	}

	// Get profile data
	profileRes := ProfileResponse{}
	if err := httpClient.Get("/profile.json", &profileRes); err != nil {
		LogFatalError("Failed to get profile data", err)
	}
	if profileRes.ResultCode != "SUCCESS" {
		LogFatalApiError("Failed to get profile data", profileRes.ApiError)
	}

	if verbose {
		log.Println(profileRes.Profile.WelcomeMessage)
		log.Println(profileRes.Profile.Agents[0].Name)
		log.Println("Possible durations:")
		for _, product := range profileRes.Profile.Agents[0].KioskProducts {
			log.Println(product.ProductName)
		}
	}

	// Get plate data
	plateRes := PlateResponse{}
	if err := httpClient.Get("/platenumber/lookup/"+regNr+".json", &plateRes); err != nil {
		LogFatalError("Failed to get plate data", err)
	}
	if plateRes.ResultCode != "SUCCESS" {
		LogFatalApiError("Failed to get plate data", plateRes.ApiError)
	}
	if verbose {
		log.Printf("Reg. nr. %q (%q)\n", regNr, plateRes.Message)
	}

	// Status
	statusRes := StatusResponse{}
	if err := httpClient.Get("/parking/status.json?kioskAgentId="+profileRes.Profile.Agents[0].Id+"&plateNumber="+regNr, &statusRes); err != nil {
		LogFatalError("Failed to get status", err)
	}
	if statusRes.ResultCode != "SUCCESS" {
		LogFatalApiError("Failed to get status", statusRes.ApiError)
	}

	if statusRes.StartTime != "" && statusRes.EndTime != "" {
		now := time.Now()
		startTime, _ := time.Parse("2006-01-02 15:04:05", statusRes.StartTime)
		endTime, _ := time.Parse("2006-01-02 15:04:05", statusRes.EndTime)
		endQuarantineTime, _ := time.Parse("2006-01-02 15:04:05", statusRes.EndOfWaitingPeriod)
		duration := endTime.Sub(startTime)
		if now.Before(endTime) {
			log.Printf("Parking (%.1f hours) for %s is still active until %s\n", duration.Hours(), regNr, endTime.Format("2006-01-02 15:04:05"))
			os.Exit(0)
			return
		}
		if now.Before(endQuarantineTime) {
			log.Printf("Failed to register parking (%.1f hours) for %s. Quarantined until %s\n", duration.Hours(), regNr, endQuarantineTime.Format("2006-01-02 15:04:05"))
			os.Exit(1)
			return
		}
	}

	productId := profileRes.Profile.Agents[0].KioskProducts[0].Id

	if extended {
		productId = profileRes.Profile.Agents[0].KioskProducts[1].Id
	}

	registerRes := RegisterResponse{}
	httpClient.Post("/parking/register.json", &registerRes, RegisterRequest{
		PlateNumber:                  regNr,
		KioskAgentId:                 profileRes.Profile.Agents[0].Id,
		ProductId:                    productId,
		TermsAccepted:                false,
		AcceptedTermsAndConditionsId: nil,
		QualificationForm:            []string{},
	})
	if registerRes.ResultCode != "SUCCESS" {
		LogFatalApiError("Failed to register", registerRes.ApiError)
	}

	log.Printf("Successfully registered at %s parking for %q, from %s to %s\n", registerRes.Parking.AgentName, regNr, registerRes.Parking.StartTime, registerRes.Parking.EndTime)
	os.Exit(0)
}
