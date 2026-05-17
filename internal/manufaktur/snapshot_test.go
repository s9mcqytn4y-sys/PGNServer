package manufaktur

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTanganiSnapshotMasterData_Sukses(t *testing.T) {
	// Set Gin ke mode pengetesan
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockRepositoriManufaktur)
	layanan := KonstruksiLayananBaru(mockRepo)
	handler := KonstruksiPenangananBaru(layanan)

	// Persiapkan Gin Engine
	router := gin.New()
	router.Use(gin.Recovery())

	// Daftarkan endpoint snapshot
	router.GET("/api/v1/master-data/snapshot", handler.TanganiSnapshotMasterData)

	// Data mock material dari repository
	materialsMock := []Material{
		{
			KodeSKU:      "SKU-MAT-001",
			NamaMaterial: "Baja Lembaran Galvanis",
		},
		{
			KodeSKU:      "SKU-MAT-002",
			NamaMaterial: "Tembaga Murni Rod",
		},
	}
	materialsMock[0].ID = 101
	materialsMock[1].ID = 102

	// Setup expectation pada mock repository
	mockRepo.On("AmbilSnapshotMasterData").Return(materialsMock, nil).Once()

	// Lakukan HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/master-data/snapshot", nil)
	router.ServeHTTP(w, req)

	// Verifikasi response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal respons
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// 1. Verifikasi Double-Envelope Keamanan API (sukses & success)
	assert.True(t, resp["sukses"].(bool))
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "Snapshot master data berhasil diambil", resp["pesan"])
	assert.Equal(t, "Snapshot master data berhasil diambil", resp["message"])

	// 2. Verifikasi Struktur Data Snapshot
	assert.NotNil(t, resp["data"])
	dataMap, ok := resp["data"].(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, "v1.0.0", dataMap["versiMasterData"])

	// Verifikasi material array
	materialsVal, hasMaterials := dataMap["material"]
	assert.True(t, hasMaterials)
	materialsList, isSlice := materialsVal.([]interface{})
	assert.True(t, isSlice)
	assert.Len(t, materialsList, 2)

	// Check detail material pertama
	m1 := materialsList[0].(map[string]interface{})
	assert.Equal(t, "101", m1["id"])
	assert.Equal(t, "SKU-MAT-001", m1["kodeSKU"])
	assert.Equal(t, "Baja Lembaran Galvanis", m1["namaMaterial"])
	assert.True(t, m1["aktif"].(bool))

	// Check detail material kedua
	m2 := materialsList[1].(map[string]interface{})
	assert.Equal(t, "102", m2["id"])
	assert.Equal(t, "SKU-MAT-002", m2["kodeSKU"])
	assert.Equal(t, "Tembaga Murni Rod", m2["namaMaterial"])
	assert.True(t, m2["aktif"].(bool))

	// 3. Verifikasi Keberadaan Metadata & Kesesuaian Nilai
	assert.NotNil(t, resp["metadata"])
	metaMap, okMeta := resp["metadata"].(map[string]interface{})
	assert.True(t, okMeta)

	assert.Equal(t, float64(2), metaMap["jumlahMaterial"])
	assert.Equal(t, float64(3), metaMap["jumlahShiftOperasional"])
	assert.Equal(t, float64(0), metaMap["jumlahLineProduksi"])

	// 4. Verifikasi Legacy metadata envelope (meta)
	assert.NotNil(t, resp["meta"])
	legacyMetaMap, okLegacyMeta := resp["meta"].(map[string]interface{})
	assert.True(t, okLegacyMeta)
	assert.Equal(t, float64(2), legacyMetaMap["jumlahMaterial"])

	mockRepo.AssertExpectations(t)
}
