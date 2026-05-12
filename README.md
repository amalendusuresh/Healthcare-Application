# Healthcare Application — Hyperledger Fabric Chaincode

A **Hyperledger Fabric chaincode** (smart contract) written in Go that records doctor-patient appointments, prescriptions, and consultation notes on a permissioned blockchain ledger. Designed as a reference implementation for tamper-evident medical record-keeping with simple, on-chain doctor authorization.

![Hyperledger Fabric](https://img.shields.io/badge/Network-Hyperledger%20Fabric-2F3134?logo=hyperledger&logoColor=white)
![Go](https://img.shields.io/badge/Chaincode-Go-00ADD8?logo=go&logoColor=white)
![Sector](https://img.shields.io/badge/Sector-Healthcare-FF4081)
![Status](https://img.shields.io/badge/status-Reference%20Impl-success)

---

## ✨ Features

- ✅ **Schedule appointments** between patients and doctors
- ✅ **Retrieve appointment** details by ID
- ✅ **Prescribe medicine** during an appointment (doctor-authorized)
- ✅ **Provide consultation notes** during an appointment (doctor-authorized)
- ✅ **On-chain access control** — only the assigned doctor can update an appointment
- ✅ **JSON-serialized state** on the Fabric ledger for easy querying

---

## 🏗️ Architecture

```
   ┌─────────────────────────────────────────────────────┐
   │                Fabric Client / SDK                   │
   │   (CLI, Node.js, Go, or Java app invoking chaincode) │
   └────────────────────────┬─────────────────────────────┘
                            │ invoke / query
                            ▼
   ┌─────────────────────────────────────────────────────┐
   │            HealthcareChaincode (Go)                  │
   │                                                      │
   │   ┌─────────────────────────────────────────────┐    │
   │   │  scheduleAppointment(id, patient, doc, t)   │    │
   │   │  getAppointment(id)                         │    │
   │   │  prescribeMedicine(apptID, drug, docID)     │    │
   │   │  provideConsultation(apptID, notes, docID)  │    │
   │   └─────────────────────────────────────────────┘    │
   └────────────────────────┬─────────────────────────────┘
                            │ PutState / GetState
                            ▼
   ┌─────────────────────────────────────────────────────┐
   │           Fabric Ledger (per-channel)                │
   │     Appointment JSON keyed by appointment ID         │
   └─────────────────────────────────────────────────────┘
```

---

## 📦 Data Model

### Patient

```go
type Patient struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

### Doctor

```go
type Doctor struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

### Appointment

```go
type Appointment struct {
    ID           string `json:"id"`
    PatientID    string `json:"patient_id"`
    DoctorID     string `json:"doctor_id"`
    Time         string `json:"time"`
    Medicine     string `json:"medicine,omitempty"`
    Consultation string `json:"consultation,omitempty"`
}
```

---

## 🔧 Chaincode Functions

| Function | Args | Description |
|---|---|---|
| `scheduleAppointment` | `id`, `patientID`, `doctorID`, `time` | Creates a new appointment record on the ledger |
| `getAppointment` | `id` | Returns the JSON-serialized appointment for the given ID |
| `prescribeMedicine` | `appointmentID`, `medicine`, `doctorID` | Doctor records a prescription. Reverts if `doctorID` ≠ appointment's doctor |
| `provideConsultation` | `appointmentID`, `consultation`, `doctorID` | Doctor records consultation notes. Reverts if `doctorID` ≠ appointment's doctor |

### Access Control

Each state-changing function (prescribe / consult) reads the appointment, unmarshals it, and **verifies that the caller's `doctorID` matches the appointment's `DoctorID`** before writing. This is a simple chaincode-level check — for production, this should be replaced with **MSP-based identity** via `stub.GetCreator()` so doctors are authenticated via their X.509 certificates instead of a string argument.

---

## 📂 Project Structure

```
Healthcare-Application/
├── chaincode.go      # Main chaincode — structs, Init, Invoke, all functions
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites

- **Hyperledger Fabric 2.x** test network (use [`fabric-samples`](https://github.com/hyperledger/fabric-samples))
- **Go** ≥ 1.18
- **Docker** & **Docker Compose**
- **Fabric binaries** (`peer`, `orderer`, `fabric-ca-client`)

### Quick Start (test-network)

Clone Fabric samples and start the test network:

```bash
git clone https://github.com/hyperledger/fabric-samples.git
cd fabric-samples/test-network

# Boot up the network with a default channel
./network.sh up createChannel -c mychannel -ca
```

Clone this chaincode into the samples directory:

```bash
cd ../
git clone https://github.com/amalendusuresh/Healthcare-Application.git healthcare
```

Package, install, approve, and commit the chaincode:

```bash
cd test-network

./network.sh deployCC \
  -ccn healthcare \
  -ccp ../healthcare \
  -ccl go \
  -c mychannel
```

### Invoke from the CLI

Set up the peer CLI environment, then:

```bash
# Schedule an appointment
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile $ORDERER_CA \
  -C mychannel \
  -n healthcare \
  -c '{"function":"scheduleAppointment","Args":["A001","P001","D001","2026-05-20T10:00:00Z"]}'

# Query the appointment
peer chaincode query \
  -C mychannel \
  -n healthcare \
  -c '{"function":"getAppointment","Args":["A001"]}'

# Doctor prescribes medicine
peer chaincode invoke \
  -C mychannel \
  -n healthcare \
  -c '{"function":"prescribeMedicine","Args":["A001","Amoxicillin 500mg","D001"]}'

# Doctor records consultation
peer chaincode invoke \
  -C mychannel \
  -n healthcare \
  -c '{"function":"provideConsultation","Args":["A001","Patient reported mild fever. Vitals normal. Follow-up in 7 days.","D001"]}'
```

---

## 🔐 Security Considerations

- **Doctor authorization** is enforced inside the chaincode — only the appointment's assigned doctor can prescribe or update
- **Ledger immutability** — all writes produce permanent, tamper-evident history
- **Permissioned network** — only enrolled organizations & identities can participate
- **JSON marshaling errors** are caught and returned as chaincode errors (no partial state writes)

> ⚠️ This is a **reference implementation**. For real clinical deployment, see the roadmap below.

---

## 🗺️ Roadmap

- [ ] **MSP-based identity** — replace string `doctorID` arguments with `stub.GetCreator()` so callers are identified by their X.509 cert
- [ ] **Patient consent** — patient must approve a doctor before that doctor can view/update records
- [ ] **Private Data Collections (PDCs)** — keep PHI off the main channel; only hashes on-chain
- [ ] **Audit trail queries** — return full history of an appointment via `GetHistoryForKey`
- [ ] **Migrate to `contractapi`** — newer Fabric chaincode API for cleaner code & auto-generated metadata
- [ ] **Unit tests** with `MockStub`
- [ ] **Sample Node.js / Go SDK client** for end-to-end demo
- [ ] **HL7 FHIR mapping** for interoperability with hospital EHR systems

---

## 📚 Tech Stack

- **Language:** Go
- **Framework:** Hyperledger Fabric 2.x · Fabric Shim
- **API:** Legacy `shim` interface (`github.com/hyperledger/fabric/core/chaincode/shim`)
- **State:** Key-value store via `PutState` / `GetState`
- **Serialization:** JSON

---

## 📄 License

MIT © [Amalendu Suresh](https://github.com/amalendusuresh)

---

## 🤝 Contact

**Amalendu Suresh** — Blockchain Engineer

- 💼 **LinkedIn:** [amalendu-blockchain](https://www.linkedin.com/in/amalendu-blockchain/)
- ✍️ **Medium:** [@amalenduvishnu](https://medium.com/@amalenduvishnu)
- 📧 **Email:** amalendusuresh95@gmail.com

If you find this project useful, please ⭐ star the repo!
