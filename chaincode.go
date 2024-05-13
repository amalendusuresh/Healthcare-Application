package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// HealthcareChaincode represents the chaincode for the Healthcare Application
type HealthcareChaincode struct {
}

// Patient represents a patient in the system
type Patient struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Add more fields as needed
}

// Doctor represents a doctor in the system
type Doctor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Add more fields as needed
}

// Appointment represents an appointment between a patient and a doctor
type Appointment struct {
	ID        string `json:"id"`
	PatientID string `json:"patient_id"`
	DoctorID  string `json:"doctor_id"`
	Time      string `json:"time"`
	// Add more fields as needed
}

// Init initializes the chaincode
func (cc *HealthcareChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke handles invocation of functions
func (cc *HealthcareChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "scheduleAppointment":
		return cc.scheduleAppointment(stub, args)
	case "getAppointment":
		return cc.getAppointment(stub, args)
	case "prescribeMedicine":
		return cc.prescribeMedicine(stub, args)
	case "provideConsultation":
		return cc.provideConsultation(stub, args)
	// Add more cases for other functions
	default:
		return shim.Error("Invalid function name.")
	}
}

// ScheduleAppointment schedules an appointment between a patient and a doctor
func (cc *HealthcareChaincode) scheduleAppointment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4.")
	}

	appointment := Appointment{
		ID:        args[0],
		PatientID: args[1],
		DoctorID:  args[2],
		Time:      args[3],
	}

	appointmentJSON, err := json.Marshal(appointment)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to marshal appointment JSON: %s", err))
	}

	err = stub.PutState(appointment.ID, appointmentJSON)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to save appointment to ledger: %s", err))
	}

	return shim.Success(nil)
}

// GetAppointment retrieves an appointment by ID
func (cc *HealthcareChaincode) getAppointment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}

	appointmentJSON, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to read appointment from ledger: %s", err))
	}
	if appointmentJSON == nil {
		return shim.Error("Appointment not found.")
	}

	return shim.Success(appointmentJSON)
}

// PrescribeMedicine allows a doctor to prescribe medicine during an appointment
func (cc *HealthcareChaincode) prescribeMedicine(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3.")
	}

	appointmentID := args[0]
	medicine := args[1]
	doctorID := args[2]

	appointmentJSON, err := stub.GetState(appointmentID)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to read appointment from ledger: %s", err))
	}
	if appointmentJSON == nil {
		return shim.Error("Appointment not found.")
	}

	var appointment Appointment
	err = json.Unmarshal(appointmentJSON, &appointment)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to unmarshal appointment JSON: %s", err))
	}

	if appointment.DoctorID != doctorID {
		return shim.Error("Doctor is not authorized to prescribe medicine for this appointment.")
	}

	appointment.Medicine = medicine

	updatedAppointmentJSON, err := json.Marshal(appointment)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to marshal updated appointment JSON: %s", err))
	}

	err = stub.PutState(appointmentID, updatedAppointmentJSON)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to update appointment in ledger: %s", err))
	}

	return shim.Success(nil)
}

// ProvideConsultation allows a doctor to provide consultation during an appointment
func (cc *HealthcareChaincode) provideConsultation(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3.")
	}

	appointmentID := args[0]
	consultation := args[1]
	doctorID := args[2]

	appointmentJSON, err := stub.GetState(appointmentID)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to read appointment from ledger: %s", err))
	}
	if appointmentJSON == nil {
		return shim.Error("Appointment not found.")
	}

	var appointment Appointment
	err = json.Unmarshal(appointmentJSON, &appointment)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to unmarshal appointment JSON: %s", err))
	}

	if appointment.DoctorID != doctorID {
		return shim.Error("Doctor is not authorized to provide consultation for this appointment.")
	}

	appointment.Consultation = consultation

	updatedAppointmentJSON, err := json.Marshal(appointment)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to marshal updated appointment JSON: %s", err))
	}

	err = stub.PutState(appointmentID, updatedAppointmentJSON)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to update appointment in ledger: %s", err))
	}

	return shim.Success(nil)
}


func main() {
	err := shim.Start(new(HealthcareChaincode))
	if err != nil {
		fmt.Printf("Error starting HealthcareChaincode: %s", err)
	}
}
