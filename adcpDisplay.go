package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ricorx7/go-rti"
)

const (
	adcpListID        = "AdcpList"
	adcpEnsembleID    = "Ensemble"
	velBeamDataID     = "VelBeamData"
	profileAmpID      = "ProfileAmpData"
	profileCorrID     = "ProfileCorrData"
	profileID         = "ProfileData"
	profileRickshawID = "ProfileRickshawData"
	profileC3ID       = "ProfileC3Data"
	profileEpochID    = "ProfileEpochData"
	hprID             = "HprData"
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

// profileBeamData will store the data for a bin.
type profileBeamData struct {
	Values [][]float32 `json:"values"` // JSON Data
	Color  string      `json:"color"`  // Color option
	Key    string      `json:"key"`    // Data key
	Area   bool        `json:"area"`   // Flag if data should be area plot
}

// timeSeriesData will store the data for a time series.
type timeSeriesData struct {
	Values [][]float32 `json:"values"` // JSON Data
	Color  string      `json:"color"`  // Color option
	Key    string      `json:"key"`    // Data key
	Area   bool        `json:"area"`   // Flag if data should be area plot
}

// profileData will store the profile data.
type profileData struct {
	ID        string            // Data ID
	SerialNum string            // Serial Number
	CepoIndex uint8             // Subystem configuration
	AmpData   []profileBeamData // Array of all the Amplitude series
	CorrData  []profileBeamData // Array of all the Correlation series
}

// profileRickshawData will store the profile data.
type profileRickshawData struct {
	ID         string           // Data ID
	SerialNum  string           // Serial Number
	CepoIndex  uint8            // Subystem configuration
	AmpB0Data  lineRickshawData // Array of all the Amplitude B0 series
	AmpB1Data  lineRickshawData // Array of all the Amplitude B1 series
	AmpB2Data  lineRickshawData // Array of all the Amplitude B2 series
	AmpB3Data  lineRickshawData // Array of all the Amplitude B3 series
	CorrB0Data lineRickshawData // Array of all the Correlation B0 series
	CorrB1Data lineRickshawData // Array of all the Correlation B1 series
	CorrB2Data lineRickshawData // Array of all the Correlation B2 series
	CorrB3Data lineRickshawData // Array of all the Correlation B3 series
}

// pointRickshawData is the point data with x and y.
type pointRickshawData struct {
	X float32 `json:"x"` // X data
	Y float32 `json:"y"` // Y data
}

// lineRickshawData is the line.
type lineRickshawData struct {
	Data  []pointRickshawData `json:"data"`  // JSON Data
	Color string              `json:"color"` // Color option
	Key   string              `json:"key"`   // Data key
	Area  bool                `json:"area"`  // Flag if data should be area plot
}

// profileC3Data will store the profile data.
type profileC3Data struct {
	ID         string    // Data ID
	SerialNum  string    // Serial Number
	CepoIndex  uint8     // Subystem configuration
	AmpXAxis   []float32 // Amplitude X Axis Labels
	AmpB0Data  []float32 // Array of all the Amplitude B0 series
	AmpB1Data  []float32 // Array of all the Amplitude B1 series
	AmpB2Data  []float32 // Array of all the Amplitude B2 series
	AmpB3Data  []float32 // Array of all the Amplitude B3 series
	CorrXAxis  []float32 // Correlation X Axis Labels
	CorrB0Data []float32 // Array of all the Correlation B0 series
	CorrB1Data []float32 // Array of all the Correlation B1 series
	CorrB2Data []float32 // Array of all the Correlation B2 series
	CorrB3Data []float32 // Array of all the Correlation B3 series
}

// hprData will store the heading, pitch and roll data.
type hprData struct {
	ID        string           // Data ID
	SerialNum string           // Serial Number
	CepoIndex uint8            // Subystem configuration
	HprData   []timeSeriesData // Array of all the hpr series
}

// heading will accumulate the heading values.
var heading = &timeSeriesData{
	Color: "#ff7f0e", // Color of the plot
	Key:   "Heading", // Key for the data
	Area:  false,     // Flag for area plot
}

// pitch will accumulate the pitch values.
var pitch = &timeSeriesData{
	Color: "#2ca02c", // Color of the plot
	Key:   "Pitch",   // Key for the data
	Area:  false,     // Flag for area plot
}

// roll will accumulate the rol values.
var roll = &timeSeriesData{
	Color: "#7777ff", // Color of the plot
	Key:   "Roll",    // Key for the data
	Area:  false,     // Flag for area plot
}

// pointEpochData is the point data with x and y.
type pointEpochPointData struct {
	X float32 `json:"x"` // X data
	Y float32 `json:"y"` // Y data
}

// pointEpochData is the point data with x and y.
type pointEpochTimeData struct {
	Time float32 `json:"time"` // Time data
	Y    float32 `json:"y"`    // Y data
}

// pointEpochHeatmapData is the histogram data with bin and average value.
type seriesEpochHeatmapData struct {
	Time      int64     `json:"time"`      // Time data
	Histogram []float64 `json:"histogram"` // Histgram data
}

// seriesEpochData holds the Epoch series data.
type seriesEpochData struct {
	Label  string                `json:"label"`  // Series Label
	Values []pointEpochPointData `json:"values"` // Values
}

// seriesEpochRealtimeData holds the Epoch series data.
type seriesEpochRealtimeData struct {
	Label  string               `json:"label"`  // Series Label
	Values []pointEpochTimeData `json:"values"` // Values
}

// profileEpochData will store the Epoch plot data.
type profileEpochData struct {
	ID                   string                    // Data ID
	SerialNum            string                    // Serial Number
	CepoIndex            uint8                     // Subystem configuration
	Data                 []seriesEpochData         // Data
	RealtimeData         []seriesEpochRealtimeData // Realtime data
	HeatmapMagData       seriesEpochHeatmapData    // Heatmap Magnitude data
	HeatmapDirYNorthData seriesEpochHeatmapData    // Heatmap Direction Y North data
	HeatmapDirXNorthData seriesEpochHeatmapData    // Heatmap Direction X North data
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

	// Send Profile data
	sendProfilePlotData(ens)

	// Send Profile Rickshaw data
	sendProfileRickshawPlotData(ens)

	// Send Profile C3 data
	sendProfileC3PlotData(ens)

	// Send Profile Epoch data
	sendProfileEpochPlotData(ens)

	// Send HPR data
	sendHprPlotData(ens)
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

// sendProfilePlotData will accumulate the amplitude and correlation data
// to pass to the display.
func sendProfilePlotData(ens rti.Ensemble) {

	// Create the data struct
	ampProfileB0 := &profileBeamData{
		Color: "#ff7f0e", // Color of the plot
		Key:   "B0",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	ampProfileB1 := &profileBeamData{
		Color: "#2ca02c", // Color of the plot
		Key:   "B1",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	ampProfileB2 := &profileBeamData{
		Color: "#7777ff", // Color of the plot
		Key:   "B2",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	ampProfileB3 := &profileBeamData{
		Color: "#d67777", // Color of the plot
		Key:   "B3",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	corrProfileB0 := &profileBeamData{
		Color: "#ff7f0e", // Color of the plot
		Key:   "B0",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	corrProfileB1 := &profileBeamData{
		Color: "#2ca02c", // Color of the plot
		Key:   "B1",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	corrProfileB2 := &profileBeamData{
		Color: "#7777ff", // Color of the plot
		Key:   "B2",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	corrProfileB3 := &profileBeamData{
		Color: "#d67777", // Color of the plot
		Key:   "B3",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	profData := &profileData{
		ID:        profileID,                                  // ID
		SerialNum: ens.EnsembleData.SerialNumber.SerialNumber, // Serial Data
		CepoIndex: ens.EnsembleData.SubsystemConfig.CepoIndex, // Subsystem Config
	}

	//for bin := int(ens.AmplitudeData.Base.NumElements - 1); bin >= 0; bin-- {
	for bin := 0; bin < int(ens.AmplitudeData.Base.NumElements); bin++ {
		for beam := 0; beam < int(ens.AmplitudeData.Base.ElementMultiplier); beam++ {
			if beam == 0 {
				var ampData = []float32{ens.AmplitudeData.Amplitude[bin][beam], float32(bin)} // Create point to hold [key, value]
				ampProfileB0.Values = append(ampProfileB0.Values, ampData)                    // Add the key and value array to array

				var corrData = []float32{ens.CorrelationData.Correlation[bin][beam], float32(bin)} // Create point to hold [key, value]
				corrProfileB0.Values = append(corrProfileB0.Values, corrData)                      // Add the key and value array to array
			}
			if beam == 1 {
				var ampData = []float32{ens.AmplitudeData.Amplitude[bin][beam], float32(bin)} // Create point to hold [key, value]
				ampProfileB1.Values = append(ampProfileB1.Values, ampData)                    // Add the key and value array to array

				var corrData = []float32{ens.CorrelationData.Correlation[bin][beam], float32(bin)} // Create point to hold [key, value]
				corrProfileB1.Values = append(corrProfileB1.Values, corrData)                      // Add the key and value array to array
			}
			if beam == 2 {
				var ampData = []float32{ens.AmplitudeData.Amplitude[bin][beam], float32(bin)} // Create point to hold [key, value]
				ampProfileB2.Values = append(ampProfileB2.Values, ampData)                    // Add the key and value array to array

				var corrData = []float32{ens.CorrelationData.Correlation[bin][beam], float32(bin)} // Create point to hold [key, value]
				corrProfileB2.Values = append(corrProfileB2.Values, corrData)                      // Add the key and value array to array
			}
			if beam == 3 {
				var ampData = []float32{ens.AmplitudeData.Amplitude[bin][beam], float32(bin)} // Create point to hold [key, value]
				ampProfileB3.Values = append(ampProfileB3.Values, ampData)                    // Add the key and value array to array

				var corrData = []float32{ens.CorrelationData.Correlation[bin][beam], float32(bin)} // Create point to hold [key, value]
				corrProfileB3.Values = append(corrProfileB3.Values, corrData)                      // Add the key and value array to array
			}
		}
	}

	// Set the data
	profData.AmpData = append(profData.AmpData, *ampProfileB0)
	profData.AmpData = append(profData.AmpData, *ampProfileB1)
	profData.AmpData = append(profData.AmpData, *ampProfileB2)
	profData.AmpData = append(profData.AmpData, *ampProfileB3)

	profData.CorrData = append(profData.CorrData, *corrProfileB0)
	profData.CorrData = append(profData.CorrData, *corrProfileB1)
	profData.CorrData = append(profData.CorrData, *corrProfileB2)
	profData.CorrData = append(profData.CorrData, *corrProfileB3)

	// Convert the JSON to byte array
	b, err := json.Marshal(profData)
	if err != nil {
		log.Println(err)
		return
	}

	// Send the data to the display
	sendDataToDisplays(b)
}

// sendProfileRickshawPlotData will accumulate the amplitude and correlation data
// to pass to the display.
func sendProfileRickshawPlotData(ens rti.Ensemble) {

	// Create the data struct
	ampProfileB0 := &lineRickshawData{
		Color: "#ff7f0e", // Color of the plot
		Key:   "B0",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	ampProfileB1 := &lineRickshawData{
		Color: "#2ca02c", // Color of the plot
		Key:   "B1",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	ampProfileB2 := &lineRickshawData{
		Color: "#7777ff", // Color of the plot
		Key:   "B2",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	ampProfileB3 := &lineRickshawData{
		Color: "#d67777", // Color of the plot
		Key:   "B3",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	corrProfileB0 := &lineRickshawData{
		Color: "#ff7f0e", // Color of the plot
		Key:   "B0",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	corrProfileB1 := &lineRickshawData{
		Color: "#2ca02c", // Color of the plot
		Key:   "B1",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	corrProfileB2 := &lineRickshawData{
		Color: "#7777ff", // Color of the plot
		Key:   "B2",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	corrProfileB3 := &lineRickshawData{
		Color: "#d67777", // Color of the plot
		Key:   "B3",      // Key for the data
		Area:  false,     // Flag for area plot
	}

	profData := &profileRickshawData{
		ID:        profileRickshawID,                          // ID
		SerialNum: ens.EnsembleData.SerialNumber.SerialNumber, // Serial Data
		CepoIndex: ens.EnsembleData.SubsystemConfig.CepoIndex, // Subsystem Config
	}

	//for bin := int(ens.AmplitudeData.Base.NumElements - 1); bin >= 0; bin-- {
	for bin := 0; bin < int(ens.AmplitudeData.Base.NumElements); bin++ {
		for beam := 0; beam < int(ens.AmplitudeData.Base.ElementMultiplier); beam++ {
			if beam == 0 {
				ampData := pointRickshawData{X: ens.AmplitudeData.Amplitude[bin][beam], Y: float32(bin)} // Create point to hold [key, value]
				ampProfileB0.Data = append(ampProfileB0.Data, ampData)                                   // Add the key and value array to array

				corrData := pointRickshawData{X: ens.CorrelationData.Correlation[bin][beam], Y: float32(bin)} // Create point to hold [key, value]
				corrProfileB0.Data = append(corrProfileB0.Data, corrData)                                     // Add the key and value array to array
			}
			if beam == 1 {
				ampData := pointRickshawData{X: ens.AmplitudeData.Amplitude[bin][beam], Y: float32(bin)} // Create point to hold [key, value]
				ampProfileB1.Data = append(ampProfileB1.Data, ampData)                                   // Add the key and value array to array

				corrData := pointRickshawData{X: ens.CorrelationData.Correlation[bin][beam], Y: float32(bin)} // Create point to hold [key, value]
				corrProfileB1.Data = append(corrProfileB1.Data, corrData)                                     // Add the key and value array to array
			}
			if beam == 2 {
				ampData := pointRickshawData{X: ens.AmplitudeData.Amplitude[bin][beam], Y: float32(bin)} // Create point to hold [key, value]
				ampProfileB2.Data = append(ampProfileB2.Data, ampData)                                   // Add the key and value array to array

				corrData := pointRickshawData{X: ens.CorrelationData.Correlation[bin][beam], Y: float32(bin)} // Create point to hold [key, value]
				corrProfileB2.Data = append(corrProfileB2.Data, corrData)                                     // Add the key and value array to array
			}
			if beam == 3 {
				ampData := pointRickshawData{X: ens.AmplitudeData.Amplitude[bin][beam], Y: float32(bin)} // Create point to hold [key, value]
				ampProfileB3.Data = append(ampProfileB3.Data, ampData)                                   // Add the key and value array to array

				corrData := pointRickshawData{X: ens.CorrelationData.Correlation[bin][beam], Y: float32(bin)} // Create point to hold [key, value]
				corrProfileB3.Data = append(corrProfileB3.Data, corrData)                                     // Add the key and value array to array
			}
		}
	}

	profData.AmpB0Data = *ampProfileB0
	profData.AmpB1Data = *ampProfileB1
	profData.AmpB2Data = *ampProfileB2
	profData.AmpB3Data = *ampProfileB3
	profData.CorrB0Data = *corrProfileB0
	profData.CorrB1Data = *corrProfileB1
	profData.CorrB2Data = *corrProfileB2
	profData.CorrB3Data = *corrProfileB3

	// Convert the JSON to byte array
	b, err := json.Marshal(profData)
	if err != nil {
		log.Println(err)
		return
	}

	// Send the data to the display
	sendDataToDisplays(b)
}

// sendProfileC3PlotData will accumulate the amplitude and correlation data
// to pass to the display.
func sendProfileC3PlotData(ens rti.Ensemble) {

	profData := &profileC3Data{
		ID:        profileC3ID,                                // ID
		SerialNum: ens.EnsembleData.SerialNumber.SerialNumber, // Serial Data
		CepoIndex: ens.EnsembleData.SubsystemConfig.CepoIndex, // Subsystem Config
	}

	//for bin := int(ens.AmplitudeData.Base.NumElements - 1); bin >= 0; bin-- {
	for bin := 0; bin < int(ens.AmplitudeData.Base.NumElements); bin++ {
		// First bin location
		firstBinLoc := ens.AncillaryData.FirstBinRange + (ens.AncillaryData.BinSize * float32(bin))
		profData.AmpXAxis = append(profData.AmpXAxis, float32(bin))
		profData.CorrXAxis = append(profData.CorrXAxis, firstBinLoc)

		for beam := 0; beam < int(ens.AmplitudeData.Base.ElementMultiplier); beam++ {
			if beam == 0 {
				profData.AmpB0Data = append(profData.AmpB0Data, ens.AmplitudeData.Amplitude[bin][beam])
				profData.CorrB0Data = append(profData.CorrB0Data, ens.CorrelationData.Correlation[bin][beam])
			}
			if beam == 1 {
				profData.AmpB1Data = append(profData.AmpB1Data, ens.AmplitudeData.Amplitude[bin][beam])
				profData.CorrB1Data = append(profData.CorrB1Data, ens.CorrelationData.Correlation[bin][beam])
			}
			if beam == 2 {
				profData.AmpB2Data = append(profData.AmpB2Data, ens.AmplitudeData.Amplitude[bin][beam])
				profData.CorrB2Data = append(profData.CorrB2Data, ens.CorrelationData.Correlation[bin][beam])
			}
			if beam == 3 {
				profData.AmpB3Data = append(profData.AmpB3Data, ens.AmplitudeData.Amplitude[bin][beam])
				profData.CorrB3Data = append(profData.CorrB3Data, ens.CorrelationData.Correlation[bin][beam])
			}
		}
	}

	// Convert the JSON to byte array
	b, err := json.Marshal(profData)
	if err != nil {
		log.Println(err)
		return
	}

	// Send the data to the display
	sendDataToDisplays(b)
}

// sendHprPlotData will accumulate the heading, pitch and roll data
// to pass to the display.
func sendHprPlotData(ens rti.Ensemble) {

	// Accumulate heading
	heading.Values = append(heading.Values, []float32{float32(ens.EnsembleData.EnsembleNumber), ens.AncillaryData.Heading})
	if len(heading.Values) > 20 {
		heading.Values = heading.Values[1:]
	}

	// Accumulate pitch
	pitch.Values = append(pitch.Values, []float32{float32(ens.EnsembleData.EnsembleNumber), ens.AncillaryData.Pitch})
	if len(pitch.Values) > 20 {
		pitch.Values = pitch.Values[1:]
	}

	// Accumulate roll
	roll.Values = append(roll.Values, []float32{float32(ens.EnsembleData.EnsembleNumber), ens.AncillaryData.Roll})
	if len(roll.Values) > 20 {
		roll.Values = roll.Values[1:]
	}

	hpr := &hprData{
		ID:        hprID,                                      // ID
		SerialNum: ens.EnsembleData.SerialNumber.SerialNumber, // Serial Data
		CepoIndex: ens.EnsembleData.SubsystemConfig.CepoIndex, // Subsystem Config
	}

	hpr.HprData = append(hpr.HprData, *heading)
	hpr.HprData = append(hpr.HprData, *pitch)
	hpr.HprData = append(hpr.HprData, *roll)

	// Convert the JSON to byte array
	b, err := json.Marshal(hpr)
	if err != nil {
		log.Println(err)
		return
	}

	// Send the data to the display
	sendDataToDisplays(b)
}

// sendProfileEpochPlotData will accumulate the amplitude and correlation data
// to pass to the display.
func sendProfileEpochPlotData(ens rti.Ensemble) {

	profData := &profileEpochData{
		ID:        profileEpochID,                             // ID
		SerialNum: ens.EnsembleData.SerialNumber.SerialNumber, // Serial Data
		CepoIndex: ens.EnsembleData.SubsystemConfig.CepoIndex, // Subsystem Config
	}

	// Get the time
	var year = int(ens.EnsembleData.Year + 2000)
	var month = time.Month(ens.EnsembleData.Month)
	var day = int(ens.EnsembleData.Day)
	var hour = int(ens.EnsembleData.Hour)
	var min = int(ens.EnsembleData.Minute)
	var sec = int(ens.EnsembleData.Second)
	//var nsec = int(ens.EnsembleData.HSec * 0.0000001)
	var nsec = 0
	var unixTime = time.Date(year, month, day, hour, min, sec, nsec, time.Local).Unix()

	// Init the maps
	profData.HeatmapMagData.Histogram = make([]float64, int(ens.EarthVelocityData.Base.NumElements))
	profData.HeatmapDirXNorthData.Histogram = make([]float64, int(ens.EarthVelocityData.Base.NumElements))
	profData.HeatmapDirYNorthData.Histogram = make([]float64, int(ens.EarthVelocityData.Base.NumElements))

	// Set the time
	profData.HeatmapMagData.Time = unixTime
	profData.HeatmapDirXNorthData.Time = unixTime
	profData.HeatmapDirYNorthData.Time = unixTime

	// Set the heatmap data
	for bin := 0; bin < int(ens.EarthVelocityData.Base.NumElements); bin++ {
		profData.HeatmapMagData.Histogram[bin] = ens.EarthVelocityData.Vectors[bin].Magnitude
		profData.HeatmapDirXNorthData.Histogram[bin] = ens.EarthVelocityData.Vectors[bin].DirectionXNorth
		profData.HeatmapDirYNorthData.Histogram[bin] = ens.EarthVelocityData.Vectors[bin].DirectionYNorth
	}

	// Convert the JSON to byte array
	b, err := json.Marshal(profData)
	if err != nil {
		log.Println(err)
		return
	}

	// Send the data to the display
	sendDataToDisplays(b)
}
