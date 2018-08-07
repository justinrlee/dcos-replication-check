package main

import (
	"encoding/json"
	// "io"
	// "io/ioutil"
	// "log"
	"net/http"
	"github.com/Sirupsen/logrus"
)

type MetricsSnapshotResponse struct {
	OverlayRecovered 		float64 `json:"overlay/log/recovered"`
	// OverlayRecovered 		float64 `json:"master/slave_registrations"`
	MesosRecovered     	float64 `json:"registrar/log/recovered"`
}

//Index for slash, returns version
func Index(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// TODO This isn't actually currently JSON
	w.Write([]byte("Health check tool running"))
	// json.NewEncoder(w).Encode("Health check tool running")
}

// TODO lots of refactoring - DRY

func MesosRepLog(w http.ResponseWriter, r *http.Request) {
	var snapshot MetricsSnapshotResponse
	// TODO This isn't actually currently JSON
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := client.Request("GET", 5050, "metrics/snapshot")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Mesos Replog state not accessible`))
		return
	}

	err = json.Unmarshal(body, &snapshot)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Mesos Replog state not accessible`))
		return
	}
	// logrus.Infoln(snapshot)

	if snapshot.MesosRecovered == 1.0 {
		w.Write([]byte(`Mesos Replog Healthy: {"registrar/log/recovered": 1.0}`))
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
		// TODO this should have the actual value, although I do think this is a binary.
		w.Write([]byte(`Mesos Replog Unhealthy: {"registrar/log/recovered": 0.0}`))
	}
}

func OverlayRepLog(w http.ResponseWriter, r *http.Request) {
	var snapshot MetricsSnapshotResponse
	// TODO This isn't actually currently JSON
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := client.Request("GET", 5050, "metrics/snapshot")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Overlay Replog state not accessible`))
		return
	}

	err = json.Unmarshal(body, &snapshot)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Overlay Replog state not accessible`))
		return
	}
	// logrus.Infoln(snapshot)

	if snapshot.OverlayRecovered == 1.0 {
		w.Write([]byte(`Overlay Replog Healthy: {"overlay/log/recovered": 1.0}`))
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
		// TODO this should have the actual value, although I do think this is a binary.
		w.Write([]byte(`Overlay Replog Unhealthy: {"overlay/log/recovered": 0.0}`))
	}
}

func RepLog(w http.ResponseWriter, r *http.Request) {
	var snapshot MetricsSnapshotResponse
	// TODO This isn't actually currently JSON
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := client.Request("GET", 5050, "metrics/snapshot")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Replog state not accessible`))
		return
	}

	err = json.Unmarshal(body, &snapshot)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Replog state not accessible`))
		return
	}
	// logrus.Infoln(snapshot)

	if (snapshot.OverlayRecovered == 1.0 && snapshot.MesosRecovered == 1.0) {
		w.Write([]byte(`Mesos and Overlay Replogs Healthy: {"registrar/log/recovered": 1.0, "overlay/log/recovered": 1.0}`))
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
		// TODO this should have the actual value, although I do think this is a binary.
		w.Write([]byte(`Replog unhealthy`))
	}
}

type ExhibitorStatus struct {
	Hostname 				string  `json:"hostname"`
	Status     			string  `json:"description"`
	StatusCode     	float64 `json:"code"`
	Leader					bool    `json:"isLeader"`
}

func Exhibitor(w http.ResponseWriter, r *http.Request) {
	// var snapshot MetricsSnapshotResponse
	// TODO This isn't actually currently JSON
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := client.Request("GET", 443, "exhibitor/exhibitor/v1/cluster/status")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Exhibitor state not accessible`))
		return
	}

	var exhibitorStatuses []ExhibitorStatus

	json.Unmarshal(body, &exhibitorStatuses)

	// logrus.Infoln(exhibitorStatuses)

	healthy := true
	leaderCount := 0
	for _, status := range exhibitorStatuses {
		if status.Status != "serving" {
			healthy = false
		}
		if status.StatusCode != 3 {
			healthy = false
		}
		if status.Leader {
			leaderCount += 1
		}
	}

	// w.Write([]byte(`Exhibitor state hmm accessible`))

	if (leaderCount == 1 && healthy) {
		w.Write([]byte("Exhibitor appears healthy:\n"))
		w.Write(body)
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("Exhibitor appears unhealthy:\n"))
		w.Write(body)
	}
}

func Cockroach(w http.ResponseWriter, r *http.Request) {
	// var snapshot MetricsSnapshotResponse
	// TODO This isn't actually currently JSON
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := client.Request("GET", 443, "cockroachdb/_status/vars")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Cockroachdb state not accessible`))
		return
	}

	logrus.Infoln(string(body))

	w.Write([]byte(`Cockroachdb state hmm accessible`))

	// if (snapshot.OverlayRecovered == 1.0 && snapshot.MesosRecovered == 1.0) {
	// 	w.Write([]byte(`Mesos and Overlay Replogs Healthy: {"registrar/log/recovered": 1.0, "overlay/log/recovered": 1.0}`))
	// } else {
	// 	w.WriteHeader(http.StatusExpectationFailed)
	// 	// TODO this should have the actual value, although I do think this is a binary.
	// 	w.Write([]byte(`Replog unhealthy`))
	// }
}