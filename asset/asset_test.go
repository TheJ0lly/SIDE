package asset

import "testing"

func TestCreateNewAsset(t *testing.T) {
	a := CreateNewAsset("Test", []byte{1})

	if a == nil || a.mName != "Test" || a.mData[0] != 1 {
		t.Errorf("Error: Failed to create asset\n")
		t.Fail()
		return
	}
}

func TestAsset_GetAssetBytes(t *testing.T) {
	a := CreateNewAsset("Test", []byte{1})

	if a.GetAssetBytes()[0] != 1 {
		t.Errorf("Error: Failed to get asset bytes\n")
		t.Fail()
		return
	}
}

func TestAsset_GetName(t *testing.T) {
	a := CreateNewAsset("Test", []byte{1})

	if a.GetName() != "Test" {
		t.Errorf("Error: Failed to get asset name\n")
		t.Fail()
		return
	}
}

func TestAsset_GetAssetCopy(t *testing.T) {
	a := CreateNewAsset("Test", []byte{1})

	b := a.GetAssetCopy()

	if a.GetName() != b.GetName() && a.GetAssetBytes()[0] != b.GetAssetBytes()[0] {
		t.Errorf("Error: Failed to make a copy\n")
		t.Fail()
		return
	}
}

func TestAsset_GetAssetSize(t *testing.T) {
	a := CreateNewAsset("Test", []byte{1})

	if a.GetAssetSize() != 1 {
		t.Errorf("Error: Failed to get asset size\n")
		t.Fail()
		return
	}
}
