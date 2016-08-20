package goscaleio

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/emccode/goscaleio/types/v1"
)

func TestVolumeGetID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/types/Volume/instances/action/queryIdByKey":
				if !requestAuthOK(resp, req) {
					return
				}

				var param = new(types.VolumeQeryIdByKeyParam)
				if err := json.NewDecoder(req.Body).Decode(param); err != nil {
					t.Fatalf("Failed to decode param %v", err)
				}
				if param.Name != "test-vol" {
					resp.WriteHeader(http.StatusInternalServerError)
					resp.Write([]byte(`{"message":"Not found","httpStatusCode":500,"errorCode":3}`))
					return
				}
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte("a2b7cc6300000000"))
			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)
	id, err := client.GetVolumeID("test-vol")
	if err != nil {
		t.Fatal(err)
	}

	if id != "a2b7cc6300000000" {
		t.Fatal("Expected volume id a2b7cc6300000000, but got ", id)
	}

	_, err = client.GetVolumeID("bad-vol-name")
	if err == nil {
		t.Fatal("Expected failure for bad volume name")
	}
}

func TestVolumeGetByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/instances/Volume::a2b7cc6300000000":
				if !requestAuthOK(resp, req) {
					return
				}

				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(jsonVolume))

			case "/api/instances/Volume::bad-vol-id":
				if !requestAuthOK(resp, req) {
					return
				}
				resp.WriteHeader(http.StatusInternalServerError)
				resp.Write([]byte(`{"message": "Could not find the volume",
                    "httpStatusCode": 500,
                    "errorCode": 79
                    }`))
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(jsonVolume))

			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)
	vol, err := client.GetVolumeByID("a2b7cc6300000000")
	if err != nil {
		t.Fatal(err)
	}

	if vol == nil {
		t.Fatal("Expecting a volume, but got nil")
	}

	if vol.ID != "a2b7cc6300000000" {
		t.Fatal("Expecting a volume id a2b7cc6300000000, got ", vol.ID)
	}

	_, err = client.GetVolumeByID("bad-vol-id")
	if err == nil {
		t.Fatal("Expecting error here, did not get it")
	}
}

func TestVolumeGetVolumes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/types/Volume/instances":
				if !requestAuthOK(resp, req) {
					return
				}

				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(jsonVolumes))
			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)
	vols, err := client.GetVolumes()
	if err != nil {
		t.Fatal(err)
	}

	if vols == nil {
		t.Fatal("Expecting volumes, but got nil")
	}
	if len(vols) != 3 {
		t.Fatal("Expected 3 volumes, got ", len(vols))
	}
}

func TestVolumeGetVolumesByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/types/Volume/instances":
				if !requestAuthOK(resp, req) {
					return
				}

				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(jsonVolumes))
			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)
	vol, err := client.GetVolumesByName("test-vol-001")
	if err != nil {
		t.Fatal(err)
	}

	if vol == nil {
		t.Fatal("Expecting volume, but got nil")
	}
	if vol.Name != "test-vol-001" {
		t.Fatal("Unexpected volume data: ", vol.Name)
	}
}

func TestVolumeGetByStoragePoolID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/instances/StoragePool::800b247900000000/relationships/Volume":
				if !requestAuthOK(resp, req) {
					return
				}

				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(jsonVolumes))

			case "/api/instances/StoragePool::bad-storage-pool-id/relationships/Volume":
				if !requestAuthOK(resp, req) {
					return
				}

				resp.WriteHeader(http.StatusInternalServerError)
				resp.Write([]byte(`{
                    "message": "Error in get relationship Volume",
                    "httpStatusCode": 500,
                    "errorCode": 0
                    }
                `))
			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)
	vols, err := client.GetVolumesByStoragePoolID("800b247900000000")
	if err != nil {
		t.Fatal(err)
	}

	if vols == nil {
		t.Fatal("Expecting volumes, but got nil")
	}
	if len(vols) != 3 {
		t.Fatal("Expected 3 volumes, got ", len(vols))
	}

	_, err = client.GetVolumesByStoragePoolID("bad-storage-pool-id")
	if err == nil {
		t.Fatal("Expected error for bad storage pool id")
	}
}

func TestVolumeCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/types/Volume/instances":
				if !requestAuthOK(resp, req) {
					return
				}

				if req.Method != "POST" {
					t.Fatal("Expecting http method POST, got ", req.Method)
				}

				contentType := req.Header.Get("Content-Type")
				if contentType != "application/json;version=2.0" {
					t.Fatal("Expecting content-type=application/json;version=2.0, got", contentType)
				}

				var volParam *types.VolumeParam
				if err := json.NewDecoder(req.Body).Decode(&volParam); err != nil {
					t.Fatal(err)
				}
				if volParam.VolumeSizeInKb == "" {
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{
                        "message": "Request message is not valid: The following parameter(s) must be part of the request body: volumeSizeInKb.",
                        "httpStatusCode": 400,
                        "errorCode": 0
                        }
                    `))
					return
				}
				if volParam.StoragePoolID == "" {
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{
                        "message": "Request message is not valid: The following parameter(s) must be part of the request body: storagePoolId.",
                        "httpStatusCode": 400,
                        "errorCode": 0
                        }
                    `))
					return
				}

				data, err := json.Marshal(types.VolumeResp{ID: "a2b7cc6700000004"})
				if err != nil {
					t.Fatal(err)
				}
				resp.WriteHeader(http.StatusOK)
				resp.Write(data)

			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}

		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)
	id, err := client.CreateVolume(&types.VolumeParam{
		Name:           "test-vol",
		VolumeSizeInKb: "1000000000",
		StoragePoolID:  "800b247900000000",
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "a2b7cc6700000004" {
		t.Fatal("Unexpected ID, expecting a2b7cc6700000004, got ", id)
	}

	_, err = client.CreateVolume(&types.VolumeParam{
		Name: "test-vol",
	})
	if err == nil {
		t.Fatal("Expecting error for missing VolumeParam attributes.")
	}
}

func TestVolumeRemove(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))

			case "/api/login":
				handleAuthToken(resp, req)

			case "/api/instances/Volume::a2b7cc6700000004/action/removeVolume":
				if !requestAuthOK(resp, req) {
					return
				}

				if req.Method != "POST" {
					t.Fatal("Expecting http method POST, got ", req.Method)
				}

				var param *types.RemoveVolumeParam
				if err := json.NewDecoder(req.Body).Decode(&param); err != nil {
					t.Fatal(err)
				}
				if param.RemoveMode != RemoveMode.OnlyMe {
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{
                        "message": "removeMode should get one of the following values: ONLY_ME, INCLUDING_DESCENDANTS, DESCENDANTS_ONLY, WHOLE_VTREE, but its value is Uknown.",
                        "httpStatusCode": 400,
                        "errorCode": 0
                        }
                    `))
					return
				}

				resp.WriteHeader(http.StatusOK)

			case "/api/instances/Volume::bad-vol-id/action/removeVolume":
				if !requestAuthOK(resp, req) {
					return
				}
				resp.WriteHeader(http.StatusInternalServerError)
				resp.Write([]byte(`{
                    "message": "Could not find the volume",
                    "httpStatusCode": 500,
                    "errorCode": 79
                    }
                `))
				return
			default:
				t.Fatal("Unexpected RequestURI", req.RequestURI)
			}
		},
	))
	defer server.Close()

	client := setupClient(t, server.URL)
	err := client.RemoveVolume("a2b7cc6700000004", RemoveMode.OnlyMe)
	if err != nil {
		t.Fatal(err)
	}

	err = client.RemoveVolume("bad-vol-id", RemoveMode.OnlyMe)
	if err == nil {
		t.Fatal("Expected error here for bad volume id")
	}

}

var jsonVolume = `
{
  "mappingToAllSdcsEnabled": false,
  "mappedSdcInfo": null,
  "creationTime": 1470105021,
  "useRmcache": false,
  "isObfuscated": false,
  "volumeType": "ThinProvisioned",
  "consistencyGroupId": null,
  "vtreeId": "1d5e2c0100000000",
  "ancestorVolumeId": null,
  "storagePoolId": "800b247900000000",
  "sizeInKb": 8388608,
  "name": "local-vol",
  "id": "a2b7cc6300000000",
  "links": [
    {
      "rel": "self",
      "href": "/api/instances/Volume::a2b7cc6300000000"
    },
    {
      "rel": "/api/Volume/relationship/Statistics",
      "href": "/api/instances/Volume::a2b7cc6300000000/relationships/Statistics"
    },
    {
      "rel": "/api/parent/relationship/vtreeId",
      "href": "/api/instances/VTree::1d5e2c0100000000"
    },
    {
      "rel": "/api/parent/relationship/storagePoolId",
      "href": "/api/instances/StoragePool::800b247900000000"
    }
  ]
}
`
var jsonVolumes = `
[
  {
    "mappingToAllSdcsEnabled": false,
    "mappedSdcInfo": null,
    "creationTime": 1470105021,
    "useRmcache": false,
    "isObfuscated": false,
    "volumeType": "ThinProvisioned",
    "consistencyGroupId": null,
    "vtreeId": "1d5e2c0100000000",
    "ancestorVolumeId": null,
    "storagePoolId": "800b247900000000",
    "sizeInKb": 8388608,
    "name": "local-vol",
    "id": "a2b7cc6300000000",
    "links": [
      {
        "rel": "self",
        "href": "/api/instances/Volume::a2b7cc6300000000"
      },
      {
        "rel": "/api/Volume/relationship/Statistics",
        "href": "/api/instances/Volume::a2b7cc6300000000/relationships/Statistics"
      },
      {
        "rel": "/api/parent/relationship/vtreeId",
        "href": "/api/instances/VTree::1d5e2c0100000000"
      },
      {
        "rel": "/api/parent/relationship/storagePoolId",
        "href": "/api/instances/StoragePool::800b247900000000"
      }
    ]
  },
  {
    "mappingToAllSdcsEnabled": false,
    "mappedSdcInfo": null,
    "creationTime": 1470105375,
    "useRmcache": true,
    "isObfuscated": false,
    "volumeType": "ThickProvisioned",
    "consistencyGroupId": null,
    "vtreeId": "1d5e2c0200000001",
    "ancestorVolumeId": null,
    "storagePoolId": "800b247900000000",
    "sizeInKb": 8388608,
    "name": "test-vol-001",
    "id": "a2b7cc6400000001",
    "links": [
      {
        "rel": "self",
        "href": "/api/instances/Volume::a2b7cc6400000001"
      },
      {
        "rel": "/api/Volume/relationship/Statistics",
        "href": "/api/instances/Volume::a2b7cc6400000001/relationships/Statistics"
      },
      {
        "rel": "/api/parent/relationship/vtreeId",
        "href": "/api/instances/VTree::1d5e2c0200000001"
      },
      {
        "rel": "/api/parent/relationship/storagePoolId",
        "href": "/api/instances/StoragePool::800b247900000000"
      }
    ]
  },
  {
    "mappingToAllSdcsEnabled": false,
    "mappedSdcInfo": null,
    "creationTime": 1470112764,
    "useRmcache": false,
    "isObfuscated": false,
    "volumeType": "Snapshot",
    "consistencyGroupId": "77fc9ce100000001",
    "vtreeId": "1d5e2c0200000001",
    "ancestorVolumeId": "a2b7cc6400000001",
    "storagePoolId": "800b247900000000",
    "sizeInKb": 8388608,
    "name": "test-vol-001-snap",
    "id": "a2b7cc6500000002",
    "links": [
      {
        "rel": "self",
        "href": "/api/instances/Volume::a2b7cc6500000002"
      },
      {
        "rel": "/api/Volume/relationship/Statistics",
        "href": "/api/instances/Volume::a2b7cc6500000002/relationships/Statistics"
      },
      {
        "rel": "/api/parent/relationship/ancestorVolumeId",
        "href": "/api/instances/Volume::a2b7cc6400000001"
      },
      {
        "rel": "/api/parent/relationship/vtreeId",
        "href": "/api/instances/VTree::1d5e2c0200000001"
      },
      {
        "rel": "/api/parent/relationship/storagePoolId",
        "href": "/api/instances/StoragePool::800b247900000000"
      }
    ]
  }
]
`
