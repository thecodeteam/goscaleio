package goscaleio

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSystemGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/instances/System::788d9efb0a8f20cb":
				if !requestAuthOK(resp, req) {
					return
				}
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(jsonSystem))
			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)

	sys, err := client.GetSystemByID("788d9efb0a8f20cb")
	if err != nil {
		t.Fatal(err)
	}

	if sys == nil {
		t.Fatal("System is nil")
	}
	if sys.ID != "788d9efb0a8f20cb" {
		t.Fatalf("Expected system id %v, got %v", "788d9efb0a8f20cb", sys.ID)
	}
}

func TestSystemGetSystems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/types/System/instances":
				if !requestAuthOK(resp, req) {
					return
				}
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(jsonSystems))
			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)

	systems, err := client.GetSystems()
	if err != nil {
		t.Fatal(err)
	}

	if systems == nil {
		t.Fatal("Systems is nil")
	}
	if len(systems) != 2 {
		t.Fatalf("Expecting len(systems) = 2, got %d", len(systems))
	}
}

func TestSystemGetSystemsByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/types/System/instances":
				if !requestAuthOK(resp, req) {
					return
				}
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(jsonSystems))
			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)

	system, err := client.GetSystemByName("scaleio0")
	if err != nil {
		t.Fatal(err)
	}

	if system == nil {
		t.Fatal("Systems is nil")
	}
	if system.Name != "scaleio0" {
		t.Fatal("Unexpected system name:", system.Name)
	}
}

func TestSystemGetProtectionDomains(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/types/ProtectionDomain/instances":
				if !requestAuthOK(resp, req) {
					return
				}
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(jsonProtectionDomains))
			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)

	domains, err := client.GetProtectionDomains()
	if err != nil {
		t.Fatal(err)
	}

	if domains == nil {
		t.Fatal("Protection domains is nil")
	}
	if len(domains) != 1 {
		t.Fatalf("Expecting len(systems) = 1, got %d", len(domains))
	}
}

func TestSystemGetProtectionDomainID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/instances/ProtectionDomain::7042970f00000000":
				if !requestAuthOK(resp, req) {
					return
				}

				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(jsonProtectionDomain))

			case "/api/instances/ProtectionDomain::bad-vol-id":
				if !requestAuthOK(resp, req) {
					return
				}
				resp.WriteHeader(http.StatusInternalServerError)
				resp.Write([]byte(`{
  					"message": "Could not find Protection Domain",
  					"httpStatusCode": 500,
  					"errorCode": 142
				}`))
				resp.WriteHeader(http.StatusInternalServerError)
				resp.Write([]byte(jsonVolume))

			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)
	domain, err := client.GetProtectionDomainByID("7042970f00000000")
	if err != nil {
		t.Fatal(err)
	}

	if domain.ID != "7042970f00000000" {
		t.Fatal("Expecting a protection domain id 7042970f00000000, got ", domain.ID)
	}

	_, err = client.GetProtectionDomainByID("bad-vol-id")
	if err == nil {
		t.Fatal("Expecting error here, did not get it")
	}
}

var jsonSystem = `
{
  "defaultIsVolumeObfuscated": false,
  "restrictedSdcModeEnabled": false,
  "capacityTimeLeftInDays": "Unlimited",
  "enterpriseFeaturesEnabled": true,
  "isInitialLicense": true,
  "swid": "",
  "daysInstalled": 5,
  "maxCapacityInGb": "Unlimited",
  "systemVersionName": "EMC ScaleIO Version: R2_0.6035.0",
  "installId": "68eb1ae058e2b67b",
  "capacityAlertHighThresholdPercent": 80,
  "capacityAlertCriticalThresholdPercent": 90,
  "remoteReadOnlyLimitState": false,
  "upgradeState": "NoUpgrade",
  "performanceParameters": {
    "perfProfile": "HighPerformance",
    "mdmNumberSdcReceiveUmt": 5,
    "mdmNumberSdsReceiveUmt": 10,
    "mdmNumberSdsSendUmt": 10,
    "mdmNumberSdsKeepaliveReceiveUmt": 10,
    "mdmSdsCapacityCountersUpdateInterval": 1,
    "mdmSdsCapacityCountersPollingInterval": 5,
    "mdmSdsVolumeSizePollingInterval": 15,
    "mdmSdsVolumeSizePollingRetryInterval": 5,
    "mdmNumberSdsTasksUmt": 1024,
    "mdmInitialSdsSnapshotCapacity": 1024,
    "mdmSdsSnapshotCapacityChunkSize": 5120,
    "mdmSdsSnapshotUsedCapacityThreshold": 50,
    "mdmSdsSnapshotFreeCapacityThreshold": 200
  },
  "currentProfilePerformanceParameters": {
    "perfProfile": null,
    "mdmNumberSdcReceiveUmt": 5,
    "mdmNumberSdsReceiveUmt": 10,
    "mdmNumberSdsSendUmt": 10,
    "mdmNumberSdsKeepaliveReceiveUmt": 10,
    "mdmSdsCapacityCountersUpdateInterval": 1,
    "mdmSdsCapacityCountersPollingInterval": 5,
    "mdmSdsVolumeSizePollingInterval": 15,
    "mdmSdsVolumeSizePollingRetryInterval": 5,
    "mdmNumberSdsTasksUmt": 1024,
    "mdmInitialSdsSnapshotCapacity": 1024,
    "mdmSdsSnapshotCapacityChunkSize": 5120,
    "mdmSdsSnapshotUsedCapacityThreshold": 50,
    "mdmSdsSnapshotFreeCapacityThreshold": 200
  },
  "sdcMdmNetworkDisconnectionsCounterParameters": {
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "sdcSdsNetworkDisconnectionsCounterParameters": {
    "shortWindow": {
      "threshold": 800,
      "windowSizeInSec": 60
    },
    "mediumWindow": {
      "threshold": 4000,
      "windowSizeInSec": 3600
    },
    "longWindow": {
      "threshold": 20000,
      "windowSizeInSec": 86400
    }
  },
  "sdcMemoryAllocationFailuresCounterParameters": {
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "sdcSocketAllocationFailuresCounterParameters": {
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "sdcLongOperationsCounterParameters": {
    "shortWindow": {
      "threshold": 10000,
      "windowSizeInSec": 60
    },
    "mediumWindow": {
      "threshold": 100000,
      "windowSizeInSec": 3600
    },
    "longWindow": {
      "threshold": 1000000,
      "windowSizeInSec": 86400
    }
  },
  "cliPasswordAllowed": true,
  "managementClientSecureCommunicationEnabled": true,
  "tlsVersion": "TLSv1.2",
  "showGuid": true,
  "authenticationMethod": "Native",
  "mdmToSdsPolicy": "None",
  "mdmCluster": {
    "clusterState": "ClusteredNormal",
    "tieBreakers": [
      {
        "managementIPs": [
          "192.168.99.201"
        ],
        "versionInfo": "R2_0.6035.0",
        "ips": [
          "192.168.99.201"
        ],
        "role": "TieBreaker",
        "status": "Normal",
        "name": "tb",
        "id": "6f9bac2315975de1",
        "port": 9011
      }
    ],
    "goodNodesNum": 3,
    "goodReplicasNum": 2,
    "clusterMode": "ThreeNodes",
    "master": {
      "managementIPs": [
        "192.168.99.202"
      ],
      "versionInfo": "R2_0.6035.0",
      "ips": [
        "192.168.99.202"
      ],
      "role": "Manager",
      "name": "mdm1",
      "id": "4cb067b24847d420",
      "port": 9011
    },
    "slaves": [
      {
        "managementIPs": [
          "192.168.99.203"
        ],
        "versionInfo": "R2_0.6035.0",
        "ips": [
          "192.168.99.203"
        ],
        "role": "Manager",
        "status": "Normal",
        "name": "mdm2",
        "id": "08a2e6a00dfaf432",
        "port": 9011
      }
    ],
    "name": "scaleio",
    "id": "8686774057318686923"
  },
  "name": "scaleio",
  "id": "788d9efb0a8f20cb",
  "links": [
    {
      "rel": "self",
      "href": "/api/instances/System::788d9efb0a8f20cb"
    },
    {
      "rel": "/api/System/relationship/Statistics",
      "href": "/api/instances/System::788d9efb0a8f20cb/relationships/Statistics"
    },
    {
      "rel": "/api/System/relationship/ProtectionDomain",
      "href": "/api/instances/System::788d9efb0a8f20cb/relationships/ProtectionDomain"
    },
    {
      "rel": "/api/System/relationship/Sdc",
      "href": "/api/instances/System::788d9efb0a8f20cb/relationships/Sdc"
    },
    {
      "rel": "/api/System/relationship/User",
      "href": "/api/instances/System::788d9efb0a8f20cb/relationships/User"
    }
  ]
}
`

var jsonSystems = `
[
  {
    "defaultIsVolumeObfuscated": false,
    "restrictedSdcModeEnabled": false,
    "capacityTimeLeftInDays": "Unlimited",
    "enterpriseFeaturesEnabled": true,
    "isInitialLicense": true,
    "swid": "",
    "daysInstalled": 5,
    "maxCapacityInGb": "Unlimited",
    "systemVersionName": "EMC ScaleIO Version: R2_0.6035.0",
    "installId": "68eb1ae058e2b67b",
    "capacityAlertHighThresholdPercent": 80,
    "capacityAlertCriticalThresholdPercent": 90,
    "remoteReadOnlyLimitState": false,
    "upgradeState": "NoUpgrade",
    "performanceParameters": {
      "perfProfile": "HighPerformance",
      "mdmNumberSdcReceiveUmt": 5,
      "mdmNumberSdsReceiveUmt": 10,
      "mdmNumberSdsSendUmt": 10,
      "mdmNumberSdsKeepaliveReceiveUmt": 10,
      "mdmSdsCapacityCountersUpdateInterval": 1,
      "mdmSdsCapacityCountersPollingInterval": 5,
      "mdmSdsVolumeSizePollingInterval": 15,
      "mdmSdsVolumeSizePollingRetryInterval": 5,
      "mdmNumberSdsTasksUmt": 1024,
      "mdmInitialSdsSnapshotCapacity": 1024,
      "mdmSdsSnapshotCapacityChunkSize": 5120,
      "mdmSdsSnapshotUsedCapacityThreshold": 50,
      "mdmSdsSnapshotFreeCapacityThreshold": 200
    },
    "currentProfilePerformanceParameters": {
      "perfProfile": null,
      "mdmNumberSdcReceiveUmt": 5,
      "mdmNumberSdsReceiveUmt": 10,
      "mdmNumberSdsSendUmt": 10,
      "mdmNumberSdsKeepaliveReceiveUmt": 10,
      "mdmSdsCapacityCountersUpdateInterval": 1,
      "mdmSdsCapacityCountersPollingInterval": 5,
      "mdmSdsVolumeSizePollingInterval": 15,
      "mdmSdsVolumeSizePollingRetryInterval": 5,
      "mdmNumberSdsTasksUmt": 1024,
      "mdmInitialSdsSnapshotCapacity": 1024,
      "mdmSdsSnapshotCapacityChunkSize": 5120,
      "mdmSdsSnapshotUsedCapacityThreshold": 50,
      "mdmSdsSnapshotFreeCapacityThreshold": 200
    },
    "sdcMdmNetworkDisconnectionsCounterParameters": {
      "shortWindow": {
        "threshold": 300,
        "windowSizeInSec": 60
      },
      "mediumWindow": {
        "threshold": 500,
        "windowSizeInSec": 3600
      },
      "longWindow": {
        "threshold": 700,
        "windowSizeInSec": 86400
      }
    },
    "sdcSdsNetworkDisconnectionsCounterParameters": {
      "shortWindow": {
        "threshold": 800,
        "windowSizeInSec": 60
      },
      "mediumWindow": {
        "threshold": 4000,
        "windowSizeInSec": 3600
      },
      "longWindow": {
        "threshold": 20000,
        "windowSizeInSec": 86400
      }
    },
    "sdcMemoryAllocationFailuresCounterParameters": {
      "shortWindow": {
        "threshold": 300,
        "windowSizeInSec": 60
      },
      "mediumWindow": {
        "threshold": 500,
        "windowSizeInSec": 3600
      },
      "longWindow": {
        "threshold": 700,
        "windowSizeInSec": 86400
      }
    },
    "sdcSocketAllocationFailuresCounterParameters": {
      "shortWindow": {
        "threshold": 300,
        "windowSizeInSec": 60
      },
      "mediumWindow": {
        "threshold": 500,
        "windowSizeInSec": 3600
      },
      "longWindow": {
        "threshold": 700,
        "windowSizeInSec": 86400
      }
    },
    "sdcLongOperationsCounterParameters": {
      "shortWindow": {
        "threshold": 10000,
        "windowSizeInSec": 60
      },
      "mediumWindow": {
        "threshold": 100000,
        "windowSizeInSec": 3600
      },
      "longWindow": {
        "threshold": 1000000,
        "windowSizeInSec": 86400
      }
    },
    "cliPasswordAllowed": true,
    "managementClientSecureCommunicationEnabled": true,
    "tlsVersion": "TLSv1.2",
    "showGuid": true,
    "authenticationMethod": "Native",
    "mdmToSdsPolicy": "None",
    "mdmCluster": {
      "clusterState": "ClusteredNormal",
      "tieBreakers": [
        {
          "managementIPs": [
            "192.168.99.201"
          ],
          "versionInfo": "R2_0.6035.0",
          "ips": [
            "192.168.99.201"
          ],
          "role": "TieBreaker",
          "status": "Normal",
          "name": "tb",
          "id": "6f9bac2315975de1",
          "port": 9011
        }
      ],
      "goodNodesNum": 3,
      "goodReplicasNum": 2,
      "clusterMode": "ThreeNodes",
      "master": {
        "managementIPs": [
          "192.168.99.202"
        ],
        "versionInfo": "R2_0.6035.0",
        "ips": [
          "192.168.99.202"
        ],
        "role": "Manager",
        "name": "mdm1",
        "id": "4cb067b24847d420",
        "port": 9011
      },
      "slaves": [
        {
          "managementIPs": [
            "192.168.99.203"
          ],
          "versionInfo": "R2_0.6035.0",
          "ips": [
            "192.168.99.203"
          ],
          "role": "Manager",
          "status": "Normal",
          "name": "mdm2",
          "id": "08a2e6a00dfaf432",
          "port": 9011
        }
      ],
      "name": "scaleio",
      "id": "8686774057318686923"
    },
    "name": "scaleio0",
    "id": "788d9efb0a8f20cb",
    "links": [
      {
        "rel": "self",
        "href": "/api/instances/System::788d9efb0a8f20cb"
      },
      {
        "rel": "/api/System/relationship/Statistics",
        "href": "/api/instances/System::788d9efb0a8f20cb/relationships/Statistics"
      },
      {
        "rel": "/api/System/relationship/ProtectionDomain",
        "href": "/api/instances/System::788d9efb0a8f20cb/relationships/ProtectionDomain"
      },
      {
        "rel": "/api/System/relationship/Sdc",
        "href": "/api/instances/System::788d9efb0a8f20cb/relationships/Sdc"
      },
      {
        "rel": "/api/System/relationship/User",
        "href": "/api/instances/System::788d9efb0a8f20cb/relationships/User"
      }
    ]
},

{
  "defaultIsVolumeObfuscated": false,
  "restrictedSdcModeEnabled": false,
  "capacityTimeLeftInDays": "Unlimited",
  "enterpriseFeaturesEnabled": true,
  "isInitialLicense": true,
  "swid": "",
  "daysInstalled": 5,
  "maxCapacityInGb": "Unlimited",
  "systemVersionName": "EMC ScaleIO Version: R2_0.6035.0",
  "installId": "68eb1ae058e2b67b",
  "capacityAlertHighThresholdPercent": 80,
  "capacityAlertCriticalThresholdPercent": 90,
  "remoteReadOnlyLimitState": false,
  "upgradeState": "NoUpgrade",
  "performanceParameters": {
    "perfProfile": "HighPerformance",
    "mdmNumberSdcReceiveUmt": 5,
    "mdmNumberSdsReceiveUmt": 10,
    "mdmNumberSdsSendUmt": 10,
    "mdmNumberSdsKeepaliveReceiveUmt": 10,
    "mdmSdsCapacityCountersUpdateInterval": 1,
    "mdmSdsCapacityCountersPollingInterval": 5,
    "mdmSdsVolumeSizePollingInterval": 15,
    "mdmSdsVolumeSizePollingRetryInterval": 5,
    "mdmNumberSdsTasksUmt": 1024,
    "mdmInitialSdsSnapshotCapacity": 1024,
    "mdmSdsSnapshotCapacityChunkSize": 5120,
    "mdmSdsSnapshotUsedCapacityThreshold": 50,
    "mdmSdsSnapshotFreeCapacityThreshold": 200
  },
  "currentProfilePerformanceParameters": {
    "perfProfile": null,
    "mdmNumberSdcReceiveUmt": 5,
    "mdmNumberSdsReceiveUmt": 10,
    "mdmNumberSdsSendUmt": 10,
    "mdmNumberSdsKeepaliveReceiveUmt": 10,
    "mdmSdsCapacityCountersUpdateInterval": 1,
    "mdmSdsCapacityCountersPollingInterval": 5,
    "mdmSdsVolumeSizePollingInterval": 15,
    "mdmSdsVolumeSizePollingRetryInterval": 5,
    "mdmNumberSdsTasksUmt": 1024,
    "mdmInitialSdsSnapshotCapacity": 1024,
    "mdmSdsSnapshotCapacityChunkSize": 5120,
    "mdmSdsSnapshotUsedCapacityThreshold": 50,
    "mdmSdsSnapshotFreeCapacityThreshold": 200
  },
  "sdcMdmNetworkDisconnectionsCounterParameters": {
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "sdcSdsNetworkDisconnectionsCounterParameters": {
    "shortWindow": {
      "threshold": 800,
      "windowSizeInSec": 60
    },
    "mediumWindow": {
      "threshold": 4000,
      "windowSizeInSec": 3600
    },
    "longWindow": {
      "threshold": 20000,
      "windowSizeInSec": 86400
    }
  },
  "sdcMemoryAllocationFailuresCounterParameters": {
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "sdcSocketAllocationFailuresCounterParameters": {
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "sdcLongOperationsCounterParameters": {
    "shortWindow": {
      "threshold": 10000,
      "windowSizeInSec": 60
    },
    "mediumWindow": {
      "threshold": 100000,
      "windowSizeInSec": 3600
    },
    "longWindow": {
      "threshold": 1000000,
      "windowSizeInSec": 86400
    }
  },
  "cliPasswordAllowed": true,
  "managementClientSecureCommunicationEnabled": true,
  "tlsVersion": "TLSv1.2",
  "showGuid": true,
  "authenticationMethod": "Native",
  "mdmToSdsPolicy": "None",
  "mdmCluster": {
    "clusterState": "ClusteredNormal",
    "tieBreakers": [
      {
        "managementIPs": [
          "192.168.99.201"
        ],
        "versionInfo": "R2_0.6035.0",
        "ips": [
          "192.168.99.201"
        ],
        "role": "TieBreaker",
        "status": "Normal",
        "name": "tb",
        "id": "6f9bac2315975de1",
        "port": 9011
      }
    ],
    "goodNodesNum": 3,
    "goodReplicasNum": 2,
    "clusterMode": "ThreeNodes",
    "master": {
      "managementIPs": [
        "192.168.99.202"
      ],
      "versionInfo": "R2_0.6035.0",
      "ips": [
        "192.168.99.202"
      ],
      "role": "Manager",
      "name": "mdm1",
      "id": "4cb067b24847d420",
      "port": 9011
    },
    "slaves": [
      {
        "managementIPs": [
          "192.168.99.203"
        ],
        "versionInfo": "R2_0.6035.0",
        "ips": [
          "192.168.99.203"
        ],
        "role": "Manager",
        "status": "Normal",
        "name": "mdm2",
        "id": "08a2e6a00dfaf432",
        "port": 9011
      }
    ],
    "name": "scaleio",
    "id": "8686774057318686923"
  },
  "name": "scaleio1",
  "id": "788d9efb0a8f20cb",
  "links": [
    {
      "rel": "self",
      "href": "/api/instances/System::788d9efb0a8f20cb"
    },
    {
      "rel": "/api/System/relationship/Statistics",
      "href": "/api/instances/System::788d9efb0a8f20cb/relationships/Statistics"
    },
    {
      "rel": "/api/System/relationship/ProtectionDomain",
      "href": "/api/instances/System::788d9efb0a8f20cb/relationships/ProtectionDomain"
    },
    {
      "rel": "/api/System/relationship/Sdc",
      "href": "/api/instances/System::788d9efb0a8f20cb/relationships/Sdc"
    },
    {
      "rel": "/api/System/relationship/User",
      "href": "/api/instances/System::788d9efb0a8f20cb/relationships/User"
    }
  ]
}
]
`

var jsonProtectionDomain = `
{
  "systemId": "788d9efb0a8f20cb",
  "protectionDomainState": "Active",
  "sdsDecoupledCounterParameters": {
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "rebuildNetworkThrottlingInKbps": null,
  "rebalanceNetworkThrottlingInKbps": null,
  "overallIoNetworkThrottlingInKbps": null,
  "sdsConfigurationFailureCounterParameters": {
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "mdmSdsNetworkDisconnectionsCounterParameters": {
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "sdsSdsNetworkDisconnectionsCounterParameters": {
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "sdsReceiveBufferAllocationFailuresCounterParameters": {
    "mediumWindow": {
      "threshold": 200000,
      "windowSizeInSec": 3600
    },
    "shortWindow": {
      "threshold": 20000,
      "windowSizeInSec": 60
    },
    "longWindow": {
      "threshold": 2000000,
      "windowSizeInSec": 86400
    }
  },
  "rebuildNetworkThrottlingEnabled": false,
  "rebalanceNetworkThrottlingEnabled": false,
  "overallIoNetworkThrottlingEnabled": false,
  "name": "default",
  "id": "7042970f00000000",
  "links": [
    {
      "rel": "self",
      "href": "/api/instances/ProtectionDomain::7042970f00000000"
    },
    {
      "rel": "/api/ProtectionDomain/relationship/Statistics",
      "href": "/api/instances/ProtectionDomain::7042970f00000000/relationships/Statistics"
    },
    {
      "rel": "/api/ProtectionDomain/relationship/StoragePool",
      "href": "/api/instances/ProtectionDomain::7042970f00000000/relationships/StoragePool"
    },
    {
      "rel": "/api/ProtectionDomain/relationship/Sds",
      "href": "/api/instances/ProtectionDomain::7042970f00000000/relationships/Sds"
    },
    {
      "rel": "/api/ProtectionDomain/relationship/FaultSet",
      "href": "/api/instances/ProtectionDomain::7042970f00000000/relationships/FaultSet"
    },
    {
      "rel": "/api/parent/relationship/systemId",
      "href": "/api/instances/System::788d9efb0a8f20cb"
    }
  ]
}
`

var jsonProtectionDomains = `
[
{
  "systemId": "788d9efb0a8f20cb",
  "protectionDomainState": "Active",
  "sdsDecoupledCounterParameters": {
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "rebuildNetworkThrottlingInKbps": null,
  "rebalanceNetworkThrottlingInKbps": null,
  "overallIoNetworkThrottlingInKbps": null,
  "sdsConfigurationFailureCounterParameters": {
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "mdmSdsNetworkDisconnectionsCounterParameters": {
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "sdsSdsNetworkDisconnectionsCounterParameters": {
    "mediumWindow": {
      "threshold": 500,
      "windowSizeInSec": 3600
    },
    "shortWindow": {
      "threshold": 300,
      "windowSizeInSec": 60
    },
    "longWindow": {
      "threshold": 700,
      "windowSizeInSec": 86400
    }
  },
  "sdsReceiveBufferAllocationFailuresCounterParameters": {
    "mediumWindow": {
      "threshold": 200000,
      "windowSizeInSec": 3600
    },
    "shortWindow": {
      "threshold": 20000,
      "windowSizeInSec": 60
    },
    "longWindow": {
      "threshold": 2000000,
      "windowSizeInSec": 86400
    }
  },
  "rebuildNetworkThrottlingEnabled": false,
  "rebalanceNetworkThrottlingEnabled": false,
  "overallIoNetworkThrottlingEnabled": false,
  "name": "default",
  "id": "7042970f00000000",
  "links": [
    {
      "rel": "self",
      "href": "/api/instances/ProtectionDomain::7042970f00000000"
    },
    {
      "rel": "/api/ProtectionDomain/relationship/Statistics",
      "href": "/api/instances/ProtectionDomain::7042970f00000000/relationships/Statistics"
    },
    {
      "rel": "/api/ProtectionDomain/relationship/StoragePool",
      "href": "/api/instances/ProtectionDomain::7042970f00000000/relationships/StoragePool"
    },
    {
      "rel": "/api/ProtectionDomain/relationship/Sds",
      "href": "/api/instances/ProtectionDomain::7042970f00000000/relationships/Sds"
    },
    {
      "rel": "/api/ProtectionDomain/relationship/FaultSet",
      "href": "/api/instances/ProtectionDomain::7042970f00000000/relationships/FaultSet"
    },
    {
      "rel": "/api/parent/relationship/systemId",
      "href": "/api/instances/System::788d9efb0a8f20cb"
    }
  ]
}
]
`
