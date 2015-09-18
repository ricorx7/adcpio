package main

import (
	"encoding/json"
	"log"

	"github.com/ricorx7/go-rti"
)

const (
	adcpListID     = "AdcpList"
	adcpEnsembleID = "Ensemble"
	velBeamDataID  = "VelBeamData"
	profileAmpID   = "ProfileAmpData"
	profileCorrID  = "ProfileCorrData"
)

// adcp will store all the ADCP it is monitoring and also the last ensemble.
type adcpList struct {
	ID            string   // Data ID
	SerialNumList []string // List of Serial number
}

// adcp will store the last ensemble.
type adcpEnsemble struct {
	ID              string       // Data ID
	SerialNum       string       // Serial number
	SubsystemConfig string       // Subsystem configuration
	Ens             rti.Ensemble // Last ensemble
}

// velBeamData will store the beam velocity data.
type velBeamData struct {
	ID              string    // Data ID
	SerialNum       string    // Serial Number
	SubsystemConfig string    // Subystem configuration
	Labels          []string  // Seprate each ensemble's data with this label
	Beam0Vel        []float32 // Beam 0 velocity
	Beam1Vel        []float32 // Beam 1 velocity
	Beam2Vel        []float32 // Beam 2 velocity
	Beam3Vel        []float32 // Beam 3 velocity
}

// profileData will store the amplitude and correlation data.
type profileData struct {
	ID              string    // Data ID
	SerialNum       string    // Serial Number
	SubsystemConfig string    // Subystem configuration
	Labels          []string  // Seprate each ensemble's data with this label
	Beam0Amp        []float32 // Beam 0 Amplitude
	Beam1Amp        []float32 // Beam 1 Amplitude
	Beam2Amp        []float32 // Beam 2 Amplitude
	Beam3Amp        []float32 // Beam 3 Amplitude
	Beam0Corr       []float32 // Beam 0 Correlation
	Beam1Corr       []float32 // Beam 1 Correlation
	Beam2Corr       []float32 // Beam 2 Correlation
	Beam3Corr       []float32 // Beam 3 Correlation
}

// processEnsemble will take the ensemble data and add it to the map.
// It will then set the latest ensemble.
func processEnsemble(server *adcpIO, ens rti.Ensemble) {
	// See if the serial number exist in the map
	if val, ok := server.adcp[ens.EnsembleData.SerialNumber.SerialNumber]; ok {
		val.lastEns = ens
		log.Print("ADCP exist")
	} else {
		var data adcp
		data.lastEns = ens
		server.adcp[ens.EnsembleData.SerialNumber.SerialNumber] = &data
		log.Print("ADCP does not exist")

		// Send a new list of ADCP
		sendAdcpList()
	}

	// Send last ensemble to display
	sendRawEnsemble(ens)

	// Set the velocity data
	sendProfilePlotData(ens)
}

// sendRawEnsemble will send the ensemble to the registered displays through
// the websocket connection.
func sendRawEnsemble(ens rti.Ensemble) {
	// Create the data struct
	adcpEns := &adcpEnsemble{
		ID:              adcpEnsembleID,                             // ID
		SerialNum:       ens.EnsembleData.SerialNumber.SerialNumber, // Serial Data
		SubsystemConfig: "",
		Ens:             ens,
	}

	// Convert the JSON to byte array
	b, err := json.Marshal(adcpEns)
	if err != nil {
		log.Println(err)
		return
	}

	// Send the data to the display
	sendDataToDisplays(b)
}

type profileBinData struct {
	X  float32 `json:"x"` // X values
	B0 float32 // Beam 0 value
	B1 float32 // Beam 1 value
	B2 float32 // Beam 2 value
	B3 float32 // Beam 3 value
}

type profileData1 struct {
	ID              string // Data ID
	SerialNum       string // Serial Number
	SubsystemConfig string // Subystem configuration
	Data            string // JSON Data
	Options         string // JSON Options
}

func sendProfilePlotData(ens rti.Ensemble) {

	//var ampData = make([]profileBinData, ens.AmplitudeData.Base.NumElements)
	//var corrData = make([]profileBinData, ens.AmplitudeData.Base.NumElements)
	var ampData []profileBinData
	var corrData []profileBinData

	for bin := 0; bin < int(ens.AmplitudeData.Base.NumElements); bin++ {
		ampBinData := profileBinData{}
		ampBinData.X = float32(bin)
		corrBinData := profileBinData{}
		corrBinData.X = float32(bin)

		for beam := 0; beam < int(ens.AmplitudeData.Base.ElementMultiplier); beam++ {
			if beam == 0 {
				ampBinData.B0 = ens.AmplitudeData.Amplitude[bin][beam]
				corrBinData.B0 = ens.CorrelationData.Correlation[bin][beam]
			}
			if beam == 1 {
				ampBinData.B1 = ens.AmplitudeData.Amplitude[bin][beam]
				corrBinData.B1 = ens.CorrelationData.Correlation[bin][beam]
			}
			if beam == 2 {
				ampBinData.B2 = ens.AmplitudeData.Amplitude[bin][beam]
				corrBinData.B2 = ens.CorrelationData.Correlation[bin][beam]
			}
			if beam == 3 {
				ampBinData.B3 = ens.AmplitudeData.Amplitude[bin][beam]
				corrBinData.B3 = ens.CorrelationData.Correlation[bin][beam]
			}
		}

		ampData = append(ampData, ampBinData)
		corrData = append(corrData, corrBinData)
	}

	// Options string
	options := "{ series: [ {y: 'B0', color: 'steelblue', thickness: '2px', type: 'line', striped: true, label: 'B0'}, {y: 'B1', color: '#ff7f0e', thickness: '2px', type: 'line', striped: true, label: 'B1'}, {y: 'B2', color: '#2ca02c', thickness: '2px', type: 'line', striped: true, label: 'B2'}, {y: 'B3', color: 'lightsteelblue', thickness: '2px', type: 'line', striped: true, label: 'B3'}], lineMode: 'linear',}"
	if ens.EnsembleData.NumBeams == 1 {
		options = "{ series: [ {y: 'B0', color: 'steelblue', thickness: '2px', type: 'line', striped: true, label: 'B0'}], lineMode: 'linear',}"
	}

	// Convert the JSON to byte array
	ampByte, err := json.Marshal(ampData)
	if err != nil {
		log.Println(err)
		return
	}

	// Convert the JSON to byte array
	corrByte, err1 := json.Marshal(corrData)
	if err1 != nil {
		log.Println(err1)
		return
	}

	ampProfile := &profileData1{
		ID:              profileAmpID,
		SerialNum:       ens.EnsembleData.SerialNumber.SerialNumber, // Serial Data
		SubsystemConfig: "",                                         // Subsystem Config
		Data:            string(ampByte),                            // Amplitude data
		Options:         options,                                    // Options
	}

	corrProfile := &profileData1{
		ID:              profileCorrID,
		SerialNum:       ens.EnsembleData.SerialNumber.SerialNumber, // Serial Data
		SubsystemConfig: "",                                         // Subsystem Config
		Data:            string(corrByte),                           // Correlation data
		Options:         options,                                    // Options
	}

	// Convert the JSON to byte array
	ampProfileJSON, err := json.Marshal(ampProfile)
	if err != nil {
		log.Println(err)
		return
	}

	// Convert the JSON to byte array
	corrProfileJSON, err := json.Marshal(corrProfile)
	if err != nil {
		log.Println(err)
		return
	}

	//log.Printf("%s", ampProfile.Data)

	// Send the data to the display
	sendDataToDisplays(ampProfileJSON)
	sendDataToDisplays(corrProfileJSON)

}

func sendProfile(ens rti.Ensemble) {
	// Create the data struct
	profile := &profileData{
		ID:              adcpEnsembleID,                             // ID
		SerialNum:       ens.EnsembleData.SerialNumber.SerialNumber, // Serial Data
		SubsystemConfig: "",                                         // Subsystem Config
		Labels:          make([]string, ens.AmplitudeData.Base.NumElements),
		Beam0Amp:        make([]float32, ens.AmplitudeData.Base.NumElements),
		Beam1Amp:        make([]float32, ens.AmplitudeData.Base.NumElements),
		Beam2Amp:        make([]float32, ens.AmplitudeData.Base.NumElements),
		Beam3Amp:        make([]float32, ens.AmplitudeData.Base.NumElements),
		Beam0Corr:       make([]float32, ens.AmplitudeData.Base.NumElements),
		Beam1Corr:       make([]float32, ens.AmplitudeData.Base.NumElements),
		Beam2Corr:       make([]float32, ens.AmplitudeData.Base.NumElements),
		Beam3Corr:       make([]float32, ens.AmplitudeData.Base.NumElements),
	}

	//for bin := 0; bin < int(ens.AmplitudeData.Base.NumElements); bin++ {
	for beam := 0; beam < int(ens.AmplitudeData.Base.ElementMultiplier); beam++ {
		if beam == 0 {
			profile.Beam0Amp = ens.AmplitudeData.Amplitude[0:ens.AmplitudeData.Base.NumElements][0]
			profile.Beam0Corr = ens.CorrelationData.Correlation[0:ens.CorrelationData.Base.NumElements][0]
		}
		if beam == 1 {
			profile.Beam1Amp = ens.AmplitudeData.Amplitude[0:ens.AmplitudeData.Base.NumElements][1]
			profile.Beam1Corr = ens.CorrelationData.Correlation[0:ens.CorrelationData.Base.NumElements][1]
		}
		if beam == 2 {
			profile.Beam2Amp = ens.AmplitudeData.Amplitude[0:ens.AmplitudeData.Base.NumElements][2]
			profile.Beam2Corr = ens.CorrelationData.Correlation[0:ens.CorrelationData.Base.NumElements][2]
		}
		if beam == 3 {
			profile.Beam3Amp = ens.AmplitudeData.Amplitude[0:ens.AmplitudeData.Base.NumElements][3]
			profile.Beam3Corr = ens.CorrelationData.Correlation[0:ens.CorrelationData.Base.NumElements][3]
		}
	}

	// Convert the JSON to byte array
	b, err := json.Marshal(profile)
	if err != nil {
		log.Println(err)
		return
	}

	// Send the data to the display
	sendDataToDisplays(b)
}
